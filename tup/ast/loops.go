package ast

import "strings"

// ForHeader represents the header of a for loop (initializer; condition; step)
type ForHeader struct {
	BaseNode
	Initializer Node // The initializer expression
	Condition   Node // The condition expression (may be nil)
	StepExpr    Node // The step expression (may be nil)
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

type ForInHeader struct {
	BaseNode
	Initializer Node // The initializer expression (may be nil)
	LoopVar     Node // The loop variable
	Iterable    Node // The iterable expression
	StepExpr    Node // The step expression (may be nil)
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

// IterableHeader represents the iterable part of a for-in loop
type IterableHeader struct {
	BaseNode
	LoopVar  Node // The loop variable
	Iterable Node // The iterable expression
}

func (i *IterableHeader) String() string {
	return i.LoopVar.String() + " in " + i.Iterable.String()
}

type ForExpression struct {
	BaseNode
	Header Node   // ForHeader, ForInHeader, or nil
	Block  *Block // The loop body
}

func NewForExpression(header Node, block *Block) *ForExpression {
	return &ForExpression{
		BaseNode: BaseNode{Type: NodeForExpression},
		Header:   header,
		Block:    block,
	}
}

func (f *ForExpression) String() string {
	var builder strings.Builder
	builder.WriteString("for ")

	if f.Header != nil {
		builder.WriteString(f.Header.String())
		builder.WriteString(" ")
	}

	builder.WriteString(f.Block.String())
	return builder.String()
}

type InlineForExpression struct {
	BaseNode
	Header *ForInHeader // The for-in header
	Block  *Block       // The loop body
}

func NewInlineForExpression(header *ForInHeader, block *Block) *InlineForExpression {
	return &InlineForExpression{
		BaseNode: BaseNode{Type: NodeInlineForExpression},
		Header:   header,
		Block:    block,
	}
}

func (i *InlineForExpression) String() string {
	return "inline for " + i.Header.String() + " " + i.Block.String()
}
