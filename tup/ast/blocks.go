package ast

import "strings"

// Block represents a code block enclosed in braces
type Block struct {
	BaseNode
	Body *BlockBody
}

// NewBlock creates a new Block node
func NewBlock(body *BlockBody) *Block {
	return &Block{
		BaseNode: BaseNode{Type: NodeBlock},
		Body:     body,
	}
}

// String returns a textual representation of the block
func (b *Block) String() string {
	if b.Body == nil {
		return "{}"
	}

	body := b.Body.String()
	if body == "" {
		return "{\n}"
	}

	var builder strings.Builder
	builder.WriteString("{\n")
	builder.WriteString(indentString(body))
	builder.WriteString("\n}")
	return builder.String()
}

// BlockParameters represents the parameters of a block (e.g., |x, y|)
type BlockParameters struct {
	BaseNode
	Parameters AssignmentLHS // The block parameters (identifiers)
}

// NewBlockParameters creates a new BlockParameters node
func NewBlockParameters(parameters AssignmentLHS) *BlockParameters {
	return &BlockParameters{
		BaseNode:   BaseNode{Type: NodeBlockParameters},
		Parameters: parameters,
	}
}

// String returns a textual representation of the block parameters
func (b *BlockParameters) String() string {
	var builder strings.Builder
	builder.WriteString("|")

	builder.WriteString(b.Parameters.String())

	builder.WriteString("|")
	return builder.String()
}

// BlockBody represents the body of a block
type BlockBody struct {
	BaseNode
	Statements []Statement // The statements in the block body
	Expression Expression  // The final expression (may be nil)
}

func NewBlockBody(statements []Statement, expression Expression) *BlockBody {
	return &BlockBody{
		BaseNode:   BaseNode{Type: NodeBlockBody},
		Statements: statements,
		Expression: expression,
	}
}

func (b *BlockBody) String() string {
	var builder strings.Builder

	for _, stmt := range b.Statements {
		builder.WriteString(stmt.String())
		builder.WriteString("\n")
	}

	if b.Expression != nil {
		builder.WriteString(b.Expression.String())
	}

	return builder.String()
}
