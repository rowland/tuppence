package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestAssignment(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenTypes []tok.TokenType
		want       *ast.Assignment
		wantErr    bool
	}{
		{
			name:       "simple assignment",
			input:      "x = 1",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokOpAssign, tok.TokDecLit, tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
				ast.Immutable,
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
			),
		},
		{
			name:       "simple mut assignment",
			input:      "x = mut 1",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokOpAssign, tok.TokKwMut, tok.TokDecLit, tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
				ast.Mutable,
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
			),
		},
		{
			name:       "tuple assignment",
			input:      "x = (1, 2)",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokOpAssign, tok.TokOpenParen, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
				ast.Immutable,
				ast.NewTupleLiteral(false, []*ast.TupleMember{
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				}),
			),
		},
		{
			name:       "ordinal assignment",
			input:      "x, y = (1, 2)",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokID, tok.TokOpAssign, tok.TokOpenParen, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
					ast.NewIdentifier("x", nil, 0, 1),
					ast.NewIdentifier("y", nil, 0, 1),
				}, nil),
				ast.Immutable,
				ast.NewTupleLiteral(false, []*ast.TupleMember{
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				},
				),
			),
		},
		{
			name:       "ordinal mut assignment",
			input:      "x, y = mut (1, 2)",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokID, tok.TokOpAssign, tok.TokKwMut, tok.TokOpenParen, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
					ast.NewIdentifier("x", nil, 0, 1),
					ast.NewIdentifier("y", nil, 0, 1),
				}, nil),
				ast.Mutable,
				ast.NewTupleLiteral(false, []*ast.TupleMember{
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				},
				),
			),
		},
		{
			name:       "ordinal assignment with rest binding",
			input:      "x, ...rest = (1, 2, 3)",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokID, tok.TokOpAssign, tok.TokOpenParen, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, ast.NewRestOperator(ast.NewIdentifier("...", nil, 0, 3))),
				ast.Immutable,
				ast.NewTupleLiteral(false, []*ast.TupleMember{
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("3", 3, nil, 0, 1)),
				}),
			),
		},
		{
			name:       "ordinal assignment with (ignored) rest binding",
			input:      "x, y, ... = (1, 2, 3)",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokOpAssign, tok.TokOpenParen, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
					ast.NewIdentifier("x", nil, 0, 1),
					ast.NewIdentifier("y", nil, 0, 1),
				}, ast.NewRestOperator(ast.NewIdentifier("...", nil, 0, 3))),
				ast.Immutable,
				ast.NewTupleLiteral(false, []*ast.TupleMember{
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
					ast.NewTupleMember(nil, ast.NewDecimalLiteral("3", 3, nil, 0, 1)),
				},
				),
			),
		},
		{
			name:       "ordinal assignment with multiple (ignored) rest bindings (invalid)",
			input:      "x, ..., ... = (1, 2, 3)",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokComma, tok.TokOpRest, tok.TokOpAssign, tok.TokOpenParen, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "labeled assignment",
			input:      "(x, y) = (x: 1, y: 2)",
			tokenTypes: []tok.TokenType{tok.TokOpenParen, tok.TokID, tok.TokComma, tok.TokID, tok.TokCloseParen, tok.TokOpAssign, tok.TokOpenParen, tok.TokID, tok.TokColon, tok.TokDecLit, tok.TokComma, tok.TokID, tok.TokColon, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewLabeledAssignmentLHS(
					[]ast.Rename{
						ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), nil),
						ast.NewRenameIdentifier(ast.NewIdentifier("y", nil, 0, 1), nil),
					},
				),
				ast.Immutable,
				ast.NewTupleLiteral(true, []*ast.TupleMember{
					ast.NewTupleMember(ast.NewIdentifier("x", nil, 0, 1), ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewTupleMember(ast.NewIdentifier("y", nil, 0, 1), ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				},
				),
			),
		},
		{
			name:  "labeled mut assignment",
			input: "(x, y) = mut (x: 1, y: 2)",
			tokenTypes: []tok.TokenType{
				tok.TokOpenParen, tok.TokID, tok.TokComma, tok.TokID, tok.TokCloseParen,
				tok.TokOpAssign, tok.TokKwMut,
				tok.TokOpenParen, tok.TokID, tok.TokColon, tok.TokDecLit, tok.TokComma, tok.TokID, tok.TokColon, tok.TokDecLit, tok.TokCloseParen,
				tok.TokEOF},
			want: ast.NewAssignment(
				ast.NewLabeledAssignmentLHS(
					[]ast.Rename{
						ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), nil),
						ast.NewRenameIdentifier(ast.NewIdentifier("y", nil, 0, 1), nil),
					},
				),
				ast.Mutable,
				ast.NewTupleLiteral(true, []*ast.TupleMember{
					ast.NewTupleMember(ast.NewIdentifier("x", nil, 0, 1), ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewTupleMember(ast.NewIdentifier("y", nil, 0, 1), ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				},
				),
			),
		},
		{
			name:  "labeled assignment with rest operator (invalid)",
			input: "(x, y, ...) = (x: 1, y: 2, z: 3)",
			tokenTypes: []tok.TokenType{
				tok.TokOpenParen, tok.TokID, tok.TokComma, tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokCloseParen,
				tok.TokOpAssign,
				tok.TokOpenParen, tok.TokID, tok.TokColon, tok.TokDecLit, tok.TokComma, tok.TokID, tok.TokColon, tok.TokDecLit, tok.TokComma, tok.TokID, tok.TokColon, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want:    nil,
			wantErr: true,
		},
		{
			name:       "assignment to it is invalid",
			input:      "it = 1",
			tokenTypes: []tok.TokenType{tok.TokKwIt, tok.TokOpAssign, tok.TokDecLit, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunParseTestExt(t, tt.name, tt.input, tt.want, tt.wantErr, "Assignment", Assignment, StringerCheck[*ast.Assignment], tt.tokenTypes)
		})
	}
}

// ordinal_assignment_lhs = identifier { "," identifier } [ "," rest_operator ] .

func TestOrdinalAssignmentLHS(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenTypes []tok.TokenType
		want       *ast.OrdinalAssignmentLHS
		wantErr    bool
	}{
		{
			name:       "empty",
			input:      "",
			tokenTypes: []tok.TokenType{tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "single",
			input:      "x",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokEOF},
			want:       ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
		},
		{
			name:       "two identifiers",
			input:      "x, y",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokID, tok.TokEOF},
			want:       ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1), ast.NewIdentifier("y", nil, 0, 1)}, nil),
		},
		{
			name:       "with rest binding",
			input:      "x, ...rest",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokID, tok.TokEOF},
			want:       ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, ast.NewRestOperator(ast.NewIdentifier("...", nil, 0, 3))),
		},
		{
			name:       "with ignored rest binding",
			input:      "x, y, ... = (1, 2, 3)",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokOpAssign, tok.TokOpenParen, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokComma, tok.TokDecLit, tok.TokCloseParen, tok.TokEOF},
			want:       ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1), ast.NewIdentifier("y", nil, 0, 1)}, ast.NewRestOperator(ast.NewIdentifier("...", nil, 0, 3))),
		},
		{
			name:       "with multiple rest bindings (invalid)",
			input:      "x, ...rest, ...",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunParseTestExt(t, tt.name, tt.input, tt.want, tt.wantErr, "OrdinalAssignmentLHS", ordinalAssignmentLHS, StringerCheck[*ast.OrdinalAssignmentLHS], tt.tokenTypes)
		})
	}
}

// labeled_assignment_lhs = "(" ( rename_identifier | rename_type ) { "," ( rename_identifier | rename_type ) } ")" .

func TestLabeledAssignmentLHS(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenTypes []tok.TokenType
		want       *ast.LabeledAssignmentLHS
		wantErr    bool
	}{
		{
			name:       "empty",
			input:      "",
			tokenTypes: []tok.TokenType{tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "single",
			input:      "(x)",
			tokenTypes: []tok.TokenType{tok.TokOpenParen, tok.TokID, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewLabeledAssignmentLHS(
				[]ast.Rename{
					ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), nil),
				}),
		},
		{
			name:       "two",
			input:      "(x, y)",
			tokenTypes: []tok.TokenType{tok.TokOpenParen, tok.TokID, tok.TokComma, tok.TokID, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewLabeledAssignmentLHS(
				[]ast.Rename{
					ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), nil),
					ast.NewRenameIdentifier(ast.NewIdentifier("y", nil, 0, 1), nil),
				}),
		},
		{
			name:       "with rest operator (invalid)",
			input:      "(x, y, ...) ",
			tokenTypes: []tok.TokenType{tok.TokOpenParen, tok.TokID, tok.TokComma, tok.TokID, tok.TokComma, tok.TokOpRest, tok.TokCloseParen, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "two renames",
			input:      "(x: foo, y: bar)",
			tokenTypes: []tok.TokenType{tok.TokOpenParen, tok.TokID, tok.TokColon, tok.TokID, tok.TokComma, tok.TokID, tok.TokColon, tok.TokID, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewLabeledAssignmentLHS(
				[]ast.Rename{
					ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), ast.NewIdentifier("foo", nil, 0, 3)),
					ast.NewRenameIdentifier(ast.NewIdentifier("y", nil, 0, 1), ast.NewIdentifier("bar", nil, 0, 3)),
				}),
		},
		{
			name:       "mixed rename and identifier",
			input:      "(x: foo, y)",
			tokenTypes: []tok.TokenType{tok.TokOpenParen, tok.TokID, tok.TokColon, tok.TokID, tok.TokComma, tok.TokID, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewLabeledAssignmentLHS(
				[]ast.Rename{
					ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), ast.NewIdentifier("foo", nil, 0, 3)),
					ast.NewRenameIdentifier(ast.NewIdentifier("y", nil, 0, 1), nil),
				}),
		},
		{
			name:       "mixed rename and type identifier",
			input:      "(x: foo, Y: Bar)",
			tokenTypes: []tok.TokenType{tok.TokOpenParen, tok.TokID, tok.TokColon, tok.TokID, tok.TokComma, tok.TokTypeID, tok.TokColon, tok.TokTypeID, tok.TokCloseParen, tok.TokEOF},
			want: ast.NewLabeledAssignmentLHS(
				[]ast.Rename{
					ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), ast.NewIdentifier("foo", nil, 0, 3)),
					ast.NewRenameType(ast.NewTypeIdentifier("Y", nil, 0, 1), ast.NewTypeIdentifier("Bar", nil, 0, 3)),
				}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunParseTestExt(t, tt.name, tt.input, tt.want, tt.wantErr, "LabeledAssignmentLHS", labeledAssignmentLHS, StringerCheck[*ast.LabeledAssignmentLHS], tt.tokenTypes)
		})
	}
}

// rename_identifier | rename_type

func TestRename(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenTypes []tok.TokenType
		want       ast.Rename
		wantErr    bool
	}{
		{
			name:       "empty",
			input:      "",
			tokenTypes: []tok.TokenType{tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "not renamed",
			input:      "x",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokEOF},
			want:       ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), nil),
		},
		{
			name:       "not renamed type",
			input:      "Foo",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokEOF},
			want:       ast.NewRenameType(ast.NewTypeIdentifier("Foo", nil, 0, 3), nil),
		},
		{
			name:       "renamed identifier",
			input:      "x: y",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokColon, tok.TokID, tok.TokEOF},
			want: ast.NewRenameIdentifier(
				ast.NewIdentifier("x", nil, 0, 1),
				ast.NewIdentifier("y", nil, 0, 1),
			),
		},
		{
			name:       "renamed type",
			input:      "Foo: Bar",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokColon, tok.TokTypeID, tok.TokEOF},
			want: ast.NewRenameType(
				ast.NewTypeIdentifier("Foo", nil, 0, 3),
				ast.NewTypeIdentifier("Bar", nil, 0, 3),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunParseTestExt(t, tt.name, tt.input, tt.want, tt.wantErr, "Rename", Rename, StringerCheck[ast.Rename], tt.tokenTypes)
		})
	}
}

// "," ( rename_identifier | rename_type )

func TestCommaRename(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenTypes []tok.TokenType
		want       ast.Rename
		wantErr    bool
	}{
		{
			name:       "empty",
			input:      "",
			tokenTypes: []tok.TokenType{tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "comma not renamed",
			input:      ", x",
			tokenTypes: []tok.TokenType{tok.TokComma, tok.TokID, tok.TokEOF},
			want:       ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), nil),
		},
		{
			name:       "comma renamed identifier",
			input:      ", x: y",
			tokenTypes: []tok.TokenType{tok.TokComma, tok.TokID, tok.TokColon, tok.TokID, tok.TokEOF},
			want: ast.NewRenameIdentifier(
				ast.NewIdentifier("x", nil, 0, 1),
				ast.NewIdentifier("y", nil, 0, 1),
			),
		},
		{
			name:       "comma renamed type identifier",
			input:      ", Foo: Bar",
			tokenTypes: []tok.TokenType{tok.TokComma, tok.TokTypeID, tok.TokColon, tok.TokTypeID, tok.TokEOF},
			want:       ast.NewRenameType(ast.NewTypeIdentifier("Foo", nil, 0, 3), ast.NewTypeIdentifier("Bar", nil, 0, 3)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunParseTestExt(t, tt.name, tt.input, tt.want, tt.wantErr, "CommaRename", commaRename, StringerCheck[ast.Rename], tt.tokenTypes)
		})
	}
}
