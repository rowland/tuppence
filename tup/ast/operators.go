package ast

// Operator kinds for classification
type OperatorKind uint8

const (
	OpKindAdditive       OperatorKind = iota // +, -, |
	OpKindMultiplicative                     // *, /, %, &, <<, >>
	OpKindRelational                         // ==, <, >, <=, >=, !=
	OpKindType                               // is
	OpKindLogical                            // &&, ||
	OpKindChecked                            // ?+, ?-, ?*, ?/, ?%
	OpKindCompound                           // +=, -=, *=, /=, <<=, >>=
)

// AddSubOp represents additive operators (+, -, |)
type AddSubOp struct {
	BaseNode
	Symbol string
	Kind   OperatorKind
}

// NewAddSubOp creates a new AddSubOp node
func NewAddSubOp(symbol string) *AddSubOp {
	return &AddSubOp{
		BaseNode: BaseNode{Type: NodeAddSubOp},
		Symbol:   symbol,
		Kind:     OpKindAdditive,
	}
}

// String returns a textual representation of the operator
func (o *AddSubOp) String() string {
	return o.Symbol
}

// MulDivOp represents multiplicative operators (*, /, %, &, <<, >>)
type MulDivOp struct {
	BaseNode
	Symbol string
	Kind   OperatorKind
}

// NewMulDivOp creates a new MulDivOp node
func NewMulDivOp(symbol string) *MulDivOp {
	return &MulDivOp{
		BaseNode: BaseNode{Type: NodeMulDivOp},
		Symbol:   symbol,
		Kind:     OpKindMultiplicative,
	}
}

// String returns a textual representation of the operator
func (o *MulDivOp) String() string {
	return o.Symbol
}

// RelOp represents relational operators (==, <, >, <=, >=, !=)
type RelOp struct {
	BaseNode
	Symbol string
	Kind   OperatorKind
}

// NewRelOp creates a new RelOp node
func NewRelOp(symbol string) *RelOp {
	return &RelOp{
		BaseNode: BaseNode{Type: NodeRelOp},
		Symbol:   symbol,
		Kind:     OpKindRelational,
	}
}

// String returns a textual representation of the operator
func (o *RelOp) String() string {
	return o.Symbol
}

// IsOp represents the 'is' operator for type checking
type IsOp struct {
	BaseNode
	Kind OperatorKind
}

// NewIsOp creates a new IsOp node
func NewIsOp() *IsOp {
	return &IsOp{
		BaseNode: BaseNode{Type: NodeIsOp},
		Kind:     OpKindType,
	}
}

// String returns a textual representation of the operator
func (o *IsOp) String() string {
	return "is"
}

// LogicalOp represents logical operators (&&, ||)
type LogicalOp struct {
	BaseNode
	Symbol       string
	Kind         OperatorKind
	ShortCircuit bool // Whether the operator short-circuits evaluation
}

// NewLogicalOp creates a new LogicalOp node
func NewLogicalOp(symbol string, shortCircuit bool) *LogicalOp {
	return &LogicalOp{
		BaseNode:     BaseNode{Type: NodeLogicalOp},
		Symbol:       symbol,
		Kind:         OpKindLogical,
		ShortCircuit: shortCircuit,
	}
}

// String returns a textual representation of the operator
func (o *LogicalOp) String() string {
	return o.Symbol
}

// ShortCircuitOp represents operators that short-circuit evaluation
type ShortCircuitOp struct {
	BaseNode
	Symbol string
	Kind   OperatorKind
}

// NewShortCircuitOp creates a new ShortCircuitOp node
func NewShortCircuitOp(symbol string) *ShortCircuitOp {
	return &ShortCircuitOp{
		BaseNode: BaseNode{Type: NodeShortCircuitOp},
		Symbol:   symbol,
		Kind:     OpKindLogical,
	}
}

// String returns a textual representation of the operator
func (o *ShortCircuitOp) String() string {
	return o.Symbol
}

// CheckedArithmeticOp represents checked arithmetic operators (?+, ?-, ?*, ?/, ?%)
type CheckedArithmeticOp struct {
	BaseNode
	Symbol string
	Kind   OperatorKind
}

// NewCheckedArithmeticOp creates a new CheckedArithmeticOp node
func NewCheckedArithmeticOp(symbol string) *CheckedArithmeticOp {
	return &CheckedArithmeticOp{
		BaseNode: BaseNode{Type: NodeCheckedArithmeticOp},
		Symbol:   symbol,
		Kind:     OpKindChecked,
	}
}

// String returns a textual representation of the operator
func (o *CheckedArithmeticOp) String() string {
	return o.Symbol
}

// CompoundAssignmentOp represents compound assignment operators (+=, -=, *=, /=, <<=, >>=)
type CompoundAssignmentOp struct {
	BaseNode
	Symbol string
	Kind   OperatorKind
}

// NewCompoundAssignmentOp creates a new CompoundAssignmentOp node
func NewCompoundAssignmentOp(symbol string) *CompoundAssignmentOp {
	return &CompoundAssignmentOp{
		BaseNode: BaseNode{Type: NodeCompoundAssignmentOp},
		Symbol:   symbol,
		Kind:     OpKindCompound,
	}
}

// String returns a textual representation of the operator
func (o *CompoundAssignmentOp) String() string {
	return o.Symbol
}

// CompoundAssignment represents a compound assignment (e.g., x += y)
type CompoundAssignment struct {
	BaseNode
	Left     Node                  // The left operand (target of assignment)
	Operator *CompoundAssignmentOp // The compound assignment operator
	Right    Node                  // The right operand
}

// NewCompoundAssignment creates a new CompoundAssignment node
func NewCompoundAssignment(left Node, operator *CompoundAssignmentOp, right Node) *CompoundAssignment {
	return &CompoundAssignment{
		BaseNode: BaseNode{Type: NodeCompoundAssignmentOp},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// String returns a textual representation of the compound assignment
func (c *CompoundAssignment) String() string {
	return c.Left.String() + " " + c.Operator.String() + " " + c.Right.String()
}

// Standard operators as constants for easy reference
var (
	// Additive
	OpAdd   = NewAddSubOp("+")
	OpSub   = NewAddSubOp("-")
	OpBitOr = NewAddSubOp("|")

	// Multiplicative
	OpMul    = NewMulDivOp("*")
	OpDiv    = NewMulDivOp("/")
	OpMod    = NewMulDivOp("%")
	OpBitAnd = NewMulDivOp("&")
	OpShiftL = NewMulDivOp("<<")
	OpShiftR = NewMulDivOp(">>")

	// Relational
	OpEq  = NewRelOp("==")
	OpLt  = NewRelOp("<")
	OpGt  = NewRelOp(">")
	OpLte = NewRelOp("<=")
	OpGte = NewRelOp(">=")
	OpNe  = NewRelOp("!=")

	// Type
	OpIs = NewIsOp()

	// Logical
	OpAnd = NewLogicalOp("&&", true)
	OpOr  = NewLogicalOp("||", true)

	// Short-circuit
	OpShortAnd = NewShortCircuitOp("&&")
	OpShortOr  = NewShortCircuitOp("||")

	// Checked arithmetic
	OpCheckedAdd = NewCheckedArithmeticOp("?+")
	OpCheckedSub = NewCheckedArithmeticOp("?-")
	OpCheckedMul = NewCheckedArithmeticOp("?*")
	OpCheckedDiv = NewCheckedArithmeticOp("?/")
	OpCheckedMod = NewCheckedArithmeticOp("?%")

	// Compound assignment
	OpPlusEq   = NewCompoundAssignmentOp("+=")
	OpMinusEq  = NewCompoundAssignmentOp("-=")
	OpMulEq    = NewCompoundAssignmentOp("*=")
	OpDivEq    = NewCompoundAssignmentOp("/=")
	OpShiftLEq = NewCompoundAssignmentOp("<<=")
	OpShiftREq = NewCompoundAssignmentOp(">>=")
)
