package driver

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
	"k8s.io/klog"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

func createVsphereClient(
	host, username, password string, insecure bool, caBundle []byte,
) func(ctx context.Context) (*govmomi.Client, *rest.Client, func(), error) {
	return func(ctx context.Context) (*govmomi.Client, *rest.Client, func(), error) {
		parsedURL, err := url.Parse(
			fmt.Sprintf(
				"https://%s:%s@%s/sdk",
				url.PathEscape(strings.TrimSpace(username)),
				url.PathEscape(strings.TrimSpace(password)),
				url.PathEscape(strings.TrimSpace(host)),
			),
		)
		if err != nil {
			return nil, nil, nil, err
		}

		userInfo := url.UserPassword(username, password)

		soapClient := soap.NewClient(parsedURL, insecure)
		if err := setCABundleIfNeed(soapClient, insecure, caBundle); err != nil {
			return nil, nil, nil, err
		}

		vimClient, err := vim25.NewClient(ctx, soapClient)
		if err != nil {
			return nil, nil, nil, err
		}

		if !vimClient.IsVC() {
			return nil, nil, nil, errors.New("not connected to vCenter")
		}

		vcClient := &govmomi.Client{
			Client:         vimClient,
			SessionManager: session.NewManager(vimClient),
		}
		if err := vcClient.SessionManager.Login(ctx, parsedURL.User); err != nil {
			return nil, nil, nil, err
		}

		restClient := rest.NewClient(vcClient.Client)
		if err := restClient.Login(ctx, userInfo); err != nil {
			return nil, nil, nil, err
		}

		logoutFunc := func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			err = vcClient.Logout(ctx)
			if err != nil {
				klog.Warningf("can't logout from SOAP API: %s", err)
			}
			err = restClient.Logout(ctx)
			if err != nil {
				klog.Warningf("can't logout from REST API: %s", err)
			}
		}
		return vcClient, restClient, logoutFunc, nil
	}
}

func drsEnabled(ctx context.Context, resource *object.ClusterComputeResource) (bool, error) {
	conf, err := resource.Configuration(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get Cluster configuration")
	}

	return conf.DrsConfig.Enabled != nil && *conf.DrsConfig.Enabled, nil
}

func setCABundleIfNeed(soapClient *soap.Client, insecure bool, caBundle []byte) error {
	hasCABundle := caBundle != nil && len(caBundle) > 0
	if !hasCABundle {
		return nil
	}

	if insecure {
		klog.Warningf("set insecure flag to true, CA bundle will be ignored")
		return nil
	}

	klog.V(2).Infof("setting CA bundle for VC client")

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(caBundle)) {
		return errors.New("failed to parse CA bundle")
	}

	soapClient.DefaultTransport().TLSClientConfig.RootCAs = pool
	return nil
}
