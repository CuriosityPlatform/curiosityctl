package main

import (
	"github.com/urfave/cli/v2"

	"curiosity/pkg/common/infrastructure/dockerclient"
	"curiosity/pkg/curiosity/app/compose"
	"curiosity/pkg/curiosity/app/usecase"
	infracompose "curiosity/pkg/curiosity/infrastructure/compose"
	"curiosity/pkg/curiosity/infrastructure/containerwaiter"
	"curiosity/pkg/curiosity/infrastructure/servicepreparer"
)

func executeUp(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	executor, err := dockerclient.NewDockerExecutorWithStaticDir(config.PlatformRoot)
	if err != nil {
		return err
	}

	client := dockerclient.NewClient(executor)

	loader := infracompose.NewLoader()

	project, err := loader.Load(infracompose.LoadParams{
		WorkDir: config.PlatformRoot,
	})
	if err != nil {
		return err
	}

	useCase := usecase.NewUp(
		client,
		servicepreparer.NewFactory(client),
		containerwaiter.NewWaiter(client),
		compose.NewProject(project),
	)

	return useCase.Execute(ctx.Context)
}
