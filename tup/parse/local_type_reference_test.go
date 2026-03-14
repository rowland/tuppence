package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestLocalTypeReference(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.LocalTypeReference
		wantErr bool
	}{
		{
			name:  "identifier",
			input: "x",
			want:  ast.NewIdentifier("x", nil, 0, 1),
		},
		{
			name:  "type_reference",
			input: "Foo",
			want:  ast.NewTypeReference(nil, ast.NewTypeIdentifier("Foo", nil, 0, 3), nil, 0, 3),
		},
		{
			name:  "qualified type_reference",
			input: "foo.bar.Foo",
			want: ast.NewTypeReference(
				[]*ast.Identifier{ast.NewIdentifier("foo", nil, 0, 3), ast.NewIdentifier("bar", nil, 0, 3)},
				ast.NewTypeIdentifier("Foo", nil, 0, 3),
				nil, 0, 0,
			),
		},
		{
			name:    "number literal",
			input:   "1",
			want:    nil,
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
			localTypeReference, _, err := LocalTypeReference(tokens)
			if err != nil {
				if test.wantErr {
					return
				}
				t.Fatalf("LocalTypeReference(%q) = %v, want nil", test.input, err)
			}
			if localTypeReference == nil {
				t.Errorf("LocalTypeReference(%q) = nil, want not nil", test.input)
			}
			switch want := test.want.(type) {
			case *ast.TypeReference:
				got, ok := localTypeReference.(*ast.TypeReference)
				if !ok {
					t.Errorf("LocalTypeReference(%q) = %T, want %T", test.input, localTypeReference, test.want)
					return
				}
				if got.TypeIdentifier.Name != want.TypeIdentifier.Name {
					t.Errorf("LocalTypeReference(%q).TypeIdentifier.Name: got %v, want %v", test.input, got.TypeIdentifier.Name, want.TypeIdentifier.Name)
				}
			case *ast.Identifier:
				got, ok := localTypeReference.(*ast.Identifier)
				if !ok {
					t.Errorf("LocalTypeReference(%q) = %T, want %T", test.input, localTypeReference, test.want)
					return
				}
				if got.Name != want.Name {
					t.Errorf("LocalTypeReference(%q).Name: got %v, want %v", test.input, got.Name, want.Name)
				}
			}
		})
	}
}
