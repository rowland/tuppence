package ast

import "strings"

// for_block = "{" { statement } [ expression ] "}" .

type ForBlock struct {
	BaseNode
	Statements []Statement
	Expression Expression
}

func NewForBlock(statements []Statement, expression Expression) *ForBlock {
	return &ForBlock{
		BaseNode:   BaseNode{Type: NodeForBlock},
		Statements: statements,
		Expression: expression,
	}
}

func (f *ForBlock) String() string {
	var builder strings.Builder
	builder.WriteString("{")

	if len(f.Statements) == 0 && f.Expression == nil {
		builder.WriteString("}")
		return builder.String()
	}

	builder.WriteString("\n")
	for _, statement := range f.Statements {
		builder.WriteString(indentString(statement.String()))
		builder.WriteString("\n")
	}
	if f.Expression != nil {
		builder.WriteString(indentString(f.Expression.String()))
		builder.WriteString("\n")
	}
	builder.WriteString("}")
	return builder.String()
}

// initializer = assignment .

type Initializer struct {
	BaseNode
	Assignment *Assignment
}

func NewInitializer(assignment *Assignment) *Initializer {
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
	Expression Expression
}

func NewStepExpression(expression Expression) *StepExpression {
	return &StepExpression{
		BaseNode:   BaseNode{Type: NodeStepExpression},
		Expression: expression,
	}
}

func (s *StepExpression) String() string {
	return s.Expression.String()
}

// iterable = expression .

type Iterable struct {
	BaseNode
	Expression Expression
}

func NewIterable(expression Expression) *Iterable {
	return &Iterable{
		BaseNode:   BaseNode{Type: NodeIterable},
		Expression: expression,
	}
}

func (i *Iterable) String() string {
	return i.Expression.String()
}

type ForExpressionHeader interface {
	Node
	forExpressionHeaderNode()
}

func (f *ForHeader) forExpressionHeaderNode()   {}
func (f *ForInHeader) forExpressionHeaderNode() {}

// for_header = initializer [ ";" condition [ ";" step_expression ] ] .

type ForHeader struct {
	BaseNode
	Initializer *Initializer
	Condition   Expression
	StepExpr    *StepExpression
}

func NewForHeader(initializer *Initializer, condition Expression, stepExpr *StepExpression) *ForHeader {
	return &ForHeader{
		BaseNode:    BaseNode{Type: NodeForHeader},
		Initializer: initializer,
		Condition:   condition,
		StepExpr:    stepExpr,
	}
}

func (f *ForHeader) String() string {
	var builder strings.Builder
	builder.WriteString(f.Initializer.String())

	if f.Condition != nil {
		builder.WriteString("; ")
		builder.WriteString(f.Condition.String())

		if f.StepExpr != nil {
			builder.WriteString("; ")
			builder.WriteString(f.StepExpr.String())
		}
	}

	return builder.String()
}

// for_in_header = ( initializer ";" assignment_lhs "in" iterable [ ";" step_expression ] )
//               | ( assignment_lhs "in" iterable ) .

type ForInHeader struct {
	BaseNode
	Initializer *Initializer
	LoopVar     AssignmentLHS
	Iterable    *Iterable
	StepExpr    *StepExpression
}

func NewForInHeader(initializer *Initializer, loopVar AssignmentLHS, iterable *Iterable, stepExpr *StepExpression) *ForInHeader {
	return &ForInHeader{
		BaseNode:    BaseNode{Type: NodeForInHeader},
		Initializer: initializer,
		LoopVar:     loopVar,
		Iterable:    iterable,
		StepExpr:    stepExpr,
	}
}

func (f *ForInHeader) String() string {
	var builder strings.Builder

	if f.Initializer != nil {
		builder.WriteString(f.Initializer.String())
		builder.WriteString("; ")
	}

	builder.WriteString(f.LoopVar.String())
	builder.WriteString(" in ")
	builder.WriteString(f.Iterable.String())

	if f.StepExpr != nil {
		builder.WriteString("; ")
		builder.WriteString(f.StepExpr.String())
	}

	return builder.String()
}

// iterable_header = assignment_lhs "in" iterable .

type IterableHeader struct {
	BaseNode
	LoopVar  AssignmentLHS
	Iterable *Iterable
}

func NewIterableHeader(loopVar AssignmentLHS, iterable *Iterable) *IterableHeader {
	return &IterableHeader{
		BaseNode: BaseNode{Type: NodeIterableHeader},
		LoopVar:  loopVar,
		Iterable: iterable,
	}
}

func (i *IterableHeader) String() string {
	return i.LoopVar.String() + " in " + i.Iterable.String()
}

// for_expression = "for" [ for_header | for_in_header ] for_block .

type ForExpression struct {
	BaseNode
	Header ForExpressionHeader
	Block  *ForBlock
}

func NewForExpression(header ForExpressionHeader, block *ForBlock) *ForExpression {
	return &ForExpression{
		BaseNode: BaseNode{Type: NodeForExpression},
		Header:   header,
		Block:    block,
	}
}

func (f *ForExpression) String() string {
	var builder strings.Builder
	builder.WriteString("for")

	if f.Header != nil {
		builder.WriteString(" ")
		builder.WriteString(f.Header.String())
	}

	builder.WriteString(" ")
	builder.WriteString(f.Block.String())
	return builder.String()
}

// inline_for_expression = "inline" "for" for_in_header for_block .

type InlineForExpression struct {
	BaseNode
	Header *ForInHeader
	Block  *ForBlock
}

func NewInlineForExpression(header *ForInHeader, block *ForBlock) *InlineForExpression {
	return &InlineForExpression{
		BaseNode: BaseNode{Type: NodeInlineForExpression},
		Header:   header,
		Block:    block,
	}
}

func (i *InlineForExpression) String() string {
	return "inline for " + i.Header.String() + " " + i.Block.String()
}
