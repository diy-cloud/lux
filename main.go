package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "lux",
		Usage: "A simple web framework",
		Commands: []*cli.Command{
			{
				Name: "init",
				Args: true,
			},
			{
				Name:      "generate",
				Aliases:   []string{"g"},
				Usage:     "Generate a file or a directory",
				UsageText: "lux generate [command] [path/to/package]\n\n" + "Example:\n" + "  lux generate middleware middleware/acl",
				Subcommands: []*cli.Command{
					{
						Name:   "cmd",
						Usage:  "Generate a command",
						Args:   true,
						Action: generateCommand,
					},
					{
						Name:    "middleware",
						Aliases: []string{"m"},
						Usage:   "Generate a middleware",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:     "request",
								Aliases:  []string{"r"},
								Usage:    "Generate a request middleware",
								Value:    false,
								Category: "middleware",
							},
							&cli.BoolFlag{
								Name:     "response",
								Aliases:  []string{"s"},
								Usage:    "Generate a response middleware",
								Value:    false,
								Category: "middleware",
							},
						},
						Args:   true,
						Action: generateMiddlewareCommand,
					},
					{
						Name:    "service",
						Aliases: []string{"s"},
						Usage:   "Generate a service",
						Args:    true,
						Action:  generateServiceCommand,
					},
					{
						Name:    "controller",
						Aliases: []string{"c"},
						Usage:   "Generate a controller",
						Args:    true,
						Action:  generateControllerCommand,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
