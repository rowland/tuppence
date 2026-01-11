package ast

// top_level_item = ( type_qualified_function_declaration
//                  | type_qualified_declaration
//                  | type_declaration
//                  | function_declaration
//                  | assignment
//                  | export_declaration
//                  ) .

type TopLevelItem interface {
	Node
	topLevelItemNode()
}

func (n *Assignment) topLevelItemNode()                       {}
func (n *FunctionDeclaration) topLevelItemNode()              {}
func (n *TypeDeclaration) topLevelItemNode()                  {}
func (n *TypeQualifiedDeclaration) topLevelItemNode()         {}
func (n *TypeQualifiedFunctionDeclaration) topLevelItemNode() {}

// ExportDeclaration
func (n *ExportTypeQualifiedFunctionDeclaration) topLevelItemNode() {}
func (n *ExportTypeQualifiedDeclaration) topLevelItemNode()         {}
func (n *ExportTypeDeclaration) topLevelItemNode()                  {}
func (n *ExportFunctionDeclaration) topLevelItemNode()              {}
func (n *ExportAssignment) topLevelItemNode()                       {}
