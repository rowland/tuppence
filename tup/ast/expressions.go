package ast

import (
	"fmt"
	"strings"
)

// expression = try_expression
//            | binary_expression
//            | unary_expression .

type Expression interface {
	Node
	expressionNode()
}

func (n *TryExpression) expressionNode()    {}
func (n *BinaryExpression) expressionNode() {}
func (n *UnaryExpression) expressionNode()  {}

func (n *LogicalOrExpression) expressionNode()  {}
func (n *LogicalAndExpression) expressionNode() {}
func (n *AddSubExpression) expressionNode()     {}
func (n *MulDivExpression) expressionNode()     {}
func (n *PowExpression) expressionNode()        {}
func (n *TypeComparison) expressionNode()       {}
func (n *RelationalComparison) expressionNode() {}

func (n *Identifier) expressionNode()            {}
func (n *FunctionIdentifier) expressionNode()    {}
func (n *Block) expressionNode()                 {}
func (n *IfExpression) expressionNode()          {}
func (n *ForExpression) expressionNode()         {}
func (n *InlineForExpression) expressionNode()   {}
func (n *ArrayFunctionCall) expressionNode()     {}
func (n *ImportExpression) expressionNode()      {}
func (n *TypeofExpression) expressionNode()      {}
func (n *MetaExpression) expressionNode()        {}
func (n *FunctionCall) expressionNode()          {}
func (n *TypeConstructorCall) expressionNode()   {}
func (n *MemberAccess) expressionNode()          {}
func (n *TupleUpdateExpression) expressionNode() {}
func (n *SafeIndexedAccess) expressionNode()     {}
func (n *IndexedAccess) expressionNode()         {}

// func (n *ImportExpression) expressionNode() {}
// func (n *ReturnExpression) expressionNode()      {}
// func (n *BreakExpression) expressionNode()       {}
// func (n *ContinueExpression) expressionNode()    {}

func (n *FloatLiteral) expressionNode()              {}
func (n *IntegerLiteral) expressionNode()            {}
func (n *BooleanLiteral) expressionNode()            {}
func (n *StringLiteral) expressionNode()             {}
func (n *InterpolatedStringLiteral) expressionNode() {}
func (n *RawStringLiteral) expressionNode()          {}
func (n *MultiLineStringLiteral) expressionNode()    {}
func (n *TupleLiteral) expressionNode()              {}
func (n *ArrayLiteral) expressionNode()              {}
func (n *SymbolLiteral) expressionNode()             {}
func (n *RuneLiteral) expressionNode()               {}
func (n *FixedSizeArrayLiteral) expressionNode()     {}

// logical_or_expression = logical_and_expression { logical_or_op logical_and_expression } .

type LogicalOrExpression struct {
	BaseNode
	Operands []Expression
}

func NewLogicalOrExpression(operands []Expression) *LogicalOrExpression {
	return &LogicalOrExpression{
		BaseNode: BaseNode{Type: NodeLogicalOrExpression},
		Operands: operands,
	}
}

func (p *LogicalOrExpression) String() string {
	var builder strings.Builder
	for i, operand := range p.Operands {
		if i > 0 {
			builder.WriteString(" || ")
		}
		builder.WriteString(operand.String())
	}
	return builder.String()
}

// logical_and_expression = comparison_expression { logical_and_op comparison_expression } .

type LogicalAndExpression struct {
	BaseNode
	Operands []Expression
}

func NewLogicalAndExpression(operands []Expression) *LogicalAndExpression {
	return &LogicalAndExpression{
		BaseNode: BaseNode{Type: NodeLogicalAndExpression},
		Operands: operands,
	}
}

func (p *LogicalAndExpression) String() string {
	var builder strings.Builder
	for i, operand := range p.Operands {
		if i > 0 {
			builder.WriteString(" && ")
		}
		builder.WriteString(operand.String())
	}
	return builder.String()
}

// comparison_expression = type_comparison | relational_comparison .

type ComparisonExpression interface {
	Node
	comparisonExpressionNode()
}

func (n *TypeComparison) comparisonExpressionNode()       {}
func (n *RelationalComparison) comparisonExpressionNode() {}

// add_sub_expression = mul_div_expression { add_sub_op mul_div_expression } .

type AddSubExpression struct {
	BaseNode
	Left     Expression
	Operator AddSubOp
	Right    Expression
}

func NewAddSubExpression(left Expression, operator AddSubOp, right Expression) *AddSubExpression {
	return &AddSubExpression{
		BaseNode: BaseNode{Type: NodeAddSubExpression},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (p *AddSubExpression) String() string {
	return fmt.Sprintf("%s %s %s", p.Left, p.Operator, p.Right)
}

// mul_div_expression = pow_expression { mul_div_op pow_expression } .

type MulDivExpression struct {
	BaseNode
	Left     Expression
	Operator MulDivOp
	Right    Expression
}

func NewMulDivExpression(left Expression, operator MulDivOp, right Expression) *MulDivExpression {
	return &MulDivExpression{
		BaseNode: BaseNode{Type: NodeMulDivExpression},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (p *MulDivExpression) String() string {
	return fmt.Sprintf("%s %s %s", p.Left, p.Operator, p.Right)
}

// pow_expression = unary_expression { "^" unary_expression } .

type PowExpression struct {
	BaseNode
	Operands []Expression
}

func NewPowExpression(operands []Expression) *PowExpression {
	return &PowExpression{
		BaseNode: BaseNode{Type: NodePowExpression},
		Operands: operands,
	}
}

func (p *PowExpression) String() string {
	var builder strings.Builder
	for i, operand := range p.Operands {
		if i > 0 {
			builder.WriteString(" ^ ")
		}
		builder.WriteString(operand.String())
	}
	return builder.String()
}

// type_predicate = type_reference | inline_union .

type TypePredicate interface {
	Node
	typePredicateNode()
}

func (n *TypeReference) typePredicateNode() {}
func (n *InlineUnion) typePredicateNode()   {}

// type_comparison = add_sub_expression is_op type_predicate .

type TypeComparison struct {
	BaseNode
	Left  Expression
	Right TypePredicate
}

func NewTypeComparison(left Expression, right TypePredicate) *TypeComparison {
	return &TypeComparison{
		BaseNode: BaseNode{Type: NodeTypeComparison},
		Left:     left,
		Right:    right,
	}
}

func (t *TypeComparison) String() string {
	return fmt.Sprintf("%s is %s", t.Left, t.Right)
}

// relational_comparison = add_sub_expression { rel_op add_sub_expression } .
type RelationalComparison struct {
	BaseNode
	Left     Expression
	Operator RelOp
	Right    Expression
}

func NewRelationalComparison(left Expression, operator RelOp, right Expression) *RelationalComparison {
	return &RelationalComparison{
		BaseNode: BaseNode{Type: NodeRelationalComparison},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (r *RelationalComparison) String() string {
	return fmt.Sprintf("%s %s %s", r.Left, r.Operator, r.Right)
}

// binary_expression = chained_expression .

type BinaryExpression = ChainedExpression

// unary_expression = prefixed_unary_expression
//                  | primary_expression .

type UnaryExpression struct {
	BaseNode
	Operator   UnaryOp
	Expression Expression
}

func NewUnaryExpression(operator UnaryOp, expression Expression) *UnaryExpression {
	return &UnaryExpression{
		BaseNode:   BaseNode{Type: NodeUnaryExpression},
		Operator:   operator,
		Expression: expression,
	}
}

func (u *UnaryExpression) String() string {
	return u.Operator.String() + u.Expression.String()
}

// chained_expression = logical_or_expression { "|>" function_call } .

type ChainedExpression struct {
	BaseNode
	Initial       Node            // The initial expression
	FunctionCalls []*FunctionCall // The chained expressions (FunctionCall nodes)
}

func NewChainedExpression(initial Node, functionCalls []*FunctionCall) *ChainedExpression {
	return &ChainedExpression{
		BaseNode:      BaseNode{Type: NodeChainedExpression},
		Initial:       initial,
		FunctionCalls: functionCalls,
	}
}

func (c *ChainedExpression) String() string {
	var builder strings.Builder
	builder.WriteString(c.Initial.String())

	for _, functionCall := range c.FunctionCalls {
		builder.WriteString(" |> ")
		builder.WriteString(functionCall.String())
	}

	return builder.String()
}
