package ir

import "testing"

func TestIndexSourceFilesAssignsStableIDsAndPreservesPaths(t *testing.T) {
	table, invalid := IndexSourceFiles([]string{
		"math.vectors--arm64.tup",
		"math.tup",
		"math.vectors--amd64.tup",
		"bad.txt",
	}, "linux", "amd64")
	if len(invalid) != 1 || invalid[0] != "bad.txt" {
		t.Fatalf("unexpected invalid paths: %v", invalid)
	}
	if !table.Valid() || len(table) != 2 {
		t.Fatalf("unexpected source file table: %#v", table)
	}
	if table[0].ID != FileID(1) || table[1].ID != FileID(2) || table[0].Path != "math.tup" || table[1].Path != "math.vectors--amd64.tup" {
		t.Fatalf("source files must be deterministic: %#v", table)
	}
}
