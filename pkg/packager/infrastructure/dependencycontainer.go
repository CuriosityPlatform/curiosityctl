package infrastructure

import (
	"curiosity/pkg/common/infrastructure/dockerclient"
	"curiosity/pkg/packager/api"
	"curiosity/pkg/packager/app"
	"curiosity/pkg/packager/infrastructure/packager"
	"curiosity/pkg/packager/infrastructure/project"
)

func BuildAPI() (api.API, error) {
	executor, err := dockerclient.NewDockerExecutor()
	if err != nil {
		return nil, err
	}
	return api.NewAPI(app.NewService(
		project.NewProjectDetector([]project.DetectionStrategy{&project.GoProjectDetectionStrategy{}}),
		packager.NewPackagerFactory(dockerclient.NewClient(executor)),
	)), nil
}
