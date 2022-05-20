package project

import (
	"bufio"
	stderrors "errors"
	"io/fs"
	"regexp"

	"github.com/pkg/errors"

	"curiosity/pkg/packager/app"
)

var (
	ErrProjectDoesNotSupported = stderrors.New("project does not supported")
)

type DetectionStrategy interface {
	Detect(fs fs.FS) (app.Project, error)
}

type GoProjectDetectionStrategy struct{}

func (strategy *GoProjectDetectionStrategy) Detect(projectFS fs.FS) (app.Project, error) {
	const goMod = "go.mod"

	const searchGroup = "project"
	projectNameRegexp := regexp.MustCompile(`module ([a-z]+\.[a-z]+)?(/)?(?P<project>.+)`)

	goModFile, err := projectFS.Open(goMod)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = errors.WithStack(ErrProjectDoesNotSupported)
		}
		return app.Project{}, err
	}

	line, err := strategy.getFileFirstLine(goModFile)
	if err != nil {
		return app.Project{}, err
	}

	projectName := findFromRegByGroupName(projectNameRegexp, line, searchGroup)
	if projectName == nil {
		return app.Project{}, errors.Errorf("failed to detect project name from %s from go.mod", line)
	}

	return app.Project{
		Name: *projectName,
	}, nil
}

func (strategy GoProjectDetectionStrategy) getFileFirstLine(file fs.File) (line string, err error) {
	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)
	scanned := fileScanner.Scan()
	if !scanned {
		return "", errors.New("empty project file")
	}
	line = fileScanner.Text()
	err = fileScanner.Err()
	return
}

func findFromRegByGroupName(reg *regexp.Regexp, s, searchGroup string) *string {
	groupNames := reg.SubexpNames()
	for _, match := range reg.FindAllStringSubmatch(s, -1) {
		for groupIdx, value := range match {
			groupName := groupNames[groupIdx]
			if groupName == "" {
				groupName = "*"
			}
			if groupName == searchGroup {
				return &value
			}
		}
	}
	return nil
}
