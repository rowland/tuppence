package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestArrayFunctionCall(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ArrayFunctionCall
		wantErr bool
	}{
		{
			name:  "type only",
			input: "array(Int)",
			want: ast.NewArrayFunctionCall(
				ast.NewTypeIdentifier("Int", nil, 0, 3),
				nil,
			),
		},
		{
			name:  "type and size",
			input: "array(Int, 10)",
			want: ast.NewArrayFunctionCall(
				ast.NewTypeIdentifier("Int", nil, 0, 3),
				ast.NewDecimalLiteral("10", 10, nil, 0, 2),
			),
		},
		{
			name:    "missing type",
			input:   "array()",
			wantErr: true,
		},
		{
			name:    "lowercase type is invalid",
			input:   "array(int)",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		RunParseTest(
			t,
			tt.name,
			tt.input,
			tt.want,
			tt.wantErr,
			"ArrayFunctionCall",
			ArrayFunctionCall,
			StringerCheck[*ast.ArrayFunctionCall],
		)
	}
}
