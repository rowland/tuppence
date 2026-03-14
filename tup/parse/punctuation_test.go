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
		{input: "@"},                 // sees at
		{input: "@x"},                // sees at followed by identifier
		{input: "x@", wantErr: true}, // sees identifier followed by at
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			_, found := At(tokens)
			if test.wantErr && found {
				t.Errorf("At(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("At(%q) = %v, want nil", test.input, found)
			}
		})
	}
}

func TestColon(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{input: ":"},                 // sees colon
		{input: ": "},                // sees colon
		{input: ":\t"},               // sees colon
		{input: ":\r"},               // sees colon
		{input: ":\n"},               // sees colon
		{input: ":x"},                // sees colon followed by identifier
		{input: "x:", wantErr: true}, // sees identifier followed by colon
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			_, found := Colon(tokens)
			if test.wantErr && found {
				t.Errorf("Colon(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("Colon(%q) = %v, want nil", test.input, found)
			}
		})
	}
}

func TestDot(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{input: "."},                 // sees dot
		{input: ".x"},                // sees dot followed by identifier
		{input: "x.", wantErr: true}, // sees identifier followed by dot
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			_, found := Dot(tokens)
			if test.wantErr && found {
				t.Errorf("Dot(%q) = %v, want error", test.input, found)
			}
			if !test.wantErr && !found {
				t.Fatalf("Dot(%q) = %v, want nil", test.input, found)
			}
		})
	}
}
