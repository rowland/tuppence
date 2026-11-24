package ir

import "testing"

func TestIndexedProgramValidity(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
		Bindings: BindingTable{
			{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: testOrigin()},
		},
		Diagnostics: []Diagnostic{
			{Severity: SeverityWarning, Code: "IR001", Message: "example warning", Origin: testOrigin()},
		},
	}
	if !program.Valid() {
		t.Fatal("complete indexed program must be valid")
	}
}

func TestIndexedProgramRejectsInvalidTables(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main"}},
	}
	if program.Valid() {
		t.Fatal("indexed program with an invalid module table must be invalid")
	}

	program = IndexedProgram{
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1)},
		},
	}
	if program.Valid() {
		t.Fatal("indexed program with an invalid table must be invalid")
	}
}

func TestIndexedProgramRejectsInvalidDiagnostics(t *testing.T) {
	program := IndexedProgram{
		Diagnostics: []Diagnostic{{Severity: SeverityError, Code: "IR001"}},
	}
	if program.Valid() {
		t.Fatal("indexed program with an invalid diagnostic must be invalid")
	}
}

func TestIndexedProgramRejectsDanglingScopeModules(t *testing.T) {
	program := IndexedProgram{
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(9), Kind: ScopeModule},
		},
	}
	if program.Valid() {
		t.Fatal("scope referring to an absent module must be invalid")
	}
}
