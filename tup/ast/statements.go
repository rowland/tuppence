package ast

// statement = ( type_qualified_function_declaration
// 	           | type_qualified_declaration
// 	           | type_declaration
// 	           | function_declaration
// 	           | compound_assignment
// 	           | assignment
// 	           | expression
// 	           ) .

type Statement interface {
	Node
	statementNode()
}

func (s *TypeQualifiedFunctionDeclaration) statementNode() {}
func (s *TypeQualifiedDeclaration) statementNode()         {}
func (s *TypeDeclaration) statementNode()                  {}
func (s *FunctionDeclaration) statementNode()              {}
func (s *CompoundAssignment) statementNode()               {}
func (s *Assignment) statementNode()                       {}

// func (s *Expression) statementNode()                       {}
