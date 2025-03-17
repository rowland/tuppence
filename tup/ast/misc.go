package ast

import (
	"fmt"
)

// ErrorNode represents a syntax error in the AST
type ErrorNode struct {
	BaseNode
	Message string // Error message
	Code    int    // Error code
}

// NewErrorNode creates a new ErrorNode
func NewErrorNode(message string, code int) *ErrorNode {
	return &ErrorNode{
		BaseNode: BaseNode{NodeType: NodeErrorNode},
		Message:  message,
		Code:     code,
	}
}

// String returns a textual representation of the error node
func (e *ErrorNode) String() string {
	return fmt.Sprintf("ERROR(%d): %s", e.Code, e.Message)
}

// Metadata represents metadata about a module or node
type Metadata struct {
	BaseNode
	Properties map[string]interface{} // Metadata properties
}

// NewMetadata creates a new Metadata node
func NewMetadata() *Metadata {
	return &Metadata{
		BaseNode:   BaseNode{NodeType: NodeMetadata},
		Properties: make(map[string]interface{}),
	}
}

// Set sets a metadata property
func (m *Metadata) Set(key string, value interface{}) {
	m.Properties[key] = value
}

// Get gets a metadata property
func (m *Metadata) Get(key string) interface{} {
	return m.Properties[key]
}

// String returns a textual representation of the metadata
func (m *Metadata) String() string {
	result := "Metadata{"
	first := true
	for k, v := range m.Properties {
		if !first {
			result += ", "
		}
		first = false
		result += k + ": " + StringifyValue(v)
	}
	result += "}"
	return result
}

// StringifyValue converts a value to a string representation
func StringifyValue(v interface{}) string {
	// This is a simple implementation; in a real implementation,
	// you would want to handle different types differently
	return fmt.Sprintf("%v", v)
}

// Directive represents a compiler directive
type Directive struct {
	BaseNode
	Name  string // Directive name
	Value string // Directive value
}

// NewDirective creates a new Directive node
func NewDirective(name, value string) *Directive {
	return &Directive{
		BaseNode: BaseNode{NodeType: NodeDirective},
		Name:     name,
		Value:    value,
	}
}

// String returns a textual representation of the directive
func (d *Directive) String() string {
	return "#" + d.Name + "(" + d.Value + ")"
}

// SyntaxTree represents the root of the entire syntax tree
type SyntaxTree struct {
	BaseNode
	Modules  []*Module    // All modules in the program
	Errors   []*ErrorNode // All errors in the program
	Metadata *Metadata    // Metadata about the syntax tree
}

// NewSyntaxTree creates a new SyntaxTree
func NewSyntaxTree() *SyntaxTree {
	return &SyntaxTree{
		BaseNode: BaseNode{NodeType: NodeSyntaxTree},
		Modules:  []*Module{},
		Errors:   []*ErrorNode{},
	}
}

// AddModule adds a module to the syntax tree
func (t *SyntaxTree) AddModule(module *Module) {
	t.Modules = append(t.Modules, module)
}

// AddError adds an error to the syntax tree
func (t *SyntaxTree) AddError(err *ErrorNode) {
	t.Errors = append(t.Errors, err)
}

// SetMetadata sets the syntax tree metadata
func (t *SyntaxTree) SetMetadata(metadata *Metadata) {
	t.Metadata = metadata
}

// String returns a textual representation of the syntax tree
func (t *SyntaxTree) String() string {
	result := "SyntaxTree{\n"

	if len(t.Errors) > 0 {
		result += "  Errors: [\n"
		for _, err := range t.Errors {
			result += "    " + err.String() + "\n"
		}
		result += "  ]\n"
	}

	result += "  Modules: [\n"
	for _, module := range t.Modules {
		result += "    " + module.Path + "\n"
	}
	result += "  ]\n"

	result += "}"
	return result
}

// Children returns the child nodes
func (t *SyntaxTree) Children() []Node {
	var children []Node

	for _, module := range t.Modules {
		children = append(children, module)
	}

	for _, err := range t.Errors {
		children = append(children, err)
	}

	if t.Metadata != nil {
		children = append(children, t.Metadata)
	}

	return children
}

// EmptyStatement represents an empty statement (semicolon with no expression)
type EmptyStatement struct {
	BaseNode
}

// NewEmptyStatement creates a new EmptyStatement node
func NewEmptyStatement() *EmptyStatement {
	return &EmptyStatement{
		BaseNode: BaseNode{NodeType: NodeEmptyStatement},
	}
}

// String returns a textual representation of the empty statement
func (s *EmptyStatement) String() string {
	return ";"
}
