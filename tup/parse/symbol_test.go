package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestSymbolLiteral(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       *ast.SymbolLiteral
		wantErr    bool
		tokenTypes []tok.TokenType
	}{
		{
			name:       "symbol literal",
			input:      ":hello",
			want:       ast.NewSymbolLiteral(":hello", nil, 0, 6),
			tokenTypes: []tok.TokenType{tok.TokColonNoSpace, tok.TokID, tok.TokEOF},
		},
		{
			name:       "symbol literal with spaced colon does not parse",
			input:      ": hello",
			want:       nil,
			wantErr:    true,
			tokenTypes: []tok.TokenType{tok.TokColon, tok.TokID, tok.TokEOF},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTestExt(t, test.name, test.input, test.want, test.wantErr,
				"SymbolLiteral", SymbolLiteral, StringerCheck[*ast.SymbolLiteral], test.tokenTypes)
		})
	}
}
