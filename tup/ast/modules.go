package ast

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
