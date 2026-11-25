package ir

import "testing"

func TestCallExpressionValidity(t *testing.T) {
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(2), Name: "parse", Arguments: []CallArgument{{Type: TypeID(3)}}, Origin: testOrigin()}
	if !call.Valid() {
		t.Fatal("complete call expression must be valid")
	}
}

func TestCallExpressionAllowsPositionalPrefixAndNamedSuffix(t *testing.T) {
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(2), Name: "parse", Arguments: []CallArgument{{Type: TypeID(3)}, {Name: "base", Type: TypeID(4)}}, Curry: true, Origin: testOrigin()}
	if !call.Valid() {
		t.Fatal("call with positional prefix, named suffix, and curry marker must be valid")
	}
}

func TestCallExpressionRejectsPositionalAfterNamed(t *testing.T) {
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(2), Name: "parse", Arguments: []CallArgument{{Name: "base", Type: TypeID(3)}, {Type: TypeID(4)}}, Origin: testOrigin()}
	if call.Valid() {
		t.Fatal("positional argument after named argument must be invalid")
	}
}

func TestCallExpressionRejectsDuplicateNamedArguments(t *testing.T) {
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(2), Name: "parse", Arguments: []CallArgument{{Name: "base", Type: TypeID(3)}, {Name: "base", Type: TypeID(4)}}, Origin: testOrigin()}
	if call.Valid() {
		t.Fatal("duplicate named arguments must be invalid")
	}
}
func TestCallExpressionAllowsCalleeStage(t *testing.T) {
	call := CallExpr{ID: ExprID(2), Scope: ScopeID(2), Callee: ExprID(1), Arguments: []CallArgument{{Type: TypeID(3)}}, Origin: testOrigin()}
	if !call.Valid() {
		t.Fatal("call stage with a callee expression must be valid")
	}
}
func TestCallExpressionRejectsMissingName(t *testing.T) {
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(2), Origin: testOrigin()}
	if call.Valid() {
		t.Fatal("call expression without a name must be invalid")
	}
}
