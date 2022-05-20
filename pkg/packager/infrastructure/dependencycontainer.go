package infrastructure

import (
	"curiosity/pkg/packager/api"
	"curiosity/pkg/packager/app"
	"curiosity/pkg/packager/infrastructure/packager"
	"curiosity/pkg/packager/infrastructure/project"
)

func BuildAPI() api.API {
	return api.NewAPI(app.NewService(
		project.NewProjectDetector([]project.DetectionStrategy{&project.GoProjectDetectionStrategy{}}),
		packager.NewPackagerFactory(),
	))
}
