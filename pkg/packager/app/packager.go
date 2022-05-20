package app

import "context"

type PackagerType int

const (
	Makefile PackagerType = iota
)

type PackagerFactory interface {
	PackagerByType(packagerType PackagerType) (Packager, error)
}

type PackageParams struct {
	Registry string
	Push     bool
}

type Packager interface {
	Package(ctx context.Context, project Project, params PackageParams) (Package, error)
}
