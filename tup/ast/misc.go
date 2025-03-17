package ast

import (
	"fmt"
)

// Module represents a complete module (file) in the AST
type Module struct {
	BaseNode
	Path         string    // File path of the module
	Namespace    string    // Module namespace
	Declarations []Node    // Top-level declarations in the module
	Comments     []Node    // Comments in the module
	Imports      []Node    // Import statements
	Directives   []Node    // Module directives
	Metadata     *Metadata // Module metadata
}

// NewModule creates a new Module node
func NewModule(path string, namespace string) *Module {
	return &Module{
		BaseNode:     BaseNode{NodeType: NodeModule},
		Path:         path,
		Namespace:    namespace,
		Declarations: []Node{},
		Comments:     []Node{},
		Imports:      []Node{},
		Directives:   []Node{},
	}
}

// AddDeclaration adds a declaration to the module
func (m *Module) AddDeclaration(decl Node) {
	m.Declarations = append(m.Declarations, decl)
}

// AddImport adds an import to the module
func (m *Module) AddImport(imp Node) {
	m.Imports = append(m.Imports, imp)
}

// AddComment adds a comment to the module
func (m *Module) AddComment(comment Node) {
	m.Comments = append(m.Comments, comment)
}

// AddDirective adds a directive to the module
func (m *Module) AddDirective(directive Node) {
	m.Directives = append(m.Directives, directive)
}

// SetMetadata sets the module metadata
func (m *Module) SetMetadata(metadata *Metadata) {
	m.Metadata = metadata
}

// String returns a textual representation of the module
func (m *Module) String() string {
	result := "// Module: " + m.Path + "\n"

	if m.Namespace != "" {
		result += "namespace " + m.Namespace + ";\n\n"
	}

	// Output imports
	for _, imp := range m.Imports {
		result += imp.String() + "\n"
	}

	if len(m.Imports) > 0 {
		result += "\n"
	}

	// Output declarations
	for _, decl := range m.Declarations {
		result += decl.String() + "\n\n"
	}

	return result
}

// Children returns the child nodes
func (m *Module) Children() []Node {
	var children []Node

	children = append(children, m.Declarations...)
	children = append(children, m.Imports...)
	children = append(children, m.Comments...)
	children = append(children, m.Directives...)

	if m.Metadata != nil {
		children = append(children, m.Metadata)
	}

	return children
}

// Comment is the base type for all comment nodes
type Comment struct {
	BaseNode
	Text string // Comment text
}

// NewComment creates a new Comment node
func NewComment(text string) *Comment {
	return &Comment{
		BaseNode: BaseNode{NodeType: NodeComment},
		Text:     text,
	}
}

// String returns a textual representation of the comment
func (c *Comment) String() string {
	return "// " + c.Text
}

// DocComment represents a documentation comment
type DocComment struct {
	BaseNode
	Text string // Comment text
}

// NewDocComment creates a new DocComment node
func NewDocComment(text string) *DocComment {
	return &DocComment{
		BaseNode: BaseNode{NodeType: NodeDocComment},
		Text:     text,
	}
}

// String returns a textual representation of the documentation comment
func (c *DocComment) String() string {
	return "/// " + c.Text
}

// LineComment represents a single-line comment
type LineComment struct {
	BaseNode
	Text string // Comment text
}

// NewLineComment creates a new LineComment node
func NewLineComment(text string) *LineComment {
	return &LineComment{
		BaseNode: BaseNode{NodeType: NodeLineComment},
		Text:     text,
	}
}

// String returns a textual representation of the line comment
func (c *LineComment) String() string {
	return "// " + c.Text
}

// BlockComment represents a multi-line block comment
type BlockComment struct {
	BaseNode
	Text string // Comment text
}

// NewBlockComment creates a new BlockComment node
func NewBlockComment(text string) *BlockComment {
	return &BlockComment{
		BaseNode: BaseNode{NodeType: NodeBlockComment},
		Text:     text,
	}
}

// String returns a textual representation of the block comment
func (c *BlockComment) String() string {
	return "/* " + c.Text + " */"
}

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
