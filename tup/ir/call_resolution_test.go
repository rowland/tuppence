package ir

import "testing"

func TestIndexedProgramResolvesLiteralConstructor(t *testing.T) {
	program := IndexedProgram{
		Constructors: TypeConstructorTable{
			{Type: TypeID(20), Storage: StorageF32, Origin: testOrigin()},
			{Type: TypeID(20), Storage: StorageI32, Origin: testOrigin()},
		},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(20)}}, Effect: EffectPure, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "take", Arguments: []CallArgument{{Literal: LiteralInteger, LiteralTemplate: true}}, Origin: testOrigin()}
	plan, ok := program.ResolveCallPlan(call, FunctionID(1))
	if !ok || len(plan.Arguments) != 1 || plan.Arguments[0].Adaptation != AdaptLiteral || plan.Arguments[0].Storage != StorageI32 {
		t.Fatal("literal call plan must select the first admissible storage constructor")
	}
}

func TestIndexedProgramRejectsLiteralWithoutConstructor(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(20)}}, Effect: EffectPure, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "take", Arguments: []CallArgument{{Literal: LiteralInteger}}, Origin: testOrigin()}
	if _, ok := program.ResolveCallPlan(call, FunctionID(1)); ok {
		t.Fatal("literal call plan must reject a parameter without a storage constructor")
	}
}
