package packager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"curiosity/pkg/common/infrastructure/executor"
	"curiosity/pkg/packager/app"
)

func NewPackagerFactory() app.PackagerFactory {
	return &packagerFactory{}
}

type packagerFactory struct {
}

func (factory *packagerFactory) PackagerByType(packagerType app.PackagerType) (app.Packager, error) {
	switch packagerType {
	case app.Makefile:
		return newMakefilePackager()
	default:
		return nil, errors.New("unknown packager type")
	}
}

func newMakefilePackager() (*MakefilePackager, error) {
	e, err := executor.New("make")
	if err != nil {
		return nil, err
	}

	return &MakefilePackager{executor: e}, nil
}

type MakefilePackager struct {
	executor executor.Executor
}

func (m *MakefilePackager) Package(ctx context.Context, project app.Project, params app.PackageParams) (app.Package, error) {
	image := fmt.Sprintf("%s/%s", params.Registry, project.Name)
	err := m.executor.PTY(ctx, []string{fmt.Sprintf("IMAGE=%s", image), "build-image"})
	if err != nil {
		return app.Package{}, errors.WithStack(err)
	}

	return app.Package{
		Image: app.Image(image),
	}, nil
}
