package main

import (
	"path"

	"github.com/urfave/cli/v2"

	"curiosity/pkg/common/infrastructure/kubeclient"
	"curiosity/pkg/curiosity/app/usecase"
	"curiosity/pkg/curiosity/infrastructure/deployer"
)

func executeDeploy(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	client, err := kubeclient.NewKubeClient()
	if err != nil {
		return err
	}

	useCase := usecase.NewDeploy(deployer.NewDeployer(
		client,
		path.Join(config.PlatformRoot, config.KubernetesManifestsBasePath),
	))

	return useCase.Execute(ctx.Context, emptyStringToPtr(ctx.String("app")))
}
