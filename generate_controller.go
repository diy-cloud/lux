package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func generateController(name string) string {
	builder := strings.Builder{}
	builder.WriteString("package ")
	builder.WriteString(name)
	builder.WriteString("\n\n")

	builder.WriteString("import (\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/context\"\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/lux\"\n")
	builder.WriteString(")\n\n")

	builder.WriteString("type ")
	builder.WriteString(name)
	builder.WriteString("Controller struct{}\n\n")

	builder.WriteString("func NewController() *")
	builder.WriteString(name)
	builder.WriteString("Controller {\n")
	builder.WriteString("\treturn &")
	builder.WriteString(name)
	builder.WriteString("Controller{}\n")
	builder.WriteString("}\n")

	builder.WriteString("func RegisterController(c *")
	builder.WriteString(name)
	builder.WriteString("Controller, l *lux.Lux) {\n")
	builder.WriteString("\tl.RegisterController(c)\n")
	builder.WriteString("}\n")

	return builder.String()
}

func generateControllerCommand(ctx *cli.Context) error {
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

	packageName := elems[len(elems)-1]
	path := strings.Join(elems, "/")

	controller := generateController(packageName)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filepath.Join(path, "controller.go"))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(controller); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
