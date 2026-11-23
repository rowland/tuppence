package ir

import "testing"

func TestScopeEntryValidity(t *testing.T) {
	entry := ScopeEntry{
		Name:        "answer",
		Declaration: DeclID(3),
		Offset:      7,
	}
	if !entry.Valid() {
		t.Fatal("complete scope entry must be valid")
	}
	if (ScopeEntry{Name: "answer", Declaration: DeclID(3)}).Valid() != true {
		t.Fatal("zero source offset should be valid")
	}
}

func TestScopeEntryRejectsMissingNameOrDeclaration(t *testing.T) {
	for name, entry := range map[string]ScopeEntry{
		"missing name":        {Declaration: DeclID(1)},
		"missing declaration": {Name: "answer"},
	} {
		if entry.Valid() {
			t.Errorf("%s scope entry must be invalid", name)
		}
	}
}
