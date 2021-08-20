package containerwaiter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/common/infrastructure/progress"
	"curiosity/pkg/curiosity/app/containerwaiter"
)

const (
	healthyStatus = "healthy"
)

func NewWaiter(client dockerclient.Client) containerwaiter.Waiter {
	return &waiter{client: client}
}

type waiter struct {
	client dockerclient.Client
}

func (w *waiter) WaitFor(ctx context.Context, containers ...string) error {
	er, _ := errgroup.WithContext(ctx)
	writer := progress.ContextWriter(ctx)

	for _, container := range containers {
		er.Go(func() error {
			eventID := fmt.Sprintf("Container %s", container)
			writer.Event(progress.WaitingEvent(eventID))

			ticker := time.NewTicker(1 * time.Second)

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					inspectResult, err := w.client.Inspect("{{.State.Health.Status}}", container)
					if err != nil {
						writer.Event(progress.ErrorEvent(eventID))
						return err
					}

					inspectResultStr := strings.TrimSpace(string(inspectResult))

					if inspectResultStr == healthyStatus {
						writer.Event(progress.HealthyEvent(eventID))
						return nil
					}
				}
			}
		})
	}

	return er.Wait()
}
