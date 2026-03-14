package parse

import (
	"fmt"
	"reflect"
	"slices"
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

// StringerCheck is a check function for RunParseTest that compares got and want
// using String(). Use it when the parsed type implements fmt.Stringer.
func StringerCheck[T fmt.Stringer](t *testing.T, input, parserName string, got, want T) {
	t.Helper()
	if isNil(want) {
		return
	}
	if isNil(got) {
		t.Fatalf("%s(%q) = nil, want non-nil", parserName, input)
	}
	if got.String() != want.String() {
		t.Errorf("%s(%q) = %v, want %v", parserName, input, got.String(), want.String())
	}
}

// RunParseTest runs tokenize, parse, and wantErr handling, then calls check
// to compare results. check is only called when !wantErr and parse succeeded.
// check receives concrete got/want so tests stay explicit and debuggable.
func RunParseTest[T any](
	t *testing.T,
	subtestName string,
	input string,
	want T,
	wantErr bool,
	parserName string,
	parse func([]tok.Token) (T, []tok.Token, error),
	check func(t *testing.T, input, parserName string, got, want T),
) {
	t.Helper()
	RunParseTestExt(t, subtestName, input, want, wantErr, parserName, parse, check, nil)
}

func RunParseTestExt[T any](
	t *testing.T,
	subtestName string,
	input string,
	want T,
	wantErr bool,
	parserName string,
	parse func([]tok.Token) (T, []tok.Token, error),
	check func(t *testing.T, input, parserName string, got, want T),
	tokenTypes []tok.TokenType,
) {
	t.Helper()
	src := source.NewSource([]byte(input), "test.tup")
	tokens, err := tok.Tokenize(src.Contents, src.Filename)
	if err != nil {
		t.Errorf("Tokenize(%q) = %v", input, err)
		return
	}
	if tokenTypes != nil && !slices.Equal(tok.Types(tokens), tokenTypes) {
		t.Fatalf("Tokens: %v, want %v", tok.Types(tokens), tokenTypes)
	}
	got, _, err := parse(tokens)
	if wantErr {
		if err == nil {
			t.Errorf("%s(%q) = %v, want error", parserName, input, got)
		}
		return
	}
	if err != nil {
		t.Fatalf("%s(%q) = %v", parserName, input, err)
	}
	if isNil(want) && !isNil(got) {
		t.Fatalf("%s(%q) = %v, want nil", parserName, input, got)
	}
	check(t, input, parserName, got, want)
}
