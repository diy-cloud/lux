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

	keys := map[string]struct{}{}
	for key := range ps.Constructors {
		keys[key] = struct{}{}
	}
	for key := range ps.Updaters {
		keys[key] = struct{}{}
	}

	writer := bytes.Buffer{}
	for name := range keys {
		log.Printf("Generating %s\n", name)
		writer.Reset()
		cons := ps.Constructors[name]
		if cons == nil {
			cons = []parser.Component{}
		}

		upds := ps.Updaters[name]
		if upds == nil {
			upds = []parser.Component{}
		}

		if err := buildComponentSet(&writer, name, cons, upds); err != nil {
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

func buildComponentSet(writer *bytes.Buffer, setName string, cons []parser.Component, upds []parser.Component) error {
	pkgName := map[string]string{}
	pkgPathMap := map[string]string{}

	for i, comp := range cons {
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

	for i, comp := range upds {
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
	writer.WriteString("\tcons []any\n")
	writer.WriteString("\tupds []any\n")
	writer.WriteString("\n}\n\n")

	writer.WriteString("func New")
	writer.WriteString(structName)
	writer.WriteString("() *")
	writer.WriteString(structName)
	writer.WriteString(" {\n")
	writer.WriteString("\treturn &")
	writer.WriteString(structName)
	writer.WriteString("{\n")
	writer.WriteString("\t\tcons: []any{\n")
	for _, comp := range cons {
		writer.WriteString("\t\t\t")
		writer.WriteString(comp.PackageName)
		writer.WriteString(".")
		writer.WriteString(comp.FunctionName)
		writer.WriteString(",\n")
	}
	writer.WriteString("\t\t},\n")
	writer.WriteString("\t\tupds: []any{\n")
	for _, comp := range upds {
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
	writer.WriteString(") Constructors() []any {\n")
	writer.WriteString("\tclone := make([]any, len(s.cons))\n")
	writer.WriteString("\tcopy(clone, s.cons)\n")
	writer.WriteString("\treturn clone\n")
	writer.WriteString("}\n")

	writer.WriteString("func (s *")
	writer.WriteString(structName)
	writer.WriteString(") Updaters() []any {\n")
	writer.WriteString("\tclone := make([]any, len(s.upds))\n")
	writer.WriteString("\tcopy(clone, s.upds)\n")
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
