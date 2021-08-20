package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

const (
	appID = "curiosityctl"
)

func main() {
	ctx := context.Background()

	ctx = subscribeForKillSignals(ctx)

	err := runApp(ctx, os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(ctx context.Context, args []string) error {
	app := &cli.App{
		Name:    appID,
		Version: "1.0",
		Commands: []*cli.Command{
			{
				Name:   "up",
				Action: executeUp,
			},
			{
				Name:   "down",
				Action: executeDown,
			},
			{
				Name:   "restart",
				Action: executeRestart,
			},
		},
	}

	return app.RunContext(ctx, args)
}

func subscribeForKillSignals(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		<-ch
		fmt.Println("Cancelled")
		cancel()
	}()

	return ctx
}
