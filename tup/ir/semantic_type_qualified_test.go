package ir

import "testing"

func TestIndexedProgramAllowsSameScopeTypeQualifiedDeclaration(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:  ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "Foo", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "default", Kind: DeclFunction, OwnerType: DeclID(1), Origin: testOrigin()},
		},
		Functions: FunctionTable{{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()}},
	}
	if !program.Valid() {
		t.Fatal("same-scope type-qualified declaration must be valid")
	}
}

func TestIndexedProgramRejectsMisownedTypeQualifiedDeclaration(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(1), Kind: ScopeFunction, Origin: testOrigin()},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "Foo", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(2), Name: "default", Kind: DeclFunction, OwnerType: DeclID(1), Origin: testOrigin()},
		},
		Functions: FunctionTable{{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("type-qualified declaration from another scope must be invalid")
	}
}

func TestIndexedProgramRejectsNonTypeQualifiedOwner(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:  ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "value", Kind: DeclValue, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "default", Kind: DeclFunction, OwnerType: DeclID(1), Origin: testOrigin()},
		},
		Functions: FunctionTable{{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("type-qualified declaration with a non-type owner must be invalid")
	}
}
