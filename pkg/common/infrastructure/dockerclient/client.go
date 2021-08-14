package dockerclient

import (
	"fmt"
	"strings"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/common/infrastructure/executor"
)

func NewClient(dockerExecutor DockerExecutor) dockerclient.Client {
	return &client{
		executor:      dockerExecutor,
		composeClient: &composeClient{executor: dockerExecutor},
	}
}

type client struct {
	executor      DockerExecutor
	composeClient *composeClient
}

func (client *client) Exec(param dockerclient.ExecParam) ([]byte, error) {
	args := []string{
		"exec",
		param.Service,
	}

	args = append(args, strings.Split(param.Command, " ")...)

	var opts []executor.Opt

	if param.Reader != nil {
		opts = append(opts, executor.WithStdin(param.Reader))
		prevParams := append([]string{"--interactive"}, args[1:]...)
		args = append(args[0:1], prevParams...)
	}

	return client.executor.Output(args, opts...)
}

func (client *client) Inspect(format, containerName string) ([]byte, error) {
	return client.executor.Output([]string{
		"inspect",
		fmt.Sprintf("--format=%s", format),
		containerName,
	})
}

func (client *client) Compose() dockerclient.Compose {
	return client.composeClient
}
