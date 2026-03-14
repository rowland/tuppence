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
			name:  "plus",
			input: "+",
			want:  ast.OpAdd,
		},
		{
			name:  "checked_add",
			input: "?+",
			want:  ast.OpCheckedAdd,
		},
		{
			name:  "minus",
			input: "-",
			want:  ast.OpSub,
		},
		{
			name:  "checked_sub",
			input: "?-",
			want:  ast.OpCheckedSub,
		},
		{
			name:  "bit_or",
			input: "|",
			want:  ast.OpBitOr,
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
			op, remainder, match := AddSubOp(tokens)
			if test.wantErr && match {
				t.Fatalf("AddSubOp(%q): match == true, want false", test.input)
			}
			if !test.wantErr && !match {
				t.Fatalf("AddSubOp(%q): match == false, want true", test.input)
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
			name:  "mul",
			input: "*",
			want:  ast.OpMul,
		},
		{
			name:  "checked_mul",
			input: "?*",
			want:  ast.OpCheckedMul,
		},
		{
			name:  "div",
			input: "/",
			want:  ast.OpDiv,
		},
		{
			name:  "checked_div",
			input: "?/",
			want:  ast.OpCheckedDiv,
		},
		{
			name:  "mod",
			input: "%",
			want:  ast.OpMod,
		},
		{
			name:  "checked_mod",
			input: "?%",
			want:  ast.OpCheckedMod,
		},
		{
			name:  "bit_and",
			input: "&",
			want:  ast.OpBitAnd,
		},
		{
			name:  "shift_left",
			input: "<<",
			want:  ast.OpShiftLeft,
		},
		{
			name:  "shift_right",
			input: ">>",
			want:  ast.OpShiftRight,
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
			op, remainder, match := MulDivOp(tokens)
			if test.wantErr && match {
				t.Fatalf("MulDivOp(%q): match == true, want false", test.input)
			}
			if !test.wantErr && !match {
				t.Fatalf("MulDivOp(%q): match == false, want true", test.input)
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
			name:  "eq",
			input: "==",
			want:  ast.OpEq,
		},
		{
			name:  "neq",
			input: "!=",
			want:  ast.OpNeq,
		},
		{
			name:  "lt",
			input: "<",
			want:  ast.OpLt,
		},
		{
			name:  "lte",
			input: "<=",
			want:  ast.OpLte,
		},
		{
			name:  "gt",
			input: ">",
			want:  ast.OpGt,
		},
		{
			name:  "gte",
			input: ">=",
			want:  ast.OpGte,
		},
		{
			name:  "match",
			input: "=~",
			want:  ast.OpMatch,
		},
		{
			name:  "compare",
			input: "<=>",
			want:  ast.OpCompare,
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
			op, remainder, match := RelOp(tokens)
			if test.wantErr && match {
				t.Fatalf("RelOp(%q): match == true, want false", test.input)
			}
			if !test.wantErr && !match {
				t.Fatalf("RelOp(%q): match == false, want true", test.input)
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
			name:  "plus_eq",
			input: "+=",
			want:  ast.OpPlusEq,
		},
		{
			name:  "minus_eq",
			input: "-=",
			want:  ast.OpMinusEq,
		},
		{
			name:  "mul_eq",
			input: "*=",
			want:  ast.OpMulEq,
		},
		{
			name:  "div_eq",
			input: "/=",
			want:  ast.OpDivEq,
		},
		{
			name:  "pow_eq",
			input: "^=",
			want:  ast.OpPowEq,
		},
		{
			name:  "shift_left_eq",
			input: "<<=",
			want:  ast.OpShiftLeftEq,
		},
		{
			name:  "shift_right_eq",
			input: ">>=",
			want:  ast.OpShiftRightEq,
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
			op, remainder, match := CompoundAssignmentOp(tokens)
			if test.wantErr && match {
				t.Fatalf("CompoundAssignmentOp(%q): match == true, want false", test.input)
			}
			if !test.wantErr && !match {
				t.Fatalf("CompoundAssignmentOp(%q): match == false, want true", test.input)
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
			name:  "positive",
			input: "+",
			want:  ast.OpPosSign,
		},
		{
			name:  "negative",
			input: "-",
			want:  ast.OpNegSign,
		},
		{
			name:  "logical not",
			input: "!",
			want:  ast.OpLogicalNot,
		},
		{
			name:  "bit not",
			input: "~",
			want:  ast.OpBitNot,
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
			op, remainder, match := UnaryOp(tokens)
			if test.wantErr && match {
				t.Fatalf("UnaryOp(%q): match == true, want false", test.input)
			}
			if !test.wantErr && !match {
				t.Fatalf("UnaryOp(%q): match == false, want true", test.input)
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
			name:  "logical or",
			input: "||",
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
			remainder, found := LogicalOrOp(tokens)
			if test.wantErr && found {
				t.Fatalf("LogicalOrOp(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("LogicalOrOp(%q) = %v, want nil", test.input, found)
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
			name:  "logical and",
			input: "&&",
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
			remainder, found := LogicalAndOp(tokens)
			if test.wantErr && found {
				t.Fatalf("LogicalAndOp(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("LogicalAndOp(%q) = %v, want nil", test.input, found)
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
			name:  "is",
			input: "is",
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
			remainder, found := IsOp(tokens)
			if test.wantErr && found {
				t.Fatalf("IsOp(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("IsOp(%q) = %v, want nil", test.input, found)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("IsOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestPipeOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "pipe operator",
			input: "|>",
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
			remainder, found := PipeOp(tokens)
			if test.wantErr && found {
				t.Fatalf("PipeOp(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("PipeOp(%q) = %v, want nil", test.input, found)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("PipeOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestPowOp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "pow operator",
			input: "^",
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
			remainder, found := PowOp(tokens)
			if test.wantErr && found {
				t.Fatalf("PowOp(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("PowOp(%q) = %v, want nil", test.input, found)
			}
			if !test.wantErr && (len(remainder) != 1 || remainder[0].Type != tok.TokEOF) {
				t.Errorf("PowOp(%q) remainder = %v, want 1 token (EOF)", test.input, remainder)
			}
		})
	}
}

func TestPartialApplication(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{input: ",*"},                 // sees partial application
		{input: ",x", wantErr: true},  // sees partial application followed by identifier
		{input: "x,*", wantErr: true}, // sees identifier followed by partial application
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			_, found := PartialApplication(tokens)
			if test.wantErr && found {
				t.Errorf("PartialApplication(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("PartialApplication(%q) = %v, want nil", test.input, found)
			}
		})
	}
}
