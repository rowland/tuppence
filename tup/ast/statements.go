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
func (s *TryExpression) statementNode()                    {}
func (s *BinaryExpression) statementNode()                 {}
func (s *UnaryExpression) statementNode()                  {}
func (s *LogicalOrExpression) statementNode()              {}
func (s *LogicalAndExpression) statementNode()             {}
func (s *AddSubExpression) statementNode()                 {}
func (s *MulDivExpression) statementNode()                 {}
func (s *PowExpression) statementNode()                    {}
func (s *TypeComparison) statementNode()                   {}
func (s *RelationalComparison) statementNode()             {}
func (s *Identifier) statementNode()                       {}
func (s *FunctionIdentifier) statementNode()               {}
func (s *Block) statementNode()                            {}
func (s *IfExpression) statementNode()                     {}
func (s *ForExpression) statementNode()                    {}
func (s *InlineForExpression) statementNode()              {}
func (s *ArrayFunctionCall) statementNode()                {}
func (s *ImportExpression) statementNode()                 {}
func (s *TypeofExpression) statementNode()                 {}
func (s *MetaExpression) statementNode()                   {}
func (s *FunctionCall) statementNode()                     {}
func (s *TypeConstructorCall) statementNode()              {}
func (s *MemberAccess) statementNode()                     {}
func (s *TupleUpdateExpression) statementNode()            {}
func (s *SafeIndexedAccess) statementNode()                {}
func (s *IndexedAccess) statementNode()                    {}
func (s *FloatLiteral) statementNode()                     {}
func (s *IntegerLiteral) statementNode()                   {}
func (s *BooleanLiteral) statementNode()                   {}
func (s *StringLiteral) statementNode()                    {}
func (s *InterpolatedStringLiteral) statementNode()        {}
func (s *RawStringLiteral) statementNode()                 {}
func (s *MultiLineStringLiteral) statementNode()           {}
func (s *TupleLiteral) statementNode()                     {}
func (s *ArrayLiteral) statementNode()                     {}
func (s *SymbolLiteral) statementNode()                    {}
func (s *RuneLiteral) statementNode()                      {}
func (s *FixedSizeArrayLiteral) statementNode()            {}
