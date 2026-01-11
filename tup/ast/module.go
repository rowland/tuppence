package ast

import "github.com/rowland/tuppence/tup/source"

// module = { top_level_item } .

// Module represents a complete module (composed of one or more files) in the AST
type Module struct {
	Name          string           // Module name (derived from the file name)
	Sources       []*source.Source // Source files that make up the module
	TopLevelItems []TopLevelItem   // Top-level items in the module
}

// NewModule creates a new Module node
func NewModule(name string) *Module {
	return &Module{
		Name:          name,
		Sources:       []*source.Source{},
		TopLevelItems: []TopLevelItem{},
	}
}

func (m *Module) AddSource(source *source.Source) {
	m.Sources = append(m.Sources, source)
}

func (m *Module) AddTopLevelItem(item TopLevelItem) {
	m.TopLevelItems = append(m.TopLevelItems, item)
}

func (m *Module) Type() NodeType {
	return NodeModule
}

// String returns a textual representation of the module
func (m *Module) String() string {
	result := "# module " + m.Name + "\n"

	return result
}
