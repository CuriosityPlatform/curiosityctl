package usecase

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/curiosity/app/containerwaiter"
	"curiosity/pkg/curiosity/app/servicepreparer"
)

func NewUp(
	dockerClient dockerclient.Client,
	preparerFactory servicepreparer.Factory,
	waiter containerwaiter.Waiter,
) *Up {
	return &Up{
		dockerClient:    dockerClient,
		preparerFactory: preparerFactory,
		waiter:          waiter,
	}
}

type Up struct {
	dockerClient    dockerclient.Client
	preparerFactory servicepreparer.Factory
	waiter          containerwaiter.Waiter
}

func (c *Up) Execute(ctx context.Context) error {
	fmt.Println("UP")
	err := c.dockerClient.Compose().Up(ctx, []string{"db"})
	if err != nil {
		return err
	}

	fmt.Println("Wait for DB")
	err = c.waiter.WaitFor("services-db")
	if err != nil {
		downErr := c.dockerClient.Compose().Down(ctx, nil)
		if downErr != nil {
			err = errors.Wrap(err, downErr.Error())
		}
		return err
	}

	fmt.Println("Prepare db")
	preparer, err := c.preparerFactory.Preparer("services-db")
	if err != nil {
		return err
	}

	err = preparer.Prepare("services-db")
	if err != nil {
		downErr := c.dockerClient.Compose().Down(ctx, nil)
		if downErr != nil {
			err = errors.Wrap(err, downErr.Error())
		}
		return err
	}

	fmt.Println("Prepare db")
	err = c.dockerClient.Compose().Up(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Println("DONE")
	return nil
}
