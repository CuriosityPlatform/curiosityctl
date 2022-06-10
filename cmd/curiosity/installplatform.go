package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"curiosity/pkg/common/infrastructure/git"
	"curiosity/pkg/curiosity/app/usecase"
)

const (
	platformRepoURL     = "git@github.com:CuriosityPlatform/platform.git"
	platformPathEnvName = "CURIOSITYCTL_PLATFORM_ROOT"
)

func executeInstallPlatform(ctx *cli.Context) error {
	vsc, err := git.NewGitVcs()
	if err != nil {
		return err
	}

	useCase := usecase.NewInstallPlatform(vsc)

	platformPath, err := useCase.Execute(ctx.Context, usecase.InstallPlatformExecuteParams{
		OutPath:       emptyStringToPtr(ctx.String("output")),
		RepoRemoteURL: platformRepoURL,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Please setup %s to %s as global environment\n", platformPathEnvName, platformPath)
	fmt.Println("Platform installed")

	return nil
}
