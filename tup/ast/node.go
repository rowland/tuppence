package ast

import (
	"fmt"
)

// Position represents a position in the source code
type Position struct {
	Filename string // Source filename
	Offset   int    // Byte offset, starting at 0
	Line     int    // Line number, starting at 1
	Column   int    // Column number, starting at 1 (in characters, not bytes)
}

// String returns a string representation of a Position
func (pos Position) String() string {
	if pos.Filename == "" {
		return fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	}
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column)
}

// NodeType represents the type of AST node
type NodeType string

// Node is the interface implemented by all AST nodes
type Node interface {
	// Pos returns the position of the first character belonging to the node
	Pos() Position
	// End returns the position of the first character immediately after the node
	End() Position
	// Type returns the type of the node
	Type() NodeType
	// String returns a textual representation of the node for debugging
	String() string
	// Children returns all the child nodes
	Children() []Node
}

// BaseNode provides the common implementation for AST nodes
type BaseNode struct {
	// StartPos is the position of the first character belonging to the node
	StartPos Position
	// EndPos is the position of the first character immediately after the node
	EndPos Position
	// NodeType identifies the specific node type
	NodeType NodeType
}

// Pos returns the position of the first character belonging to the node
func (n *BaseNode) Pos() Position {
	return n.StartPos
}

// End returns the position of the first character immediately after the node
func (n *BaseNode) End() Position {
	return n.EndPos
}

// Type returns the type of the node
func (n *BaseNode) Type() NodeType {
	return n.NodeType
}

// String returns a textual representation of the node for debugging
func (n *BaseNode) String() string {
	return string(n.NodeType)
}

// Children returns all the child nodes (empty for BaseNode)
func (n *BaseNode) Children() []Node {
	return nil
}

// SetPos sets the start and end positions for the node
func (n *BaseNode) SetPos(start, end Position) {
	n.StartPos = start
	n.EndPos = end
}
