package ast

import (
	"encoding/json"
	"fmt"
	"strings"
)

func FormatFunctions(functions []FunctionSignature) string {
	if len(functions) == 0 {
		return "No functions found"
	}

	var sb strings.Builder
	sb.WriteString("Functions:\n\n")

	for _, fn := range functions {
		// Format function signature
		signature := fn.Name
		if fn.Receiver != "" {
			signature = fmt.Sprintf("(%s) %s", fn.Receiver, fn.Name)
		}

		signature += "("
		params := make([]string, len(fn.Parameters))
		for i, param := range fn.Parameters {
			if param.Name != "" {
				params[i] = fmt.Sprintf("%s %s", param.Name, param.Type)
			} else {
				params[i] = param.Type
			}
		}
		signature += strings.Join(params, ", ")
		signature += ")"

		if len(fn.Returns) > 0 {
			if len(fn.Returns) == 1 {
				signature += " " + fn.Returns[0]
			} else {
				signature += " (" + strings.Join(fn.Returns, ", ") + ")"
			}
		}

		sb.WriteString(fmt.Sprintf("Line %d: %s\n", fn.LineNumber, signature))

		if fn.DocComment != "" {
			sb.WriteString(fmt.Sprintf("  Doc: %s\n", strings.TrimSpace(fn.DocComment)))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func FormatTypes(types []TypeDefinition) string {
	if len(types) == 0 {
		return "No types found"
	}

	var sb strings.Builder
	sb.WriteString("Types:\n\n")

	for _, t := range types {
		sb.WriteString(fmt.Sprintf("Line %d: %s %s\n", t.LineNumber, t.Kind, t.Name))

		if t.DocComment != "" {
			sb.WriteString(fmt.Sprintf("  Doc: %s\n", strings.TrimSpace(t.DocComment)))
		}

		if len(t.Fields) > 0 {
			sb.WriteString("  Fields:\n")
			for _, field := range t.Fields {
				sb.WriteString(fmt.Sprintf("    %s %s", field.Name, field.Type))
				if field.Tag != "" {
					sb.WriteString(fmt.Sprintf(" %s", field.Tag))
				}
				sb.WriteString("\n")
			}
		}

		if len(t.Methods) > 0 {
			sb.WriteString("  Methods:\n")
			for _, method := range t.Methods {
				sb.WriteString(fmt.Sprintf("    %s\n", method))
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

func FormatPackageAnalysis(analysis *PackageAnalysis) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Package: %s\n\n", analysis.PackageName))

	if len(analysis.Imports) > 0 {
		sb.WriteString("Imports:\n")
		for _, imp := range analysis.Imports {
			sb.WriteString(fmt.Sprintf("  - %s\n", imp))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(FormatTypes(analysis.Types))
	sb.WriteString(FormatFunctions(analysis.Functions))

	return sb.String()
}

func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
