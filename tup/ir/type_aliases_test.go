package ir

import "testing"

func TestTypeAliasValidity(t *testing.T) {
	alias := TypeAlias{
		Declaration: DeclID(1),
		Target:      TypeID(2),
		Origin:      testOrigin(),
	}
	if !alias.Valid() {
		t.Fatal("complete type alias must be valid")
	}
}

func TestTypeAliasRejectsMissingFields(t *testing.T) {
	for name, alias := range map[string]TypeAlias{
		"missing declaration": {Target: TypeID(2), Origin: testOrigin()},
		"missing target":      {Declaration: DeclID(1), Origin: testOrigin()},
		"missing origin":      {Declaration: DeclID(1), Target: TypeID(2)},
	} {
		if alias.Valid() {
			t.Errorf("%s alias must be invalid", name)
		}
	}
}
