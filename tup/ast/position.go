package ast

import "fmt"

// Position represents a position in the source code
type Position struct {
	Filename string // Source filename
	Offset   int    // Byte offset, starting at 0
	Line     int    // Line number, starting at 1
	Column   int    // Column number, starting at 1 (currently in bytes, TODO: should be in graphemes)
}

// String returns a string representation of a Position
func (pos Position) String() string {
	if pos.Filename == "" {
		return fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	}
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column)
}
