package containerwaiter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/common/infrastructure/progress"
	"curiosity/pkg/curiosity/app/containerwaiter"
)

const (
	healthyStatus = "healthy"

	waitSeconds = 60
)

func NewWaiter(client dockerclient.Client) containerwaiter.Waiter {
	return &waiter{client: client}
}

type waiter struct {
	client dockerclient.Client
}

func (w *waiter) WaitFor(ctx context.Context, containers ...string) error {
	eg, _ := errgroup.WithContext(ctx)
	writer := progress.ContextWriter(ctx)

	deadlineContext, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*waitSeconds))

	for _, container := range containers {
		eg.Go(func() error {
			eventID := fmt.Sprintf("Container %s", container)
			writer.Event(progress.WaitingEvent(eventID))

			ticker := time.NewTicker(1 * time.Second)

			var inspectResultStr string

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-deadlineContext.Done():
					if inspectResultStr == "" {
						writer.Event(progress.ErrorMessageEvent(eventID, "Check status manually"))
						return errors.New("check status manually")
					}

					writer.Event(progress.ErrorMessageEvent(eventID, fmt.Sprintf("State: %s", inspectResultStr)))
					return errors.New(fmt.Sprintf("current status: %s, check container logs manually", inspectResultStr))
				case <-ticker.C:
					inspectResult, err := w.client.Inspect("{{.State.Health.Status}}", container)
					if err != nil {
						writer.Event(progress.ErrorEvent(eventID))
						return err
					}

					inspectResultStr = strings.TrimSpace(string(inspectResult))

					if inspectResultStr == healthyStatus {
						writer.Event(progress.HealthyEvent(eventID))
						return nil
					}
				}
			}
		})
	}

	err := eg.Wait()
	cancel()
	return err
}
