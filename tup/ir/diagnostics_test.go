package ir

import "testing"

func TestDiagnosticRequiresStableCodeMessageAndOrigin(t *testing.T) {
	origin := Origin{Primary: Span{File: FileID(1), Start: 2, End: 5}}
	diagnostic := Diagnostic{
		Severity: SeverityError,
		Code:     "IR001",
		Message:  "invalid reference",
		Origin:   origin,
	}
	if !diagnostic.Valid() {
		t.Fatal("complete diagnostic must be valid")
	}

	for name, incomplete := range map[string]Diagnostic{
		"missing severity": {Code: "IR001", Message: "message", Origin: origin},
		"missing code":     {Severity: SeverityError, Message: "message", Origin: origin},
		"missing message":  {Severity: SeverityError, Code: "IR001", Origin: origin},
		"missing origin":   {Severity: SeverityError, Code: "IR001", Message: "message"},
	} {
		if incomplete.Valid() {
			t.Errorf("%s diagnostic must be invalid", name)
		}
	}
}

func TestDiagnosticRejectsInvalidRelatedOrigins(t *testing.T) {
	diagnostic := Diagnostic{
		Severity: SeverityWarning,
		Code:     "IR002",
		Message:  "ambiguous reference",
		Origin:   Origin{Primary: Span{File: FileID(1), Start: 0, End: 1}},
		Related:  []OriginID{OriginID(2), OriginID(0)},
	}
	if diagnostic.Valid() {
		t.Fatal("diagnostic with an invalid related origin must be invalid")
	}
}

func TestDiagnosticRejectsImportCandidatesWithoutPath(t *testing.T) {
	diagnostic := Diagnostic{
		Severity:         SeverityError,
		Code:             "IR004",
		Message:          "module not found",
		Origin:           testOrigin(),
		ImportCandidates: []string{"app/foo"},
	}
	if diagnostic.Valid() {
		t.Fatal("import candidates require a requested import path")
	}
}

func TestDiagnosticCanonicalIdentityMustBeComplete(t *testing.T) {
	diagnostic := Diagnostic{
		Severity:        SeverityError,
		Code:            "IR003",
		Message:         "missing symbol",
		Origin:          testOrigin(),
		CanonicalModule: ModuleID(1),
	}
	if diagnostic.Valid() {
		t.Fatal("diagnostic must provide both canonical module and declaration")
	}
}
