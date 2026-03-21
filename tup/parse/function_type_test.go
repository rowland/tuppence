package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestParameter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Parameter
		wantErr bool
	}{
		{
			name:  "simple parameter",
			input: "Int",
			want: ast.NewParameter(
				ast.NewAnnotations(nil),
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
			),
		},
		{
			name:  "annotated nested function parameter",
			input: "@pure\nfn(Int) String",
			want: ast.NewParameter(
				ast.NewAnnotations([]ast.Annotation{ast.NewSimpleAnnotation("pure")}),
				ast.NewFunctionType(
					false,
					[]ast.FunctionTypeParameter{
						ast.NewParameter(
							ast.NewAnnotations(nil),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
					},
					ast.NewReturnType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Parameter", Parameter, StringerCheck[*ast.Parameter])
		})
	}
}

func TestLabeledParameter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.LabeledParameter
		wantErr bool
	}{
		{
			name:  "simple labeled parameter",
			input: "value: String",
			want: ast.NewLabeledParameter(
				ast.NewAnnotations(nil),
				ast.NewIdentifier("value", nil, 0, 5),
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"LabeledParameter", LabeledParameter, StringerCheck[*ast.LabeledParameter])
		})
	}
}

func TestLabeledParameters(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []ast.FunctionTypeParameter
		wantErr bool
	}{
		{
			name:  "two labeled parameters",
			input: "value: String, count: Int",
			want: []ast.FunctionTypeParameter{
				ast.NewLabeledParameter(
					ast.NewAnnotations(nil),
					ast.NewIdentifier("value", nil, 0, 5),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				),
				ast.NewLabeledParameter(
					ast.NewAnnotations(nil),
					ast.NewIdentifier("count", nil, 0, 5),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
			},
		},
		{
			name:  "labeled parameters with trailing comma",
			input: "value: String, count: Int,",
			want: []ast.FunctionTypeParameter{
				ast.NewLabeledParameter(
					ast.NewAnnotations(nil),
					ast.NewIdentifier("value", nil, 0, 5),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				),
				ast.NewLabeledParameter(
					ast.NewAnnotations(nil),
					ast.NewIdentifier("count", nil, 0, 5),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
			},
		},
		{
			name:  "labeled rest parameter followed by callable parameter",
			input: "args: ...Int, transform: fn(Int) Int",
			want: []ast.FunctionTypeParameter{
				ast.NewLabeledRestParameter(
					ast.NewAnnotations(nil),
					ast.NewIdentifier("args", nil, 0, 4),
					ast.NewRestParameter(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					),
				),
				ast.NewLabeledParameter(
					ast.NewAnnotations(nil),
					ast.NewIdentifier("transform", nil, 0, 9),
					ast.NewFunctionType(
						false,
						[]ast.FunctionTypeParameter{
							ast.NewParameter(
								ast.NewAnnotations(nil),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							),
						},
						ast.NewReturnType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
					),
				),
			},
		},
	}

	check := func(t *testing.T, input, parserName string, got, want []ast.FunctionTypeParameter) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("%s(%q) len = %d, want %d", parserName, input, len(got), len(want))
		}
		for i := range want {
			if got[i].String() != want[i].String() {
				t.Fatalf("%s(%q)[%d] = %q, want %q", parserName, input, i, got[i].String(), want[i].String())
			}
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"LabeledParameters", LabeledParameters, check)
		})
	}
}

func TestParameters(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []ast.FunctionTypeParameter
		wantErr bool
	}{
		{
			name:  "two positional parameters",
			input: "String, Int",
			want: []ast.FunctionTypeParameter{
				ast.NewParameter(
					ast.NewAnnotations(nil),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				),
				ast.NewParameter(
					ast.NewAnnotations(nil),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
			},
		},
		{
			name:  "parameters with trailing comma",
			input: "String, Int,",
			want: []ast.FunctionTypeParameter{
				ast.NewParameter(
					ast.NewAnnotations(nil),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				),
				ast.NewParameter(
					ast.NewAnnotations(nil),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
			},
		},
		{
			name:  "rest parameter followed by callable parameter",
			input: "...Int, fn(Int) Int",
			want: []ast.FunctionTypeParameter{
				ast.NewRestParameter(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
				ast.NewParameter(
					ast.NewAnnotations(nil),
					ast.NewFunctionType(
						false,
						[]ast.FunctionTypeParameter{
							ast.NewParameter(
								ast.NewAnnotations(nil),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							),
						},
						ast.NewReturnType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
					),
				),
			},
		},
	}

	check := func(t *testing.T, input, parserName string, got, want []ast.FunctionTypeParameter) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("%s(%q) len = %d, want %d", parserName, input, len(got), len(want))
		}
		for i := range want {
			if got[i].String() != want[i].String() {
				t.Fatalf("%s(%q)[%d] = %q, want %q", parserName, input, i, got[i].String(), want[i].String())
			}
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Parameters", Parameters, check)
		})
	}
}

func TestRestParameter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.RestParameter
		wantErr bool
	}{
		{
			name:  "simple rest parameter",
			input: "...Int",
			want: ast.NewRestParameter(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
			),
		},
		{
			name:  "function rest parameter",
			input: "...fn(Int) String",
			want: ast.NewRestParameter(
				ast.NewFunctionType(
					false,
					[]ast.FunctionTypeParameter{
						ast.NewParameter(
							ast.NewAnnotations(nil),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
					},
					ast.NewReturnType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"RestParameter", RestParameter, StringerCheck[*ast.RestParameter])
		})
	}
}

func TestLabeledRestParameter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.LabeledRestParameter
		wantErr bool
	}{
		{
			name:  "simple labeled rest parameter",
			input: "args: ...Int",
			want: ast.NewLabeledRestParameter(
				ast.NewAnnotations(nil),
				ast.NewIdentifier("args", nil, 0, 4),
				ast.NewRestParameter(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"LabeledRestParameter", LabeledRestParameter, StringerCheck[*ast.LabeledRestParameter])
		})
	}
}

func TestFunctionType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionType
		wantErr bool
	}{
		{
			name:  "simple function type",
			input: "fn(Int) String",
			want: ast.NewFunctionType(
				false,
				[]ast.FunctionTypeParameter{
					ast.NewParameter(
						ast.NewAnnotations(nil),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					),
				},
				ast.NewReturnType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				),
			),
		},
		{
			name:  "simple fx function type",
			input: "fx(message: String) error",
			want: ast.NewFunctionType(
				true,
				[]ast.FunctionTypeParameter{
					ast.NewLabeledParameter(
						ast.NewAnnotations(nil),
						ast.NewIdentifier("message", nil, 0, 7),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					),
				},
				ast.NewReturnType(ast.NewInferredErrorType()),
			),
		},
		{
			name:  "nested function type parameter",
			input: "fn(input: a, process: fn(a) b) b",
			want: ast.NewFunctionType(
				false,
				[]ast.FunctionTypeParameter{
					ast.NewLabeledParameter(
						ast.NewAnnotations(nil),
						ast.NewIdentifier("input", nil, 0, 5),
						ast.NewIdentifier("a", nil, 0, 1),
					),
					ast.NewLabeledParameter(
						ast.NewAnnotations(nil),
						ast.NewIdentifier("process", nil, 0, 7),
						ast.NewFunctionType(
							false,
							[]ast.FunctionTypeParameter{
								ast.NewParameter(
									ast.NewAnnotations(nil),
									ast.NewIdentifier("a", nil, 0, 1),
								),
							},
							ast.NewReturnType(ast.NewIdentifier("b", nil, 0, 1)),
						),
					),
				},
				ast.NewReturnType(ast.NewIdentifier("b", nil, 0, 1)),
			),
		},
		{
			name:  "rest parameter followed by callable parameter",
			input: "fn(args: ...Int, transform: fn(Int) Int) Int",
			want: ast.NewFunctionType(
				false,
				[]ast.FunctionTypeParameter{
					ast.NewLabeledRestParameter(
						ast.NewAnnotations(nil),
						ast.NewIdentifier("args", nil, 0, 4),
						ast.NewRestParameter(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
					),
					ast.NewLabeledParameter(
						ast.NewAnnotations(nil),
						ast.NewIdentifier("transform", nil, 0, 9),
						ast.NewFunctionType(
							false,
							[]ast.FunctionTypeParameter{
								ast.NewParameter(
									ast.NewAnnotations(nil),
									ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
								),
							},
							ast.NewReturnType(
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							),
						),
					),
				},
				ast.NewReturnType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionType", FunctionType, StringerCheck[*ast.FunctionType])
		})
	}
}
