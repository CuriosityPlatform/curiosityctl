package dockerclient

import (
	"context"
	"io"
)

type ExecParam struct {
	Service string
	Command string
	Reader  io.Reader
}

type Client interface {
	Exec(param ExecParam) ([]byte, error)
	Inspect(format string, containerName string) ([]byte, error)

	Compose() Compose
}

type Compose interface {
	Up(ctx context.Context, services []string) error
	Down(ctx context.Context, services []string) error
}
