package ast

import (
	"github.com/rowland/tuppence/tup/source"
)

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
	// Source is the reference to the source file
	Source *source.Source
	// StartOffset is the byte offset of the first character belonging to the node
	StartOffset int
	// Length is the length of the node in bytes
	Length int
	// NodeType identifies the specific node type
	NodeType NodeType
}

// Pos returns the position of the first character belonging to the node
func (n *BaseNode) Pos() Position {
	if n.Source == nil {
		return Position{}
	}

	line := n.Source.Line(n.StartOffset)
	column := n.Source.Column(n.StartOffset)

	// TODO: Convert column from byte offset to grapheme count for proper
	// display to users. Currently, this will be incorrect for non-ASCII text.
	return Position{
		Filename: n.Source.Filename,
		Offset:   n.StartOffset,
		Line:     line + 1,   // Convert from 0-based to 1-based
		Column:   column + 1, // Convert from 0-based to 1-based (currently in bytes)
	}
}

// End returns the position of the first character immediately after the node
func (n *BaseNode) End() Position {
	if n.Source == nil {
		return Position{}
	}

	endOffset := n.StartOffset + n.Length
	line := n.Source.Line(endOffset)
	column := n.Source.Column(endOffset)

	// TODO: Convert column from byte offset to grapheme count for proper
	// display to users. Currently, this will be incorrect for non-ASCII text.
	return Position{
		Filename: n.Source.Filename,
		Offset:   endOffset,
		Line:     line + 1,   // Convert from 0-based to 1-based
		Column:   column + 1, // Convert from 0-based to 1-based (currently in bytes)
	}
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

// SetPos sets the source, start offset, and length for the node
func (n *BaseNode) SetPos(source *source.Source, startOffset int, length int) {
	n.Source = source
	n.StartOffset = startOffset
	n.Length = length
}
