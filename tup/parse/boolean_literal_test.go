package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType ast.NodeType
		wantErr  bool
	}{
		{"true", "true", ast.NodeBooleanLiteral, false},
		{"false", "false", ast.NodeBooleanLiteral, false},
		{"invalid", "invalid", ast.NodeBooleanLiteral, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tok.Tokenize([]byte(tt.input), "test.tup")
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", tt.input, err)
			}
			literal, _, err := BooleanLiteral(tokens)
			if err != nil && !tt.wantErr {
				t.Errorf("BooleanLiteral(%q) = %v, want nil", tt.input, err)
			}
			if err == nil && literal == nil {
				t.Errorf("BooleanLiteral(%q) = nil, want not nil", tt.input)
			}
			if err == nil && literal.NodeType().String() != ast.NodeTypes[ast.NodeType(tt.wantType)] {
				t.Errorf("BooleanLiteral(%q) = %v, want %v", tt.input, literal.NodeType().String(), ast.NodeTypes[ast.NodeType(tt.wantType)])
			}
		})
	}
}
