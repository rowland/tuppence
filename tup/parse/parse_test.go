package parse

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

// isNil reports whether v is nil. It only checks types that can be nil (pointer,
// slice, map, chan, func, interface); value types always return false.
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}

// RunParseTest runs a tokenize-and-parse test. The parse function consumes tokens
// and returns (result, remainder, error). Results are compared via String().
func RunParseTest[T fmt.Stringer](
	t *testing.T,
	input string,
	want T,
	wantErr bool,
	parse func([]tok.Token) (T, []tok.Token, error)) {
	t.Helper()
	src := source.NewSource([]byte(input), "test.tup")
	tokens, err := tok.Tokenize(src.Contents, src.Filename)
	if err != nil {
		t.Errorf("Tokenize(%q) = %v", input, err)
		return
	}
	got, _, err := parse(tokens)
	if wantErr {
		if err == nil {
			t.Fatalf("parse() = nil, want error")
		}
		return
	}
	if err != nil {
		t.Fatalf("parse() error = %v, want nil", err)
	}
	if isNil(got) {
		t.Fatalf("parse() = nil, want non-nil")
	}
	if got.String() != want.String() {
		t.Errorf("parse() = %v, want %v", got.String(), want.String())
	}
}
