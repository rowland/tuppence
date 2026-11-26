package ir

import "testing"

func TestSelectModuleFilesGroupsMatchingTargetDeterministically(t *testing.T) {
	selected, invalid := SelectModuleFiles([]string{
		"zeta.more.tup",
		"math.vectors--arm64.tup",
		"math.tup",
		"math.vectors--amd64.tup",
		"filesystem--linux.tup",
		"invalid.txt",
	}, "linux", "amd64")
	if len(invalid) != 1 || invalid[0] != "invalid.txt" {
		t.Fatalf("unexpected invalid files: %v", invalid)
	}
	if len(selected) != 3 {
		t.Fatalf("unexpected selected groups: %#v", selected)
	}
	if selected[0].ModuleName != "filesystem" || selected[1].ModuleName != "math" || selected[2].ModuleName != "zeta" {
		t.Fatalf("groups must be sorted by module name: %#v", selected)
	}
	if len(selected[1].Files) != 2 || selected[1].Files[0] != "math.tup" || selected[1].Files[1] != "math.vectors--amd64.tup" {
		t.Fatalf("matching files must be sorted and filtered: %#v", selected[1].Files)
	}
}
