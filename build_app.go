package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func buildAppCommand(ctx *cli.Context) error {
	packageName := ctx.Args().First()
	if packageName == "" {
		return fmt.Errorf("package name is required")
	}

	packageName = strings.TrimLeft(packageName, "./")
	packageName = strings.TrimRight(packageName, "./")

	base := filepath.Base(packageName)

	if err := generateComponentsSetCommand(ctx); err != nil {
		return fmt.Errorf("failed to generate components set: %w", err)
	}

	if err := exec.Command("go", "build", "-o", base, "./"+packageName+"/.").Run(); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}
