package parser

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Component struct {
	PackagePath  string
	PackageName  string
	FunctionName string
}

type Parser struct {
	RootPath   string
	ModuleName string
	Components map[string][]Component
}

func New() *Parser {
	return &Parser{
		Components: map[string][]Component{},
	}
}

func (p *Parser) ReadModule() error {
	cp, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current path: %w", err)
	}

	p.RootPath = cp

	f, err := os.Open(filepath.Join(p.RootPath, "/go.mod"))
	if err != nil {
		return fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && strings.HasPrefix(line, "module ") {
			p.ModuleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			break
		}
	}

	if p.ModuleName == "" {
		return fmt.Errorf("failed to get module name")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	return nil
}

func (p *Parser) ParseFromRoot() error {
	if err := p.ReadModule(); err != nil {
		return err
	}

	return p.ParseFromPath(p.RootPath)
}

func (p *Parser) ParseFromPath(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		packagePath := strings.TrimPrefix(path, p.RootPath)
		packagePath = strings.TrimPrefix(packagePath, "/")
		packagePath = filepath.Dir(packagePath)

		comps, err := p.ParseFile(path, packagePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		for k, v := range comps {
			p.Components[k] = append(p.Components[k], v...)
		}

		return nil
	})
}

func (p *Parser) ParseFile(path string, packagePath string) (map[string][]Component, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	packageName := ""
	components := map[string][]Component{}

	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		switch x := n.(type) {
		case *ast.File:
			packageName = x.Name.Name
		case *ast.FuncDecl:
			if x.Recv == nil {
				isComp := false
				compSets := []string{}
				if x.Doc == nil {
					return true
				}

				for _, c := range x.Doc.List {
					c.Text = strings.ToLower(strings.TrimSpace(c.Text))
					if strings.HasPrefix(c.Text, "// lux:comp") {
						isComp = true
						sp := strings.Split(strings.TrimPrefix(c.Text, "// lux:comp"), " ")
						for i := range sp {
							sp[i] = strings.TrimSpace(sp[i])
						}
						compSets = append(compSets, sp...)
						break
					}
				}

				if !isComp {
					return true
				}

				for _, compSet := range compSets {
					if compSet == "" {
						continue
					}

					components[compSet] = append(components[compSet], Component{
						PackagePath:  packagePath,
						PackageName:  packageName,
						FunctionName: x.Name.Name,
					})
				}
			}
		}

		return true
	})

	return components, nil
}
