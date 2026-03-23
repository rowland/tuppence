package ast

// primary_expression = "(" expression ")"
//	                  | block
//	                  | if_expression
//	                  | switch_expression
//	                  | for_expression
//	                  | inline_for_expression
//	                  | array_function_call
//	                  | import_expression
//	                  | typeof_expression
//	                  | meta_expression
//	                  | function_identifier
//	                  | type_constructor_call
//	                  | return_expression
//	                  | break_expression
//	                  | continue_expression
//	                  | member_access
//	                  | tuple_update_expression
//	                  | safe_indexed_access
//	                  | indexed_access
//	                  | range
//	                  | identifier
//	                  | literal .

type PrimaryExpression interface {
	Expression
	primaryExpressionNode()
}

func (n *Block) primaryExpressionNode()                 {}
func (n *IfExpression) primaryExpressionNode()          {}
func (n *SwitchExpression) primaryExpressionNode()      {}
func (n *ForExpression) primaryExpressionNode()         {}
func (n *InlineForExpression) primaryExpressionNode()   {}
func (n *ArrayFunctionCall) primaryExpressionNode()     {}
func (n *ImportExpression) primaryExpressionNode()      {}
func (n *TypeofExpression) primaryExpressionNode()      {}
func (n *MetaExpression) primaryExpressionNode()        {}
func (n *FunctionCall) primaryExpressionNode()          {}
func (n *TypeConstructorCall) primaryExpressionNode()   {}
func (n *ReturnExpression) primaryExpressionNode()      {}
func (n *BreakExpression) primaryExpressionNode()       {}
func (n *ContinueExpression) primaryExpressionNode()    {}
func (n *MemberAccess) primaryExpressionNode()          {}
func (n *TupleUpdateExpression) primaryExpressionNode() {}
func (n *SafeIndexedAccess) primaryExpressionNode()     {}
func (n *IndexedAccess) primaryExpressionNode()         {}
func (n *Range) primaryExpressionNode()                 {}
func (n *FunctionIdentifier) primaryExpressionNode()    {}
func (n *Identifier) primaryExpressionNode()            {}
