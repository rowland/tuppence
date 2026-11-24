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

func TestIndexedProgramRejectsDanglingDeclarationScopes(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(9), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("declaration referring to an absent scope must be invalid")
	}
}

func TestIndexedProgramRejectsDanglingBindingScopes(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Bindings: BindingTable{
			{ID: BindingID(1), Scope: ScopeID(9), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("binding referring to an absent scope must be invalid")
	}
}

func TestIndexedProgramRejectsDanglingBindingVersions(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule},
		},
		Bindings: BindingTable{
			{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", Previous: BindingID(9), TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("binding referring to an absent previous version must be invalid")
	}
}

func TestIndexedProgramRejectsInconsistentBindingVersions(t *testing.T) {
	base := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule},
			{ID: ScopeID(2), Module: ModuleID(1), Kind: ScopeBlock},
		},
	}
	for name, current := range map[string]Binding{
		"different name":  {ID: BindingID(2), Scope: ScopeID(1), Name: "other", Previous: BindingID(1), TypeSlot: TypeVarID(1), Value: ExprID(2), Origin: testOrigin()},
		"different scope": {ID: BindingID(2), Scope: ScopeID(2), Name: "answer", Previous: BindingID(1), TypeSlot: TypeVarID(1), Value: ExprID(2), Origin: testOrigin()},
	} {
		program := base
		program.Bindings = BindingTable{
			{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: testOrigin()},
			current,
		}
		if program.Valid() {
			t.Errorf("%s binding version must be invalid", name)
		}
	}
}

func TestIndexedProgramRejectsDanglingScopeParents(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Parent: ScopeID(9), Module: ModuleID(1), Kind: ScopeBlock},
		},
	}
	if program.Valid() {
		t.Fatal("scope referring to an absent parent must be invalid")
	}
}

func TestIndexedProgramRejectsCrossModuleScopeParents(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
			{ID: ModuleID(2), Name: "lib", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule},
			{ID: ScopeID(2), Parent: ScopeID(1), Module: ModuleID(2), Kind: ScopeBlock},
		},
	}
	if program.Valid() {
		t.Fatal("scope with a parent from another module must be invalid")
	}
}

func TestIndexedProgramRejectsScopeParentCycles(t *testing.T) {
	base := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
	}
	for name, scopes := range map[string]ScopeTable{
		"self cycle": {
			{ID: ScopeID(1), Parent: ScopeID(1), Module: ModuleID(1), Kind: ScopeBlock},
		},
		"two-scope cycle": {
			{ID: ScopeID(1), Parent: ScopeID(2), Module: ModuleID(1), Kind: ScopeBlock},
			{ID: ScopeID(2), Parent: ScopeID(1), Module: ModuleID(1), Kind: ScopeBlock},
		},
	} {
		program := base
		program.Scopes = scopes
		if program.Valid() {
			t.Errorf("%s must be invalid", name)
		}
	}
}

func TestIndexedProgramRejectsDanglingScopeEntries(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Entries: []ScopeEntry{{Name: "answer", Declaration: DeclID(9)}}},
		},
	}
	if program.Valid() {
		t.Fatal("scope entry referring to an absent declaration must be invalid")
	}
}

func TestIndexedProgramRejectsMisownedScopeEntries(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Entries: []ScopeEntry{{Name: "answer", Declaration: DeclID(1)}}},
			{ID: ScopeID(2), Module: ModuleID(1), Kind: ScopeBlock},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(2), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("scope entry owned by another scope must be invalid")
	}
}
