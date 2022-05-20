package app

import (
	"context"
	"io/fs"
)

type PackageProjectParams struct {
	// If nil name will be detected automatically
	PackageName *string
	Packager    PackagerType
	Registry    string
	Push        bool
	ProjectFS   fs.FS
}

type Image string

type Package struct {
	Image Image
}

func NewService(projectDetector ProjectDetector, packagerFactory PackagerFactory) Service {
	return &service{projectDetector: projectDetector, packagerFactory: packagerFactory}
}

type Service interface {
	PackageProject(ctx context.Context, params PackageProjectParams) (Package, error)
}

type service struct {
	projectDetector ProjectDetector
	packagerFactory PackagerFactory
}

func (service *service) PackageProject(ctx context.Context, params PackageProjectParams) (Package, error) {
	project, err := service.projectDetector.FindProject(params.ProjectFS)
	if err != nil {
		return Package{}, err
	}

	packager, err := service.packagerFactory.PackagerByType(params.Packager)
	if err != nil {
		return Package{}, err
	}

	return packager.Package(ctx, project, PackageParams{
		Registry: params.Registry,
		Push:     params.Push,
	})
}
