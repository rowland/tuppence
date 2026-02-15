package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestAt(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{input: "@", wantErr: false},  // sees at
		{input: "@x", wantErr: false}, // sees at followed by identifier
		{input: "x@", wantErr: true},  // sees identifier followed by at
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			_, err = At(tokens)
			if test.wantErr && err == nil {
				t.Errorf("At(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("At(%q) = %v, want nil", test.input, err)
			}
		})
	}
}

func TestColon(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{input: ":", wantErr: false},   // sees colon
		{input: ": ", wantErr: false},  // sees colon
		{input: ":\t", wantErr: false}, // sees colon
		{input: ":\r", wantErr: false}, // sees colon
		{input: ":\n", wantErr: false}, // sees colon
		{input: ":x", wantErr: false},  // sees colon followed by identifier
		{input: "x:", wantErr: true},   // sees identifier followed by colon
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			_, err = Colon(tokens)
			if test.wantErr && err == nil {
				t.Errorf("Colon(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Colon(%q) = %v, want nil", test.input, err)
			}
		})
	}
}

func TestDot(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{input: ".", wantErr: false},  // sees dot
		{input: ".x", wantErr: false}, // sees dot followed by identifier
		{input: "x.", wantErr: true},  // sees identifier followed by dot
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			_, err = Dot(tokens)
			if test.wantErr && err == nil {
				t.Errorf("Dot(%q) = %v, want error", test.input, err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Dot(%q) = %v, want nil", test.input, err)
			}
		})
	}
}
