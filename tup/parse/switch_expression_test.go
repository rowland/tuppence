package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func switchBody(expr ast.Expression) *ast.FunctionBlock {
	return ast.NewFunctionBlock(nil, ast.NewBlockBody(nil, expr))
}

func TestMatchElement(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.MatchElement
		wantErr bool
	}{
		{
			name:  "constant",
			input: "1",
			want:  ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
		},
		{
			name:  "range",
			input: "1..10",
			want: ast.NewRange(
				ast.NewRangeBound(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewRangeBound(ast.NewDecimalLiteral("10", 10, nil, 0, 2)),
			),
		},
		{
			name:  "type reference",
			input: "Int",
			want:  ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
		},
		{
			name:  "inferred error type",
			input: "error",
			want:  ast.NewInferredErrorType(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"MatchElement", MatchElement, StringerCheck[ast.MatchElement])
		})
	}
}

func TestListMatch(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ListMatch
		wantErr bool
	}{
		{
			name:  "constants and ranges",
			input: "1..10, 15, 20..30",
			want: ast.NewListMatch([]ast.MatchElement{
				ast.NewRange(
					ast.NewRangeBound(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewRangeBound(ast.NewDecimalLiteral("10", 10, nil, 0, 2)),
				),
				ast.NewConstant(ast.NewDecimalLiteral("15", 15, nil, 0, 2)),
				ast.NewRange(
					ast.NewRangeBound(ast.NewDecimalLiteral("20", 20, nil, 0, 2)),
					ast.NewRangeBound(ast.NewDecimalLiteral("30", 30, nil, 0, 2)),
				),
			}),
		},
		{
			name:  "type list",
			input: "Int, Int8",
			want: ast.NewListMatch([]ast.MatchElement{
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int8", nil, 0, 4), nil, 0, 4),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ListMatch", ListMatch, StringerCheck[*ast.ListMatch])
		})
	}
}

func TestPattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.Pattern
		wantErr bool
	}{
		{
			name:  "wildcard",
			input: "_",
			want:  ast.NewWildcardPattern(ast.NewIdentifier("_", nil, 0, 1)),
		},
		{
			name:  "constant pattern",
			input: "1",
			want:  ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
		},
		{
			name:  "type pattern without destructure",
			input: "String",
			want:  ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
		},
		{
			name:  "typed pattern",
			input: "Hearts(3..10)",
			want: ast.NewTypedPattern(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Hearts", nil, 0, 6), nil, 0, 6),
				ast.NewTuplePattern([]ast.Pattern{
					ast.NewRange(
						ast.NewRangeBound(ast.NewDecimalLiteral("3", 3, nil, 0, 1)),
						ast.NewRangeBound(ast.NewDecimalLiteral("10", 10, nil, 0, 2)),
					),
				}),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Pattern", Pattern, StringerCheck[ast.Pattern])
		})
	}
}

func TestPatternMatch(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.Pattern
		wantErr bool
	}{
		{
			name:  "structured tuple",
			input: "(1, 2)",
			want: ast.NewTuplePattern([]ast.Pattern{
				ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewConstant(ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
			}),
		},
		{
			name:  "typed labeled pattern",
			input: "Point(x: 0, y: 0)",
			want: ast.NewTypedPattern(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Point", nil, 0, 5), nil, 0, 5),
				ast.NewLabeledPattern([]*ast.LabeledPatternMember{
					ast.NewLabeledPatternMember(ast.NewIdentifier("x", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1))),
					ast.NewLabeledPatternMember(ast.NewIdentifier("y", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1))),
				}),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"PatternMatch", PatternMatch, StringerCheck[ast.Pattern])
		})
	}
}

func TestMatchCondition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.MatchCondition
		wantErr bool
	}{
		{
			name:  "list match",
			input: "1, 2, 3",
			want: ast.NewListMatch([]ast.MatchElement{
				ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewConstant(ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				ast.NewConstant(ast.NewDecimalLiteral("3", 3, nil, 0, 1)),
			}),
		},
		{
			name:  "pattern",
			input: "Point(x: 0, y: 0)",
			want: ast.NewTypedPattern(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Point", nil, 0, 5), nil, 0, 5),
				ast.NewLabeledPattern([]*ast.LabeledPatternMember{
					ast.NewLabeledPatternMember(ast.NewIdentifier("x", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1))),
					ast.NewLabeledPatternMember(ast.NewIdentifier("y", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1))),
				}),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"MatchCondition", MatchCondition, StringerCheck[ast.MatchCondition])
		})
	}
}

func TestStructuredMatch(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.Pattern
		wantErr bool
	}{
		{
			name:  "labeled",
			input: "(x: 1, y: 2)",
			want: ast.NewLabeledPattern([]*ast.LabeledPatternMember{
				ast.NewLabeledPatternMember(ast.NewIdentifier("x", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1))),
				ast.NewLabeledPatternMember(ast.NewIdentifier("y", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("2", 2, nil, 0, 1))),
			}),
		},
		{
			name:  "array",
			input: "[_, ...]",
			want: ast.NewArrayPattern([]ast.Pattern{
				ast.NewWildcardPattern(ast.NewIdentifier("_", nil, 0, 1)),
			}, true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"StructuredMatch", StructuredMatch, StringerCheck[ast.Pattern])
		})
	}
}

func TestTuplePattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TuplePattern
		wantErr bool
	}{
		{
			name:  "tuple pattern",
			input: "(1, _)",
			want: ast.NewTuplePattern([]ast.Pattern{
				ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewWildcardPattern(ast.NewIdentifier("_", nil, 0, 1)),
			}),
		},
		{
			name:  "tuple pattern with identifiers",
			input: "(x, y)",
			want: ast.NewTuplePattern([]ast.Pattern{
				ast.NewConstant(ast.NewScopedIdentifier([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)})),
				ast.NewConstant(ast.NewScopedIdentifier([]*ast.Identifier{ast.NewIdentifier("y", nil, 0, 1)})),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TuplePattern", TuplePattern, StringerCheck[*ast.TuplePattern])
		})
	}
}

func TestLabeledPattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.LabeledPattern
		wantErr bool
	}{
		{
			name:  "labeled pattern",
			input: "(x: 1, y: _)",
			want: ast.NewLabeledPattern([]*ast.LabeledPatternMember{
				ast.NewLabeledPatternMember(
					ast.NewIdentifier("x", nil, 0, 1),
					ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				),
				ast.NewLabeledPatternMember(
					ast.NewIdentifier("y", nil, 0, 1),
					ast.NewWildcardPattern(ast.NewIdentifier("_", nil, 0, 1)),
				),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"LabeledPattern", LabeledPattern, StringerCheck[*ast.LabeledPattern])
		})
	}
}

func TestArrayPattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ArrayPattern
		wantErr bool
	}{
		{
			name:  "specific elements",
			input: "[1, 2, 3]",
			want: ast.NewArrayPattern([]ast.Pattern{
				ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewConstant(ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				ast.NewConstant(ast.NewDecimalLiteral("3", 3, nil, 0, 1)),
			}, false),
		},
		{
			name:  "rest marker",
			input: "[_, ...]",
			want: ast.NewArrayPattern([]ast.Pattern{
				ast.NewWildcardPattern(ast.NewIdentifier("_", nil, 0, 1)),
			}, true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ArrayPattern", ArrayPattern, StringerCheck[*ast.ArrayPattern])
		})
	}
}

func TestTypedPattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypedPattern
		wantErr bool
	}{
		{
			name:  "labeled typed pattern",
			input: "Point(x: 0, y: 0)",
			want: ast.NewTypedPattern(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Point", nil, 0, 5), nil, 0, 5),
				ast.NewLabeledPattern([]*ast.LabeledPatternMember{
					ast.NewLabeledPatternMember(
						ast.NewIdentifier("x", nil, 0, 1),
						ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1)),
					),
					ast.NewLabeledPatternMember(
						ast.NewIdentifier("y", nil, 0, 1),
						ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1)),
					),
				}),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypedPattern", TypedPattern, StringerCheck[*ast.TypedPattern])
		})
	}
}

func TestSwitchCase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.SwitchCase
		wantErr bool
	}{
		{
			name:  "constant list case",
			input: `1, 2 { "small" }`,
			want: ast.NewSwitchCase(
				ast.NewListMatch([]ast.MatchElement{
					ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
					ast.NewConstant(ast.NewDecimalLiteral("2", 2, nil, 0, 1)),
				}),
				switchBody(ast.NewStringLiteral(`"small"`, "small", nil, 0, 7)),
			),
		},
		{
			name:  "typed pattern case",
			input: `Point(x: 0, y: 0) { "origin" }`,
			want: ast.NewSwitchCase(
				ast.NewTypedPattern(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Point", nil, 0, 5), nil, 0, 5),
					ast.NewLabeledPattern([]*ast.LabeledPatternMember{
						ast.NewLabeledPatternMember(ast.NewIdentifier("x", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1))),
						ast.NewLabeledPatternMember(ast.NewIdentifier("y", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1))),
					}),
				),
				switchBody(ast.NewStringLiteral(`"origin"`, "origin", nil, 0, 8)),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"SwitchCase", SwitchCase, StringerCheck[*ast.SwitchCase])
		})
	}
}

func TestSwitchExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.SwitchExpression
		wantErr bool
	}{
		{
			name: "constant and else",
			input: `switch value {
    1 { "one" }
    else { "other" }
}`,
			want: ast.NewSwitchExpression(
				ast.NewIdentifier("value", nil, 0, 5),
				[]*ast.SwitchCase{
					ast.NewSwitchCase(
						ast.NewConstant(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
						switchBody(ast.NewStringLiteral(`"one"`, "one", nil, 0, 5)),
					),
				},
				switchBody(ast.NewStringLiteral(`"other"`, "other", nil, 0, 7)),
			),
		},
		{
			name: "try style error switch",
			input: `switch result {
    error { return it }
    else { it }
}`,
			want: ast.NewSwitchExpression(
				ast.NewIdentifier("result", nil, 0, 6),
				[]*ast.SwitchCase{
					ast.NewSwitchCase(
						ast.NewInferredErrorType(),
						switchBody(ast.NewReturnExpression(ast.NewItExpression(nil, 0, 2))),
					),
				},
				switchBody(ast.NewItExpression(nil, 0, 2)),
			),
		},
		{
			name: "typed pattern and type list",
			input: `switch value {
    Int, Int8 { |i| i }
    Point(x: 0, y: 0) { "origin" }
}`,
			want: ast.NewSwitchExpression(
				ast.NewIdentifier("value", nil, 0, 5),
				[]*ast.SwitchCase{
					ast.NewSwitchCase(
						ast.NewListMatch([]ast.MatchElement{
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int8", nil, 0, 4), nil, 0, 4),
						}),
						ast.NewFunctionBlock(
							ast.NewBlockParameters(ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("i", nil, 0, 1)}, nil)),
							ast.NewBlockBody(nil, ast.NewIdentifier("i", nil, 0, 1)),
						),
					),
					ast.NewSwitchCase(
						ast.NewTypedPattern(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Point", nil, 0, 5), nil, 0, 5),
							ast.NewLabeledPattern([]*ast.LabeledPatternMember{
								ast.NewLabeledPatternMember(ast.NewIdentifier("x", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1))),
								ast.NewLabeledPatternMember(ast.NewIdentifier("y", nil, 0, 1), ast.NewConstant(ast.NewDecimalLiteral("0", 0, nil, 0, 1))),
							}),
						),
						switchBody(ast.NewStringLiteral(`"origin"`, "origin", nil, 0, 8)),
					),
				},
				nil,
			),
		},
		{
			name: "semicolon between cases is invalid",
			input: `switch value {
    1 { "one" };
    else { "other" }
}`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"SwitchExpression", SwitchExpression, StringerCheck[*ast.SwitchExpression])
		})
	}
}
