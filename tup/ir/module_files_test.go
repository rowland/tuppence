package ir

import "testing"

func TestParseModuleFileName(t *testing.T) {
	cases := []struct {
		filename    string
		module      string
		constraints []string
	}{
		{"math.tup", "math", nil},
		{"math.vectors.tup", "math", nil},
		{"math.vectors--amd64.tup", "math", []string{"amd64"}},
		{"filesystem--linux.tup", "filesystem", []string{"linux"}},
		{"inference--darwin-arm64.tup", "inference", []string{"darwin", "arm64"}},
	}
	for _, test := range cases {
		got, ok := ParseModuleFileName(test.filename)
		if !ok || got.ModuleName != test.module {
			t.Fatalf("ParseModuleFileName(%q) = %#v, %t", test.filename, got, ok)
		}
		if len(got.Constraints) != len(test.constraints) {
			t.Fatalf("ParseModuleFileName(%q) constraints = %v", test.filename, got.Constraints)
		}
		for i := range test.constraints {
			if got.Constraints[i] != test.constraints[i] {
				t.Fatalf("ParseModuleFileName(%q) constraints = %v", test.filename, got.Constraints)
			}
		}
	}
}

func TestModuleFileNameMatchesTarget(t *testing.T) {
	file, ok := ParseModuleFileName("inference--darwin-arm64.tup")
	if !ok || !file.Matches("darwin", "arm64") || file.Matches("linux", "arm64") {
		t.Fatal("OS and architecture suffixes must constrain file selection")
	}
}
