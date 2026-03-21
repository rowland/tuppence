package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestFunctionDeclarationLHS(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionDeclarationLHS
		wantErr bool
	}{
		{
			name:  "simple lhs",
			input: "atoi",
			want: ast.NewFunctionDeclarationLHS(
				ast.NewFunctionIdentifier("atoi", nil, 0, 4),
				nil,
			),
		},
		{
			name:  "lhs with selector types",
			input: "read[String]",
			want: ast.NewFunctionDeclarationLHS(
				ast.NewFunctionIdentifier("read", nil, 0, 4),
				ast.NewFunctionParameterTypes([]ast.FunctionParameterType{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				}),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionDeclarationLHS", FunctionDeclarationLHS, StringerCheck[*ast.FunctionDeclarationLHS])
		})
	}
}

func TestFunctionDeclarationLHSStringRoundTrip(t *testing.T) {
	input := "read[?String, !Handle]"

	src := source.NewSource([]byte(input), "test.tup")
	tokens, err := tok.Tokenize(src.Contents, src.Filename)
	if err != nil {
		t.Fatalf("Tokenize(%q) = %v", input, err)
	}

	got, remainder, err := FunctionDeclarationLHS(tokens)
	if err != nil {
		t.Fatalf("FunctionDeclarationLHS(%q) unexpected error: %v", input, err)
	}
	if len(remainder) != 1 || remainder[0].Type != tok.TokEOF {
		t.Fatalf("FunctionDeclarationLHS(%q) left remainder: %v", input, remainder)
	}
	if got.String() != input {
		t.Fatalf("FunctionDeclarationLHS(%q) = %q, want %q", input, got.String(), input)
	}
}

func TestFunctionDeclarationType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionDeclarationType
		wantErr bool
	}{
		{
			name:  "pure declaration type with explicit return",
			input: "fn(v: String) !Int",
			want: ast.NewFunctionDeclarationType(
				false,
				[]ast.FunctionTypeParameter{
					ast.NewLabeledParameter(
						ast.NewAnnotations(nil),
						ast.NewIdentifier("v", nil, 0, 1),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					),
				},
				ast.NewReturnType(ast.NewUnionWithError([]ast.UnionMemberType{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				}, true)),
				false,
			),
		},
		{
			name:  "pure declaration type with inferred return",
			input: "fn() _",
			want: ast.NewFunctionDeclarationType(
				false,
				nil,
				nil,
				true,
			),
		},
		{
			name:  "effectful declaration type with omitted return",
			input: "fx(message: String)",
			want: ast.NewFunctionDeclarationType(
				true,
				[]ast.FunctionTypeParameter{
					ast.NewLabeledParameter(
						ast.NewAnnotations(nil),
						ast.NewIdentifier("message", nil, 0, 7),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					),
				},
				nil,
				false,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionDeclarationType", FunctionDeclarationType, StringerCheck[*ast.FunctionDeclarationType])
		})
	}
}

func TestFunctionDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionDeclaration
		wantErr bool
	}{
		{
			name:  "simple pure function declaration",
			input: "sqr = fn(i: Int) Int { i * i }",
			want: ast.NewFunctionDeclaration(
				nil,
				ast.NewFunctionDeclarationLHS(
					ast.NewFunctionIdentifier("sqr", nil, 0, 3),
					nil,
				),
				ast.NewFunctionDeclarationType(
					false,
					[]ast.FunctionTypeParameter{
						ast.NewLabeledParameter(
							ast.NewAnnotations(nil),
							ast.NewIdentifier("i", nil, 0, 1),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
					},
					ast.NewReturnType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					),
					false,
				),
				ast.NewBlock(
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewMulDivExpression(
							ast.NewIdentifier("i", nil, 0, 1),
							ast.OpMul,
							ast.NewIdentifier("i", nil, 0, 1),
						),
					),
				),
			),
		},
		{
			name:  "effectful function declaration with omitted return",
			input: "log = fx(message: String) { message }",
			want: ast.NewFunctionDeclaration(
				nil,
				ast.NewFunctionDeclarationLHS(
					ast.NewFunctionIdentifier("log", nil, 0, 3),
					nil,
				),
				ast.NewFunctionDeclarationType(
					true,
					[]ast.FunctionTypeParameter{
						ast.NewLabeledParameter(
							ast.NewAnnotations(nil),
							ast.NewIdentifier("message", nil, 0, 7),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						),
					},
					nil,
					false,
				),
				ast.NewBlock(
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewIdentifier("message", nil, 0, 7),
					),
				),
			),
		},
		{
			name:  "function declaration with selector types and inferred return",
			input: "read[!String] = fn(handle: Handle) _ { handle }",
			want: ast.NewFunctionDeclaration(
				nil,
				ast.NewFunctionDeclarationLHS(
					ast.NewFunctionIdentifier("read", nil, 0, 4),
					ast.NewFunctionParameterTypes([]ast.FunctionParameterType{
						ast.NewFallibleType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						),
					}),
				),
				ast.NewFunctionDeclarationType(
					false,
					[]ast.FunctionTypeParameter{
						ast.NewLabeledParameter(
							ast.NewAnnotations(nil),
							ast.NewIdentifier("handle", nil, 0, 6),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Handle", nil, 0, 6), nil, 0, 6),
						),
					},
					nil,
					true,
				),
				ast.NewBlock(
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewIdentifier("handle", nil, 0, 6),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionDeclaration", FunctionDeclaration, StringerCheck[*ast.FunctionDeclaration])
		})
	}
}
