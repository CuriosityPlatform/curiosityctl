package usecase

import (
	"context"

	"github.com/pkg/errors"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/common/infrastructure/progress"
	"curiosity/pkg/curiosity/app/compose"
	"curiosity/pkg/curiosity/app/containerwaiter"
	"curiosity/pkg/curiosity/app/servicepreparer"
)

func NewUp(
	dockerClient dockerclient.Client,
	preparerFactory servicepreparer.Factory,
	waiter containerwaiter.Waiter,
	project compose.Project,
) *Up {
	return &Up{
		dockerClient:    dockerClient,
		preparerFactory: preparerFactory,
		waiter:          waiter,
		project:         project,
	}
}

type Up struct {
	dockerClient    dockerclient.Client
	preparerFactory servicepreparer.Factory
	waiter          containerwaiter.Waiter
	project         compose.Project
}

func (c *Up) Execute(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			downErr := c.dockerClient.Compose().Down(ctx, nil)
			if downErr != nil {
				err = errors.Wrap(err, downErr.Error())
			}
		}
	}()

	var awaitableServices []compose.Service
	awaitableServices, err = c.project.AwaitableServices()
	if err != nil {
		return err
	}

	if len(awaitableServices) > 0 {
		err = c.waitServices(ctx, awaitableServices)
		if err != nil {
			return err
		}
	}

	var bootableServices []compose.Service
	bootableServices, err = c.project.BootableServices()
	if err != nil {
		return err
	}

	if len(bootableServices) > 0 {
		err = c.bootServices(ctx, bootableServices)
		if err != nil {
			return err
		}
	}

	err = c.dockerClient.Compose().Up(ctx, nil)
	return err
}

func (c *Up) waitServices(ctx context.Context, awaitableServices []compose.Service) (err error) {
	servicesNames := make([]string, 0, len(awaitableServices))
	containerNames := make([]string, 0, len(awaitableServices))
	for _, service := range awaitableServices {
		servicesNames = append(servicesNames, service.Name())
		containerNames = append(containerNames, service.ContainerName())
	}

	err = c.dockerClient.Compose().Up(ctx, servicesNames)
	if err != nil {
		return err
	}

	return progress.Run(ctx, func(ctx context.Context) error {
		return c.waiter.WaitFor(ctx, containerNames...)
	})
}

func (c Up) bootServices(ctx context.Context, bootableServices []compose.Service) error {
	containerNames := make([]string, 0, len(bootableServices))
	for _, service := range bootableServices {
		containerNames = append(containerNames, service.ContainerName())
	}

	return progress.Run(ctx, func(ctx context.Context) error {
		g := servicepreparer.NewPreparersGroup(c.preparerFactory)

		return g.Start(ctx, containerNames)
	})
}
