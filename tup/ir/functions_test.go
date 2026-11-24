package ir

import "testing"

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
