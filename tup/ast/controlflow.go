package ast

import (
	"strings"
)

// Block represents a code block enclosed in braces
type Block struct {
	BaseNode
	Body       *BlockBody       // The block's body
	Parameters *BlockParameters // Optional block parameters (may be nil)
}

// NewBlock creates a new Block node
func NewBlock(body *BlockBody, parameters *BlockParameters) *Block {
	return &Block{
		BaseNode:   BaseNode{NodeType: NodeBlock},
		Body:       body,
		Parameters: parameters,
	}
}

// String returns a textual representation of the block
func (b *Block) String() string {
	var builder strings.Builder
	builder.WriteString("{")

	if b.Parameters != nil {
		builder.WriteString(" ")
		builder.WriteString(b.Parameters.String())
		builder.WriteString(" ")
	}

	if b.Body != nil {
		if b.Parameters != nil {
			builder.WriteString("\n")
		} else {
			builder.WriteString(" ")
		}

		builder.WriteString(b.Body.String())
	}

	builder.WriteString(" }")
	return builder.String()
}

// Children returns the child nodes
func (b *Block) Children() []Node {
	var children []Node
	if b.Parameters != nil {
		children = append(children, b.Parameters)
	}
	if b.Body != nil {
		children = append(children, b.Body)
	}
	return children
}

// BlockParameters represents the parameters of a block (e.g., |x, y|)
type BlockParameters struct {
	BaseNode
	Parameters []Node // The block parameters (identifiers)
}

// NewBlockParameters creates a new BlockParameters node
func NewBlockParameters(parameters []Node) *BlockParameters {
	return &BlockParameters{
		BaseNode:   BaseNode{NodeType: NodeBlockParameters},
		Parameters: parameters,
	}
}

// String returns a textual representation of the block parameters
func (b *BlockParameters) String() string {
	var builder strings.Builder
	builder.WriteString("|")

	for i, param := range b.Parameters {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(param.String())
	}

	builder.WriteString("|")
	return builder.String()
}

// Children returns the child nodes
func (b *BlockParameters) Children() []Node {
	return b.Parameters
}

// BlockBody represents the body of a block
type BlockBody struct {
	BaseNode
	Statements []Node // The statements in the block body
	Expression Node   // The final expression (may be nil)
}

// NewBlockBody creates a new BlockBody node
func NewBlockBody(statements []Node, expression Node) *BlockBody {
	return &BlockBody{
		BaseNode:   BaseNode{NodeType: NodeBlockBody},
		Statements: statements,
		Expression: expression,
	}
}

// String returns a textual representation of the block body
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

// Children returns the child nodes
func (b *BlockBody) Children() []Node {
	children := make([]Node, 0, len(b.Statements)+1)
	children = append(children, b.Statements...)
	if b.Expression != nil {
		children = append(children, b.Expression)
	}
	return children
}

// ForHeader represents the header of a for loop (initializer; condition; step)
type ForHeader struct {
	BaseNode
	Initializer Node // The initializer expression
	Condition   Node // The condition expression (may be nil)
	StepExpr    Node // The step expression (may be nil)
}

// String returns a textual representation of the for header
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

// Children returns the child nodes
func (f *ForHeader) Children() []Node {
	var children []Node

	children = append(children, f.Initializer)
	if f.Condition != nil {
		children = append(children, f.Condition)
	}
	if f.StepExpr != nil {
		children = append(children, f.StepExpr)
	}

	return children
}

// ForInHeader represents the header of a for-in loop
type ForInHeader struct {
	BaseNode
	Initializer Node // The initializer expression (may be nil)
	LoopVar     Node // The loop variable
	Iterable    Node // The iterable expression
	StepExpr    Node // The step expression (may be nil)
}

// String returns a textual representation of the for-in header
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

// Children returns the child nodes
func (f *ForInHeader) Children() []Node {
	var children []Node

	if f.Initializer != nil {
		children = append(children, f.Initializer)
	}

	children = append(children, f.LoopVar, f.Iterable)

	if f.StepExpr != nil {
		children = append(children, f.StepExpr)
	}

	return children
}

// IterableHeader represents the iterable part of a for-in loop
type IterableHeader struct {
	BaseNode
	LoopVar  Node // The loop variable
	Iterable Node // The iterable expression
}

// String returns a textual representation of the iterable header
func (i *IterableHeader) String() string {
	return i.LoopVar.String() + " in " + i.Iterable.String()
}

// Children returns the child nodes
func (i *IterableHeader) Children() []Node {
	return []Node{i.LoopVar, i.Iterable}
}

// ForExpression represents a for loop expression
type ForExpression struct {
	BaseNode
	Header Node   // ForHeader, ForInHeader, or nil
	Block  *Block // The loop body
}

// NewForExpression creates a new ForExpression node
func NewForExpression(header Node, block *Block) *ForExpression {
	return &ForExpression{
		BaseNode: BaseNode{NodeType: NodeForExpression},
		Header:   header,
		Block:    block,
	}
}

// String returns a textual representation of the for expression
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

// Children returns the child nodes
func (f *ForExpression) Children() []Node {
	var children []Node

	if f.Header != nil {
		children = append(children, f.Header)
	}

	children = append(children, f.Block)
	return children
}

// InlineForExpression represents an inline for loop expression
type InlineForExpression struct {
	BaseNode
	Header *ForInHeader // The for-in header
	Block  *Block       // The loop body
}

// NewInlineForExpression creates a new InlineForExpression node
func NewInlineForExpression(header *ForInHeader, block *Block) *InlineForExpression {
	return &InlineForExpression{
		BaseNode: BaseNode{NodeType: NodeInlineForExpression},
		Header:   header,
		Block:    block,
	}
}

// String returns a textual representation of the inline for expression
func (i *InlineForExpression) String() string {
	return "inline for " + i.Header.String() + " " + i.Block.String()
}

// Children returns the child nodes
func (i *InlineForExpression) Children() []Node {
	return []Node{i.Header, i.Block}
}

// IfExpression represents an if expression
type IfExpression struct {
	BaseNode
	Conditions []Node   // The conditions (one for the if, and one for each else if)
	Blocks     []*Block // The blocks (one for the if, one for each else if, and optionally one for else)
	HasElse    bool     // Whether there's a final else block
}

// NewIfExpression creates a new IfExpression node
func NewIfExpression(conditions []Node, blocks []*Block, hasElse bool) *IfExpression {
	return &IfExpression{
		BaseNode:   BaseNode{NodeType: NodeIfExpression},
		Conditions: conditions,
		Blocks:     blocks,
		HasElse:    hasElse,
	}
}

// String returns a textual representation of the if expression
func (i *IfExpression) String() string {
	var builder strings.Builder

	// First condition and block (if)
	builder.WriteString("if ")
	builder.WriteString(i.Conditions[0].String())
	builder.WriteString(" ")
	builder.WriteString(i.Blocks[0].String())

	// Else if conditions and blocks
	for j := 1; j < len(i.Conditions); j++ {
		builder.WriteString(" else if ")
		builder.WriteString(i.Conditions[j].String())
		builder.WriteString(" ")
		builder.WriteString(i.Blocks[j].String())
	}

	// Final else block (if present)
	if i.HasElse {
		builder.WriteString(" else ")
		builder.WriteString(i.Blocks[len(i.Blocks)-1].String())
	}

	return builder.String()
}

// Children returns the child nodes
func (i *IfExpression) Children() []Node {
	children := make([]Node, 0, len(i.Conditions)+len(i.Blocks))
	children = append(children, i.Conditions...)

	for _, block := range i.Blocks {
		children = append(children, block)
	}

	return children
}

// CaseBlock represents a case block in a switch statement
type CaseBlock struct {
	BaseNode
	Condition Node   // The case condition
	Block     *Block // The case body
}

// NewCaseBlock creates a new CaseBlock node
func NewCaseBlock(condition Node, block *Block) *CaseBlock {
	return &CaseBlock{
		BaseNode:  BaseNode{NodeType: NodeCaseBlock},
		Condition: condition,
		Block:     block,
	}
}

// String returns a textual representation of the case block
func (c *CaseBlock) String() string {
	return c.Condition.String() + " " + c.Block.String()
}

// Children returns the child nodes
func (c *CaseBlock) Children() []Node {
	return []Node{c.Condition, c.Block}
}

// ElseBlock represents an else block in a switch statement
type ElseBlock struct {
	BaseNode
	Block *Block // The else body
}

// NewElseBlock creates a new ElseBlock node
func NewElseBlock(block *Block) *ElseBlock {
	return &ElseBlock{
		BaseNode: BaseNode{NodeType: NodeElseBlock},
		Block:    block,
	}
}

// String returns a textual representation of the else block
func (e *ElseBlock) String() string {
	return "else " + e.Block.String()
}

// Children returns the child nodes
func (e *ElseBlock) Children() []Node {
	return []Node{e.Block}
}

// SwitchStatement represents a switch statement
type SwitchStatement struct {
	BaseNode
	Expression Node         // The expression being switched on
	Cases      []*CaseBlock // The case blocks
	ElseBlock  *ElseBlock   // The optional else block (may be nil)
}

// NewSwitchStatement creates a new SwitchStatement node
func NewSwitchStatement(expression Node, cases []*CaseBlock, elseBlock *ElseBlock) *SwitchStatement {
	return &SwitchStatement{
		BaseNode:   BaseNode{NodeType: NodeSwitchStatement},
		Expression: expression,
		Cases:      cases,
		ElseBlock:  elseBlock,
	}
}

// String returns a textual representation of the switch statement
func (s *SwitchStatement) String() string {
	var builder strings.Builder
	builder.WriteString("switch ")
	builder.WriteString(s.Expression.String())
	builder.WriteString(" {")

	for _, caseBlock := range s.Cases {
		builder.WriteString("\n  ")
		builder.WriteString(caseBlock.String())
	}

	if s.ElseBlock != nil {
		builder.WriteString("\n  ")
		builder.WriteString(s.ElseBlock.String())
	}

	builder.WriteString("\n}")
	return builder.String()
}

// Children returns the child nodes
func (s *SwitchStatement) Children() []Node {
	children := make([]Node, 0, len(s.Cases)+2)
	children = append(children, s.Expression)

	for _, caseBlock := range s.Cases {
		children = append(children, caseBlock)
	}

	if s.ElseBlock != nil {
		children = append(children, s.ElseBlock)
	}

	return children
}
