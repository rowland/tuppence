package parse

import (
	"slices"
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestSize(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenTypes []tok.TokenType
		want       ast.Size
		wantErr    bool
	}{
		{
			name:       "simple",
			input:      "1",
			tokenTypes: []tok.TokenType{tok.TokDecLit, tok.TokEOF},
			want:       ast.NewDecimalLiteral("1", 1, nil, 0, 1),
			wantErr:    false,
		},
		{
			name:       "identifier",
			input:      "x",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokEOF},
			want:       ast.NewIdentifier("x", nil, 0, 1),
			wantErr:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			if !slices.Equal(tok.Types(tokens), test.tokenTypes) {
				t.Fatalf("tokenTypes: %v, want %v", tok.Types(tokens), test.tokenTypes)
			}

			got, _, err := Size(tokens)

			if test.wantErr {
				if err == nil {
					t.Errorf("Size(%q): want error", test.input)
				}
				return
			}

			if !test.wantErr && err != nil {
				t.Errorf("Size(%q): got error %v, want nil", test.input, err)
			}

			if err != nil {
				t.Errorf("Size(%q) = %v", test.input, err)
			}

			switch got := got.(type) {
			case *ast.IntegerLiteral:
				want, ok := test.want.(*ast.IntegerLiteral)
				if !ok {
					t.Fatalf("Size(%q).want: got %T, want *ast.IntegerLiteral", test.input, test.want)
				}
				if got.Value != want.Value {
					t.Fatalf("Size(%q).Value: got %v, want %v", test.input, got.Value, want.Value)
				}
				if got.IntegerValue != want.IntegerValue {
					t.Fatalf("Size(%q).IntegerValue: got %v, want %v", test.input, got.IntegerValue, want.IntegerValue)
				}
				if got.Base != want.Base {
					t.Fatalf("Size(%q).Base: got %v, want %v", test.input, got.Base, want.Base)
				}
			case *ast.Identifier:
				want, ok := test.want.(*ast.Identifier)
				if !ok {
					t.Fatalf("Size(%q).want: got %T, want *ast.Identifier", test.input, test.want)
				}
				if got.Name != want.Name {
					t.Fatalf("Size(%q).Name: got %v, want %v", test.input, got.Name, want.Name)
				}
			default:
				t.Errorf("Size(%q): got %T, want *ast.IntegerLiteral or *ast.Identifier", test.input, got)
			}
		})
	}
}
