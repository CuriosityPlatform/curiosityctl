package api

import (
	"context"
	"io/fs"

	"curiosity/pkg/packager/app"
)

type Packager int

const (
	Makefile Packager = iota
)

type PackageParams struct {
	// If nil name will be detected automatically
	PackageName *string
	Packager    Packager
	Registry    string
	Push        bool
	ProjectFS   fs.FS
}

type Image string

type Package struct {
	Image Image
}

func NewAPI(service app.Service) API {
	return &api{service: service}
}

type API interface {
	Package(ctx context.Context, params PackageParams) (Package, error)
}

type api struct {
	service app.Service
}

func (a *api) Package(ctx context.Context, params PackageParams) (Package, error) {
	p, err := a.service.PackageProject(ctx, app.PackageProjectParams{
		PackageName: params.PackageName,
		Packager:    app.PackagerType(params.Packager),
		Registry:    params.Registry,
		Push:        params.Push,
		ProjectFS:   params.ProjectFS,
	})
	if err != nil {
		return Package{}, err
	}

	return Package{
		Image: Image(p.Image),
	}, err
}
