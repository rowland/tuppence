package ir

import "testing"

func TestReferenceExpressionValidity(t *testing.T) {
	cases := []ReferenceExpr{
		{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(2), Origin: testOrigin()},
		{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(2), Function: FunctionID(3), Origin: testOrigin()},
		{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(2), Value: ValueID(4), Origin: testOrigin()},
	}
	for _, reference := range cases {
		if !reference.Valid() {
			t.Errorf("reference expression %#v must be valid", reference)
		}
	}
}

func TestReferenceExpressionRejectsConflictingTargets(t *testing.T) {
	for name, reference := range map[string]ReferenceExpr{
		"missing id":          {Scope: ScopeID(1), Declaration: DeclID(2), Origin: testOrigin()},
		"missing scope":       {ID: ExprID(1), Declaration: DeclID(2), Origin: testOrigin()},
		"missing declaration": {ID: ExprID(1), Scope: ScopeID(1), Origin: testOrigin()},
		"missing origin":      {ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(2)},
		"function and value":  {ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(2), Function: FunctionID(3), Value: ValueID(4), Origin: testOrigin()},
	} {
		if reference.Valid() {
			t.Errorf("%s reference expression must be invalid", name)
		}
	}
}
