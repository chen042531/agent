package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) ExtractFunctions(filePath string) ([]FunctionSignature, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	var functions []FunctionSignature
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			sig := a.extractFunctionSignature(fn, fset)
			functions = append(functions, sig)
		}
		return true
	})

	return functions, nil
}

func (a *Analyzer) extractFunctionSignature(fn *ast.FuncDecl, fset *token.FileSet) FunctionSignature {
	sig := FunctionSignature{
		Name:       fn.Name.Name,
		LineNumber: fset.Position(fn.Pos()).Line,
	}

	// Extract receiver (for methods)
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		sig.Receiver = a.exprToString(fn.Recv.List[0].Type)
	}

	// Extract parameters
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			typeStr := a.exprToString(field.Type)
			if len(field.Names) == 0 {
				sig.Parameters = append(sig.Parameters, Parameter{
					Name: "",
					Type: typeStr,
				})
			} else {
				for _, name := range field.Names {
					sig.Parameters = append(sig.Parameters, Parameter{
						Name: name.Name,
						Type: typeStr,
					})
				}
			}
		}
	}

	// Extract return types
	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			typeStr := a.exprToString(field.Type)
			sig.Returns = append(sig.Returns, typeStr)
		}
	}

	// Extract doc comment
	if fn.Doc != nil {
		sig.DocComment = fn.Doc.Text()
	}

	return sig
}

func (a *Analyzer) ExtractTypes(filePath string) ([]TypeDefinition, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	var types []TypeDefinition
	ast.Inspect(node, func(n ast.Node) bool {
		if gen, ok := n.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					typeDef := a.extractTypeDefinition(typeSpec, gen, fset)
					types = append(types, typeDef)
				}
			}
		}
		return true
	})

	return types, nil
}

func (a *Analyzer) extractTypeDefinition(typeSpec *ast.TypeSpec, gen *ast.GenDecl, fset *token.FileSet) TypeDefinition {
	typeDef := TypeDefinition{
		Name:       typeSpec.Name.Name,
		LineNumber: fset.Position(typeSpec.Pos()).Line,
	}

	// Extract doc comment
	if gen.Doc != nil {
		typeDef.DocComment = gen.Doc.Text()
	}

	switch t := typeSpec.Type.(type) {
	case *ast.StructType:
		typeDef.Kind = "struct"
		if t.Fields != nil {
			for _, field := range t.Fields.List {
				typeStr := a.exprToString(field.Type)
				var tag string
				if field.Tag != nil {
					tag = field.Tag.Value
				}

				if len(field.Names) == 0 {
					// Embedded field
					typeDef.Fields = append(typeDef.Fields, FieldInfo{
						Name: typeStr,
						Type: typeStr,
						Tag:  tag,
					})
				} else {
					for _, name := range field.Names {
						typeDef.Fields = append(typeDef.Fields, FieldInfo{
							Name: name.Name,
							Type: typeStr,
							Tag:  tag,
						})
					}
				}
			}
		}

	case *ast.InterfaceType:
		typeDef.Kind = "interface"
		if t.Methods != nil {
			for _, method := range t.Methods.List {
				if len(method.Names) > 0 {
					typeDef.Methods = append(typeDef.Methods, method.Names[0].Name)
				}
			}
		}

	default:
		typeDef.Kind = "alias"
	}

	return typeDef
}

func (a *Analyzer) AnalyzePackage(packagePath string) (*PackageAnalysis, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse directory error: %w", err)
	}

	analysis := &PackageAnalysis{
		Functions: make([]FunctionSignature, 0),
		Types:     make([]TypeDefinition, 0),
		Imports:   make([]string, 0),
	}

	for pkgName, pkg := range pkgs {
		analysis.PackageName = pkgName

		for _, file := range pkg.Files {
			// Extract imports
			for _, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, "\"")
				analysis.Imports = append(analysis.Imports, importPath)
			}

			// Extract functions
			ast.Inspect(file, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok {
					// Only include exported functions
					if fn.Name.IsExported() {
						sig := a.extractFunctionSignature(fn, fset)
						analysis.Functions = append(analysis.Functions, sig)
					}
				}
				return true
			})

			// Extract types
			ast.Inspect(file, func(n ast.Node) bool {
				if gen, ok := n.(*ast.GenDecl); ok {
					for _, spec := range gen.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							// Only include exported types
							if typeSpec.Name.IsExported() {
								typeDef := a.extractTypeDefinition(typeSpec, gen, fset)
								analysis.Types = append(analysis.Types, typeDef)
							}
						}
					}
				}
				return true
			})
		}
	}

	return analysis, nil
}

func (a *Analyzer) exprToString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + a.exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + a.exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + a.exprToString(t.Key) + "]" + a.exprToString(t.Value)
	case *ast.SelectorExpr:
		return a.exprToString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func"
	case *ast.ChanType:
		return "chan " + a.exprToString(t.Value)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func GetRelativePath(basePath, targetPath string) string {
	relPath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return targetPath
	}
	return relPath
}
