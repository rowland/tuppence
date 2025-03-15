package source

import (
	"testing"
)

func TestNewSource(t *testing.T) {
	source := []byte("line 1\nline 2\nline 3\nline 4\nline 5\n")
	filename := "test.txt"
	s := NewSource(source, filename)

	if s.Filename != filename {
		t.Errorf("Expected filename %s, got %s", filename, s.Filename)
	}
	if len(s.bol) != 6 { // 5 lines + 1 for the start
		t.Errorf("Expected 6 BOL indices, got %d", len(s.bol))
	}
}

func TestLine(t *testing.T) {
	source := []byte("line 1\nline 2\nline 3\nline 4\nline 5\n")
	s := NewSource(source, "test.txt")

	tests := []struct {
		index int
		line  int
	}{
		{0, 0},  // index 0 is in line 0
		{10, 1}, // index 10 falls in line 1 (since line 1 starts at index 7)
		{11, 1}, // still line 1
		{21, 3}, // index 21 is the start of line 3
		{22, 3}, // index 22 is still in line 3
		{32, 4}, // index 32 is in line 4 (line 4 starts at index 28)
		{33, 4}, // index 33 is also in line 4
	}

	for _, test := range tests {
		result := s.Line(test.index)
		if result != test.line {
			t.Errorf("Line(%d) = %d; want %d", test.index, result, test.line)
		}
	}
}

func TestColumn(t *testing.T) {
	source := []byte("line 1\nline 2\nline 3\nline 4\nline 5\n")
	s := NewSource(source, "test.txt")

	tests := []struct {
		index  int
		column int
	}{
		{0, 0},  // index 0: line 0 starts at index 0, so column = 0 - 0 = 0
		{10, 3}, // index 10: line 1 starts at index 7, so column = 10 - 7 = 3
		{11, 4}, // index 11: line 1 starts at index 7, so column = 11 - 7 = 4
		{21, 0}, // index 21: line 3 starts at index 21, so column = 21 - 21 = 0
		{22, 1}, // index 22: line 3 starts at index 21, so column = 22 - 21 = 1
		{32, 4}, // index 32: line 4 starts at index 28, so column = 32 - 28 = 4
		{33, 5}, // index 33: line 4 starts at index 28, so column = 33 - 28 = 5
	}

	for _, test := range tests {
		result := s.Column(test.index)
		if result != test.column {
			t.Errorf("Column(%d) = %d; want %d", test.index, result, test.column)
		}
	}
}

func TestPosition(t *testing.T) {
	source := []byte("line 1\nline 2\nline 3\nline 4\nline 5\n")
	s := NewSource(source, "test.txt")

	tests := []struct {
		index  int
		line   int
		column int
	}{
		{0, 0, 0},  // index 0: line 0 starts at index 0, so column = 0 - 0 = 0
		{10, 1, 3}, // index 10: falls in line 1 (line 1 starts at index 7), so column = 10 - 7 = 3
		{11, 1, 4}, // index 11: still in line 1, column = 11 - 7 = 4
		{21, 3, 0}, // index 21: is the first character of line 3 (line 3 starts at index 21), so column = 21 - 21 = 0
		{22, 3, 1}, // index 22: in line 3, column = 22 - 21 = 1
		{32, 4, 4}, // index 32: falls in line 4 (line 4 starts at index 28), so column = 32 - 28 = 4
		{33, 4, 5}, // index 33: still in line 4, column = 33 - 28 = 5
	}

	for _, test := range tests {
		line, column := s.Position(test.index)
		if line != test.line || column != test.column {
			t.Errorf("Position(%d) = (%d, %d); want (%d, %d)", test.index, line, column, test.line, test.column)
		}
	}
}
