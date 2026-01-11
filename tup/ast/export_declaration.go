package ast

// export_declaration = ( export_type_qualified_function_declaration
// 	| export_type_qualified_declaration
// 	| export_type_declaration
// 	| export_function_declaration
// 	| export_assignment ) .

type ExportDeclaration interface {
	// Node
	TopLevelItem
	exportDeclarationNode()
}

func (n *ExportAssignment) exportDeclarationNode()                       {}
func (n *ExportFunctionDeclaration) exportDeclarationNode()              {}
func (n *ExportTypeDeclaration) exportDeclarationNode()                  {}
func (n *ExportTypeQualifiedDeclaration) exportDeclarationNode()         {}
func (n *ExportTypeQualifiedFunctionDeclaration) exportDeclarationNode() {}
