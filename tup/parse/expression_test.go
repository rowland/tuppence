package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.Expression
		wantErr bool
	}{
		{
			name:    "empty expression",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		// boolean
		{
			name:  "true",
			input: "true",
			want:  &ast.BooleanLiteral{BooleanValue: true},
		},
		{
			name:  "false",
			input: "false",
			want:  &ast.BooleanLiteral{BooleanValue: false},
		},
		// binary
		{
			name:  "0b1010",
			input: "0b1010",
			want:  &ast.IntegerLiteral{IntegerValue: 10, Base: 2},
		},
		// octal
		{
			name:  "0o12",
			input: "0o12",
			want:  &ast.IntegerLiteral{IntegerValue: 10, Base: 8},
		},
		// decimal
		{
			name:  "123",
			input: "123",
			want:  &ast.IntegerLiteral{IntegerValue: 123},
		},
		// hexadecimal
		{
			name:  "0x1A",
			input: "0x1A",
			want:  &ast.IntegerLiteral{IntegerValue: 26, Base: 16},
		},
		// float
		{
			name:  "1.0",
			input: "1.0",
			want:  &ast.FloatLiteral{FloatValue: 1.0},
		},
		{
			name:  "1.0e10",
			input: "1.0e10",
			want:  &ast.FloatLiteral{FloatValue: 1.0e10},
		},
		// string
		{
			name:  "\"hello\"",
			input: "\"hello\"",
			want:  &ast.StringLiteral{StringValue: "hello"},
		},
		// raw string
		{
			name:  "`hello`",
			input: "`hello`",
			want:  &ast.RawStringLiteral{StringValue: "hello"},
		},
		// logical or
		{
			name:  "true || false",
			input: "true || false",
			want: ast.NewLogicalOrExpression([]ast.Expression{
				ast.NewBooleanLiteral("true", true, nil, 0, 0),
				ast.NewBooleanLiteral("false", false, nil, 0, 0),
			}),
		},
		{
			name:  "false || true",
			input: "false || true",
			want: ast.NewLogicalOrExpression([]ast.Expression{
				ast.NewBooleanLiteral("false", false, nil, 0, 0),
				ast.NewBooleanLiteral("true", true, nil, 0, 0),
			}),
		},
		// // add sub
		{
			name:  "1 + 2",
			input: "1 + 2",
			want: ast.NewAddSubExpression(
				ast.NewDecimalLiteral("1", 1, nil, 0, 0),
				ast.OpAdd,
				ast.NewDecimalLiteral("2", 2, nil, 0, 0)),
		},
		{
			name:  "2 - 1",
			input: "2 - 1",
			want: ast.NewAddSubExpression(
				ast.NewDecimalLiteral("2", 2, nil, 0, 0),
				ast.OpSub,
				ast.NewDecimalLiteral("1", 1, nil, 0, 0)),
		},
		{
			name:  "1 | 2",
			input: "1 | 2",
			want: ast.NewAddSubExpression(
				ast.NewDecimalLiteral("1", 1, nil, 0, 0),
				ast.OpBitOr,
				ast.NewDecimalLiteral("2", 2, nil, 0, 0)),
		},
		// mul div
		{
			name:  "2 * 3",
			input: "2 * 3",
			want: ast.NewMulDivExpression(
				ast.NewDecimalLiteral("2", 2, nil, 0, 0),
				ast.OpMul,
				ast.NewDecimalLiteral("3", 3, nil, 0, 0)),
		},
		// pow
		{
			name:  "3 ^ 4",
			input: "3 ^ 4",
			want: ast.NewPowExpression([]ast.Expression{
				ast.NewDecimalLiteral("3", 3, nil, 0, 0),
				ast.NewDecimalLiteral("4", 4, nil, 0, 0),
			}),
		},
		{
			name:    "invalid expression",
			input:   "1 +",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "tuple expression",
			input: "(1, 2)",
			want: ast.NewTupleLiteral(false, []*ast.TupleMember{
				ast.NewTupleMember(nil, ast.NewDecimalLiteral("1", 1, nil, 0, 0)),
				ast.NewTupleMember(nil, ast.NewDecimalLiteral("2", 2, nil, 0, 0)),
			}),
		},
		{
			name:  "function call",
			input: "foo(1, 2)",
			want: ast.NewFunctionCall(
				ast.NewFunctionIdentifier("foo", nil, 0, 3),
				nil,
				ast.NewFunctionArguments(
					// args
					ast.NewArguments([]*ast.Argument{
						ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 0), false),
						ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 0), false),
					}),
					// labeledArgs
					nil,
					// partialApplication
					false,
				),
				nil,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tok.Tokenize([]byte(tt.input), "test.tup")
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", tt.input, err)
			}
			expression, _, err := Expression(tokens)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expression(%q): want error, got nil", tt.input)
				}
				return
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Expression(%q): got error %v, want nil", tt.input, err)
			}
			if expression == nil {
				t.Fatalf("Expression(%q) = nil, want not nil", tt.input)
			}
			switch want := tt.want.(type) {
			case *ast.IntegerLiteral:
				got, ok := expression.(*ast.IntegerLiteral)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.IntegerValue != want.IntegerValue {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got.IntegerValue, want.IntegerValue)
				}
			case *ast.FloatLiteral:
				got, ok := expression.(*ast.FloatLiteral)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.FloatValue != want.FloatValue {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got.FloatValue, want.FloatValue)
				}
			case *ast.StringLiteral:
				got, ok := expression.(*ast.StringLiteral)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.StringValue != want.StringValue {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got.StringValue, want.StringValue)
				}
			case *ast.RawStringLiteral:
				got, ok := expression.(*ast.RawStringLiteral)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.StringValue != want.StringValue {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got.StringValue, want.StringValue)
				}
			case *ast.LogicalOrExpression:
				got, ok := expression.(*ast.LogicalOrExpression)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.String() != want.String() {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got, want)
					return
				}
			case *ast.AddSubExpression:
				got, ok := expression.(*ast.AddSubExpression)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.String() != want.String() {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got, want)
					return
				}
			case *ast.LogicalAndExpression:
				got, ok := expression.(*ast.LogicalAndExpression)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.String() != want.String() {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got, want)
					return
				}
			case *ast.TypeComparison:
				got, ok := expression.(*ast.TypeComparison)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.String() != want.String() {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got, want)
					return
				}
			case *ast.PowExpression:
				got, ok := expression.(*ast.PowExpression)
				if !ok {
					t.Errorf("Expression(%q) = %T, want %T", tt.input, expression, tt.want)
					return
				}
				if got.String() != want.String() {
					t.Errorf("Expression(%q) = %v, want %v", tt.input, got, want)
					return
				}
			}
		})
	}
}
