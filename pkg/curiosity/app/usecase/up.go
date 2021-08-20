package usecase

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/common/infrastructure/progress"
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
	err := c.dockerClient.Compose().Up(ctx, []string{"db"})
	if err != nil {
		return err
	}

	err = progress.Run(ctx, func(ctx2 context.Context) error {
		err = c.waiter.WaitFor(ctx, "services-db")
		if err != nil {
			downErr := c.dockerClient.Compose().Down(ctx2, nil)
			if downErr != nil {
				err = errors.Wrap(err, downErr.Error())
			}
			return err
		}

		fmt.Println("Prepare db")
		preparer, err2 := c.preparerFactory.Preparer("services-db")
		if err2 != nil {
			return err2
		}

		err2 = preparer.Prepare(ctx2, "services-db")
		if err2 != nil {
			downErr := c.dockerClient.Compose().Down(ctx2, nil)
			if downErr != nil {
				err2 = errors.Wrap(err2, downErr.Error())
			}
			return err2
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = c.dockerClient.Compose().Up(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Println("DONE")
	return nil
}
