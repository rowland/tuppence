package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestCompoundAssignment(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.CompoundAssignment
		wantErr bool
	}{
		{
			name:  "add equals",
			input: "x += 1",
			want: ast.NewCompoundAssignment(
				ast.NewIdentifier("x", nil, 0, 1),
				ast.OpPlusEq,
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
			),
		},
		{
			name:  "shift right equals with expression",
			input: "bits >>= count + 1",
			want: ast.NewCompoundAssignment(
				ast.NewIdentifier("bits", nil, 0, 4),
				ast.OpShiftRightEq,
				ast.NewAddSubExpression(
					ast.NewIdentifier("count", nil, 0, 5),
					ast.OpAdd,
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				),
			),
		},
		{
			name:    "plain assignment is not a compound assignment",
			input:   "x = 1",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := tok.Tokenize([]byte(test.input), "test.tup")
			if err != nil {
				t.Fatalf("Tokenize(%q) error = %v", test.input, err)
			}

			got, remainder, err := CompoundAssignment(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("CompoundAssignment(%q): want error, got nil", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("CompoundAssignment(%q): got error %v, want nil", test.input, err)
			}
			if got == nil {
				t.Fatalf("CompoundAssignment(%q) = nil, want not nil", test.input)
			}
			if got.String() != test.want.String() {
				t.Fatalf("CompoundAssignment(%q) = %q, want %q", test.input, got.String(), test.want.String())
			}
			if len(remainder) != 1 {
				t.Fatalf("CompoundAssignment(%q) remainder = %v, want EOF only", test.input, remainder)
			}
		})
	}
}
