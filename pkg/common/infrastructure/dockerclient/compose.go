package dockerclient

import (
	"context"

	"github.com/pkg/errors"
)

const (
	composeCommand = "compose"
)

type composeClient struct {
	executor DockerExecutor
}

func (client *composeClient) Up(ctx context.Context, services []string) error {
	args := client.buildArgs([]string{"up", "-d"})

	if len(services) != 0 {
		args = append(args, services...)
	}

	return errors.WithStack(client.executor.PTY(ctx, args))
}

func (client *composeClient) Down(ctx context.Context, services []string) error {
	args := client.buildArgs([]string{"down"})

	if len(services) != 0 {
		args = append(args, services...)
	}

	return errors.WithStack(client.executor.PTY(ctx, args))
}

func (client *composeClient) buildArgs(args []string) []string {
	return append([]string{composeCommand}, args...)
}
