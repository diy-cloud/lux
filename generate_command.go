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

	packageName := elems[len(elems)-1]
	path := strings.Join(elems, "/")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create package directory: %w", err)
	}

	builder := strings.Builder{}
	builder.WriteString("package ")
	builder.WriteString(packageName)
	builder.WriteString("\n\n")

	builder.WriteString("import (\n")
	builder.WriteString("\t\"context\"\n")
	builder.WriteString("\t\"os/signal\"\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/lux\"\n")
	builder.WriteString("\t\"github.com/snowmerak/lux/provider\"\n")
	builder.WriteString(")\n\n")

	builder.WriteString("func main() {\n")
	builder.WriteString("\tctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)\n")
	builder.WriteString("\tdefer cancel()\n")
	builder.WriteString("\tp := provider.New()\n")
	builder.WriteString("\tp.Register(lux.New)\n\n")
	builder.WriteString("\tif err := p.Construct(ctx); err != nil {\n")
	builder.WriteString("\t\tlog.Fatal(err)\n")
	builder.WriteString("\t}\n\n")
	builder.WriteString("\tif err := provider.JustRun(p, lux.ListenAndServe1); err != nil {\n")
	builder.WriteString("\t\tlog.Fatal(err)\n")
	builder.WriteString("\t}\n\n")
	builder.WriteString("\t<-ctx.Done()\n")
	builder.WriteString("}\n")

	if err := os.WriteFile(path+"/main.go", []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	return nil
}
