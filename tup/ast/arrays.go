package ast

// size = decimal_literal | identifier .

type Size interface {
	Node
	sizeNode()
}

func (n *IntegerLiteral) sizeNode() {}
func (n *Identifier) sizeNode()     {}
