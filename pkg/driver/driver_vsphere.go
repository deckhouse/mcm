// Package driver contains the cloud provider specific implementations to manage machines

package driver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/vmware/govmomi/task"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"k8s.io/klog"

	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vapi/rest"

	"github.com/vmware/govmomi/vim25/mo"

	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"

	"github.com/vmware/govmomi/object"

	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/view"
	corev1 "k8s.io/api/core/v1"
)

//noinspection GoUnusedConst
const (
	tagCategoryCardinalitySingle   string = "SINGLE"
	tagCategoryCardinalityMultiple string = "MULTIPLE"
	hwVersion                             = 15 // recommended in https://cloud-provider-vsphere.sigs.k8s.io/tutorials/kubernetes-on-vsphere-with-kubeadm.html
)

type VsphereInfo struct {
	sync.Mutex

	loginFunc  func(ctx context.Context) (*govmomi.Client, *rest.Client, func(), error)
	tagManager *TagManager
}

var vsphereInfo *VsphereInfo

// PacketDriver is the driver struct for holding Packet machine information
type VsphereDriver struct {
	// TODO: Init with persistent search results
	VsphereMachineClass *v1alpha1.VsphereMachineClass
	SecretData          map[string][]byte
	UserData            string
	MachineID           string
	MachineName         string
}

func NewVsphereDriver(machineClass *v1alpha1.VsphereMachineClass, secretData map[string][]byte, userData, machineID, machineName string) *VsphereDriver {
	if vsphereInfo == nil {
		vsphereInfo = &VsphereInfo{}
		vsphereInfo.Lock()
		defer vsphereInfo.Unlock()
		host, ok := secretData[v1alpha1.VsphereHost]
		if !ok {
			panic(fmt.Errorf("missing %s in secret", v1alpha1.VsphereHost))
		}
		username, ok := secretData[v1alpha1.VsphereUsername]
		if !ok {
			panic(fmt.Errorf("missing %s in secret", v1alpha1.VsphereUsername))
		}
		password, ok := secretData[v1alpha1.VspherePassword]
		if !ok {
			panic(fmt.Errorf("missing %s in secret", v1alpha1.VspherePassword))
		}
		insecureBytes, ok := secretData[v1alpha1.VsphereInsecure]
		if !ok {
			panic(fmt.Errorf("missing %s in secret", v1alpha1.VsphereInsecure))
		}
		insecure, err := strconv.ParseBool(string(insecureBytes))
		if err != nil {
			panic(fmt.Errorf("not bool: %s", insecureBytes))
		}
		regionTagCategoryName, ok := secretData[v1alpha1.VsphereRegionTagCategory]
		if !ok {
			panic(fmt.Errorf("missing %s in secret", v1alpha1.VsphereRegionTagCategory))
		}
		zoneTagCategoryName, ok := secretData[v1alpha1.VsphereZoneTagCategory]
		if !ok {
			panic(fmt.Errorf("missing %s in secret", v1alpha1.VsphereZoneTagCategory))
		}

		vsphereInfo.loginFunc = createVsphereClient(string(host), string(username), string(password), insecure)
		vsphereInfo.tagManager = NewTagManager(string(regionTagCategoryName), string(zoneTagCategoryName))
	}

	return &VsphereDriver{machineClass, secretData, userData, machineID, machineName}
}

func (d *VsphereDriver) Create() (string, string, error) {
	var machineID string
	var sideEffectsPresent bool

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, restClient, logoutFunc, err := vsphereInfo.loginFunc(ctx)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}
	defer logoutFunc()

	datacenter, err := vsphereInfo.tagManager.GetDcByRegion(ctx, client, restClient, d.VsphereMachineClass.Spec.Region)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	cluster, err := vsphereInfo.tagManager.GetClusterByZone(ctx, client, restClient, d.VsphereMachineClass.Spec.Zone)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	clusterDrsEnabled, err := drsEnabled(ctx, cluster)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	hosts, err := cluster.Hosts(ctx)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	finder := find.NewFinder(client.Client, true)

	// find resources by their paths
	vmFolder, err := finder.Folder(ctx, path.Join("/", datacenter.Name(), "vm", d.VsphereMachineClass.Spec.VirtualMachineFolder))
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	var resourcePool *object.ResourcePool
	if len(d.VsphereMachineClass.Spec.ResourcePool) != 0 {
		resourcePool, err = finder.ResourcePool(ctx, path.Join(cluster.InventoryPath, "Resources", d.VsphereMachineClass.Spec.ResourcePool))
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}
	} else {
		resourcePool, err = finder.ResourcePool(ctx, path.Join(cluster.InventoryPath, "Resources"))
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}
	}

	mainNetwork, err := finder.Network(ctx, path.Join("/", datacenter.Name(), "network", d.VsphereMachineClass.Spec.MainNetwork))
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	var additionalNetworks []object.NetworkReference
	for _, addNetPath := range d.VsphereMachineClass.Spec.AdditionalNetworks {
		addNet, err := finder.Network(ctx, path.Join("/", datacenter.Name(), "network", addNetPath))
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}

		additionalNetworks = append(additionalNetworks, addNet)
	}

	dsPath := path.Join("/", datacenter.Name(), "datastore", d.VsphereMachineClass.Spec.Datastore)

	var (
		datastore        *object.Datastore
		datastoreCluster *object.StoragePod
	)
	datastoreCluster, err = finder.DatastoreCluster(ctx, dsPath)
	if _, ok := err.(*find.NotFoundError); ok {
		klog.V(2).Infof("DatastoreCluster by path %q not found, trying to find a Datastore...", dsPath)
		datastore, err = finder.Datastore(ctx, path.Join("/", datacenter.Name(), "datastore", d.VsphereMachineClass.Spec.Datastore))
		// TODO: remove at earliest convenience
		if _, ok := err.(*find.NotFoundError); ok {
			hackyDscPath := path.Join("/", datacenter.Name(), "datastore", "3par_4_k8s")
			klog.V(2).Infof("Can't find either DatastoreCluster nor Datastore by name %q, falling to the deepest level of hell: %s", dsPath, hackyDscPath)

			datastoreCluster, err = finder.DatastoreCluster(ctx, hackyDscPath)
			if err != nil {
				return d.encodeMachineID(machineID), "", fmt.Errorf("even the strongest magic incantation failed, all hope is lost: %s", NewCreateFailErr(err, sideEffectsPresent))
			}
		} else if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}
	}

	vmTemplate, err := finder.VirtualMachine(ctx, path.Join("/", datacenter.Name(), "vm", d.VsphereMachineClass.Spec.Template))
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	err = vsphereInfo.tagManager.EnsureTagCategory(ctx, restClient, string(d.SecretData[v1alpha1.VsphereClusterNameTagCategory]), "", tagCategoryCardinalitySingle, []string{"VirtualMachine"})
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	err = vsphereInfo.tagManager.EnsureTagCategory(ctx, restClient, string(d.SecretData[v1alpha1.VsphereNodeRoleTagCategory]), "", tagCategoryCardinalitySingle, []string{"VirtualMachine"})
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	err = vsphereInfo.tagManager.EnsureTag(ctx, restClient, d.VsphereMachineClass.Spec.ClusterNameTag, string(d.SecretData[v1alpha1.VsphereClusterNameTagCategory]), "")
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	err = vsphereInfo.tagManager.EnsureTag(ctx, restClient, d.VsphereMachineClass.Spec.NodeRoleTag, string(d.SecretData[v1alpha1.VsphereNodeRoleTagCategory]), "")
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	var oldNetworkDevicesConfigSpecs []types.BaseVirtualDeviceConfigSpec
	var newNetworkDevicesConfigSpecs []types.BaseVirtualDeviceConfigSpec
	var clonedVM *object.VirtualMachine

	existingVM, err := finder.VirtualMachine(ctx, path.Join(vmFolder.InventoryPath, d.MachineName))
	if _, ok := err.(*find.NotFoundError); ok {
		// VM does not exist, clone it from the template

		templateDeviceList, err := vmTemplate.Device(ctx)
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}

		var oldEthernetCardsList object.VirtualDeviceList
		for _, device := range templateDeviceList {
			if _, ok := device.(types.BaseVirtualEthernetCard); ok {
				oldEthernetCardsList = append(oldEthernetCardsList, device)
			}
		}

		var newEthernetCardsList object.VirtualDeviceList
		networkBackingInfo, err := mainNetwork.EthernetCardBackingInfo(ctx)
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}
		mainEthernetCard, err := templateDeviceList.CreateEthernetCard("vmxnet3", networkBackingInfo)
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}
		newEthernetCardsList = append(newEthernetCardsList, mainEthernetCard)

		for _, addNet := range additionalNetworks {
			backing, err := addNet.EthernetCardBackingInfo(ctx)
			if err != nil {
				return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
			}
			ethernetCard, err := templateDeviceList.CreateEthernetCard("vmxnet3", backing)
			if err != nil {
				return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
			}
			newEthernetCardsList = append(newEthernetCardsList, ethernetCard)
		}

		// govmomi generates "-1" key for each device,
		// vSphere 7.0 treats these keys as duplicates and refuses to add new devices
		for i, newNetDev := range newEthernetCardsList {
			key := 10000 + i*1000
			newNetDev.(*types.VirtualVmxnet3).Key = int32(key)
		}

		oldEthernetConfigSpec, err := oldEthernetCardsList.ConfigSpec(types.VirtualDeviceConfigSpecOperationRemove)
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}

		newEthernetConfigSpec, err := newEthernetCardsList.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}

		oldNetworkDevicesConfigSpecs = append(oldNetworkDevicesConfigSpecs, oldEthernetConfigSpec...)
		newNetworkDevicesConfigSpecs = append(newNetworkDevicesConfigSpecs, newEthernetConfigSpec...)

		// construct MO references
		rpRef := resourcePool.Reference()
		folderRef := vmFolder.Reference()
		vmTemplateRef := vmTemplate.Reference()

		var cloneSpec = &types.VirtualMachineCloneSpec{
			Location: types.VirtualMachineRelocateSpec{
				Pool:   &rpRef,
				Folder: &folderRef,
			},
			PowerOn: false,
		}
		if !clusterDrsEnabled {
			host, err := d.selectHost(ctx, client, restClient, d.VsphereMachineClass.Spec.ClusterNameTag, d.VsphereMachineClass.Spec.NodeRoleTag, hosts)
			if err != nil {
				return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
			}
			hostRef := host.Reference()

			cloneSpec.Location.Host = &hostRef
		}

		if datastoreCluster != nil {
			dsClusterRef := datastoreCluster.Reference()

			srm := object.NewStorageResourceManager(client.Client)
			recommendedDatastores, err := srm.RecommendDatastores(ctx, types.StoragePlacementSpec{
				Type: string(types.StoragePlacementSpecPlacementTypeClone),
				PodSelectionSpec: types.StorageDrsPodSelectionSpec{
					StoragePod: &dsClusterRef,
				},
				Vm:        &vmTemplateRef,
				CloneName: d.MachineName,
				Folder:    &folderRef,
				CloneSpec: cloneSpec,
			})
			if err != nil {
				return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
			}
			var recommendationKeys []string
			for _, recommendation := range recommendedDatastores.Recommendations {
				recommendationKeys = append(recommendationKeys, recommendation.Key)
			}

			recommendTask, err := srm.ApplyStorageDrsRecommendation(ctx, recommendationKeys)
			if err != nil {
				return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
			}

			recommendTaskResult, err := recommendTask.WaitForResult(ctx)
			if err != nil {
				return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
			}

			clonedVM = object.NewVirtualMachine(client.Client, *recommendTaskResult.Result.(types.ApplyStorageRecommendationResult).Vm)
		} else {
			dsRef := datastore.Reference()
			cloneSpec.Location.Datastore = &dsRef

			cloneTask, err := vmTemplate.Clone(ctx, vmFolder, d.MachineName, *cloneSpec)
			if err != nil {
				return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
			}

			cloneTaskResult, err := cloneTask.WaitForResult(ctx)
			if err != nil {
				return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
			}

			clonedVM = object.NewVirtualMachine(client.Client, cloneTaskResult.Result.(types.ManagedObjectReference))
		}
	} else {
		clonedVM = existingVM
	}

	machineID = clonedVM.UUID(ctx)
	sideEffectsPresent = true

	err = vsphereInfo.tagManager.EnsureAttachedTagToObject(ctx, restClient, d.VsphereMachineClass.Spec.ClusterNameTag, string(d.SecretData[v1alpha1.VsphereClusterNameTagCategory]), clonedVM.Reference())
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	err = vsphereInfo.tagManager.EnsureAttachedTagToObject(ctx, restClient, d.VsphereMachineClass.Spec.NodeRoleTag, string(d.SecretData[v1alpha1.VsphereNodeRoleTagCategory]), clonedVM.Reference())
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	// fill up the extraConfig field of a new VM
	type CloudInitMetadata struct {
		LocalHostname  string `json:"local-hostname"`
		PublicKeysData string `json:"public-keys-data"`
	}

	guestInfoMetadata, err := json.Marshal(
		CloudInitMetadata{
			LocalHostname:  d.MachineName,
			PublicKeysData: strings.Join(d.VsphereMachineClass.Spec.SshKeys, "\n"),
		},
	)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	var options []*types.OptionValue
	options = append(options, []*types.OptionValue{
		{
			Key:   "disk.EnableUUID",
			Value: "TRUE",
		},
		{
			Key:   "guestinfo.metadata",
			Value: base64.StdEncoding.EncodeToString(guestInfoMetadata),
		},
		{
			Key:   "guestinfo.metadata.encoding",
			Value: "base64",
		},
		{
			Key:   "guestinfo.userdata",
			Value: base64.StdEncoding.EncodeToString([]byte(d.UserData)),
		},
		{
			Key:   "guestinfo.userdata.encoding",
			Value: "base64",
		},
	}...,
	)

	if d.VsphereMachineClass.Spec.DisableTimesync {
		options = append(options, []*types.OptionValue{
			{
				Key:   "tools.syncTime",
				Value: "0",
			},
			{
				Key:   "time.synchronize.continue",
				Value: "0",
			},
			{
				Key:   "time.synchronize.restore",
				Value: "0",
			},
			{
				Key:   "time.synchronize.resume.disk",
				Value: "0",
			},
			{
				Key:   "time.synchronize.shrink",
				Value: "0",
			},
			{
				Key:   "time.synchronize.startup",
				Value: "0",
			},
			{
				Key:   "time.synchronize.enable",
				Value: "0",
			},
			{
				Key:   "time.synchronize.resume.host",
				Value: "0",
			},
		}...,
		)
	}

	for k, v := range d.VsphereMachineClass.Spec.ExtraConfig {
		options = append(options, &types.OptionValue{
			Key:   k,
			Value: v,
		})
	}

	devices, err := clonedVM.Device(ctx)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	var rootDisk *types.VirtualDisk
	for _, device := range devices {
		switch md := device.(type) {
		case *types.VirtualDisk:
			rootDisk = md
			break
		default:
			continue
		}
	}
	if rootDisk == nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(fmt.Errorf("template \"%s\" doesn't have any disks", vmTemplate.InventoryPath), sideEffectsPresent)
	}
	rootDisk.CapacityInBytes = int64(d.VsphereMachineClass.Spec.RootDiskSize * 1024 * 1024 * 1024)

	var deviceChangeSpec []types.BaseVirtualDeviceConfigSpec
	deviceChangeSpec = append(deviceChangeSpec, oldNetworkDevicesConfigSpecs...)
	deviceChangeSpec = append(deviceChangeSpec, newNetworkDevicesConfigSpecs...)
	deviceChangeSpec = append(deviceChangeSpec, &types.VirtualDeviceConfigSpec{
		Operation: types.VirtualDeviceConfigSpecOperationEdit,
		Device:    rootDisk,
	})

	configSpec := types.VirtualMachineConfigSpec{
		Name:     d.MachineName,
		NumCPUs:  int32(d.VsphereMachineClass.Spec.NumCPUs),
		MemoryMB: int64(d.VsphereMachineClass.Spec.Memory),
		CpuAllocation: &types.ResourceAllocationInfo{
			Shares: func(shares *int32) *types.SharesInfo {
				if shares != nil {
					return &types.SharesInfo{Level: types.SharesLevelCustom, Shares: *shares}
				}

				return nil
			}(d.VsphereMachineClass.Spec.RuntimeOptions.ResourceAllocationInfo.CpuShares),
			Limit:       d.VsphereMachineClass.Spec.RuntimeOptions.ResourceAllocationInfo.CpuLimit,
			Reservation: d.VsphereMachineClass.Spec.RuntimeOptions.ResourceAllocationInfo.CpuReservation,
		},
		MemoryAllocation: &types.ResourceAllocationInfo{
			Shares: func(shares *int32) *types.SharesInfo {
				if shares != nil {
					return &types.SharesInfo{Level: types.SharesLevelCustom, Shares: *shares}
				}

				return nil
			}(d.VsphereMachineClass.Spec.RuntimeOptions.ResourceAllocationInfo.MemoryShares),
			Limit:       d.VsphereMachineClass.Spec.RuntimeOptions.ResourceAllocationInfo.MemoryLimit,
			Reservation: d.VsphereMachineClass.Spec.RuntimeOptions.ResourceAllocationInfo.MemoryReservation,
		},
		DeviceChange:    deviceChangeSpec,
		NestedHVEnabled: &d.VsphereMachineClass.Spec.RuntimeOptions.NestedHardwareVirtualization,
		ExtraConfig: func() []types.BaseOptionValue {
			baseOptionValues := make([]types.BaseOptionValue, len(options))
			for index, option := range options {
				baseOptionValues[index] = option
			}
			return baseOptionValues
		}(),
		GuestId: string(types.VirtualMachineGuestOsIdentifierUbuntu64Guest),
	}

	// power off the VM (if it isn't) before reconfiguring it
	vmPowerState, err := clonedVM.PowerState(ctx)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	if vmPowerState != types.VirtualMachinePowerStatePoweredOff {
		vmPowerOffTask, err := clonedVM.PowerOff(ctx)
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}

		err = vmPowerOffTask.Wait(ctx)
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}

		timeout, _ := context.WithTimeout(ctx, 30*time.Second)
		err = clonedVM.WaitForPowerState(timeout, types.VirtualMachinePowerStatePoweredOff)
		if err != nil {
			return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
		}
	}

	reconfigureClonedVmTask, err := clonedVM.Reconfigure(ctx, configSpec)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	err = reconfigureClonedVmTask.Wait(ctx)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	err = d.upgradeHardware(ctx, clonedVM, hwVersion)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	powerOnVmTask, err := clonedVM.PowerOn(ctx)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	err = powerOnVmTask.Wait(ctx)
	if err != nil {
		return d.encodeMachineID(machineID), "", NewCreateFailErr(err, sideEffectsPresent)
	}

	return d.encodeMachineID(machineID), d.MachineName, nil
}

func (d *VsphereDriver) upgradeHardware(ctx context.Context, vm *object.VirtualMachine, version int) error {
	if version > 0 {
		// update hardware
		version := fmt.Sprintf("vmx-%02d", hwVersion)
		upgradeVmTask, err := vm.UpgradeVM(ctx, version)
		if err != nil {
			return err
		}
		err = upgradeVmTask.Wait(ctx)
		if err != nil {
			if isAlreadyUpgraded(err) {
				klog.V(4).Infof("Already upgraded: %s", err)
			} else {
				return err
			}
		}
	}
	return nil
}

func isAlreadyUpgraded(err error) bool {
	if fault, ok := err.(task.Error); ok {
		_, ok = fault.Fault().(*types.AlreadyUpgraded)
		return ok
	}

	return false
}

func (d *VsphereDriver) Delete(machineID string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vsphereClient, restClient, logoutFunc, err := vsphereInfo.loginFunc(ctx)
	if err != nil {
		return err
	}
	defer logoutFunc()

	datacenter, err := vsphereInfo.tagManager.GetDcByRegion(ctx, vsphereClient, restClient, d.VsphereMachineClass.Spec.Region)
	if err != nil {
		return err
	}

	machineID = d.decodeMachineID(machineID)

	finder := find.NewFinder(vsphereClient.Client, false)
	vmFolder, err := finder.Folder(ctx, path.Join("/", datacenter.Name(), "vm", d.VsphereMachineClass.Spec.VirtualMachineFolder))
	if err != nil {
		return err
	}
	folderRef := vmFolder.Reference()

	viewManager := view.NewManager(vsphereClient.Client)
	vmView, err := viewManager.CreateContainerView(ctx, folderRef, []string{"VirtualMachine"}, true)
	if err != nil {
		return err
	}
	defer vmView.Destroy(ctx)

	var vms []mo.VirtualMachine

	err = vmView.RetrieveWithFilter(ctx, []string{"VirtualMachine"}, []string{"config.uuid", "config.name"}, &vms, property.Filter{"config.uuid": machineID})
	// TODO: Upstream proper error handling (https://github.com/vmware/govmomi/property/collector.go:133) to skip this anti-pattern
	if err != nil && err.Error() != "object references is empty" {
		return err
	}

	if len(vms) == 0 {
		klog.V(2).Infof("Not a single VM with UUID \"%s\". Skipping \"Delete\" operation.", machineID)
		return nil
	}

	if len(vms) != 1 {
		return fmt.Errorf("number of returned \"%v\" != \"1\" for uuid \"%v\"", len(vms), machineID)
	}

	vmObject := object.NewVirtualMachine(vsphereClient.Client, vms[0].Reference())

	vmPowerState, err := vmObject.PowerState(ctx)
	if err != nil {
		return err
	}

	if vmPowerState != types.VirtualMachinePowerStatePoweredOff {
		vmPowerOffTask, err := vmObject.PowerOff(ctx)
		if err != nil {
			return err
		}

		err = vmPowerOffTask.Wait(ctx)
		if err != nil {
			return err
		}

		timeout, _ := context.WithTimeout(ctx, 30*time.Second)
		err = vmObject.WaitForPowerState(timeout, types.VirtualMachinePowerStatePoweredOff)
		if err != nil {
			return err
		}
	}

	// destroy only first disk, detach everything else
	deviceList, err := vmObject.Device(ctx)
	if err != nil {
		return err
	}

	var disksSansFirst object.VirtualDeviceList
	var diskIndex int
	for _, device := range deviceList {
		switch device.(type) {
		case *types.VirtualDisk:
			// remove this abominable monstrosity!
			if func() bool {
				diskIndex++
				return diskIndex > 1
			}() {
				disksSansFirst = append(disksSansFirst, device)
			}
		default:
			continue
		}
	}

	if len(disksSansFirst) > 0 {
		deviceConfigSpec, err := disksSansFirst.ConfigSpec(types.VirtualDeviceConfigSpecOperationRemove)
		if err != nil {
			return err
		}

		for i, deviceChange := range deviceConfigSpec {
			deviceChange.GetVirtualDeviceConfigSpec().FileOperation = ""
			deviceConfigSpec[i] = deviceChange
		}

		configSpec := types.VirtualMachineConfigSpec{
			DeviceChange: deviceConfigSpec,
		}

		reconfigureClonedVmTask, err := vmObject.Reconfigure(ctx, configSpec)
		if err != nil {
			return err
		}

		err = reconfigureClonedVmTask.Wait(ctx)
		if err != nil {
			return err
		}
	}

	vmDestroyTask, err := vmObject.Destroy(ctx)
	if err != nil {
		return err
	}

	err = vmDestroyTask.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *VsphereDriver) GetExisting() (string, error) {
	return d.MachineID, nil
}

func (d *VsphereDriver) GetVMs(machineID string) (VMs, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, restClient, logoutFunc, err := vsphereInfo.loginFunc(ctx)
	if err != nil {
		return nil, err
	}
	defer logoutFunc()

	datacenter, err := vsphereInfo.tagManager.GetDcByRegion(ctx, client, restClient, d.VsphereMachineClass.Spec.Region)
	if err != nil {
		return nil, err
	}

	listOfVMs := make(map[string]string)

	finder := find.NewFinder(client.Client)

	vmFolder, err := finder.Folder(ctx, path.Join("/", datacenter.Name(), "vm", d.VsphereMachineClass.Spec.VirtualMachineFolder))
	if err != nil {
		return nil, err
	}
	folderRef := vmFolder.Reference()

	viewManager := view.NewManager(client.Client)
	vmView, err := viewManager.CreateContainerView(ctx, folderRef, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}
	defer vmView.Destroy(ctx)

	if machineID == "" {
		clusterName := d.VsphereMachineClass.Spec.ClusterNameTag
		nodeRole := d.VsphereMachineClass.Spec.NodeRoleTag

		if clusterName == "" || nodeRole == "" {
			return listOfVMs, nil
		}

		vmsByClusterName, err := vsphereInfo.tagManager.GetObjectRefsByTag(ctx, client, restClient, clusterName)
		// TODO: remove ugliness, restore type-safe cleanliness
		if err != nil && strings.Contains(err.Error(), "404 Not Found") {
			klog.Warningf("Can't get objects by tag %q, skipping VM listing: %s", clusterName, err)
			return listOfVMs, nil
		} else if err != nil {
			return nil, err
		}

		// TODO: remove ugliness, restore type-safe cleanliness
		vmsByNodeRole, err := vsphereInfo.tagManager.GetObjectRefsByTag(ctx, client, restClient, nodeRole)
		if err != nil && strings.Contains(err.Error(), "404 Not Found") {
			klog.Warningf("Can't get objects by tag %q, skipping VM listing: %s", nodeRole, err)
			return listOfVMs, nil
		} else if err != nil {
			return nil, err
		}

		clusterNodeIntersection := intersectMOReferences(vmsByClusterName, vmsByNodeRole)

		var vms []mo.VirtualMachine
		err = vmView.Retrieve(ctx, []string{"VirtualMachine"}, []string{}, &vms)
		if err != nil {
			return nil, err
		}

		var vmObjectRefs []object.Reference
		for _, vmMO := range vms {
			vmObjectRefs = append(vmObjectRefs, vmMO.Reference())
		}

		matchingVMs := intersectMOReferences(vmObjectRefs, clusterNodeIntersection)
		for _, vmObjectRef := range matchingVMs {
			vmObject := vmObjectRef.(*object.VirtualMachine)
			var vmMO mo.VirtualMachine
			err := vmObject.Properties(ctx, vmObject.Reference(), []string{"config"}, &vmMO)
			if err != nil {
				return nil, err
			}
			listOfVMs[d.encodeMachineID(vmMO.Config.Uuid)] = vmMO.Config.Name
		}

	} else {
		machineID = d.decodeMachineID(machineID)

		var vms []mo.VirtualMachine
		err = vmView.RetrieveWithFilter(ctx, []string{"VirtualMachine"}, []string{"config.uuid", "config.name"}, &vms, property.Filter{"config.uuid": machineID})
		if err != nil {
			return nil, err
		}

		if len(vms) != 1 {
			return nil, fmt.Errorf("number of returned \"%v\" != \"1\" for uuid \"%v\"", len(vms), machineID)
		}

		listOfVMs[d.encodeMachineID(machineID)] = vms[0].Config.Name
	}

	return listOfVMs, nil
}

func (d *VsphereDriver) GetVolNames(specs []corev1.PersistentVolumeSpec) ([]string, error) {
	var names []string
	for i := range specs {
		spec := &specs[i]
		if spec.CSI == nil {
			// Not a CSI-managed volume
			continue
		}
		if spec.CSI.Driver != "vsphere.csi.vmware.com" {
			// Not a volume provisioned by vSphere CSI driver
			continue
		}
		names = append(names, spec.CSI.VolumeHandle)
	}

	return names, nil
}

//GetUserData return the used data whit which the VM will be booted
func (d *VsphereDriver) GetUserData() string {
	return d.UserData
}

//SetUserData set the used data whit which the VM will be booted
func (d *VsphereDriver) SetUserData(userData string) {
	d.UserData = userData
}

func (d *VsphereDriver) encodeMachineID(machineID string) string {
	if len(machineID) != 0 {
		return fmt.Sprintf("vsphere://%s", machineID)
	}

	return ""
}

func (d *VsphereDriver) decodeMachineID(id string) string {
	splitProviderID := strings.Split(id, "/")
	return splitProviderID[len(splitProviderID)-1]
}

func (d *VsphereDriver) selectHost(ctx context.Context, client *govmomi.Client, restClient *rest.Client, clusterName, nodeRole string, hosts []*object.HostSystem) (*object.HostSystem, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts to choose from")
	}

	vmsByClusterName, err := vsphereInfo.tagManager.GetObjectRefsByTag(ctx, client, restClient, clusterName)
	if err != nil {
		return nil, err
	}

	vmsByNodeRole, err := vsphereInfo.tagManager.GetObjectRefsByTag(ctx, client, restClient, nodeRole)
	if err != nil {
		return nil, err
	}

	matchingVMs := intersectMOReferences(vmsByClusterName, vmsByNodeRole)

	hostWeight := make(map[string]uint64)
	for _, vm := range matchingVMs {
		vmObject := vm.(*object.VirtualMachine)
		var vmWithProperties mo.VirtualMachine
		err := vmObject.Properties(ctx, vmObject.Reference(), []string{"runtime.host"}, &vmWithProperties)
		if err != nil {
			return &object.HostSystem{}, err
		}

		hostWeight[vmWithProperties.Runtime.Host.Reference().String()]++
	}

	var hostsWithProperties []mo.HostSystem
	for _, host := range hosts {
		err := host.Properties(ctx, host.Reference(), []string{"summary", "name"}, &hostsWithProperties)
		if err != nil {
			return nil, err
		}
	}
	// sort based on fairness
	// https://code.vmware.com/apis/358/vsphere/doc/vim.HostSystem.html/doc/vim.host.Summary.QuickStats.html
	sort.Slice(hostsWithProperties, func(i, j int) bool {
		if hostWeight[hostsWithProperties[i].Reference().String()] < hostWeight[hostsWithProperties[j].Reference().String()] {
			return true
		}

		if hostWeight[hostsWithProperties[i].Reference().String()] > hostWeight[hostsWithProperties[j].Reference().String()] {
			return false
		}

		iCpuFairnessDeviation := abs(int64(1000 - hostsWithProperties[i].Summary.QuickStats.DistributedCpuFairness))
		iMemoryFairnessDeviation := abs(int64(1000 - hostsWithProperties[i].Summary.QuickStats.DistributedMemoryFairness))
		jCpuFairnessDeviation := abs(int64(1000 - hostsWithProperties[j].Summary.QuickStats.DistributedCpuFairness))
		jMemoryFairnessDeviation := abs(int64(1000 - hostsWithProperties[j].Summary.QuickStats.DistributedMemoryFairness))

		if iMemoryFairnessDeviation < jMemoryFairnessDeviation {
			return true
		}
		if iMemoryFairnessDeviation > jMemoryFairnessDeviation {
			return false
		}
		// use CPU fairness to break Memory fairness ties
		return iCpuFairnessDeviation < jCpuFairnessDeviation
	})

	hostSystem := object.NewHostSystem(client.Client, hostsWithProperties[0].Reference())

	return hostSystem, nil
}

// http://cavaliercoder.com/blog/optimized-abs-for-int64-in-go.html
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func intersectMOReferences(a, b []object.Reference) (c []object.Reference) {
	m := make(map[string]struct{})

	for _, item := range a {
		m[item.Reference().String()] = struct{}{}
	}

	for _, item := range b {
		if _, ok := m[item.Reference().String()]; ok {
			c = append(c, item)
		}
	}
	return
}
