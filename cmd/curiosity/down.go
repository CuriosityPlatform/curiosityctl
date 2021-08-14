package main

import (
	"github.com/urfave/cli/v2"

	"curiosity/pkg/common/infrastructure/dockerclient"
	"curiosity/pkg/curiosity/app/usecase"
)

func executeDown(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	executor, err := dockerclient.NewDockerExecutorWithStaticDir(config.PlatformRoot)
	if err != nil {
		return err
	}

	useCase := usecase.NewDown(dockerclient.NewClient(executor))

	return useCase.Execute(ctx.Context)
}
