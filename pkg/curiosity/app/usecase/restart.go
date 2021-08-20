package usecase

import (
	"context"

	"github.com/pkg/errors"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/common/infrastructure/progress"
	"curiosity/pkg/curiosity/app/containerwaiter"
	"curiosity/pkg/curiosity/app/servicepreparer"
)

func NewRestart(
	dockerClient dockerclient.Client,
	preparerFactory servicepreparer.Factory,
	waiter containerwaiter.Waiter,
) *Restart {
	return &Restart{
		dockerClient:    dockerClient,
		preparerFactory: preparerFactory,
		waiter:          waiter,
	}
}

type Restart struct {
	dockerClient    dockerclient.Client
	preparerFactory servicepreparer.Factory
	waiter          containerwaiter.Waiter
}

func (c *Restart) Execute(ctx context.Context) (err error) {
	err = c.dockerClient.Compose().Down(ctx, nil)
	if err != nil {
		return err
	}

	err = c.dockerClient.Compose().Up(ctx, []string{"db"})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			downErr := c.dockerClient.Compose().Down(ctx, nil)
			if downErr != nil {
				err = errors.Wrap(err, downErr.Error())
			}
		}
	}()

	err = progress.Run(ctx, func(ctx context.Context) error {
		err = c.waiter.WaitFor(ctx, "services-db")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = progress.Run(ctx, func(ctx context.Context) error {
		preparer, err2 := c.preparerFactory.Preparer("services-db")
		if err2 != nil {
			return err2
		}

		err2 = preparer.Prepare(ctx, "services-db")
		if err2 != nil {
			return err2
		}

		return nil
	})

	err = c.dockerClient.Compose().Up(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
