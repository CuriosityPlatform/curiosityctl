package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"curiosity/pkg/packager/api"
	"curiosity/pkg/packager/infrastructure"
)

func Package() *cli.Command {
	return &cli.Command{
		Name:   "package",
		Usage:  "Pack project",
		Action: executePackage,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "project-path",
				Aliases: []string{"pr"},
			},
			&cli.BoolFlag{
				Name:    "push",
				Aliases: []string{"p"},
			},
		},
	}
}

func executePackage(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	projectPath := ctx.String("project-path")
	if projectPath == "" {
		var err2 error
		projectPath, err2 = os.Getwd()
		if err2 != nil {
			return err2
		}
	}

	if exists, err2 := projectDirExists(projectPath); !exists || err2 != nil {
		if err2 == nil {
			return errors.Errorf("project at %s not found", projectPath)
		}
		return err2
	}

	fs := os.DirFS(projectPath)

	packagerAPI := infrastructure.BuildAPI()
	p, err := packagerAPI.Package(ctx.Context, api.PackageParams{
		PackageName: nil,
		Packager:    api.Makefile,
		Push:        ctx.Bool("push"),
		Registry:    config.DockerRegistry,
		ProjectFS:   fs,
	})
	if err != nil {
		return err
	}

	fmt.Println("Packaged with image", p.Image)

	return nil
}

func projectDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
