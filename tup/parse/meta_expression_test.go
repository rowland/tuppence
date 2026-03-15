package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestMetaExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{
			name:  "single key",
			input: `$(file: "config.json")`,
			want: map[string]string{
				"file": `"config.json"`,
			},
		},
		{
			name:  "multiple keys with trailing comma",
			input: `$(hash: "hello", algorithm: "sha256",)`,
			want: map[string]string{
				"hash":      `"hello"`,
				"algorithm": `"sha256"`,
			},
		},
		{
			name:    "missing labeled arguments",
			input:   `$()`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(src.Contents, src.Filename)
			if err != nil {
				t.Fatalf("Tokenize(%q) = %v", test.input, err)
			}

			got, _, err := MetaExpression(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("MetaExpression(%q): want error, got nil", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("MetaExpression(%q): got error %v, want nil", test.input, err)
			}
			if len(got.KeyValues) != len(test.want) {
				t.Fatalf("MetaExpression(%q): got %d keys, want %d", test.input, len(got.KeyValues), len(test.want))
			}
			for key, wantValue := range test.want {
				value, ok := got.KeyValues[key]
				if !ok {
					t.Fatalf("MetaExpression(%q): missing key %q", test.input, key)
				}
				if value.String() != wantValue {
					t.Fatalf("MetaExpression(%q): key %q = %q, want %q", test.input, key, value.String(), wantValue)
				}
			}
		})
	}
}
