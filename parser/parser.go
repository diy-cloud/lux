package parser

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
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
	RootPath     string
	ModuleName   string
	Constructors map[string][]Component
	Updaters     map[string][]Component
}

func New() *Parser {
	return &Parser{
		Constructors: map[string][]Component{},
		Updaters:     map[string][]Component{},
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
		return fmt.Errorf("failed to read module: %w", err)
	}

	if err := p.ParseFromPath(p.RootPath); err != nil {
		return fmt.Errorf("failed to parse from path: %w", err)
	}

	for _, v := range p.Constructors {
		for i := range v {
			v[i].PackagePath = filepath.Join(p.ModuleName, v[i].PackagePath)
		}
	}

	for _, v := range p.Updaters {
		for i := range v {
			v[i].PackagePath = filepath.Join(p.ModuleName, v[i].PackagePath)
		}
	}

	return nil
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

		log.Printf("Parsing %s\n", path)

		packagePath := strings.TrimPrefix(path, p.RootPath)
		packagePath = strings.TrimPrefix(packagePath, "/")
		packagePath = filepath.Dir(packagePath)

		comps, upds, err := p.ParseFile(path, packagePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		for k, v := range comps {
			p.Constructors[k] = append(p.Constructors[k], v...)
		}

		for k, v := range upds {
			p.Updaters[k] = append(p.Updaters[k], v...)
		}

		return nil
	})
}

func (p *Parser) ParseFile(path string, packagePath string) (map[string][]Component, map[string][]Component, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file: %w", err)
	}

	packageName := ""
	constructors := map[string][]Component{}
	updaters := map[string][]Component{}

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
				consSets := []string{}
				updSets := []string{}
				if x.Doc == nil {
					return true
				}

				for _, c := range x.Doc.List {
					c.Text = strings.ToLower(strings.TrimSpace(c.Text))

					if strings.HasPrefix(c.Text, "// lux:cons") {
						isComp = true
						sp := strings.Split(strings.TrimPrefix(c.Text, "// lux:cons"), " ")
						for i := range sp {
							sp[i] = strings.TrimSpace(sp[i])
						}
						consSets = append(consSets, sp...)
						break
					}

					if strings.HasPrefix(c.Text, "// lux:upd") {
						isComp = true
						sp := strings.Split(strings.TrimPrefix(c.Text, "// lux:upd"), " ")
						for i := range sp {
							sp[i] = strings.TrimSpace(sp[i])
						}
						updSets = append(updSets, sp...)
						break
					}
				}

				if !isComp {
					return true
				}

				for _, consSet := range consSets {
					if consSet == "" {
						continue
					}

					constructors[consSet] = append(constructors[consSet], Component{
						PackagePath:  packagePath,
						PackageName:  packageName,
						FunctionName: x.Name.Name,
					})
				}

				for _, updSet := range updSets {
					if updSet == "" {
						continue
					}

					updaters[updSet] = append(updaters[updSet], Component{
						PackagePath:  packagePath,
						PackageName:  packageName,
						FunctionName: x.Name.Name,
					})
				}
			}
		}

		return true
	})

	return constructors, updaters, nil
}
