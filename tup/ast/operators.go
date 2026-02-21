package ast

type AddSubOp byte

const (
	OpAdd AddSubOp = iota
	OpCheckedAdd
	OpSub
	OpCheckedSub
	OpBitOr
)

var AddSubOpStrings = map[AddSubOp]string{
	OpAdd:        "+",
	OpCheckedAdd: "?+",
	OpSub:        "-",
	OpCheckedSub: "?-",
	OpBitOr:      "|",
}

func (o AddSubOp) String() string {
	return AddSubOpStrings[o]
}

type MulDivOp byte

const (
	OpMul MulDivOp = iota
	OpCheckedMul
	OpDiv
	OpCheckedDiv
	OpMod
	OpCheckedMod
	OpBitAnd
	OpShiftLeft
	OpShiftRight
)

var MulDivOpStrings = map[MulDivOp]string{
	OpMul:        "*",
	OpCheckedMul: "?*",
	OpDiv:        "/",
	OpCheckedDiv: "?/",
	OpMod:        "%",
	OpCheckedMod: "?%",
	OpBitAnd:     "&",
	OpShiftLeft:  "<<",
	OpShiftRight: ">>",
}

func (o MulDivOp) String() string {
	return MulDivOpStrings[o]
}

type RelOp byte

const (
	OpEq RelOp = iota
	OpNeq
	OpLt
	OpLte
	OpGt
	OpGte
	OpMatch
	OpCompare
)

var RelOpStrings = map[RelOp]string{
	OpEq:      "==",
	OpNeq:     "!=",
	OpLt:      "<",
	OpLte:     "<=",
	OpGt:      ">",
	OpGte:     ">=",
	OpMatch:   "=~",
	OpCompare: "<=>",
}

func (o RelOp) String() string {
	return RelOpStrings[o]
}

type CompoundAssignmentOp byte

const (
	OpPlusEq CompoundAssignmentOp = iota
	OpMinusEq
	OpMulEq
	OpDivEq
	OpShiftLeftEq
	OpShiftRightEq
)

var CompoundAssignmentOpStrings = map[CompoundAssignmentOp]string{
	OpPlusEq:       "+=",
	OpMinusEq:      "-=",
	OpMulEq:        "*=",
	OpDivEq:        "/=",
	OpShiftLeftEq:  "<<=",
	OpShiftRightEq: ">>=",
}

func (o CompoundAssignmentOp) String() string {
	return CompoundAssignmentOpStrings[o]
}

type UnaryOp byte

const (
	OpPosSign UnaryOp = iota
	OpNegSign
	OpLogicalNot
	OpBitNot
)

var UnaryOpStrings = map[UnaryOp]string{
	OpPosSign:    "+",
	OpNegSign:    "-",
	OpLogicalNot: "!",
	OpBitNot:     "~",
}

func (o UnaryOp) String() string {
	return UnaryOpStrings[o]
}
