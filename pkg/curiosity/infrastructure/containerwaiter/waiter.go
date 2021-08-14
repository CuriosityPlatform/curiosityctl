package containerwaiter

import (
	"fmt"
	"time"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/curiosity/app/containerwaiter"
)

func NewWaiter(client dockerclient.Client) containerwaiter.Waiter {
	return &waiter{client: client}
}

type waiter struct {
	client dockerclient.Client
}

func (w *waiter) WaitFor(container ...string) error {
	fmt.Println("Sleep for 30 sec...")
	time.Sleep(time.Second * 30)
	fmt.Println("Awake...")
	return nil
}
