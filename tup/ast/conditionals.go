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

type SwitchCase struct {
	BaseNode
	Condition MatchCondition
	Body      *FunctionBlock
}

func NewSwitchCase(condition MatchCondition, body *FunctionBlock) *SwitchCase {
	return &SwitchCase{
		BaseNode:  BaseNode{Type: NodeSwitchCase},
		Condition: condition,
		Body:      body,
	}
}

func (c *SwitchCase) String() string {
	return c.Condition.String() + " " + c.Body.String()
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

type SwitchExpression struct {
	BaseNode
	Expression Expression
	Cases      []*SwitchCase
	ElseBlock  *FunctionBlock
}

func NewSwitchExpression(expression Expression, cases []*SwitchCase, elseBlock *FunctionBlock) *SwitchExpression {
	return &SwitchExpression{
		BaseNode:   BaseNode{Type: NodeSwitchExpression},
		Expression: expression,
		Cases:      cases,
		ElseBlock:  elseBlock,
	}
}

func (s *SwitchExpression) String() string {
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
		builder.WriteString("else ")
		builder.WriteString(s.ElseBlock.String())
	}

	builder.WriteString("\n}")
	return builder.String()
}
