package ast

// LogicalOrOp represents logical operators (&&, ||)
type LogicalOrOp struct {
	BaseNode
	Symbol string
}

// NewLogicalOrOp creates a new LogicalOrOp node
func NewLogicalOrOp(symbol string) *LogicalOrOp {
	return &LogicalOrOp{
		BaseNode: BaseNode{Type: NodeLogicalOrOp},
		Symbol:   symbol,
	}
}

// String returns a textual representation of the operator
func (o *LogicalOrOp) String() string {
	return "||"
}

// LogicalAndOp represents logical operators (&&)
type LogicalAndOp struct {
	BaseNode
	Symbol string
}

// NewLogicalAndOp creates a new LogicalAndOp node
func NewLogicalAndOp(symbol string) *LogicalAndOp {
	return &LogicalAndOp{
		BaseNode: BaseNode{Type: NodeLogicalAndOp},
		Symbol:   symbol,
	}
}

// String returns a textual representation of the operator
func (o *LogicalAndOp) String() string {
	return "&&"
}
