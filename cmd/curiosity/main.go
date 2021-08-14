package main

import (
	stdlog "log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	appID = "curiosityctl"
)

func main() {
	err := runApp(os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(args []string) error {
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
		},
	}

	return app.Run(args)
}
