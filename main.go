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
				Name:      "init",
				Usage:     "Initialize a new project",
				UsageText: "lux init [repository]\n\n" + "Example:\n" + "  lux init myproject",
				Action:    initProject,
				Args:      true,
			},
			{
				Name:      "generate",
				Aliases:   []string{"g"},
				Usage:     "Generate a file or a directory",
				UsageText: "lux generate [command] [path/to/package]\n\n" + "Example:\n" + "  lux generate middleware middleware/acl",
				Subcommands: []*cli.Command{
					{
						Name:    "cmd",
						Aliases: []string{"e"},
						Usage:   "Generate a command",
						Args:    true,
						Action:  generateCommand,
					},
					{
						Name:    "component",
						Aliases: []string{"cp"},
						Usage:   "Generate a component",
						Args:    true,
						Action:  generateComponentsSetCommand,
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
						Aliases: []string{"co"},
						Usage:   "Generate a controller",
						Args:    true,
						Action:  generateControllerCommand,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "get",
								Aliases: []string{"g"},
								Usage:   "Generate a GET controller",
								Value:   false,
							},
							&cli.BoolFlag{
								Name:    "post",
								Aliases: []string{"p"},
								Usage:   "Generate a POST controller",
								Value:   false,
							},
							&cli.BoolFlag{
								Name:    "put",
								Aliases: []string{"u"},
								Usage:   "Generate a PUT controller",
								Value:   false,
							},
							&cli.BoolFlag{
								Name:    "patch",
								Aliases: []string{"a"},
								Usage:   "Generate a PATCH controller",
								Value:   false,
							},
							&cli.BoolFlag{
								Name:    "delete",
								Aliases: []string{"d"},
								Usage:   "Generate a DELETE controller",
								Value:   false,
							},
							&cli.BoolFlag{
								Name:    "socket",
								Aliases: []string{"w"},
								Usage:   "Generate a SOCKET controller",
								Value:   false,
							},
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
