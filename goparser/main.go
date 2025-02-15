package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

// PackageInfo stores information about a Go package
type PackageInfo struct {
	Name    string
	Structs map[string]*StructInfo
	Funcs   map[string]*FuncInfo
}

// StructInfo stores information about a struct
type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

// FieldInfo stores information about a struct field
type FieldInfo struct {
	Name string
	Type string
}

// FuncInfo stores information about a function
type FuncInfo struct {
	Name       string
	Receiver   string
	Parameters []ParameterInfo
	Returns    []string
}

// ParameterInfo stores information about a function parameter
type ParameterInfo struct {
	Name string
	Type string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: program <directory>")
	}

	dir := os.Args[1]
	packages := make(map[string]*PackageInfo)

	// Create a new file set for position information
	fset := token.NewFileSet()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse the Go source file
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("error parsing %s: %v", path, err)
		}

		// Get or create package info
		pkgInfo, exists := packages[node.Name.Name]
		if !exists {
			pkgInfo = &PackageInfo{
				Name:    node.Name.Name,
				Structs: make(map[string]*StructInfo),
				Funcs:   make(map[string]*FuncInfo),
			}
			packages[node.Name.Name] = pkgInfo
		}

		// Inspect the AST
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				// Handle struct types
				if structType, ok := x.Type.(*ast.StructType); ok {
					structInfo := &StructInfo{
						Name:   x.Name.Name,
						Fields: make([]FieldInfo, 0),
					}

					for _, field := range structType.Fields.List {
						fieldType := fmt.Sprintf("%s", field.Type)
						for _, name := range field.Names {
							structInfo.Fields = append(structInfo.Fields, FieldInfo{
								Name: name.Name,
								Type: fieldType,
							})
						}
					}

					pkgInfo.Structs[structInfo.Name] = structInfo
				}

			case *ast.FuncDecl:
				// Handle functions and methods
				funcInfo := &FuncInfo{
					Name:       x.Name.Name,
					Parameters: make([]ParameterInfo, 0),
					Returns:    make([]string, 0),
				}

				// Get receiver type for methods
				if x.Recv != nil && len(x.Recv.List) > 0 {
					recv := x.Recv.List[0]
					if len(recv.Names) > 0 {
						funcInfo.Receiver = fmt.Sprintf("%s", recv.Type)
					}
				}

				// Get parameters
				if x.Type.Params != nil {
					for _, param := range x.Type.Params.List {
						paramType := fmt.Sprintf("%s", param.Type)
						for _, name := range param.Names {
							funcInfo.Parameters = append(funcInfo.Parameters, ParameterInfo{
								Name: name.Name,
								Type: paramType,
							})
						}
					}
				}

				// Get return types
				if x.Type.Results != nil {
					for _, ret := range x.Type.Results.List {
						funcInfo.Returns = append(funcInfo.Returns, fmt.Sprintf("%s", ret.Type))
					}
				}

				pkgInfo.Funcs[funcInfo.Name] = funcInfo
			}
			return true
		})

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Print the collected information
	for pkgName, pkg := range packages {
		fmt.Printf("Package: %s\n", pkgName)

		fmt.Println("\nStructs:")
		for _, str := range pkg.Structs {
			fmt.Printf("  %s\n", str.Name)
			for _, field := range str.Fields {
				fmt.Printf("    %s: %s\n", field.Name, field.Type)
			}
		}

		fmt.Println("\nFunctions:")
		for _, fn := range pkg.Funcs {
			if fn.Receiver != "" {
				fmt.Printf("  func (%s) %s(", fn.Receiver, fn.Name)
			} else {
				fmt.Printf("  func %s(", fn.Name)
			}

			for i, param := range fn.Parameters {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s %s", param.Name, param.Type)
			}
			fmt.Print(")")

			if len(fn.Returns) > 0 {
				fmt.Print(" (")
				for i, ret := range fn.Returns {
					if i > 0 {
						fmt.Print(", ")
					}
					fmt.Print(ret)
				}
				fmt.Print(")")
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
