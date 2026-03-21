package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestFunctionTypeIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeIdentifier
		wantErr bool
	}{
		{
			name:  "simple type identifier",
			input: "Transformer",
			want:  ast.NewTypeIdentifier("Transformer", nil, 0, 11),
		},
		{
			name:    "lowercase identifier rejected",
			input:   "transformer",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionTypeIdentifier", FunctionTypeIdentifier, StringerCheck[*ast.TypeIdentifier])
		})
	}
}

func TestFunctionTypeDeclarationLHS(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		wantName           *ast.TypeIdentifier
		wantParameterTypes *ast.FunctionParameterTypes
		wantErr            bool
	}{
		{
			name:               "simple lhs",
			input:              "Transformer",
			wantName:           ast.NewTypeIdentifier("Transformer", nil, 0, 11),
			wantParameterTypes: nil,
		},
		{
			name:     "lhs with parameter types",
			input:    "Transformer[String, Int]",
			wantName: ast.NewTypeIdentifier("Transformer", nil, 0, 11),
			wantParameterTypes: ast.NewFunctionParameterTypes([]ast.LocalTypeReference{
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(src.Contents, src.Filename)
			if err != nil {
				t.Fatalf("Tokenize(%q) = %v", test.input, err)
			}
			name, parameterTypes, _, err := FunctionTypeDeclarationLHS(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if name.String() != test.wantName.String() {
				t.Fatalf("name mismatch: want %q, got %q", test.wantName.String(), name.String())
			}
			switch {
			case test.wantParameterTypes == nil && parameterTypes == nil:
			case test.wantParameterTypes == nil || parameterTypes == nil:
				t.Fatalf("parameter types mismatch: want %v, got %v", test.wantParameterTypes, parameterTypes)
			case parameterTypes.String() != test.wantParameterTypes.String():
				t.Fatalf("parameter types mismatch: want %q, got %q", test.wantParameterTypes.String(), parameterTypes.String())
			}
		})
	}
}

func TestFunctionTypeDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionTypeDeclaration
		wantErr bool
	}{
		{
			name:  "simple function type declaration",
			input: "Transformer = fn(Int) String",
			want: ast.NewFunctionTypeDeclaration(
				ast.NewTypeIdentifier("Transformer", nil, 0, 11),
				nil,
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
		{
			name:  "function type declaration with parameter types",
			input: "Transformer[String, Int] = fn(input: Int) String",
			want: ast.NewFunctionTypeDeclaration(
				ast.NewTypeIdentifier("Transformer", nil, 0, 11),
				ast.NewFunctionParameterTypes([]ast.LocalTypeReference{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				}),
				ast.NewFunctionType(
					false,
					[]ast.FunctionTypeParameter{
						ast.NewLabeledParameter(
							ast.NewAnnotations(nil),
							ast.NewIdentifier("input", nil, 0, 5),
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
				"FunctionTypeDeclaration", FunctionTypeDeclaration, StringerCheck[*ast.FunctionTypeDeclaration])
		})
	}
}
