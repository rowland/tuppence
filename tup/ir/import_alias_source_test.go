package ir

import "testing"

func TestImportAliasSourceValidity(t *testing.T) {
	if !ImportAliasSourceUser.Valid() || !ImportAliasSourceProlog.Valid() {
		t.Fatal("known import alias sources must be valid")
	}
	if ImportAliasSource(99).Valid() {
		t.Fatal("unknown import alias source must be invalid")
	}
}
