package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// add_sub_op = add_op | checked_add_op | sub_op | checked_sub_op | bit_or_op .

func AddSubOp(tokens []tok.Token) (op ast.AddSubOp, remainder []tok.Token, match bool) {
	remainder = skipComments(tokens)

	switch peek(remainder).Type {
	case tok.TokOpPlus:
		op = ast.OpAdd
	case tok.TokOpCheckedAdd:
		op = ast.OpCheckedAdd
	case tok.TokOpMinus:
		op = ast.OpSub
	case tok.TokOpCheckedSub:
		op = ast.OpCheckedSub
	case tok.TokOpBitOr:
		op = ast.OpBitOr
	default:
		return 0, tokens, false
	}

	return op, remainder[1:], true
}

// mul_div_op = mul_op | checked_mul_op | div_op | checked_div_op | mod_op | checked_mod_op | bit_and_op | shift_left_op | shift_right_op .

func MulDivOp(tokens []tok.Token) (op ast.MulDivOp, remainder []tok.Token, match bool) {
	remainder = skipComments(tokens)

	switch peek(remainder).Type {
	case tok.TokOpMul:
		op = ast.OpMul
	case tok.TokOpCheckedMul:
		op = ast.OpCheckedMul
	case tok.TokOpDiv:
		op = ast.OpDiv
	case tok.TokOpCheckedDiv:
		op = ast.OpCheckedDiv
	case tok.TokOpMod:
		op = ast.OpMod
	case tok.TokOpCheckedMod:
		op = ast.OpCheckedMod
	case tok.TokOpBitAnd:
		op = ast.OpBitAnd
	case tok.TokOpSHL:
		op = ast.OpShiftLeft
	case tok.TokOpSHR:
		op = ast.OpShiftRight
	default:
		return 0, tokens, false
	}

	return op, remainder[1:], true
}

// rel_op = eq_op | neq_op | lt_op | lte_op | gt_op | gte_op | match_op | compare_op.

func RelOp(tokens []tok.Token) (op ast.RelOp, remainder []tok.Token, match bool) {
	remainder = skipComments(tokens)

	switch peek(remainder).Type {
	case tok.TokOpEQ:
		op = ast.OpEq
	case tok.TokOpNE:
		op = ast.OpNeq
	case tok.TokOpLT:
		op = ast.OpLt
	case tok.TokOpLE:
		op = ast.OpLte
	case tok.TokOpGT:
		op = ast.OpGt
	case tok.TokOpGE:
		op = ast.OpGte
	case tok.TokOpMatch:
		op = ast.OpMatch
	case tok.TokOpCompare:
		op = ast.OpCompare
	default:
		return 0, tokens, false
	}

	return op, remainder[1:], true
}

// compound_assignment_op = plus_eq_op | minus_eq_op | mul_eq_op | div_eq_op | shift_left_eq_op | shift_right_eq_op .

func CompoundAssignmentOp(tokens []tok.Token) (op ast.CompoundAssignmentOp, remainder []tok.Token, match bool) {
	remainder = skipComments(tokens)

	switch peek(remainder).Type {
	case tok.TokOpPlusEQ:
		op = ast.OpPlusEq
	case tok.TokOpMinusEQ:
		op = ast.OpMinusEq
	case tok.TokOpMulEQ:
		op = ast.OpMulEq
	case tok.TokOpDivEQ:
		op = ast.OpDivEq
	case tok.TokOpPowEQ:
		op = ast.OpPowEq
	case tok.TokOpSHL_EQ:
		op = ast.OpShiftLeftEq
	case tok.TokOpSHR_EQ:
		op = ast.OpShiftRightEq
	default:
		return 0, tokens, false
	}

	return op, remainder[1:], true
}

// unary_op = add_op | sub_op | logical_not_op | bit_not_op .

func UnaryOp(tokens []tok.Token) (op ast.UnaryOp, remainder []tok.Token, match bool) {
	remainder = skipComments(tokens)

	switch peek(remainder).Type {
	case tok.TokOpPlus:
		op = ast.OpPosSign
	case tok.TokOpMinus:
		op = ast.OpNegSign
	case tok.TokOpNot:
		op = ast.OpLogicalNot
	case tok.TokOpBitNot:
		op = ast.OpBitNot
	default:
		return 0, tokens, false
	}

	return op, remainder[1:], true
}

// logical_or_op = "||" .

var LogicalOrOp = expectFunc(tok.TokOpLogOr)

// logical_and_op = "&&" .

var LogicalAndOp = expectFunc(tok.TokOpLogAnd)

// is_op = "is" .

var IsOp = expectFunc(tok.TokKwIs)

// pipe_op = "|>" .

var PipeOp = expectFunc(tok.TokOpPipe)

// partial_application = ","  "*" .

func PartialApplication(tokens []tok.Token) (remainder []tok.Token, found bool) {
	remainder = skipComments(tokens)

	if remainder, found = Comma(remainder); !found {
		return remainder, false
	}

	return Star(remainder)
}

// pow_op = "^" .

var PowOp = expectFunc(tok.TokOpPow)

// spread_op = "..." .

var SpreadOp = expectFunc(tok.TokOpRest)
