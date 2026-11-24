package ir

import "testing"

func TestIndexedProgramRejectsDanglingFunctionReferences(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "compute", Kind: DeclFunction, Origin: testOrigin()}},
		Expressions:  ExpressionTable{{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(1), Function: FunctionID(9), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("reference to an absent function must be invalid")
	}
}

func TestIndexedProgramRejectsMismatchedFunctionReferences(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main", Origin: testOrigin()}},
		Scopes:  ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule}},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "first", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "second", Kind: DeclFunction, Origin: testOrigin()},
		},
		Functions:   FunctionTable{{ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()}},
		Expressions: ExpressionTable{{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(1), Function: FunctionID(1), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("reference whose function declaration differs must be invalid")
	}
}

func TestIndexedProgramRejectsFunctionReferenceToNonFunction(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Effect: EffectPure, Origin: testOrigin()}},
		Expressions:  ExpressionTable{{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(1), Function: FunctionID(1), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("function reference to a non-function declaration must be invalid")
	}
}
