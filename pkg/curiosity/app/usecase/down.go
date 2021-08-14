package usecase

import (
	"context"

	"curiosity/pkg/common/app/dockerclient"
)

func NewDown(dockerClient dockerclient.Client) *Down {
	return &Down{dockerClient: dockerClient}
}

type Down struct {
	dockerClient dockerclient.Client
}

func (c *Down) Execute(ctx context.Context) error {
	return c.dockerClient.Compose().Down(ctx, nil)
}
