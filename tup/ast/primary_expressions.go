package ast

// primary_expression = "(" expression ")"
//	                  | block
//	                  | if_expression
//	                  | for_expression
//	                  | inline_for_expression
//	                  | array_function_call
//	                  | import_expression
//	                  | typeof_expression
//	                  | meta_expression
//	                  | function_call
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
func (n *ForExpression) primaryExpressionNode()         {}
func (n *InlineForExpression) primaryExpressionNode()   {}
func (n *ArrayFunctionCall) primaryExpressionNode()     {}
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
func (n *Identifier) primaryExpressionNode()            {}

// func (n *ImportExpression) primaryExpressionNode()      {}

// initializer = assignment .

type Initializer struct {
	BaseNode
	Assignment Assignment
}

func NewInitializer(assignment Assignment) *Initializer {
	return &Initializer{
		BaseNode:   BaseNode{Type: NodeInitializer},
		Assignment: assignment,
	}
}

func (i *Initializer) String() string {
	return i.Assignment.String()
}

// step_expression = expression .

type StepExpression struct {
	BaseNode
	Expression Node // The step expression
}

func NewStepExpression(expression Node) *StepExpression {
	return &StepExpression{
		BaseNode:   BaseNode{Type: NodeStepExpression},
		Expression: expression,
	}
}

func (s *StepExpression) String() string {
	return s.Expression.String()
}

// iterable = expression .

// Iterable represents an iterable expression used in for-in loops
type Iterable struct {
	BaseNode
	Expression Node // The expression that produces an iterable value
}

func NewIterable(expression Node) *Iterable {
	return &Iterable{
		BaseNode:   BaseNode{Type: NodeIterable},
		Expression: expression,
	}
}

func (i *Iterable) String() string {
	return i.Expression.String()
}
