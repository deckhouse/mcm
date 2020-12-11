package driver

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/vmware/govmomi/vapi/rest"
	"k8s.io/klog"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

func createVsphereClient(host, username, password string, insecure bool) func(ctx context.Context) (*govmomi.Client, *rest.Client, func(), error) {
	return func(ctx context.Context) (*govmomi.Client, *rest.Client, func(), error) {
		parsedURL, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk", url.PathEscape(strings.TrimSpace(username)), url.PathEscape(strings.TrimSpace(password)), url.PathEscape(strings.TrimSpace(host))))
		if err != nil {
			return nil, nil, nil, err
		}

		userInfo := url.UserPassword(username, password)

		vcClient, err := govmomi.NewClient(ctx, parsedURL, insecure)
		if err != nil {
			return nil, nil, nil, err
		}

		if !vcClient.IsVC() {
			return nil, nil, nil, errors.New("not connected to vCenter")
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
