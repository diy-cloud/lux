package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func generateControllerMetadata(name string) string {
	builder := strings.Builder{}
	builder.WriteString("package ")
	builder.WriteString(name)
	builder.WriteString("\n\n")

	builder.WriteString("const (\n")
	builder.WriteString("\tRoute = \"/write/your/route\"\n")
	builder.WriteString(")\n")

	return builder.String()
}

func generateController(name string, method string) string {
	method = strings.ToUpper(method[:1]) + strings.ToLower(method[1:])

	builder := strings.Builder{}
	builder.WriteString("package ")
	builder.WriteString(name)
	builder.WriteString("\n\n")

	builder.WriteString("import (\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/v3/context\"\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/v3/controller\"\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/v3/lux\"\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/v3/middleware\"\n")
	builder.WriteString(")\n\n")

	builder.WriteString("type ")
	builder.WriteString(method)
	builder.WriteString("Controller struct{\n")
	builder.WriteString("\trequestMiddlewares []middleware.Request\n")
	builder.WriteString("\tresponseMiddlewares []middleware.Response\n")
	builder.WriteString("\thandler controller.Handler")
	builder.WriteString("}\n\n")

	builder.WriteString("func New")
	builder.WriteString(method)
	builder.WriteString("Controller() *")
	builder.WriteString(method)
	builder.WriteString("Controller {\n")
	builder.WriteString("\treturn &")
	builder.WriteString(method)
	builder.WriteString("Controller{\n")
	builder.WriteString("\t\trequestMiddlewares: []middleware.Request{},\n")
	builder.WriteString("\t\tresponseMiddlewares: []middleware.Response{},\n")
	builder.WriteString("\t\thandler: func(lc *context.LuxContext) error {\n")
	builder.WriteString("\t\t\treturn lc.ReplyString(\"Hello, World!\")\n")
	builder.WriteString("\t\t},\n")
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	builder.WriteString("func RegisterController(c *")
	builder.WriteString(method)
	builder.WriteString("Controller, l *lux.Lux) {\n")
	builder.WriteString("\tl.RegisterController(Route, controller.")
	builder.WriteString(strings.ToUpper(method))
	builder.WriteString(", controller.Controller{\n")
	builder.WriteString("\t\tRequestMiddlewares: c.requestMiddlewares,\n")
	builder.WriteString("\t\tHandler: c.handler,\n")
	builder.WriteString("\t\tResponseMiddlewares: c.responseMiddlewares,\n")
	builder.WriteString("\t})\n")
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

	methods := make([]string, 0)
	if ctx.Bool("get") {
		methods = append(methods, "get")
	}
	if ctx.Bool("post") {
		methods = append(methods, "post")
	}
	if ctx.Bool("put") {
		methods = append(methods, "put")
	}
	if ctx.Bool("patch") {
		methods = append(methods, "patch")
	}
	if ctx.Bool("delete") {
		methods = append(methods, "delete")
	}

	packageName := elems[len(elems)-1]
	path := strings.Join(elems, "/")

	metadataFilePath := filepath.Join(path, "metadata.controller.go")
	if _, err := os.Stat(metadataFilePath); os.IsNotExist(err) {
		metadata := generateControllerMetadata(packageName)

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		file, err := os.Create(metadataFilePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		if _, err := file.WriteString(metadata); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	for _, method := range methods {
		controller := generateController(packageName, method)

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		controllerFilePath := filepath.Join(path, method+".controller.go")

		if _, err := os.Stat(controllerFilePath); !os.IsNotExist(err) {
			return fmt.Errorf("controller already exists")
		}

		file, err := os.Create(controllerFilePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()

		if _, err := file.WriteString(controller); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}
