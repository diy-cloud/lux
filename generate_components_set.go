package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"

	"github.com/snowmerak/lux/v3/parser"
	"github.com/urfave/cli/v2"
	"golang.org/x/tools/imports"
)

const (
	ComponentsSetDirectory = "./gen/components"
	ComponentsSetPackage   = "components"
)

func generateComponentsSetCommand(ctx *cli.Context) error {
	if err := os.MkdirAll(ComponentsSetDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create components set directory: %w", err)
	}

	ps := parser.New()
	if err := ps.ParseFromRoot(); err != nil {
		return fmt.Errorf("failed to parse module: %w", err)
	}

	writer := bytes.Buffer{}
	for name, list := range ps.Components {
		log.Printf("Generating %s\n", name)
		writer.Reset()
		if err := buildComponentSet(&writer, name, list); err != nil {
			return fmt.Errorf("failed to build component set: %w", err)
		}

		file, err := prettyFormat(writer.Bytes())
		if err != nil {
			return fmt.Errorf("failed to format source: %w", err)
		}

		filePath := fmt.Sprintf("%s/%s%s", ComponentsSetDirectory, name, GeneratedFileSuffix)
		if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		log.Printf("Generated %s\n", filePath)
	}

	return nil
}

func buildComponentSet(writer *bytes.Buffer, setName string, list []parser.Component) error {
	pkgName := map[string]string{}
	pkgPathMap := map[string]string{}

	for i, comp := range list {
		if _, ok := pkgName[comp.PackageName]; !ok {
			pkgName[comp.PackageName] = comp.PackagePath
			pkgPathMap[comp.PackagePath] = comp.PackageName
			continue
		}

		if pkgName[comp.PackageName] == comp.PackagePath {
			continue
		}

		comp.PackageName = fmt.Sprintf("%s%d", comp.PackageName, i)
		pkgName[comp.PackageName] = comp.PackagePath
		pkgPathMap[comp.PackagePath] = comp.PackageName
	}

	writer.WriteString("package ")
	writer.WriteString(ComponentsSetPackage)
	writer.WriteString("\n\n")

	writer.WriteString("import (\n")
	for name, path := range pkgName {
		writer.WriteString("\t")
		writer.WriteString(name)
		writer.WriteString(" \"")
		writer.WriteString(path)
		writer.WriteString("\"\n")
	}
	writer.WriteString(")\n\n")

	structName := strings.ToUpper(setName[:1]) + setName[1:]
	writer.WriteString("type ")
	writer.WriteString(structName)
	writer.WriteString(" struct {\n")
	writer.WriteString("\tlist []any")
	writer.WriteString("\n}\n\n")

	writer.WriteString("func New")
	writer.WriteString(structName)
	writer.WriteString("() *")
	writer.WriteString(structName)
	writer.WriteString(" {\n")
	writer.WriteString("\treturn &")
	writer.WriteString(structName)
	writer.WriteString("{\n")
	writer.WriteString("\t\tlist: []any{\n")
	for _, comp := range list {
		writer.WriteString("\t\t\t")
		writer.WriteString(comp.PackageName)
		writer.WriteString(".")
		writer.WriteString(comp.FunctionName)
		writer.WriteString(",\n")
	}
	writer.WriteString("\t\t},\n")
	writer.WriteString("\t}\n")
	writer.WriteString("}\n\n")

	writer.WriteString("func (s *")
	writer.WriteString(structName)
	writer.WriteString(") List() []any {\n")
	writer.WriteString("\tclone := make([]any, len(s.list))\n")
	writer.WriteString("\tcopy(clone, s.list)\n")
	writer.WriteString("\treturn clone\n")
	writer.WriteString("}\n")

	return nil
}

func prettyFormat(file []byte) ([]byte, error) {
	file, err := format.Source(file)
	if err != nil {
		return nil, fmt.Errorf("failed to format source: %w", err)
	}

	file, err = imports.Process("", file, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to format imports: %w", err)
	}

	return file, nil
}
