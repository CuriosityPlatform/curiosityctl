package main

import (
	"github.com/urfave/cli/v2"

	"curiosity/pkg/common/infrastructure/git"
	"curiosity/pkg/curiosity/app/usecase"
)

const (
	platformRepoURL = "git@github.com:CuriosityPlatform/platform.git"
)

func executeInstallPlatform(ctx *cli.Context) error {
	vsc, err := git.NewGitVcs()
	if err != nil {
		return err
	}

	useCase := usecase.NewInstallPlatform(vsc)

	return useCase.Execute(ctx.Context, usecase.InstallPlatformExecuteParams{
		OutPath:       emptyStringToPtr(ctx.String("output")),
		RepoRemoteURL: platformRepoURL,
	})
}
