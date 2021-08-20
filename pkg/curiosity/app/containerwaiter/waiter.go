package containerwaiter

import "context"

type Waiter interface {
	WaitFor(ctx context.Context, containers ...string) error
}
