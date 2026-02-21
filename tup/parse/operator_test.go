package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestAddSubOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.AddSubOp
		wantErr bool
	}{
		{
			name:    "plus",
			input:   "+",
			want:    ast.OpAdd,
			wantErr: false,
		},
		{
			name:    "checked_add",
			input:   "?+",
			want:    ast.OpCheckedAdd,
			wantErr: false,
		},
		{
			name:    "minus",
			input:   "-",
			want:    ast.OpSub,
			wantErr: false,
		},
		{
			name:    "checked_sub",
			input:   "?-",
			want:    ast.OpCheckedSub,
			wantErr: false,
		},
		{
			name:    "bit_or",
			input:   "|",
			want:    ast.OpBitOr,
			wantErr: false,
		},
		{
			name:    "mul not allowed",
			input:   "*",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			op, remainder, err := AddSubOp(tokens)
			if test.wantErr && err == nil {
				t.Fatalf("AddSubOp(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("AddSubOp(%q) = %v, want nil", test.input, err)
			}
			if op != test.want {
				t.Errorf("AddSubOp(%q) = %v, want %v", test.input, op, test.want)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("AddSubOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestMulDivOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.MulDivOp
		wantErr bool
	}{
		{
			name:    "mul",
			input:   "*",
			want:    ast.OpMul,
			wantErr: false,
		},
		{
			name:    "checked_mul",
			input:   "?*",
			want:    ast.OpCheckedMul,
			wantErr: false,
		},
		{
			name:    "div",
			input:   "/",
			want:    ast.OpDiv,
			wantErr: false,
		},
		{
			name:    "checked_div",
			input:   "?/",
			want:    ast.OpCheckedDiv,
			wantErr: false,
		},
		{
			name:    "mod",
			input:   "%",
			want:    ast.OpMod,
			wantErr: false,
		},
		{
			name:    "checked_mod",
			input:   "?%",
			want:    ast.OpCheckedMod,
			wantErr: false,
		},
		{
			name:    "bit_and",
			input:   "&",
			want:    ast.OpBitAnd,
			wantErr: false,
		},
		{
			name:    "shift_left",
			input:   "<<",
			want:    ast.OpShiftLeft,
			wantErr: false,
		},
		{
			name:    "shift_right",
			input:   ">>",
			want:    ast.OpShiftRight,
			wantErr: false,
		},
		{
			name:    "plus not allowed",
			input:   "+",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			op, remainder, err := MulDivOp(tokens)
			if test.wantErr && err == nil {
				t.Fatalf("MulDivOp(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("MulDivOp(%q) = %v, want nil", test.input, err)
			}
			if op != test.want {
				t.Errorf("MulDivOp(%q) = %v, want %v", test.input, op, test.want)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("MulDivOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestRelOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.RelOp
		wantErr bool
	}{
		{
			name:    "eq",
			input:   "==",
			want:    ast.OpEq,
			wantErr: false,
		},
		{
			name:    "neq",
			input:   "!=",
			want:    ast.OpNeq,
			wantErr: false,
		},
		{
			name:    "lt",
			input:   "<",
			want:    ast.OpLt,
			wantErr: false,
		},
		{
			name:    "lte",
			input:   "<=",
			want:    ast.OpLte,
			wantErr: false,
		},
		{
			name:    "gt",
			input:   ">",
			want:    ast.OpGt,
			wantErr: false,
		},
		{
			name:    "gte",
			input:   ">=",
			want:    ast.OpGte,
			wantErr: false,
		},
		{
			name:    "match",
			input:   "=~",
			want:    ast.OpMatch,
			wantErr: false,
		},
		{
			name:    "compare",
			input:   "<=>",
			want:    ast.OpCompare,
			wantErr: false,
		},
		{
			name:    "plus not allowed",
			input:   "+",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			op, remainder, err := RelOp(tokens)
			if test.wantErr && err == nil {
				t.Fatalf("RelOp(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("RelOp(%q) = %v, want nil", test.input, err)
			}
			if op != test.want {
				t.Errorf("RelOp(%q) = %v, want %v", test.input, op, test.want)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("RelOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestCompoundAssignmentOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.CompoundAssignmentOp
		wantErr bool
	}{
		{
			name:    "plus_eq",
			input:   "+=",
			want:    ast.OpPlusEq,
			wantErr: false,
		},
		{
			name:    "minus_eq",
			input:   "-=",
			want:    ast.OpMinusEq,
			wantErr: false,
		},
		{
			name:    "mul_eq",
			input:   "*=",
			want:    ast.OpMulEq,
			wantErr: false,
		},
		{
			name:    "div_eq",
			input:   "/=",
			want:    ast.OpDivEq,
			wantErr: false,
		},
		{
			name:    "shift_left_eq",
			input:   "<<=",
			want:    ast.OpShiftLeftEq,
			wantErr: false,
		},
		{
			name:    "shift_right_eq",
			input:   ">>=",
			want:    ast.OpShiftRightEq,
			wantErr: false,
		},
		{
			name:    "plus not allowed",
			input:   "+",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			op, remainder, err := CompoundAssignmentOp(tokens)
			if test.wantErr && err == nil {
				t.Fatalf("CompoundAssignmentOp(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("CompoundAssignmentOp(%q) = %v, want nil", test.input, err)
			}
			if op != test.want {
				t.Errorf("CompoundAssignmentOp(%q) = %v, want %v", test.input, op, test.want)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("CompoundAssignmentOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestUnaryOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.UnaryOp
		wantErr bool
	}{
		{
			name:    "positive",
			input:   "+",
			want:    ast.OpPosSign,
			wantErr: false,
		},
		{
			name:    "negative",
			input:   "-",
			want:    ast.OpNegSign,
			wantErr: false,
		},
		{
			name:    "logical not",
			input:   "!",
			want:    ast.OpLogicalNot,
			wantErr: false,
		},
		{
			name:    "bit not",
			input:   "~",
			want:    ast.OpBitNot,
			wantErr: false,
		},
		{
			name:    "mul not allowed",
			input:   "*",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			op, remainder, err := UnaryOp(tokens)
			if test.wantErr && err == nil {
				t.Fatalf("UnaryOp(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("UnaryOp(%q) = %v, want nil", test.input, err)
			}
			if op != test.want {
				t.Errorf("UnaryOp(%q) = %v, want %v", test.input, op, test.want)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("UnaryOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestLogicalOrOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "logical or",
			input:   "||",
			wantErr: false,
		},
		{
			name:    "plus not allowed",
			input:   "+",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			remainder, err := LogicalOrOp(tokens)
			if test.wantErr && err == nil {
				t.Fatalf("LogicalOrOp(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("LogicalOrOp(%q) = %v, want nil", test.input, err)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("LogicalOrOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestLogicalAndOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "logical and",
			input:   "&&",
			wantErr: false,
		},
		{
			name:    "plus not allowed",
			input:   "+",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			remainder, err := LogicalAndOp(tokens)
			if test.wantErr && err == nil {
				t.Fatalf("LogicalAndOp(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("LogicalAndOp(%q) = %v, want nil", test.input, err)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("LogicalAndOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestIsOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "is",
			input:   "is",
			wantErr: false,
		},
		{
			name:    "plus not allowed",
			input:   "+",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			remainder, err := IsOp(tokens)
			if test.wantErr && err == nil {
				t.Fatalf("IsOp(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("IsOp(%q) = %v, want nil", test.input, err)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("IsOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}
