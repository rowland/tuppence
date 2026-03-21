package ast

// Base type for all literals
type Literal interface {
	Expression
	literalNode()
}

func (n *FloatLiteral) literalNode()   {}
func (n *IntegerLiteral) literalNode() {}

func (n *BooleanLiteral) literalNode()            {}
func (n *StringLiteral) literalNode()             {}
func (n *InterpolatedStringLiteral) literalNode() {}
func (n *RawStringLiteral) literalNode()          {}
func (n *MultiLineStringLiteral) literalNode()    {}
func (n *TupleLiteral) literalNode()              {}
func (n *ArrayLiteral) literalNode()              {}
func (n *SymbolLiteral) literalNode()             {}
func (n *RuneLiteral) literalNode()               {}
func (n *FixedSizeArrayLiteral) literalNode()     {}

// number = float_literal | integer_literal .

// Base type for number literals
type Number interface {
	Literal
	numberNode()
}

func (n *FloatLiteral) numberNode()   {}
func (n *IntegerLiteral) numberNode() {}
