package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func generateMiddlewareRequest(name string) string {
	builder := strings.Builder{}
	builder.WriteString("package ")
	builder.WriteString(name)
	builder.WriteString("\n\n")
	builder.WriteString("import (\n")
	builder.WriteString("\t\"net/http\"\n\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/v3/context\"\n")
	builder.WriteString(")\n\n")

	builder.WriteString("func NewRequest() func(*context.LuxContext) (*context.LuxContext, int) {\n")
	builder.WriteString("\treturn func(ctx *context.LuxContext) (*context.LuxContext, int) {\n")
	builder.WriteString("\t\treturn ctx, http.StatusOK\n")
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	return builder.String()
}

func generateMiddlewareResponse(name string) string {
	builder := strings.Builder{}
	builder.WriteString("package ")
	builder.WriteString(name)
	builder.WriteString("\n\n")
	builder.WriteString("import (\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/context\"\n")
	builder.WriteString(")\n\n")

	builder.WriteString("func NewResponse() func(*context.LuxContext) (*context.LuxContext, error) {\n")
	builder.WriteString("\treturn func(ctx *context.LuxContext) (*context.LuxContext, error) {\n")
	builder.WriteString("\t\treturn ctx, nil\n")
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	return builder.String()
}

func generateMiddlewareCommand(ctx *cli.Context) error {
	packagePath := ctx.Args().First()

	isRequest := ctx.Bool("request")
	isResponse := ctx.Bool("response")
	if !isRequest && !isResponse {
		isRequest = true
		isResponse = true
	}

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

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if isRequest {
		requestContent := generateMiddlewareRequest(packageName)
		requestPath := filepath.Join(path, "request.go")
		if err := os.WriteFile(requestPath, []byte(requestContent), 0644); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}

	if isResponse {
		responseContent := generateMiddlewareResponse(packageName)
		responsePath := filepath.Join(path, "response.go")
		if err := os.WriteFile(responsePath, []byte(responseContent), 0644); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}

	return nil
}
