package ir

import "testing"

func TestIndexedProgramValidity(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Types: TypeTable{
			{ID: TypeID(1), Name: "Answer", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "Answer", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "compute", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(3), Scope: ScopeID(1), Name: "Alias", Kind: DeclType, Origin: testOrigin()},
		},
		Aliases: TypeAliasTable{
			{Declaration: DeclID(3), Target: TypeID(1), Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()},
		},
		Bindings: BindingTable{
			{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: testOrigin()},
		},
		Expressions: ExpressionTable{
			{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(1), Origin: testOrigin()},
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
			{ID: ScopeID(1), Module: ModuleID(9), Kind: ScopeModule, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("scope referring to an absent module must be invalid")
	}
}

func TestIndexedProgramRejectsDanglingDeclarationScopes(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
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
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
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
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
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
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()},
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
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Parent: ScopeID(9), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("scope referring to an absent parent must be invalid")
	}
}

func TestIndexedProgramRejectsCrossModuleScopeParents(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "lib", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
			{ID: ScopeID(2), Parent: ScopeID(1), Module: ModuleID(2), Kind: ScopeBlock, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("scope with a parent from another module must be invalid")
	}
}

func TestIndexedProgramRejectsScopeParentCycles(t *testing.T) {
	base := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
	}
	for name, scopes := range map[string]ScopeTable{
		"self cycle": {
			{ID: ScopeID(1), Parent: ScopeID(1), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()},
		},
		"two-scope cycle": {
			{ID: ScopeID(1), Parent: ScopeID(2), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()},
			{ID: ScopeID(2), Parent: ScopeID(1), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()},
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
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin(), Entries: []ScopeEntry{{Name: "answer", Declaration: DeclID(9)}}},
		},
	}
	if program.Valid() {
		t.Fatal("scope entry referring to an absent declaration must be invalid")
	}
}

func TestIndexedProgramRejectsMisownedScopeEntries(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin(), Entries: []ScopeEntry{{Name: "answer", Declaration: DeclID(1)}}},
			{ID: ScopeID(2), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(2), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("scope entry owned by another scope must be invalid")
	}
}

func TestIndexedProgramRejectsMismatchedScopeEntryNames(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin(), Entries: []ScopeEntry{{Name: "other", Declaration: DeclID(1)}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("scope entry with a mismatched name must be invalid")
	}
}

func TestIndexedProgramAllowsOverloadedFunctionNames(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin(), Entries: []ScopeEntry{
				{Name: "convert", Declaration: DeclID(1)},
				{Name: "convert", Declaration: DeclID(2)},
			}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "convert", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "convert", Kind: DeclFunction, Origin: testOrigin()},
		},
	}
	if !program.Valid() {
		t.Fatal("overloaded function names must be valid")
	}
}

func TestIndexedProgramRejectsDuplicateNonFunctionNames(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin(), Entries: []ScopeEntry{
				{Name: "answer", Declaration: DeclID(1)},
				{Name: "answer", Declaration: DeclID(2)},
			}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("duplicate non-function names must be invalid")
	}
}
