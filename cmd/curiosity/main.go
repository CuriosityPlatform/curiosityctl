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
		Usage:   "Tool to manage, init curiosity environment",
		Commands: []*cli.Command{
			{
				Name:        "up",
				Action:      executeUp,
				Usage:       "UP environment",
				Description: "UP curiosity environment via docker-compose.\nUses configuration based on docker context.\nRequired installed 'Platform'",
			},
			{
				Name:        "down",
				Action:      executeDown,
				Usage:       "DOWN environment",
				Description: "DOWN curiosity environment",
			},
			{
				Name:        "restart",
				Action:      executeRestart,
				Usage:       "RESTART environment",
				Description: "RESTART curiosity environment",
			},
			{
				Name:            "install",
				Usage:           "Install curiosity environment modules",
				Description:     "Install curiosity environment modules, like Platform and one or all services",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
					{
						Name:     "platform",
						Category: "install",
					},
				},
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
