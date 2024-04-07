package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func generateServiceTemplate(name string) string {
	builder := strings.Builder{}
	builder.WriteString("package ")
	builder.WriteString(name)
	builder.WriteString("\n\n")

	name = strings.ToUpper(name[:1]) + name[1:]

	builder.WriteString("type ")
	builder.WriteString(name)
	builder.WriteString("Service struct{}\n\n")

	builder.WriteString("func NewService() *")
	builder.WriteString(name)
	builder.WriteString("Service {\n")
	builder.WriteString("\treturn &")
	builder.WriteString(name)
	builder.WriteString("Service{}\n")
	builder.WriteString("}\n")

	return builder.String()
}

func generateServiceCommand(ctx *cli.Context) error {
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

	serviceTemplate := generateServiceTemplate(packageName)
	servicePath := filepath.Join(path, "service.go")

	if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
		return fmt.Errorf("service already exists")
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(servicePath, []byte(serviceTemplate), 0644); err != nil {
		return err
	}

	return nil
}
