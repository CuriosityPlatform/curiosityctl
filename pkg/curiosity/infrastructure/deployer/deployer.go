package deployer

import (
	"os"
	"path"

	"github.com/pkg/errors"

	"curiosity/pkg/common/infrastructure/kubeclient"
	"curiosity/pkg/curiosity/app/service"
)

func NewDeployer(client kubeclient.KubeClient, basePath string) service.Deployer {
	return &deployer{client: client, basePath: basePath}
}

type deployer struct {
	client   kubeclient.KubeClient
	basePath string
}

func (deployer *deployer) DeployAll() error {
	return deployer.client.ApplyPath(deployer.basePath)
}

func (deployer *deployer) DeployApp(name string) error {
	appPath := path.Join(deployer.basePath, name)
	if _, err := os.Stat(appPath); errors.Is(err, os.ErrNotExist) {
		return errors.Wrapf(err, "failed to deploy %s", name)
	}

	return deployer.client.ApplyPath(appPath)
}
