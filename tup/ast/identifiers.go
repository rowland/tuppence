package ast

import "github.com/rowland/tuppence/tup/source"

// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .

type Identifier struct {
	BaseNode
	Name string // The identifier name
}

func NewIdentifier(name string, source *source.Source, startOffset int32, length int32) *Identifier {
	return &Identifier{
		BaseNode: BaseNode{Type: NodeIdentifier, Source: source, StartOffset: startOffset, Length: length},
		Name:     name,
	}
}

func (i *Identifier) String() string {
	return i.Name
}

// type_identifier = uppercase_letter { letter | decimal_digit | "_" } .

// TypeIdentifier represents a type identifier (starts with uppercase)
type TypeIdentifier struct {
	BaseNode
	Name string // The type name
}

func NewTypeIdentifier(name string, source *source.Source, startOffset int32, length int32) *TypeIdentifier {
	return &TypeIdentifier{
		BaseNode: BaseNode{Type: NodeTypeIdentifier, Source: source, StartOffset: startOffset, Length: length},
		Name:     name,
	}
}

func (t *TypeIdentifier) String() string {
	return t.Name
}

// function_identifier = lowercase_letter { letter | decimal_digit | "_" } [ "?" | "!" ] .

type FunctionIdentifier struct {
	BaseNode
	Name string // The function name
}

func NewFunctionIdentifier(name string, source *source.Source, startOffset int32, length int32) *FunctionIdentifier {
	return &FunctionIdentifier{
		BaseNode: BaseNode{Type: NodeFunctionIdentifier, Source: source, StartOffset: startOffset, Length: length},
		Name:     name,
	}
}

func (f *FunctionIdentifier) String() string {
	return f.Name
}

type Rename interface {
	Node
	renameNode()
	Name() string
}

func (r *RenameIdentifier) renameNode() {}
func (r *RenameType) renameNode()       {}

// rename_identifier = identifier [ ":" identifier ] .

// RenameIdentifier represents an identifier with an optional new name for import renaming
type RenameIdentifier struct {
	BaseNode
	Identifier *Identifier
	Original   *Identifier // may be nil if not renamed
}

func NewRenameIdentifier(identifier *Identifier, original *Identifier) *RenameIdentifier {
	return &RenameIdentifier{
		BaseNode:   BaseNode{Type: NodeRenameIdentifier},
		Identifier: identifier,
		Original:   original,
	}
}

func (r *RenameIdentifier) Name() string {
	return r.Identifier.Name
}

func (r *RenameIdentifier) String() string {
	if r.Original != nil {
		return r.Identifier.String() + ": " + r.Original.String()
	}
	return r.Identifier.String()
}

// rename_type = type_identifier [ ":" type_identifier ] .

// RenameType represents a type identifier with an optional new name for import renaming
type RenameType struct {
	BaseNode
	Identifier *TypeIdentifier
	Original   *TypeIdentifier // may be nil if not renamed
}

func NewRenameType(identifier *TypeIdentifier, original *TypeIdentifier) *RenameType {
	return &RenameType{
		BaseNode:   BaseNode{Type: NodeRenameType},
		Identifier: identifier,
		Original:   original,
	}
}

func (r *RenameType) Name() string {
	return r.Identifier.Name
}

func (r *RenameType) String() string {
	if r.Original != nil {
		return r.Identifier.String() + ": " + r.Original.String()
	}
	return r.Identifier.String()
}
