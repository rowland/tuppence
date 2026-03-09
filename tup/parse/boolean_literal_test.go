package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.BooleanLiteral
		wantErr bool
	}{
		{"true", "true", ast.NewBooleanLiteral("true", true, nil, 0, 4), false},
		{"false", "false", ast.NewBooleanLiteral("false", false, nil, 0, 5), false},
		{"invalid", "invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunParseTest(t, tt.name, tt.input, tt.want, tt.wantErr, "BooleanLiteral", BooleanLiteral, StringerCheck[*ast.BooleanLiteral])
		})
	}
}
