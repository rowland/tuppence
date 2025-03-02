package main

import "sort"

type Source struct {
	source   []byte
	filename string
	bol      []int
}

// NewSource instantiates a new Source and calculates the beginning-of-line indices.
func NewSource(source []byte, filename string) *Source {
	newlines := 0
	for _, c := range source {
		if c == '\n' {
			newlines++
		}
	}
	bol := make([]int, newlines+1)
	bol[0] = 0
	line := 1
	for i, c := range source {
		if c == '\n' {
			bol[line] = i + 1
			line++
		}
	}
	return &Source{source: source, filename: filename, bol: bol}
}

// Line returns the 0-basedline number for the given index.
func (s *Source) Line(index int) int {
	return sort.SearchInts(s.bol, index+1) - 1
}

// Column returns the 0-based column number for the given index.
func (s *Source) Column(index int) int {
	return index - s.bol[s.Line(index)]
}

// Position returns the 0-based line and column numbers for the given index.
func (s *Source) Position(index int) (int, int) {
	line := s.Line(index)
	column := index - s.bol[line]
	return line, column
}
