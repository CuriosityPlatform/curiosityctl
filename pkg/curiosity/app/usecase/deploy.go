package usecase

import (
	"context"

	"curiosity/pkg/curiosity/app/service"
)

func NewDeploy(deployer service.Deployer) *Deploy {
	return &Deploy{deployer: deployer}
}

type Deploy struct {
	deployer service.Deployer
}

func (d *Deploy) Execute(ctx context.Context, appName *string) error {
	if appName != nil {
		return d.deployer.DeployApp(*appName)
	}

	return d.deployer.DeployAll()
}
