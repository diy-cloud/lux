package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func generateCommand(ctx *cli.Context) error {
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

	builder := strings.Builder{}
	builder.WriteString("package main\n\n")

	builder.WriteString("import (\n")
	builder.WriteString("\t\"context\"\n")
	builder.WriteString("\t\"os\"\n")
	builder.WriteString("\t\"log\"\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/v3/lux\"\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/v3/provider\"\n")
	builder.WriteString(")\n\n")

	builder.WriteString("func main() {\n")
	builder.WriteString("\tctx, cancel := context.WithCancel(context.Backgroud())\n")
	builder.WriteString("\tdefer cancel()\n\n")
	builder.WriteString("\tp := provider.New()\n")
	builder.WriteString("\tif err := p.Register(lux.New, lux.GenerateListenAddress(\":8080\")); err != nil {\n\n")
	builder.WriteString("\t\tlog.Fatal(err)\n")
	builder.WriteString("\t}\n\n")
	builder.WriteString("\tif err := p.Construct(ctx); err != nil {\n")
	builder.WriteString("\t\tlog.Fatal(err)\n")
	builder.WriteString("\t}\n\n")
	builder.WriteString("\tif err := provider.JustRun(p, lux.ListenAndServe1); err != nil {\n")
	builder.WriteString("\t\tlog.Fatal(err)\n")
	builder.WriteString("\t}\n\n")
	builder.WriteString("\t<-ctx.Done()\n")
	builder.WriteString("}\n")

	cmdFilePath := path + "/main.go"

	if _, err := os.Stat(cmdFilePath); !os.IsNotExist(err) {
		return fmt.Errorf("command already exists")
	}

	if err := os.WriteFile(cmdFilePath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	return nil
}
