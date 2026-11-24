package ir

import "testing"

func TestExpressionTableFindsReferencesByID(t *testing.T) {
	table := ExpressionTable{
		{ID: ExprID(1), Declaration: DeclID(2), Origin: testOrigin()},
		{ID: ExprID(2), Declaration: DeclID(3), Function: FunctionID(4), Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("complete expression table must be valid")
	}
	expression, ok := table.Find(ExprID(2))
	if !ok || expression.Function != FunctionID(4) {
		t.Fatal("expression table must find the requested reference")
	}
	if _, ok := table.Find(ExprID(9)); ok {
		t.Fatal("expression table must not find an absent expression")
	}
}

func TestExpressionTableRejectsDuplicateIDs(t *testing.T) {
	table := ExpressionTable{
		{ID: ExprID(1), Declaration: DeclID(2), Origin: testOrigin()},
		{ID: ExprID(1), Declaration: DeclID(3), Origin: testOrigin()},
	}
	if table.Valid() {
		t.Fatal("expression table with duplicate IDs must be invalid")
	}
}
