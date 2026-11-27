package ir

import "testing"

func TestIndexedProgramLookupFallsThroughToPrologAndAllowsModuleShadowing(t *testing.T) {
	prolog := ScopeID(0x80000001)
	base := IndexedProgram{
		Scopes: ScopeTable{
			{ID: prolog, Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
			{ID: ScopeID(1), Parent: prolog, Module: ModuleID(1), Kind: ScopeModule, Initial: true, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: testOrigin()},
		},
		ImportAliases: ImportAliasTable{{Scope: prolog, LocalName: "take", Module: ModuleID(2), Declaration: DeclID(1), Origin: testOrigin()}},
		Declarations:  DeclarationTable{{ID: DeclID(1), Scope: ScopeID(2), Name: "take", Kind: DeclFunction, Exported: true, Origin: testOrigin()}},
		Functions:     FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Effect: EffectPure, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "take", Arguments: []CallArgument{{Type: TypeID(10)}}, Origin: testOrigin()}
	if matches := base.FindCallCandidates(call); len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("lookup must fall through from the module scope to the prolog scope")
	}

	shadowed := base
	shadowed.Declarations = append(shadowed.Declarations, Declaration{ID: DeclID(2), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()})
	shadowed.Functions = append(shadowed.Functions, Function{ID: FunctionID(2), Declaration: DeclID(2), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Effect: EffectPure, Origin: testOrigin()})
	if matches := shadowed.FindCallCandidates(call); len(matches) != 1 || matches[0] != FunctionID(2) {
		t.Fatal("module declarations must shadow prolog aliases")
	}
}
