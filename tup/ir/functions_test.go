package ir

import "testing"

func TestFunctionPreservesSignature(t *testing.T) {
	function := Function{ID: FunctionID(1), Declaration: DeclID(2), Parameters: []FunctionParameter{{Name: "left", Type: TypeID(3)}, {Name: "right", Type: TypeID(4)}}, Result: TypeID(5), Effect: EffectPure, Origin: testOrigin()}
	if !function.Valid() || len(function.Parameters) != 2 || function.Result != TypeID(5) {
		t.Fatal("function signature must be preserved")
	}
}
func TestFunctionValidity(t *testing.T) {
	for _, effect := range []EffectKind{EffectPure, Effectful} {
		function := Function{
			ID:          FunctionID(1),
			Declaration: DeclID(2),
			Effect:      effect,
			Generic:     effect == EffectPure,
			Origin:      testOrigin(),
		}
		if !function.Valid() {
			t.Errorf("complete function with effect %d must be valid", effect)
		}
	}
}

func TestFunctionRejectsMissingFields(t *testing.T) {
	for name, function := range map[string]Function{
		"missing id":          {Declaration: DeclID(2), Effect: EffectPure, Origin: testOrigin()},
		"missing declaration": {ID: FunctionID(1), Effect: EffectPure, Origin: testOrigin()},
		"missing effect":      {ID: FunctionID(1), Declaration: DeclID(2), Origin: testOrigin()},
		"missing origin":      {ID: FunctionID(1), Declaration: DeclID(2), Effect: EffectPure},
	} {
		if function.Valid() {
			t.Errorf("%s function must be invalid", name)
		}
	}
}
