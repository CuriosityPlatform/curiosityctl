package containerwaiter

type Waiter interface {
	WaitFor(container ...string) error
}
