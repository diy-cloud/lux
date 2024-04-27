package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/snowmerak/lux/v3/parser"
	"github.com/urfave/cli/v2"
)

func generateComponentsSetCommand(ctx *cli.Context) error {
	packagePath := ctx.Args().First()
	elems := strings.Split(packagePath, "/")
	for i := range elems {
		elems[i] = strings.TrimSpace(elems[i])
		if elems[i] == "" {
			elems = append(elems[:i], elems[i+1:]...)
		}
	}

	if len(elems) == 0 {
		return fmt.Errorf("invalid package path, must not be empty")
	}

	path := strings.Join(elems, "/")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create package directory: %w", err)
	}

	ps := parser.New()
	if err := ps.ParseFromRoot(); err != nil {
		return fmt.Errorf("failed to parse module: %w", err)
	}

	for k, comp := range ps.Components {
		fmt.Println(k, comp)
	}

	return nil
}
