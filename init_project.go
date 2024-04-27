package main

import (
	"fmt"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func initProject(ctx *cli.Context) error {
	moduleName := ctx.Args().First()
	if moduleName == "" {
		return fmt.Errorf("module name is required")
	}

	if err := exec.Command("go", "mod", "init", moduleName).Run(); err != nil {
		return fmt.Errorf("failed to initialize module: %w", err)
	}

	if err := exec.Command("go", "get", "github.com/snowmerak/lux/v3").Run(); err != nil {
		return fmt.Errorf("failed to get lux: %w", err)
	}

	return nil
}
