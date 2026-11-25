package ir

import "testing"

func TestCallExpressionValidity(t *testing.T) {
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(2), Name: "parse", Arguments: []TypeID{TypeID(3)}, Origin: testOrigin()}
	if !call.Valid() {
		t.Fatal("complete call expression must be valid")
	}
}

func TestCallExpressionRejectsMissingName(t *testing.T) {
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(2), Origin: testOrigin()}
	if call.Valid() {
		t.Fatal("call expression without a name must be invalid")
	}
}
