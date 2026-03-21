package ast

import "strings"

type IfExpression struct {
	BaseNode
	Conditions []Node   // The conditions (one for the if, and one for each else if)
	Blocks     []*Block // The blocks (one for the if, one for each else if, and optionally one for else)
	HasElse    bool     // Whether there's a final else block
}

func NewIfExpression(conditions []Node, blocks []*Block, hasElse bool) *IfExpression {
	return &IfExpression{
		BaseNode:   BaseNode{Type: NodeIfExpression},
		Conditions: conditions,
		Blocks:     blocks,
		HasElse:    hasElse,
	}
}

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

type CaseBlock struct {
	BaseNode
	Condition Node   // The case condition
	Block     *Block // The case body
}

func NewCaseBlock(condition Node, block *Block) *CaseBlock {
	return &CaseBlock{
		BaseNode:  BaseNode{Type: NodeCaseBlock},
		Condition: condition,
		Block:     block,
	}
}

func (c *CaseBlock) String() string {
	return c.Condition.String() + " " + c.Block.String()
}

type ElseBlock struct {
	BaseNode
	Block *Block // The else body
}

func NewElseBlock(block *Block) *ElseBlock {
	return &ElseBlock{
		BaseNode: BaseNode{Type: NodeElseBlock},
		Block:    block,
	}
}

func (e *ElseBlock) String() string {
	return "else " + e.Block.String()
}

type SwitchStatement struct {
	BaseNode
	Expression Node         // The expression being switched on
	Cases      []*CaseBlock // The case blocks
	ElseBlock  *ElseBlock   // The optional else block (may be nil)
}

func NewSwitchStatement(expression Node, cases []*CaseBlock, elseBlock *ElseBlock) *SwitchStatement {
	return &SwitchStatement{
		BaseNode:   BaseNode{Type: NodeSwitchStatement},
		Expression: expression,
		Cases:      cases,
		ElseBlock:  elseBlock,
	}
}

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
