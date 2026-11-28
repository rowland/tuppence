package ir

import "testing"

func TestDeclarationValidity(t *testing.T) {
	declaration := Declaration{
		ID:     DeclID(1),
		Scope:  ScopeID(2),
		Name:   "answer",
		Kind:   DeclValue,
		Origin: Origin{Primary: Span{File: FileID(1), Start: 0, End: 6}},
	}
	if !declaration.Valid() {
		t.Fatal("complete declaration must be valid")
	}
}

func TestExportedDeclarationRemainsValid(t *testing.T) {
	declaration := Declaration{
		ID:       DeclID(1),
		Scope:    ScopeID(2),
		Name:     "answer",
		Kind:     DeclValue,
		Exported: true,
		Origin:   Origin{Primary: Span{File: FileID(1), Start: 0, End: 6}},
	}
	if !declaration.Valid() {
		t.Fatal("exported declaration must be valid")
	}
}

func TestDeclarationRejectsMissingFields(t *testing.T) {
	origin := Origin{Primary: Span{File: FileID(1), Start: 0, End: 1}}
	for name, declaration := range map[string]Declaration{
		"missing id":     {Scope: ScopeID(2), Name: "x", Kind: DeclValue, Origin: origin},
		"missing scope":  {ID: DeclID(1), Name: "x", Kind: DeclValue, Origin: origin},
		"missing name":   {ID: DeclID(1), Scope: ScopeID(2), Kind: DeclValue, Origin: origin},
		"missing kind":   {ID: DeclID(1), Scope: ScopeID(2), Name: "x", Origin: origin},
		"missing origin": {ID: DeclID(1), Scope: ScopeID(2), Name: "x", Kind: DeclValue},
	} {
		if declaration.Valid() {
			t.Errorf("%s declaration must be invalid", name)
		}
	}
}

func TestDeclarationKindsHaveValidValues(t *testing.T) {
	for _, kind := range []DeclKind{DeclValue, DeclFunction, DeclType, DeclParameter} {
		if !kind.Valid() {
			t.Errorf("declaration kind %d must be valid", kind)
		}
	}
}
