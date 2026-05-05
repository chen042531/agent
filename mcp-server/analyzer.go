package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type CodeAnalyzer struct {
	workspaceRoot string
	cache         map[string]interface{}
}

func NewCodeAnalyzer(workspaceRoot string) *CodeAnalyzer {
	return &CodeAnalyzer{
		workspaceRoot: workspaceRoot,
		cache:         make(map[string]interface{}),
	}
}

// ListPackages lists all Go packages with configurable detail level
func (ca *CodeAnalyzer) ListPackages(args map[string]interface{}) (*mcp.CallToolResult, error) {
	basePath := ca.getStringArg(args, "path", ".")
	level := ca.getStringArg(args, "level", "summary")

	fullPath := filepath.Join(ca.workspaceRoot, basePath)
	packages := []map[string]interface{}{}

	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() {
			return nil
		}

		// Check if directory contains Go files
		hasGoFiles := false
		entries, _ := os.ReadDir(path)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
				hasGoFiles = true
				break
			}
		}

		if !hasGoFiles {
			return nil
		}

		relPath, _ := filepath.Rel(ca.workspaceRoot, path)
		pkg := map[string]interface{}{
			"path": relPath,
		}

		if level != "minimal" {
			// Parse package to get summary info
			fset := token.NewFileSet()
			pkgs, _ := parser.ParseDir(fset, path, nil, parser.ParseComments)

			for pkgName, pkgAST := range pkgs {
				pkg["name"] = pkgName

				if level == "summary" || level == "detailed" {
					typeCount := 0
					funcCount := 0
					exports := []string{}

					for _, file := range pkgAST.Files {
						ast.Inspect(file, func(n ast.Node) bool {
							switch node := n.(type) {
							case *ast.GenDecl:
								for _, spec := range node.Specs {
									if ts, ok := spec.(*ast.TypeSpec); ok {
										typeCount++
										if ts.Name.IsExported() && level == "detailed" {
											exports = append(exports, "type "+ts.Name.Name)
										}
									}
								}
							case *ast.FuncDecl:
								funcCount++
								if node.Name.IsExported() && level == "detailed" {
									funcSig := "func " + node.Name.Name
									if node.Recv != nil {
										funcSig = "method " + node.Name.Name
									}
									exports = append(exports, funcSig)
								}
							}
							return true
						})
					}

					pkg["type_count"] = typeCount
					pkg["function_count"] = funcCount
					if level == "detailed" {
						pkg["exports"] = exports
					}
				}
				break // Only process first package in directory
			}
		}

		packages = append(packages, pkg)
		return nil
	})

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"workspace": ca.workspaceRoot,
		"packages":  packages,
		"count":     len(packages),
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// AnalyzePackage provides deep analysis of a package
func (ca *CodeAnalyzer) AnalyzePackage(args map[string]interface{}) (*mcp.CallToolResult, error) {
	packagePath := ca.getStringArg(args, "package_path", "")
	if packagePath == "" {
		return nil, fmt.Errorf("package_path is required")
	}

	include := ca.getSliceArg(args, "include", []string{"types", "functions"})
	exportedOnly := ca.getBoolArg(args, "exported_only", true)

	fullPath := filepath.Join(ca.workspaceRoot, packagePath)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, fullPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse package: %w", err)
	}

	result := map[string]interface{}{
		"package_path": packagePath,
	}

	for pkgName, pkg := range pkgs {
		result["package_name"] = pkgName

		if contains(include, "imports") {
			imports := []string{}
			for _, file := range pkg.Files {
				for _, imp := range file.Imports {
					path := strings.Trim(imp.Path.Value, `"`)
					if !contains(imports, path) {
						imports = append(imports, path)
					}
				}
			}
			result["imports"] = imports
		}

		if contains(include, "types") {
			types := []map[string]interface{}{}
			for _, file := range pkg.Files {
				ast.Inspect(file, func(n ast.Node) bool {
					if gen, ok := n.(*ast.GenDecl); ok {
						for _, spec := range gen.Specs {
							if ts, ok := spec.(*ast.TypeSpec); ok {
								if exportedOnly && !ts.Name.IsExported() {
									continue
								}

								typeInfo := map[string]interface{}{
									"name": ts.Name.Name,
									"line": fset.Position(ts.Pos()).Line,
								}

								if gen.Doc != nil {
									typeInfo["doc"] = gen.Doc.Text()
								}

								// Determine type kind
								switch t := ts.Type.(type) {
								case *ast.StructType:
									typeInfo["kind"] = "struct"
									if contains(include, "fields") && t.Fields != nil {
										fields := []map[string]interface{}{}
										for _, field := range t.Fields.List {
											fieldType := ca.exprToString(field.Type)
											var tag string
											if field.Tag != nil {
												tag = field.Tag.Value
											}
											if len(field.Names) == 0 {
												fields = append(fields, map[string]interface{}{
													"name": fieldType,
													"type": fieldType,
													"tag":  tag,
												})
											} else {
												for _, name := range field.Names {
													fields = append(fields, map[string]interface{}{
														"name": name.Name,
														"type": fieldType,
														"tag":  tag,
													})
												}
											}
										}
										typeInfo["fields"] = fields
									}
								case *ast.InterfaceType:
									typeInfo["kind"] = "interface"
									if contains(include, "methods") && t.Methods != nil {
										methods := []string{}
										for _, method := range t.Methods.List {
											if len(method.Names) > 0 {
												methods = append(methods, method.Names[0].Name)
											}
										}
										typeInfo["methods"] = methods
									}
								default:
									typeInfo["kind"] = "alias"
								}

								types = append(types, typeInfo)
							}
						}
					}
					return true
				})
			}
			result["types"] = types
		}

		if contains(include, "functions") {
			functions := []map[string]interface{}{}
			for _, file := range pkg.Files {
				ast.Inspect(file, func(n ast.Node) bool {
					if fn, ok := n.(*ast.FuncDecl); ok {
						if exportedOnly && !fn.Name.IsExported() {
							return true
						}

						funcInfo := map[string]interface{}{
							"name": fn.Name.Name,
							"line": fset.Position(fn.Pos()).Line,
						}

						if fn.Recv != nil && len(fn.Recv.List) > 0 {
							funcInfo["receiver"] = ca.exprToString(fn.Recv.List[0].Type)
						}

						// Parameters
						params := []map[string]interface{}{}
						if fn.Type.Params != nil {
							for _, field := range fn.Type.Params.List {
								typeStr := ca.exprToString(field.Type)
								if len(field.Names) == 0 {
									params = append(params, map[string]interface{}{
										"name": "",
										"type": typeStr,
									})
								} else {
									for _, name := range field.Names {
										params = append(params, map[string]interface{}{
											"name": name.Name,
											"type": typeStr,
										})
									}
								}
							}
						}
						funcInfo["parameters"] = params

						// Return types
						returns := []string{}
						if fn.Type.Results != nil {
							for _, field := range fn.Type.Results.List {
								returns = append(returns, ca.exprToString(field.Type))
							}
						}
						funcInfo["returns"] = returns

						if fn.Doc != nil {
							funcInfo["doc"] = fn.Doc.Text()
						}

						functions = append(functions, funcInfo)
					}
					return true
				})
			}
			result["functions"] = functions
		}

		break // Only process first package
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// GetFunctionDetails gets detailed info about a specific function
func (ca *CodeAnalyzer) GetFunctionDetails(args map[string]interface{}) (*mcp.CallToolResult, error) {
	filePath := ca.getStringArg(args, "file_path", "")
	functionName := ca.getStringArg(args, "function_name", "")
	includeBody := ca.getBoolArg(args, "include_body", false)

	if filePath == "" || functionName == "" {
		return nil, fmt.Errorf("file_path and function_name are required")
	}

	fullPath := filepath.Join(ca.workspaceRoot, filePath)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var result map[string]interface{}

	// Handle method notation (Type.Method)
	parts := strings.Split(functionName, ".")
	targetFunc := functionName
	targetReceiver := ""
	if len(parts) == 2 {
		targetReceiver = parts[0]
		targetFunc = parts[1]
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.Name != targetFunc {
				return true
			}

			// Check receiver if specified
			if targetReceiver != "" {
				if fn.Recv == nil || len(fn.Recv.List) == 0 {
					return true
				}
				recvType := ca.exprToString(fn.Recv.List[0].Type)
				if !strings.Contains(recvType, targetReceiver) {
					return true
				}
			}

			result = map[string]interface{}{
				"name": fn.Name.Name,
				"line": fset.Position(fn.Pos()).Line,
				"file": filePath,
			}

			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				result["receiver"] = ca.exprToString(fn.Recv.List[0].Type)
			}

			// Full signature
			sig := ca.formatFunctionSignature(fn)
			result["signature"] = sig

			if fn.Doc != nil {
				result["documentation"] = fn.Doc.Text()
			}

			// Parameters with details
			params := []map[string]interface{}{}
			if fn.Type.Params != nil {
				for _, field := range fn.Type.Params.List {
					typeStr := ca.exprToString(field.Type)
					if len(field.Names) == 0 {
						params = append(params, map[string]interface{}{
							"type": typeStr,
						})
					} else {
						for _, name := range field.Names {
							params = append(params, map[string]interface{}{
								"name": name.Name,
								"type": typeStr,
							})
						}
					}
				}
			}
			result["parameters"] = params

			// Return values
			returns := []string{}
			if fn.Type.Results != nil {
				for _, field := range fn.Type.Results.List {
					returns = append(returns, ca.exprToString(field.Type))
				}
			}
			result["returns"] = returns

			if includeBody && fn.Body != nil {
				// Extract function body source
				bodyStart := fset.Position(fn.Body.Pos())
				bodyEnd := fset.Position(fn.Body.End())

				content, _ := os.ReadFile(fullPath)
				lines := strings.Split(string(content), "\n")
				body := strings.Join(lines[bodyStart.Line-1:bodyEnd.Line], "\n")
				result["body"] = body
			}

			return false // Found it, stop searching
		}
		return true
	})

	if result == nil {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// GetTypeDetails gets detailed info about a specific type
func (ca *CodeAnalyzer) GetTypeDetails(args map[string]interface{}) (*mcp.CallToolResult, error) {
	filePath := ca.getStringArg(args, "file_path", "")
	typeName := ca.getStringArg(args, "type_name", "")
	includeMethods := ca.getBoolArg(args, "include_methods", true)

	if filePath == "" || typeName == "" {
		return nil, fmt.Errorf("file_path and type_name are required")
	}

	fullPath := filepath.Join(ca.workspaceRoot, filePath)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var result map[string]interface{}

	ast.Inspect(file, func(n ast.Node) bool {
		if gen, ok := n.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.Name == typeName {
					result = map[string]interface{}{
						"name": ts.Name.Name,
						"line": fset.Position(ts.Pos()).Line,
						"file": filePath,
					}

					if gen.Doc != nil {
						result["documentation"] = gen.Doc.Text()
					}

					switch t := ts.Type.(type) {
					case *ast.StructType:
						result["kind"] = "struct"
						if t.Fields != nil {
							fields := []map[string]interface{}{}
							for _, field := range t.Fields.List {
								fieldType := ca.exprToString(field.Type)
								var tag string
								if field.Tag != nil {
									tag = strings.Trim(field.Tag.Value, "`")
								}
								var doc string
								if field.Doc != nil {
									doc = field.Doc.Text()
								}

								if len(field.Names) == 0 {
									fields = append(fields, map[string]interface{}{
										"name":          fieldType,
										"type":          fieldType,
										"tag":           tag,
										"documentation": doc,
										"embedded":      true,
									})
								} else {
									for _, name := range field.Names {
										fields = append(fields, map[string]interface{}{
											"name":          name.Name,
											"type":          fieldType,
											"tag":           tag,
											"documentation": doc,
											"exported":      name.IsExported(),
										})
									}
								}
							}
							result["fields"] = fields
						}

					case *ast.InterfaceType:
						result["kind"] = "interface"
						if t.Methods != nil {
							methods := []map[string]interface{}{}
							for _, method := range t.Methods.List {
								if len(method.Names) > 0 {
									methodInfo := map[string]interface{}{
										"name": method.Names[0].Name,
									}
									if funcType, ok := method.Type.(*ast.FuncType); ok {
										methodInfo["signature"] = ca.formatFunctionType(funcType)
									}
									if method.Doc != nil {
										methodInfo["documentation"] = method.Doc.Text()
									}
									methods = append(methods, methodInfo)
								}
							}
							result["interface_methods"] = methods
						}

					default:
						result["kind"] = "alias"
						result["underlying_type"] = ca.exprToString(ts.Type)
					}

					return false // Found it
				}
			}
		}
		return true
	})

	if result == nil {
		return nil, fmt.Errorf("type %s not found", typeName)
	}

	// Find methods if requested
	if includeMethods {
		methods := []map[string]interface{}{}
		ast.Inspect(file, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Recv != nil && len(fn.Recv.List) > 0 {
					recvType := ca.exprToString(fn.Recv.List[0].Type)
					if strings.Contains(recvType, typeName) {
						methodInfo := map[string]interface{}{
							"name":      fn.Name.Name,
							"line":      fset.Position(fn.Pos()).Line,
							"signature": ca.formatFunctionSignature(fn),
							"exported":  fn.Name.IsExported(),
						}
						if fn.Doc != nil {
							methodInfo["documentation"] = fn.Doc.Text()
						}
						methods = append(methods, methodInfo)
					}
				}
			}
			return true
		})
		if len(methods) > 0 {
			result["methods"] = methods
		}
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// SearchCode searches for symbols in the codebase
func (ca *CodeAnalyzer) SearchCode(args map[string]interface{}) (*mcp.CallToolResult, error) {
	query := ca.getStringArg(args, "query", "")
	kind := ca.getStringArg(args, "kind", "all")
	maxResults := ca.getIntArg(args, "max_results", 20)

	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	results := []map[string]interface{}{}
	queryLower := strings.ToLower(query)

	err := filepath.Walk(ca.workspaceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(ca.workspaceRoot, path)

		ast.Inspect(file, func(n ast.Node) bool {
			if len(results) >= maxResults {
				return false
			}

			switch node := n.(type) {
			case *ast.FuncDecl:
				if (kind == "all" || kind == "function" || kind == "method") &&
					strings.Contains(strings.ToLower(node.Name.Name), queryLower) {

					resultKind := "function"
					if node.Recv != nil {
						resultKind = "method"
					}

					results = append(results, map[string]interface{}{
						"kind":      resultKind,
						"name":      node.Name.Name,
						"file":      relPath,
						"line":      fset.Position(node.Pos()).Line,
						"signature": ca.formatFunctionSignature(node),
					})
				}

			case *ast.GenDecl:
				for _, spec := range node.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if strings.Contains(strings.ToLower(ts.Name.Name), queryLower) {
							resultKind := "type"
							switch ts.Type.(type) {
							case *ast.InterfaceType:
								if kind == "all" || kind == "interface" || kind == "type" {
									resultKind = "interface"
								} else {
									return true
								}
							case *ast.StructType:
								if kind != "all" && kind != "type" {
									return true
								}
								resultKind = "struct"
							default:
								if kind != "all" && kind != "type" {
									return true
								}
							}

							results = append(results, map[string]interface{}{
								"kind": resultKind,
								"name": ts.Name.Name,
								"file": relPath,
								"line": fset.Position(ts.Pos()).Line,
							})
						}
					}
				}
			}
			return true
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	resultData := map[string]interface{}{
		"query":   query,
		"results": results,
		"count":   len(results),
	}

	data, _ := json.MarshalIndent(resultData, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// GetFileOverview provides a structured overview of a file
func (ca *CodeAnalyzer) GetFileOverview(args map[string]interface{}) (*mcp.CallToolResult, error) {
	filePath := ca.getStringArg(args, "file_path", "")
	level := ca.getStringArg(args, "level", "signatures")

	if filePath == "" {
		return nil, fmt.Errorf("file_path is required")
	}

	fullPath := filepath.Join(ca.workspaceRoot, filePath)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	result := map[string]interface{}{
		"file":         filePath,
		"package_name": file.Name.Name,
	}

	// Imports
	imports := []string{}
	for _, imp := range file.Imports {
		imports = append(imports, strings.Trim(imp.Path.Value, `"`))
	}
	result["imports"] = imports

	// Types
	types := []map[string]interface{}{}
	functions := []map[string]interface{}{}

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					typeInfo := map[string]interface{}{
						"name": ts.Name.Name,
						"line": fset.Position(ts.Pos()).Line,
					}

					if level == "signatures" || level == "full" {
						switch t := ts.Type.(type) {
						case *ast.StructType:
							typeInfo["kind"] = "struct"
							if t.Fields != nil {
								typeInfo["field_count"] = t.Fields.NumFields()
							}
						case *ast.InterfaceType:
							typeInfo["kind"] = "interface"
							if t.Methods != nil {
								typeInfo["method_count"] = t.Methods.NumFields()
							}
						default:
							typeInfo["kind"] = "alias"
						}
					}

					if level == "full" && node.Doc != nil {
						typeInfo["documentation"] = node.Doc.Text()
					}

					types = append(types, typeInfo)
				}
			}

		case *ast.FuncDecl:
			funcInfo := map[string]interface{}{
				"name": node.Name.Name,
				"line": fset.Position(node.Pos()).Line,
			}

			if level == "signatures" || level == "full" {
				funcInfo["signature"] = ca.formatFunctionSignature(node)
				if node.Recv != nil {
					funcInfo["kind"] = "method"
				} else {
					funcInfo["kind"] = "function"
				}
			}

			if level == "full" && node.Doc != nil {
				funcInfo["documentation"] = node.Doc.Text()
			}

			functions = append(functions, funcInfo)
		}
		return true
	})

	result["types"] = types
	result["functions"] = functions
	result["type_count"] = len(types)
	result["function_count"] = len(functions)

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// FindDependencies finds dependencies for a package or symbol
func (ca *CodeAnalyzer) FindDependencies(args map[string]interface{}) (*mcp.CallToolResult, error) {
	target := ca.getStringArg(args, "target", "")
	direction := ca.getStringArg(args, "direction", "reverse")

	if target == "" {
		return nil, fmt.Errorf("target is required")
	}

	result := map[string]interface{}{
		"target":    target,
		"direction": direction,
		"note":      "Dependency analysis coming soon - requires full import graph analysis",
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// Resource methods
func (ca *CodeAnalyzer) GetPackageResource(path string) (string, error) {
	result, err := ca.AnalyzePackage(map[string]interface{}{
		"package_path":  path,
		"include":       []string{"types", "functions", "imports", "methods"},
		"exported_only": false,
	})
	if err != nil {
		return "", err
	}

	// Extract text from result
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			return textContent.Text, nil
		}
	}
	return "", fmt.Errorf("unexpected result format")
}

func (ca *CodeAnalyzer) GetFileResource(path string) (string, error) {
	result, err := ca.GetFileOverview(map[string]interface{}{
		"file_path": path,
		"level":     "full",
	})
	if err != nil {
		return "", err
	}

	// Extract text from result
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			return textContent.Text, nil
		}
	}
	return "", fmt.Errorf("unexpected result format")
}

// Helper functions

func (ca *CodeAnalyzer) getStringArg(args map[string]interface{}, key, defaultVal string) string {
	if val, ok := args[key].(string); ok {
		return val
	}
	return defaultVal
}

func (ca *CodeAnalyzer) getBoolArg(args map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := args[key].(bool); ok {
		return val
	}
	return defaultVal
}

func (ca *CodeAnalyzer) getIntArg(args map[string]interface{}, key string, defaultVal int) int {
	if val, ok := args[key].(float64); ok {
		return int(val)
	}
	if val, ok := args[key].(int); ok {
		return val
	}
	return defaultVal
}

func (ca *CodeAnalyzer) getSliceArg(args map[string]interface{}, key string, defaultVal []string) []string {
	if val, ok := args[key].([]interface{}); ok {
		result := make([]string, len(val))
		for i, v := range val {
			if s, ok := v.(string); ok {
				result[i] = s
			}
		}
		return result
	}
	return defaultVal
}

func (ca *CodeAnalyzer) exprToString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + ca.exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + ca.exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + ca.exprToString(t.Key) + "]" + ca.exprToString(t.Value)
	case *ast.SelectorExpr:
		return ca.exprToString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return ca.formatFunctionType(t)
	case *ast.ChanType:
		return "chan " + ca.exprToString(t.Value)
	case *ast.Ellipsis:
		return "..." + ca.exprToString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func (ca *CodeAnalyzer) formatFunctionSignature(fn *ast.FuncDecl) string {
	sig := "func "
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		sig += "(" + ca.exprToString(fn.Recv.List[0].Type) + ") "
	}
	sig += fn.Name.Name + "("

	// Parameters
	if fn.Type.Params != nil {
		params := []string{}
		for _, field := range fn.Type.Params.List {
			typeStr := ca.exprToString(field.Type)
			if len(field.Names) == 0 {
				params = append(params, typeStr)
			} else {
				for _, name := range field.Names {
					params = append(params, name.Name+" "+typeStr)
				}
			}
		}
		sig += strings.Join(params, ", ")
	}
	sig += ")"

	// Returns
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		returns := []string{}
		for _, field := range fn.Type.Results.List {
			returns = append(returns, ca.exprToString(field.Type))
		}
		if len(returns) == 1 {
			sig += " " + returns[0]
		} else {
			sig += " (" + strings.Join(returns, ", ") + ")"
		}
	}

	return sig
}

func (ca *CodeAnalyzer) formatFunctionType(ft *ast.FuncType) string {
	sig := "func("

	// Parameters
	if ft.Params != nil {
		params := []string{}
		for _, field := range ft.Params.List {
			params = append(params, ca.exprToString(field.Type))
		}
		sig += strings.Join(params, ", ")
	}
	sig += ")"

	// Returns
	if ft.Results != nil && len(ft.Results.List) > 0 {
		returns := []string{}
		for _, field := range ft.Results.List {
			returns = append(returns, ca.exprToString(field.Type))
		}
		if len(returns) == 1 {
			sig += " " + returns[0]
		} else {
			sig += " (" + strings.Join(returns, ", ") + ")"
		}
	}

	return sig
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
