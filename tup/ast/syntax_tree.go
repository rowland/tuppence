package ast

// SyntaxTree represents the root of the entire syntax tree
type SyntaxTree struct {
	Modules []*Module // All modules in the program
}

// NewSyntaxTree creates a new SyntaxTree
func NewSyntaxTree() *SyntaxTree {
	return &SyntaxTree{
		Modules: []*Module{},
	}
}

// AddModule adds a module to the syntax tree
func (t *SyntaxTree) AddModule(module *Module) {
	t.Modules = append(t.Modules, module)
}

// String returns a textual representation of the syntax tree
func (t *SyntaxTree) String() string {
	result := "SyntaxTree{\n"

	result += "  Modules: [\n"
	for _, module := range t.Modules {
		result += "    " + module.Name + "\n"
	}
	result += "  ]\n"

	result += "}"
	return result
}
