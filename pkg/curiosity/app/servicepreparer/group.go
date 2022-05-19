package servicepreparer

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func NewPreparersGroup(factory Factory) *PreparersGroup {
	return &PreparersGroup{factory: factory}
}

type PreparersGroup struct {
	factory Factory
}

func (g *PreparersGroup) Start(ctx context.Context, services []string) error {
	eg, _ := errgroup.WithContext(ctx)

	preparers := map[string]ServicePreparer{}
	for _, service := range services {
		preparer, err := g.factory.Preparer(service)
		if err != nil {
			return err
		}

		preparers[service] = preparer
	}

	for service, preparer := range preparers {
		s, p := service, preparer
		eg.Go(func() error {
			return p.Prepare(ctx, s)
		})
	}

	return eg.Wait()
}
