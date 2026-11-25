package ir

import "testing"

func TestFunctionTableFindsFunctionsByID(t *testing.T) {
	table := FunctionTable{
		{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()},
		{ID: FunctionID(2), Declaration: DeclID(3), Effect: Effectful, Generic: true, Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("complete function table must be valid")
	}
	function, ok := table.Find(FunctionID(2))
	if !ok || !function.Generic || function.Effect != Effectful {
		t.Fatal("function table must find the requested function")
	}
	if _, ok := table.Find(FunctionID(9)); ok {
		t.Fatal("function table must not find an absent function")
	}
}

func TestFunctionTableRejectsDuplicateDeclarations(t *testing.T) {
	table := FunctionTable{
		{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()},
		{ID: FunctionID(2), Declaration: DeclID(2), Effect: Effectful, Origin: testOrigin()},
	}
	if table.Valid() {
		t.Fatal("function table with duplicate declarations must be invalid")
	}
}
func TestFunctionTableRejectsInvalidOrDuplicateFunctions(t *testing.T) {
	invalid := FunctionTable{{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure}}
	if invalid.Valid() {
		t.Fatal("function table with an invalid function must be invalid")
	}

	duplicate := FunctionTable{
		{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()},
		{ID: FunctionID(1), Declaration: DeclID(3), Effect: Effectful, Origin: testOrigin()},
	}
	if duplicate.Valid() {
		t.Fatal("function table with duplicate IDs must be invalid")
	}
}
