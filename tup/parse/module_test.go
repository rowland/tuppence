package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
)

func TestModule(t *testing.T) {
	src := source.NewSource([]byte("Answer = Int\nanswer = 42\n"), "module.tup")
	module := ast.NewModule("module")

	got, err := Module(src, module)
	if err != nil {
		t.Fatalf("Module(...) = %v, want nil", err)
	}
	if got == nil {
		t.Fatalf("Module(...) = nil, want module")
	}
	if len(got.TopLevelItems) != 2 {
		t.Fatalf("len(Module(...).TopLevelItems) = %d, want 2", len(got.TopLevelItems))
	}

	if _, ok := got.TopLevelItems[0].(*ast.TypeDeclaration); !ok {
		t.Fatalf("Module(...).TopLevelItems[0] = %T, want *ast.TypeDeclaration", got.TopLevelItems[0])
	}
	if _, ok := got.TopLevelItems[1].(*ast.Assignment); !ok {
		t.Fatalf("Module(...).TopLevelItems[1] = %T, want *ast.Assignment", got.TopLevelItems[1])
	}
}
