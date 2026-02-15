package parse

import (
	"slices"
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestAnnotation(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenTypes []tok.TokenType
		want       ast.Annotation
		wantErr    bool
	}{
		{
			name:       "simple",
			input:      "@x\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokEOL, tok.TokEOF},
			want:       ast.NewSimpleAnnotation("x"),
			wantErr:    false,
		},
		{
			name:       "namespaced with no value",
			input:      "@x:y\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokEOL, tok.TokEOF},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "namespaced with string value",
			input:      "@x:y \"z\"\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokStrLit, tok.TokEOL, tok.TokEOF},
			want:       ast.NewNamespacedAnnotation("x", "y", ast.NewStringLiteral(`"z"`, "z", nil, 0, 4)),
			wantErr:    false,
		},
		{
			name:       "namespaced with integer value",
			input:      "@x:y 1\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokDecLit, tok.TokEOL, tok.TokEOF},
			want:       ast.NewNamespacedAnnotation("x", "y", ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
			wantErr:    false,
		},
		{
			name:       "namespaced with float value",
			input:      "@x:y 1.0\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokFloatLit, tok.TokEOL, tok.TokEOF},
			want:       ast.NewNamespacedAnnotation("x", "y", ast.NewFloatLiteral("1.0", 1.0, nil, 0, 4)),
			wantErr:    false,
		},
		{
			name:       "namespaced with true value",
			input:      "@x:y true\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokBoolLit, tok.TokEOL, tok.TokEOF},
			want:       ast.NewNamespacedAnnotation("x", "y", ast.NewBooleanLiteral("true", true, nil, 0, 4)),
			wantErr:    false,
		},
		{
			name:       "namespaced with false value",
			input:      "@x:y false\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokBoolLit, tok.TokEOL, tok.TokEOF},
			want:       ast.NewNamespacedAnnotation("x", "y", ast.NewBooleanLiteral("false", false, nil, 0, 5)),
			wantErr:    false,
		},
		{
			name:       "namespaced with unqualified type reference value",
			input:      "@x:y Foo\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokTypeID, tok.TokEOL, tok.TokEOF},
			want:       ast.NewNamespacedAnnotation("x", "y", ast.NewTypeReference(nil, ast.NewTypeIdentifier("Foo", nil, 0, 3), nil, 0, 3)),
			wantErr:    false,
		},
		{
			name:       "namespaced with qualified type reference value",
			input:      "@x:y foo.Bar\n",
			tokenTypes: []tok.TokenType{tok.TokAt, tok.TokID, tok.TokColonNoSpace, tok.TokID, tok.TokID, tok.TokDot, tok.TokTypeID, tok.TokEOL, tok.TokEOF},
			want: ast.NewNamespacedAnnotation("x", "y", ast.NewTypeReference(
				[]*ast.Identifier{ast.NewIdentifier("foo", nil, 0, 3)},
				ast.NewTypeIdentifier("Bar", nil, 0, 3),
				nil, 0, 0)),
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
				return
			}
			if !slices.Equal(tok.Types(tokens), test.tokenTypes) {
				t.Fatalf("tokenTypes: %v, want %v", tok.Types(tokens), test.tokenTypes)
			}

			got, _, err := Annotation(tokens)

			if test.wantErr {
				if err == nil {
					t.Errorf("Annotation(%q) = %v, want error", test.input, err)
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Annotation(%q) = %v, want nil", test.input, err)
			}

			switch want := test.want.(type) {
			case *ast.SimpleAnnotation:
				sa, ok := got.(*ast.SimpleAnnotation)
				if !ok {
					t.Errorf("Annotation(%q) = %T, want %T", test.input, got, test.want)
					return
				}
				if sa.Identifier != want.Identifier {
					t.Errorf("Annotation(%q).Identifier: got %v, want %v", test.input, sa.Identifier, want.Identifier)
				}
			case *ast.NamespacedAnnotation:
				na, ok := got.(*ast.NamespacedAnnotation)
				if !ok {
					t.Errorf("Annotation(%q) = %T, want %T", test.input, got, test.want)
					return
				}
				if na.Namespace != want.Namespace {
					t.Errorf("Annotation(%q).Namespace: got %v, want %v", test.input, na.Namespace, want.Namespace)
				}
				if na.Identifier != want.Identifier {
					t.Errorf("Annotation(%q).Identifier: got %v, want %v", test.input, na.Identifier, want.Identifier)
				}
				if na.Value.String() != want.Value.String() {
					t.Errorf("Annotation(%q).Value: got %v, want %v", test.input, na.Value.String(), want.Value.String())
				}
			case nil:
				if got != nil {
					t.Errorf("Annotation(%q) = %v, want nil", test.input, got)
				}
			default:
				t.Errorf("Annotation(%q) = %T, want %T", test.input, got, test.want)
			}

		})
	}
}
