package project

import (
	"io/fs"

	"github.com/pkg/errors"

	"curiosity/pkg/packager/app"
)

func NewProjectDetector(detectStrategies []DetectionStrategy) app.ProjectDetector {
	return &projectDetector{detectStrategies: detectStrategies}
}

type projectDetector struct {
	detectStrategies []DetectionStrategy
}

func (detector *projectDetector) FindProject(f fs.FS) (app.Project, error) {
	for _, strategy := range detector.detectStrategies {
		project, err := strategy.Detect(f)
		if err != nil {
			if errors.Is(err, ErrProjectDoesNotSupported) {
				continue
			}
			return app.Project{}, err
		}

		return project, nil
	}

	return app.Project{}, nil
}
