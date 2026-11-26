package ir

import "testing"

func TestPrologConfigValidity(t *testing.T) {
	if !(PrologConfig{}).Valid() {
		t.Fatal("disabled prolog configuration must be valid")
	}
	if (PrologConfig{Enabled: true, Module: ModuleID(1)}).Valid() {
		t.Fatal("enabled prolog configuration requires an architecture")
	}
	if !(PrologConfig{Enabled: true, Module: ModuleID(1), Architecture: "amd64"}).Valid() {
		t.Fatal("enabled prolog configuration with module and architecture must be valid")
	}
}
