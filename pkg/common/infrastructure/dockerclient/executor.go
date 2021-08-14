package dockerclient

import (
	"context"

	"curiosity/pkg/common/infrastructure/executor"
)

const (
	dockerPath = "docker"
)

type DockerExecutor interface {
	executor.Executor
}

func NewDockerExecutor() (DockerExecutor, error) {
	return executor.New(dockerPath)
}

func NewDockerExecutorWithStaticDir(workDir string) (DockerExecutor, error) {
	dockerExecutor, err := NewDockerExecutor()
	if err != nil {
		return nil, err
	}

	return &dockerExecutorWithStaticDir{
		workDir:        workDir,
		DockerExecutor: dockerExecutor,
	}, nil
}

// dockerExecutorWithStaticDir
type dockerExecutorWithStaticDir struct {
	workDir string
	DockerExecutor
}

func (e *dockerExecutorWithStaticDir) Output(args []string, opts ...executor.Opt) ([]byte, error) {
	return e.DockerExecutor.Output(args, append([]executor.Opt{executor.WithWorkdir(e.workDir)}, opts...)...)
}

func (e *dockerExecutorWithStaticDir) Run(args []string, opts ...executor.Opt) error {
	return e.DockerExecutor.Run(args, append([]executor.Opt{executor.WithWorkdir(e.workDir)}, opts...)...)
}

func (e *dockerExecutorWithStaticDir) PTY(ctx context.Context, args []string, opts ...executor.Opt) error {
	return e.DockerExecutor.PTY(ctx, args, append([]executor.Opt{executor.WithWorkdir(e.workDir)}, opts...)...)
}
