package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"

	"github.com/yandex-cloud/go-sdk/iamkey"
	"golang.org/x/net/context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"

	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"

	ycsdk "github.com/yandex-cloud/go-sdk"
	ycsdkoperation "github.com/yandex-cloud/go-sdk/operation"
)

const (
	apiTimeout = 10 * time.Minute
)

type SecretUIDToYcSDK map[uint64]*ycsdk.SDK

type YandexSessionStore struct {
	sync.Mutex
	sdk SecretUIDToYcSDK
}

func NewYandexSession() *YandexSessionStore {
	return &YandexSessionStore{
		sdk: make(SecretUIDToYcSDK),
	}
}

func (ys *YandexSessionStore) GetSession(ctx context.Context, secretData map[string][]byte) (*ycsdk.SDK, error) {
	ys.Lock()
	defer ys.Unlock()

	var serviceAccountKey iamkey.Key
	err := json.Unmarshal(secretData[v1alpha1.YandexServiceAccountJSON], &serviceAccountKey)
	if err != nil {
		return nil, err
	}

	storeKeyHash := fnv.New64a()
	_, err = storeKeyHash.Write([]byte(serviceAccountKey.String()))
	if err != nil {
		return nil, fmt.Errorf("error while calculating hash for the Yandex.Cloud session store key: %s", err)
	}

	storeKey := storeKeyHash.Sum64()

	if sdk, ok := ys.sdk[storeKey]; ok {
		return sdk, nil
	}

	creds, err := ycsdk.ServiceAccountKey(&serviceAccountKey)
	if err != nil {
		return nil, err
	}

	sdk, err := ycsdk.Build(ctx, ycsdk.Config{Credentials: creds})
	if err != nil {
		return nil, err
	}
	ys.sdk[storeKey] = sdk

	return sdk, nil
}

// TODO: would like to use this somewhere, but I've got some reservations about shutting down in-progress operations
func (ys *YandexSessionStore) GC() {
	ys.Lock()
	defer ys.Unlock()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()

	for k, v := range ys.sdk {
		_ = v.Shutdown(ctx)
		delete(ys.sdk, k)
	}
}

var yss = NewYandexSession()

type YandexDriver struct {
	YandexMachinesClass *v1alpha1.YandexMachineClass
	SecretData          map[string][]byte
	UserData            string
	MachineID           string
	MachineName         string
	svc                 *ycsdk.SDK
}

func NewYandexDriver(machineClass *v1alpha1.YandexMachineClass, secretData map[string][]byte, userData, machineID, machineName string) *YandexDriver {
	ctx, cancelFunc := context.WithTimeout(context.Background(), apiTimeout)
	defer cancelFunc()

	svc, err := yss.GetSession(ctx, secretData)
	if err != nil {
		panic(err)
	}

	return &YandexDriver{machineClass, secretData, userData, machineID, machineName, svc}
}

func (d *YandexDriver) Create() (string, string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), apiTimeout)
	defer cancelFunc()

	var instanceMetadata = make(map[string]string)
	if d.YandexMachinesClass.Spec.Metadata != nil {
		instanceMetadata = d.YandexMachinesClass.Spec.Metadata
	}

	instanceMetadata["user-data"] = d.UserData

	var networkInterfaceSpecs []*compute.NetworkInterfaceSpec
	for _, netIf := range d.YandexMachinesClass.Spec.NetworkInterfaceSpecs {
		var natSpec *compute.OneToOneNatSpec
		if netIf.AssignPublicIPAddress {
			natSpec = &compute.OneToOneNatSpec{
				IpVersion: compute.IpVersion_IPV4,
			}
		}

		var externalIPv4Address = &compute.PrimaryAddressSpec{
			OneToOneNatSpec: natSpec,
		}

		networkInterfaceSpecs = append(networkInterfaceSpecs, &compute.NetworkInterfaceSpec{
			SubnetId:             netIf.SubnetID,
			PrimaryV4AddressSpec: externalIPv4Address,
		})
	}

	createInstanceParams := &compute.CreateInstanceRequest{
		FolderId:   string(d.SecretData[v1alpha1.YandexFolderID]),
		Name:       d.MachineName,
		Hostname:   d.MachineName,
		Labels:     d.YandexMachinesClass.Spec.Labels,
		Metadata:   instanceMetadata,
		ZoneId:     d.YandexMachinesClass.Spec.ZoneID,
		PlatformId: d.YandexMachinesClass.Spec.PlatformID,
		ResourcesSpec: &compute.ResourcesSpec{
			Memory:       d.YandexMachinesClass.Spec.ResourcesSpec.Memory,
			Cores:        d.YandexMachinesClass.Spec.ResourcesSpec.Cores,
			CoreFraction: d.YandexMachinesClass.Spec.ResourcesSpec.CoreFraction,
			Gpus:         d.YandexMachinesClass.Spec.ResourcesSpec.GPUs,
		},
		BootDiskSpec: &compute.AttachedDiskSpec{
			Mode:       compute.AttachedDiskSpec_READ_WRITE,
			AutoDelete: d.YandexMachinesClass.Spec.BootDiskSpec.AutoDelete,
			Disk: &compute.AttachedDiskSpec_DiskSpec_{
				DiskSpec: &compute.AttachedDiskSpec_DiskSpec{
					TypeId: d.YandexMachinesClass.Spec.BootDiskSpec.TypeID,
					Size:   d.YandexMachinesClass.Spec.BootDiskSpec.Size,
					Source: &compute.AttachedDiskSpec_DiskSpec_ImageId{
						ImageId: d.YandexMachinesClass.Spec.BootDiskSpec.ImageID,
					},
				},
			},
		},
		NetworkInterfaceSpecs: networkInterfaceSpecs,
		SchedulingPolicy: &compute.SchedulingPolicy{
			Preemptible: d.YandexMachinesClass.Spec.SchedulingPolicy.Preemptible,
		},
	}

	result, _, err := waitForResult(ctx, d.svc, func() (*operation.Operation, error) {
		return d.svc.Compute().Instance().Create(ctx, createInstanceParams)
	})
	if err != nil {
		return "", "", err
	}

	newInstance, ok := result.(*compute.Instance)
	if !ok {
		return "", "", fmt.Errorf("Yandex.Cloud API returned %q instead of \"*compute.Instance\". That shouldn't happen", reflect.TypeOf(result).String())
	}

	return d.encodeMachineID(newInstance.Id), d.MachineName, nil
}

func (d *YandexDriver) encodeMachineID(instanceId string) string {
	return fmt.Sprintf("yandex://%s", instanceId)
}

var regExpProviderID = regexp.MustCompile(`^yandex://([^/]+)$`)
var unexpectedMachineIdErr = errors.New("unexpected machineID")

func (d *YandexDriver) decodeMachineID(machineID string) (string, error) {
	// providerID is in the following form "${providerName}://${folderID}/${zone}/${instanceName}"
	// So for input "yandex://b1g4c2a3g6vkffp3qacq/ru-central1-a/e2e-test-node0" output will be  "b1g4c2a3g6vkffp3qacq", "ru-central1-a", "e2e-test-node0".
	matches := regExpProviderID.FindStringSubmatch(machineID)
	if len(matches) != 2 {
		return "", fmt.Errorf("%w: %s", unexpectedMachineIdErr, machineID)
	}

	return matches[1], nil
}

func (d *YandexDriver) Delete(machineID string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), apiTimeout)
	defer cancelFunc()

	var instanceID string
	instanceID, err := d.decodeMachineID(machineID)
	// folder, _, instanceName, err := d.decodeMachineID(machineID)
	// TODO: Remove migration from older Yandex.Cloud InstanceID-based MachineID
	if err != nil {
		return err
	}

	_, _, err = waitForResult(ctx, d.svc, func() (*operation.Operation, error) {
		return d.svc.Compute().Instance().Delete(ctx, &compute.DeleteInstanceRequest{InstanceId: instanceID})
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			klog.V(2).Infof("No machine matching the machineID %q found on the provider: %s", d.MachineID, err)
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (d *YandexDriver) GetExisting() (string, error) {
	return d.MachineID, nil
}

// GetVMs returns a machine matching the machineID
// If machineID is an empty string then it returns all matching instances
func (d *YandexDriver) GetVMs(machineID string) (VMs, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), apiTimeout)
	defer cancelFunc()

	listOfVMs := make(map[string]string)

	clusterName := ""
	nodeRole := ""

	for key := range d.YandexMachinesClass.Spec.Labels {
		if strings.Contains(key, "kubernetes-io-cluster-") {
			clusterName = key
		} else if strings.Contains(key, "kubernetes-io-role-") {
			nodeRole = key
		}
	}

	if clusterName == "" || nodeRole == "" {
		return listOfVMs, nil
	}

	if machineID == "" {
		// TODO: Replace this abomination with something more Go-like
		var instances []*compute.Instance
		instanceIterator := d.svc.Compute().Instance().InstanceIterator(ctx, string(d.SecretData[v1alpha1.YandexFolderID]))
	instanceIteration:
		for {
			next := instanceIterator.Next()
			switch next {
			case true:
				instances = append(instances, instanceIterator.Value())
			case false:
				break instanceIteration
			}
		}
		if instanceIterator.Error() != nil {
			return nil, fmt.Errorf("could not list instances for FolderID %q: %s", d.SecretData[v1alpha1.YandexFolderID], instanceIterator.Error())
		}

		for _, instance := range instances {
			matchedCluster := false
			matchedRole := false
			for label := range instance.Labels {
				switch label {
				case clusterName:
					matchedCluster = true
				case nodeRole:
					matchedRole = true
				}
			}
			if matchedCluster && matchedRole {
				listOfVMs[d.encodeMachineID(instance.Id)] = instance.Name
			}
		}
	} else {
		instanceID, err := d.decodeMachineID(machineID)
		if err != nil {
			return nil, err
		}

		instance, err := d.svc.Compute().Instance().Get(ctx, &compute.GetInstanceRequest{InstanceId: instanceID})
		if err != nil {
			return nil, err
		}

		if instance != nil {
			listOfVMs[machineID] = instance.Name
		}
	}
	return listOfVMs, nil
}

// GetVolNames parses volume names from pv specs
func (d *YandexDriver) GetVolNames(specs []corev1.PersistentVolumeSpec) ([]string, error) {
	var names []string
	for _, spec := range specs {
		if spec.CSI == nil {
			// Not a CSI-managed volume
			continue
		}
		if spec.CSI.Driver != "yandex.csi.flant.com" {
			// Not a volume provisioned by Yandex.Cloud CSI driver
			continue
		}
		names = append(names, spec.CSI.VolumeHandle)
	}
	return names, nil
}

//GetUserData return the used data whit which the VM will be booted
func (d *YandexDriver) GetUserData() string {
	return d.UserData
}

//SetUserData set the used data whit which the VM will be booted
func (d *YandexDriver) SetUserData(userData string) {
	d.UserData = userData
}

func waitForResult(ctx context.Context, sdk *ycsdk.SDK, origFunc func() (*operation.Operation, error)) (proto.Message, *ycsdkoperation.Operation, error) {
	op, err := sdk.WrapOperation(origFunc())
	if err != nil {
		return nil, nil, err
	}

	err = op.Wait(ctx)
	if err != nil {
		return nil, op, err
	}

	resp, err := op.Response()
	if err != nil {
		return nil, op, err
	}

	return resp, op, nil
}
