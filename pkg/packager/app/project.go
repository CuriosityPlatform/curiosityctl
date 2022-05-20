package app

import "io/fs"

type Project struct {
	Name string
}

type ProjectDetector interface {
	FindProject(fs fs.FS) (Project, error)
}
