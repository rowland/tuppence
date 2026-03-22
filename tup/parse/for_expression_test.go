package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestForBlock(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ForBlock
		wantErr bool
	}{
		{
			name:  "empty for block",
			input: "{}",
			want:  ast.NewForBlock(nil, nil),
		},
		{
			name:  "for block with final expression",
			input: "{ i + 1 }",
			want: ast.NewForBlock(
				[]ast.Statement{},
				ast.NewAddSubExpression(
					ast.NewIdentifier("i", nil, 0, 1),
					ast.OpAdd,
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				),
			),
		},
		{
			name:  "for block with statement and final expression",
			input: "{ print(i)\n i + 1 }",
			want: ast.NewForBlock(
				[]ast.Statement{
					ast.NewFunctionCall(
						ast.NewFunctionIdentifier("print", nil, 0, 5),
						nil,
						ast.NewFunctionArguments(
							ast.NewArguments([]*ast.Argument{
								ast.NewArgument(ast.NewIdentifier("i", nil, 0, 1), false),
							}),
							nil,
							false,
						),
						nil,
					),
				},
				ast.NewAddSubExpression(
					ast.NewIdentifier("i", nil, 0, 1),
					ast.OpAdd,
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ForBlock", ForBlock, StringerCheck[*ast.ForBlock])
		})
	}
}

func TestInitializer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Initializer
		wantErr bool
	}{
		{
			name:  "simple initializer",
			input: "i = 0",
			want: ast.NewInitializer(
				ast.NewAssignment(
					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
					ast.Immutable,
					ast.NewDecimalLiteral("0", 0, nil, 0, 1),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Initializer", Initializer, StringerCheck[*ast.Initializer])
		})
	}
}

func TestIterable(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Iterable
		wantErr bool
	}{
		{
			name:  "simple iterable",
			input: "items",
			want:  ast.NewIterable(ast.NewIdentifier("items", nil, 0, 5)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Iterable", Iterable, StringerCheck[*ast.Iterable])
		})
	}
}

func TestStepExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.StepExpression
		wantErr bool
	}{
		{
			name:  "simple step expression",
			input: "i + 1",
			want: ast.NewStepExpression(
				ast.NewAddSubExpression(
					ast.NewIdentifier("i", nil, 0, 1),
					ast.OpAdd,
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"StepExpression", StepExpression, StringerCheck[*ast.StepExpression])
		})
	}
}

func TestForHeader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ForHeader
		wantErr bool
	}{
		{
			name:  "initializer only",
			input: "i = 0",
			want: ast.NewForHeader(
				ast.NewInitializer(
					ast.NewAssignment(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
						ast.Immutable,
						ast.NewDecimalLiteral("0", 0, nil, 0, 1),
					),
				),
				nil,
				nil,
			),
		},
		{
			name:  "initializer and condition",
			input: "i = 0; i < 10",
			want: ast.NewForHeader(
				ast.NewInitializer(
					ast.NewAssignment(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
						ast.Immutable,
						ast.NewDecimalLiteral("0", 0, nil, 0, 1),
					),
				),
				ast.NewRelationalComparison(
					ast.NewIdentifier("i", nil, 0, 1),
					ast.OpLt,
					ast.NewDecimalLiteral("10", 10, nil, 0, 2),
				),
				nil,
			),
		},
		{
			name:  "initializer condition and step",
			input: "i = 0; i < 10; i + 1",
			want: ast.NewForHeader(
				ast.NewInitializer(
					ast.NewAssignment(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
						ast.Immutable,
						ast.NewDecimalLiteral("0", 0, nil, 0, 1),
					),
				),
				ast.NewRelationalComparison(
					ast.NewIdentifier("i", nil, 0, 1),
					ast.OpLt,
					ast.NewDecimalLiteral("10", 10, nil, 0, 2),
				),
				ast.NewStepExpression(
					ast.NewAddSubExpression(
						ast.NewIdentifier("i", nil, 0, 1),
						ast.OpAdd,
						ast.NewDecimalLiteral("1", 1, nil, 0, 1),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ForHeader", ForHeader, StringerCheck[*ast.ForHeader])
		})
	}
}

func TestIterableHeader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.IterableHeader
		wantErr bool
	}{
		{
			name:  "simple iterable header",
			input: "i in items",
			want: ast.NewIterableHeader(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
				ast.NewIterable(ast.NewIdentifier("items", nil, 0, 5)),
			),
		},
		{
			name:  "tuple iterable header",
			input: "a, b in pairs",
			want: ast.NewIterableHeader(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
					ast.NewIdentifier("a", nil, 0, 1),
					ast.NewIdentifier("b", nil, 0, 1),
				}, nil),
				ast.NewIterable(ast.NewIdentifier("pairs", nil, 0, 5)),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"IterableHeader", IterableHeader, StringerCheck[*ast.IterableHeader])
		})
	}
}

func TestForInHeader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ForInHeader
		wantErr bool
	}{
		{
			name:  "simple for-in header",
			input: "i in items",
			want: ast.NewForInHeader(
				nil,
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
				ast.NewIterable(ast.NewIdentifier("items", nil, 0, 5)),
				nil,
			),
		},
		{
			name:  "for-in header with initializer",
			input: "acc = 0; i in items",
			want: ast.NewForInHeader(
				ast.NewInitializer(
					ast.NewAssignment(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("acc", nil, 0, 3)}, nil),
						ast.Immutable,
						ast.NewDecimalLiteral("0", 0, nil, 0, 1),
					),
				),
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
				ast.NewIterable(ast.NewIdentifier("items", nil, 0, 5)),
				nil,
			),
		},
		{
			name:  "for-in header with initializer and step",
			input: "acc = 0; i in items; acc + i",
			want: ast.NewForInHeader(
				ast.NewInitializer(
					ast.NewAssignment(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("acc", nil, 0, 3)}, nil),
						ast.Immutable,
						ast.NewDecimalLiteral("0", 0, nil, 0, 1),
					),
				),
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
				ast.NewIterable(ast.NewIdentifier("items", nil, 0, 5)),
				ast.NewStepExpression(
					ast.NewAddSubExpression(
						ast.NewIdentifier("acc", nil, 0, 3),
						ast.OpAdd,
						ast.NewIdentifier("i", nil, 0, 1),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ForInHeader", ForInHeader, StringerCheck[*ast.ForInHeader])
		})
	}
}

func TestForExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ForExpression
		wantErr bool
	}{
		{
			name:  "bare for loop",
			input: "for { break }",
			want: ast.NewForExpression(
				nil,
				ast.NewForBlock(
					[]ast.Statement{
						ast.NewBreakExpression(nil),
					},
					nil,
				),
			),
		},
		{
			name:  "for loop with header and final expression",
			input: "for i = 0; i < 10 { i + 1 }",
			want: ast.NewForExpression(
				ast.NewForHeader(
					ast.NewInitializer(
						ast.NewAssignment(
							ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
							ast.Immutable,
							ast.NewDecimalLiteral("0", 0, nil, 0, 1),
						),
					),
					ast.NewRelationalComparison(
						ast.NewIdentifier("i", nil, 0, 1),
						ast.OpLt,
						ast.NewDecimalLiteral("10", 10, nil, 0, 2),
					),
					nil,
				),
				ast.NewForBlock(
					[]ast.Statement{},
					ast.NewAddSubExpression(
						ast.NewIdentifier("i", nil, 0, 1),
						ast.OpAdd,
						ast.NewDecimalLiteral("1", 1, nil, 0, 1),
					),
				),
			),
		},
		{
			name:  "for-in loop with initializer",
			input: "for acc = 0; i in items { acc + i }",
			want: ast.NewForExpression(
				ast.NewForInHeader(
					ast.NewInitializer(
						ast.NewAssignment(
							ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("acc", nil, 0, 3)}, nil),
							ast.Immutable,
							ast.NewDecimalLiteral("0", 0, nil, 0, 1),
						),
					),
					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil),
					ast.NewIterable(ast.NewIdentifier("items", nil, 0, 5)),
					nil,
				),
				ast.NewForBlock(
					[]ast.Statement{},
					ast.NewAddSubExpression(
						ast.NewIdentifier("acc", nil, 0, 3),
						ast.OpAdd,
						ast.NewIdentifier("i", nil, 0, 1),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ForExpression", ForExpression, StringerCheck[*ast.ForExpression])
		})
	}
}

func TestInlineForExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.InlineForExpression
		wantErr bool
	}{
		{
			name:  "inline for with initializer",
			input: "inline for acc = (); name, value in some_tuple { acc }",
			want: ast.NewInlineForExpression(
				ast.NewForInHeader(
					ast.NewInitializer(
						ast.NewAssignment(
							ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("acc", nil, 0, 3)}, nil),
							ast.Immutable,
							ast.NewTupleLiteral(false, []*ast.TupleMember{}),
						),
					),
					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
						ast.NewIdentifier("name", nil, 0, 4),
						ast.NewIdentifier("value", nil, 0, 5),
					}, nil),
					ast.NewIterable(ast.NewIdentifier("some_tuple", nil, 0, 10)),
					nil,
				),
				ast.NewForBlock(
					[]ast.Statement{},
					ast.NewIdentifier("acc", nil, 0, 3),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"InlineForExpression", InlineForExpression, StringerCheck[*ast.InlineForExpression])
		})
	}
}
