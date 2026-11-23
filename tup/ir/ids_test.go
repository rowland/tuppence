package ir

import "testing"

func TestIDsReserveZeroAsInvalid(t *testing.T) {
	if FileID(0).Valid() {
		t.Fatal("zero FileID must be invalid")
	}
	if ScopeID(0).Valid() {
		t.Fatal("zero ScopeID must be invalid")
	}
	if DeclID(0).Valid() {
		t.Fatal("zero DeclID must be invalid")
	}
	if ExprID(0).Valid() {
		t.Fatal("zero ExprID must be invalid")
	}
	if TypeID(0).Valid() {
		t.Fatal("zero TypeID must be invalid")
	}
	if FunctionID(0).Valid() {
		t.Fatal("zero FunctionID must be invalid")
	}
}

func TestNonZeroIDsAreValid(t *testing.T) {
	if !FileID(1).Valid() {
		t.Fatal("non-zero FileID must be valid")
	}
	if !OriginID(7).Valid() {
		t.Fatal("non-zero OriginID must be valid")
	}
}
