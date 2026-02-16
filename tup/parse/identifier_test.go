package parse

import (
	"slices"
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestIdentifier(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.Identifier
		wantErr bool
	}{
		{input: "x", want: ast.NewIdentifier("x", nil, 0, 1)},
		{input: "_x", want: ast.NewIdentifier("_x", nil, 0, 2)},
		{input: "x_", want: ast.NewIdentifier("x_", nil, 0, 2)},
		{input: "x_y", want: ast.NewIdentifier("x_y", nil, 0, 3)},
		{input: "x_y_", want: ast.NewIdentifier("x_y_", nil, 0, 4)},
		{input: "x_y_z", want: ast.NewIdentifier("x_y_z", nil, 0, 5)},
		{input: "XYZ", want: nil, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}

			got, _, err := Identifier(tokens)

			if test.wantErr {
				if err == nil {
					t.Errorf("Identifier(%q): want error", test.input)
				}
				return
			}

			if !test.wantErr && err != nil {
				t.Fatalf("Identifier(%q): got error %v, want nil", test.input, err)
			}

			if got == nil {
				t.Fatalf("Identifier(%q) = nil, want %v", test.input, test.want)
			}

			if got.Name != test.want.Name {
				t.Errorf("Identifier(%q) = %v, want %v", test.input, got, test.want)
			}
		})
	}
}

func TestTypeIdentifier(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.TypeIdentifier
		wantErr bool
	}{
		{input: "Foo", want: ast.NewTypeIdentifier("Foo", nil, 0, 3)},
		{input: "Foo_", want: ast.NewTypeIdentifier("Foo_", nil, 0, 4)},
		{input: "Foo_Bar", want: ast.NewTypeIdentifier("Foo_Bar", nil, 0, 7)},
		{input: "Foo_Bar_", want: ast.NewTypeIdentifier("Foo_Bar_", nil, 0, 8)},
		{input: "Foo_Bar_Baz", want: ast.NewTypeIdentifier("Foo_Bar_Baz", nil, 0, 11)},
		{input: "xyz", want: nil, wantErr: true},
		{input: "123", want: nil, wantErr: true},
		{input: "123Foo", want: nil, wantErr: true},
		{input: "Foo123", want: ast.NewTypeIdentifier("Foo123", nil, 0, 5)},
		{input: "Foo_123", want: ast.NewTypeIdentifier("Foo_123", nil, 0, 7)},
		{input: "Foo_123_", want: ast.NewTypeIdentifier("Foo_123_", nil, 0, 8)},
		{input: "Foo_123_Bar", want: ast.NewTypeIdentifier("Foo_123_Bar", nil, 0, 11)},
		{input: "Foo_123_Bar_", want: ast.NewTypeIdentifier("Foo_123_Bar_", nil, 0, 12)},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}

			got, _, err := TypeIdentifier(tokens)

			if test.wantErr {
				if err == nil {
					t.Errorf("TypeIdentifier(%q): want error", test.input)
				}
				return
			}

			if !test.wantErr && err != nil {
				t.Fatalf("TypeIdentifier(%q): got error %v, want nil", test.input, err)
			}

			if got == nil {
				t.Fatalf("TypeIdentifier(%q): got nil, want %v", test.input, test.want)
			}

			if got.Name != test.want.Name {
				t.Errorf("TypeIdentifier(%q).Name: got %v, want %v", test.input, got.Name, test.want.Name)
			}
		})
	}
}

func TestRenameIdentifier(t *testing.T) {
	tests := []struct {
		input      string
		tokenTypes []tok.TokenType
		want       *ast.RenameIdentifier
		wantErr    bool
	}{
		{
			input:      "x: y",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokColon, tok.TokID, tok.TokEOF},
			want:       ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), ast.NewIdentifier("y", nil, 0, 1)),
			wantErr:    false,
		},
		{
			input:      "x:y",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokEOF},
			want:       ast.NewRenameIdentifier(ast.NewIdentifier("x", nil, 0, 1), ast.NewIdentifier("y", nil, 0, 1)),
			wantErr:    false,
		},
		{
			input:      "x: Y",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokColon, tok.TokTypeID, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
		{
			input:      "x:Y",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokColonNoSpace, tok.TokTypeID, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}

			if !slices.Equal(tok.Types(tokens), test.tokenTypes) {
				t.Fatalf("tokenTypes: %v, want %v", tok.Types(tokens), test.tokenTypes)
			}

			got, _, err := RenameIdentifier(tokens)
			// fmt.Printf("got: %#v\n", got)

			if test.wantErr {
				if err == nil {
					t.Errorf("RenameIdentifier(%q): want error (E1)", test.input)
				}
				return
			}

			if !test.wantErr && err != nil {
				t.Fatalf("RenameIdentifier(%q): got error %v, want nil (E2)", test.input, err)
			}
			if got == nil {
				t.Fatalf("RenameIdentifier(%q): got nil, want %v (E3)", test.input, test.want)
			}
			if got.Identifier == nil {
				t.Fatalf("RenameIdentifier(%q).Identifier: got nil, want %v (E4)", test.input, test.want.Identifier)
			}
			if got.Identifier.Name != test.want.Identifier.Name {
				t.Errorf("RenameIdentifier(%q).Identifier.Name: got %v, want %v (E5)", test.input, got.Identifier.Name, test.want.Identifier.Name)
			}
			if test.want.Original == nil {
				if got.Original != nil {
					t.Fatalf("RenameIdentifier(%q).Original: got %v, want nil (E6)", test.input, got.Original)
				}
			} else {
				if got.Original == nil {
					t.Fatalf("RenameIdentifier(%q).Original: got nil, want %v (E7)", test.input, test.want.Original)
				}
				if got.Original.Name != test.want.Original.Name {
					t.Errorf("RenameIdentifier(%q).Original.Name: got %v, want %v (E8)", test.input, got.Original.Name, test.want.Original.Name)
				}
			}
		})
	}
}

func TestRenameType(t *testing.T) {
	tests := []struct {
		input      string
		tokenTypes []tok.TokenType
		want       *ast.RenameType
		wantErr    bool
	}{
		{
			input:      "Foo: Bar",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokColon, tok.TokTypeID, tok.TokEOF},
			want:       ast.NewRenameType(ast.NewTypeIdentifier("Foo", nil, 0, 3), ast.NewTypeIdentifier("Bar", nil, 0, 3)),
			wantErr:    false,
		},
		{
			input:      "Foo:Bar",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokColonNoSpace, tok.TokTypeID, tok.TokEOF},
			want:       ast.NewRenameType(ast.NewTypeIdentifier("Foo", nil, 0, 3), ast.NewTypeIdentifier("Bar", nil, 0, 3)),
			wantErr:    false,
		},
		{
			input:      "Foo: y",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokColon, tok.TokID, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}

			if !slices.Equal(tok.Types(tokens), test.tokenTypes) {
				t.Fatalf("tokenTypes: %v, want %v", tok.Types(tokens), test.tokenTypes)
			}

			got, _, err := RenameType(tokens)

			if test.wantErr {
				if err == nil {
					t.Errorf("RenameType(%q): want error", test.input)
				}
				return
			}

			if !test.wantErr && err != nil {
				t.Fatalf("RenameType(%q): got error %v, want nil", test.input, err)
			}

			if got == nil {
				t.Fatalf("RenameType(%q): got nil, want %v", test.input, test.want)
			}

			if got.Identifier == nil {
				t.Fatalf("RenameType(%q).Identifier: got nil, want %v", test.input, test.want.Identifier)
			}

			if got.Identifier.Name != test.want.Identifier.Name {
				t.Errorf("RenameType(%q).Identifier.Name: got %v, want %v", test.input, got.Identifier.Name, test.want.Identifier.Name)
			}

			if test.want.Original != nil {
				if got.Original == nil {
					t.Fatalf("RenameType(%q).Original: got nil, want %v", test.input, test.want.Original)
				}
				if got.Original.Name != test.want.Original.Name {
					t.Errorf("RenameType(%q).Original.Name: got %v, want %v", test.input, got.Original.Name, test.want.Original.Name)
				}
			} else {
				if got.Original != nil {
					t.Fatalf("RenameType(%q).Original: got %v, want nil", test.input, got.Original)
				}
			}
		})
	}
}

func TestTypeReference(t *testing.T) {
	tests := []struct {
		input      string
		tokenTypes []tok.TokenType
		want       *ast.TypeReference
		wantErr    bool
	}{
		{
			input:      "Foo",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokEOF},
			want: ast.NewTypeReference(
				nil, ast.NewTypeIdentifier("Foo", nil, 0, 3),
				nil, 0, 3,
			),
			wantErr: false,
		},
		{
			input:      "foo.bar.Baz",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokDot, tok.TokID, tok.TokDot, tok.TokTypeID, tok.TokEOF},
			want: ast.NewTypeReference(
				[]*ast.Identifier{
					ast.NewIdentifier("foo", nil, 0, 3),
					ast.NewIdentifier("bar", nil, 0, 3),
				},
				ast.NewTypeIdentifier("Baz", nil, 0, 3),
				nil, 0, 0,
			),
			wantErr: false,
		},
		{
			input:      "Foo.Bar.Baz",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokDot, tok.TokTypeID, tok.TokDot, tok.TokTypeID, tok.TokEOF},
			want: ast.NewTypeReference(
				nil,
				ast.NewTypeIdentifier("Foo", nil, 0, 3),
				nil, 0, 3,
			),
			wantErr: false,
		},
		{
			input:      "Foo.Bar.y",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokDot, tok.TokTypeID, tok.TokDot, tok.TokID, tok.TokEOF},
			want:       ast.NewTypeReference(nil, ast.NewTypeIdentifier("Foo", nil, 0, 3), nil, 0, 0),
			wantErr:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}

			if !slices.Equal(tok.Types(tokens), test.tokenTypes) {
				t.Fatalf("tokenTypes: %v, want %v", tok.Types(tokens), test.tokenTypes)
			}

			got, _, err := TypeReference(tokens)

			if test.wantErr {
				if err == nil {
					t.Errorf("TypeReference(%q): want error", test.input)
				}
				return
			}

			if !test.wantErr && err != nil {
				t.Fatalf("TypeReference(%q): got error %v, want nil", test.input, err)
			}

			if got == nil {
				t.Fatalf("TypeReference(%q): got nil, want %v", test.input, test.want)
			}

			if got.TypeIdentifier == nil {
				t.Fatalf("TypeReference(%q).TypeIdentifier: got nil, want %v", test.input, test.want.TypeIdentifier)
			}

			if got.TypeIdentifier.Name != test.want.TypeIdentifier.Name {
				t.Errorf("TypeReference(%q).TypeIdentifier.Name: got %v, want %v", test.input, got.TypeIdentifier.Name, test.want.TypeIdentifier.Name)
			}

			if len(got.Identifiers) != len(test.want.Identifiers) {
				t.Fatalf("TypeReference(%q).Identifiers: got %v, want %v", test.input, got.Identifiers, test.want.Identifiers)
			}
			for i, identifier := range got.Identifiers {
				if identifier.Name != test.want.Identifiers[i].Name {
					t.Errorf("TypeReference(%q).Identifiers[%d].Name: got %v, want %v", test.input, i, identifier.Name, test.want.Identifiers[i].Name)
				}
			}
		})
	}
}

func TestFunctionIdentifier(t *testing.T) {
	tests := []struct {
		input      string
		tokenTypes []tok.TokenType
		want       *ast.FunctionIdentifier
		wantErr    bool
	}{
		{
			input:      "foo",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokEOF},
			want:       ast.NewFunctionIdentifier("foo", nil, 0, 3),
			wantErr:    false,
		},
		{
			input:      "foo?",
			tokenTypes: []tok.TokenType{tok.TokFuncID, tok.TokEOF},
			want:       ast.NewFunctionIdentifier("foo?", nil, 0, 4),
			wantErr:    false,
		},
		{
			input:      "foo!",
			tokenTypes: []tok.TokenType{tok.TokFuncID, tok.TokEOF},
			want:       ast.NewFunctionIdentifier("foo!", nil, 0, 4),
			wantErr:    false,
		},
		{
			input:      "Foo",
			tokenTypes: []tok.TokenType{tok.TokTypeID, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}

			if !slices.Equal(tok.Types(tokens), test.tokenTypes) {
				t.Fatalf("tokenTypes: %v, want %v", tok.Types(tokens), test.tokenTypes)
			}

			got, _, err := FunctionIdentifier(tokens)

			if test.wantErr {
				if err == nil {
					t.Errorf("FunctionIdentifier(%q): want error", test.input)
				}
				return
			}

			if !test.wantErr && err != nil {
				t.Fatalf("FunctionIdentifier(%q): got error %v, want nil", test.input, err)
			}

			if got == nil {
				t.Fatalf("FunctionIdentifier(%q): got nil, want %v", test.input, test.want)
			}

			if got.Name != test.want.Name {
				t.Errorf("FunctionIdentifier(%q).Name: got %v, want %v", test.input, got.Name, test.want.Name)
			}
		})
	}
}
