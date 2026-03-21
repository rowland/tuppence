package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestRangeBound(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.RangeBound
		wantErr bool
	}{
		{
			name:  "identifier bound",
			input: "start",
			want: ast.NewRangeBound(
				ast.NewIdentifier("start", nil, 0, 5),
			),
		},
		{
			name:  "postfix bound",
			input: "foo.bar[1]",
			want: ast.NewRangeBound(
				ast.NewIndexedAccess(
					ast.NewMemberAccess(
						ast.NewIdentifier("foo", nil, 0, 3),
						ast.NewIdentifier("bar", nil, 0, 3),
					),
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"RangeBound", RangeBound, StringerCheck[*ast.RangeBound])
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Range
		wantErr bool
	}{
		{
			name:  "simple range",
			input: "1..10",
			want: ast.NewRange(
				ast.NewRangeBound(ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				ast.NewRangeBound(ast.NewDecimalLiteral("10", 10, nil, 0, 2)),
			),
		},
		{
			name:  "range with postfix bounds",
			input: "foo.bar..foo.qux",
			want: ast.NewRange(
				ast.NewRangeBound(
					ast.NewMemberAccess(
						ast.NewIdentifier("foo", nil, 0, 3),
						ast.NewIdentifier("bar", nil, 0, 3),
					),
				),
				ast.NewRangeBound(
					ast.NewMemberAccess(
						ast.NewIdentifier("foo", nil, 0, 3),
						ast.NewIdentifier("qux", nil, 0, 3),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Range", Range, StringerCheck[*ast.Range])
		})
	}
}
