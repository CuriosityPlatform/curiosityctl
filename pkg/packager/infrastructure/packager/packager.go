package packager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/common/infrastructure/executor"
	"curiosity/pkg/packager/app"
)

func NewPackagerFactory(docker dockerclient.Client) app.PackagerFactory {
	return &packagerFactory{docker: docker}
}

type packagerFactory struct {
	docker dockerclient.Client
}

func (factory *packagerFactory) PackagerByType(packagerType app.PackagerType) (app.Packager, error) {
	switch packagerType {
	case app.Makefile:
		return newMakefilePackager(factory.docker)
	default:
		return nil, errors.New("unknown packager type")
	}
}

func newMakefilePackager(docker dockerclient.Client) (*MakefilePackager, error) {
	e, err := executor.New("make")
	if err != nil {
		return nil, err
	}

	return &MakefilePackager{executor: e, docker: docker}, nil
}

type MakefilePackager struct {
	executor executor.Executor
	docker   dockerclient.Client
}

func (m *MakefilePackager) Package(ctx context.Context, project app.Project, params app.PackageParams) (app.Package, error) {
	image := fmt.Sprintf("%s/%s", params.Registry, project.Name)
	err := m.executor.PTY(ctx, []string{fmt.Sprintf("IMAGE=%s", image), "build-image"})
	if err != nil {
		return app.Package{}, errors.WithStack(err)
	}

	if params.Push {
		err = m.docker.Push(image)
		if err != nil {
			return app.Package{}, err
		}
	}

	return app.Package{
		Image: app.Image(image),
	}, nil
}
