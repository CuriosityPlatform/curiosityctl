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

	var outputPath *string

	if path := ctx.String("output"); path != "" {
		outputPath = &path
	}

	return useCase.Execute(ctx.Context, usecase.InstallPlatformExecuteParams{
		OutPath:       outputPath,
		RepoRemoteURL: platformRepoURL,
	})
}
