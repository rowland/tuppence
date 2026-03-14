package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestStatement(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.Statement
		wantErr bool
	}{
		{
			name:    "simple assignment",
			input:   "x = 1",
			want:    ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
			wantErr: false,
		},
		{
			name:  "positional assignment",
			input: "x, y = (1, 2)",
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS(
					[]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1), ast.NewIdentifier("y", nil, 0, 1)},
					nil,
				),
				ast.Immutable,
				ast.NewTupleLiteral(false, []*ast.TupleMember{
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				})),
			wantErr: false,
		},
		{
			name:  "labeled assignment",
			input: "(x: y) = (y: 1)",
			want: ast.NewAssignment(
				ast.NewLabeledAssignmentLHS(
					[]ast.Rename{ast.NewRenameIdentifier(
						ast.NewIdentifier("x", nil, 0, 1),
						ast.NewIdentifier("y", nil, 0, 1),
					)},
				),
				ast.Immutable,
				ast.NewTupleLiteral(true, []*ast.TupleMember{
					ast.NewTupleMember(ast.NewIdentifier("y", nil, 0, 1), ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				})),
			wantErr: false,
		},
		{
			name:  "assignment preceded by a newline",
			input: "\nx = 1",
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
				ast.Immutable,
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := tok.Tokenize([]byte(test.input), "test.tup")
			if err != nil {
				t.Errorf("Tokenize(%q) error = %v", test.input, err)
			}
			stmt, _, err := Statement(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("Statement(%q): want error, got nil", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("Statement(%q): got error %v, want nil", test.input, err)
			}
			if stmt == nil {
				t.Fatalf("Statement(%q) = nil, want not nil", test.input)
			}
			switch want := test.want.(type) {
			case *ast.Assignment:
				got, ok := stmt.(*ast.Assignment)
				if !ok {
					t.Errorf("Statement(%q) = %T, want %T", test.input, stmt, test.want)
					return
				}
				if got.String() != want.String() {
					t.Errorf("Statement(%q) = %v, want %v", test.input, got.String(), want.String())
				}
			}
		})
	}
}

func TestStatements(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []ast.Statement
		wantErr bool
	}{
		{
			name:  "simple assignment",
			input: "x = 1",
			want: []ast.Statement{
				ast.NewAssignment(
					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
					ast.Immutable,
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				),
			},
		},
		{
			name:  "two assignments",
			input: "x = 1\ny = 2",
			want: []ast.Statement{
				ast.NewAssignment(
					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
					ast.Immutable,
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				),
				ast.NewAssignment(
					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("y", nil, 0, 1)}, nil),
					ast.Immutable,
					ast.NewDecimalLiteral("2", 2, nil, 0, 1),
				),
			},
		},
		{
			name:  "two assignments with newlines",
			input: "x = 1\n\ny = 2",
			want: []ast.Statement{
				ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("y", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
			},
		},
		{
			name:  "two assignments with semicolons",
			input: "x = 1; y = 2",
			want: []ast.Statement{
				ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("y", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
			},
		},
		{
			name:  "two assignments with spaces",
			input: "x = 1 y = 2",
			want: []ast.Statement{
				ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("y", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
			},
		},
		{
			name:  "two assignments with tabs",
			input: "x = 1\ty = 2",
			want: []ast.Statement{
				ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewAssignment(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("y", nil, 0, 1)}, nil), false, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := tok.Tokenize([]byte(test.input), "test.tup")
			if err != nil {
				t.Errorf("Tokenize(%q) error = %v", test.input, err)
			}
			stmts, _, err := Statements(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("Statements(%q): want error, got nil", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("Statements(%q): got error %v, want nil", test.input, err)
			}
			if len(stmts) != len(test.want) {
				t.Fatalf("Statements(%q): got %d statements, want %d", test.input, len(stmts), len(test.want))
			}
			for i, got := range stmts {
				if got.String() != test.want[i].String() {
					t.Errorf("Statements(%q)[%d] = %v, want %v", test.input, i, got.String(), test.want[i].String())
				}
			}
		})
	}
}
