package ast

type FunctionSignature struct {
	Name       string      `json:"name"`
	Receiver   string      `json:"receiver,omitempty"`
	Parameters []Parameter `json:"parameters"`
	Returns    []string    `json:"returns"`
	DocComment string      `json:"doc_comment"`
	LineNumber int         `json:"line_number"`
}

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TypeDefinition struct {
	Name       string      `json:"name"`
	Kind       string      `json:"kind"` // struct, interface, alias
	Fields     []FieldInfo `json:"fields,omitempty"`
	Methods    []string    `json:"methods,omitempty"`
	DocComment string      `json:"doc_comment"`
	LineNumber int         `json:"line_number"`
}

type FieldInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tag  string `json:"tag,omitempty"`
}

type PackageAnalysis struct {
	PackageName string             `json:"package_name"`
	Functions   []FunctionSignature `json:"functions"`
	Types       []TypeDefinition   `json:"types"`
	Imports     []string           `json:"imports"`
}
