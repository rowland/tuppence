package ast

import "github.com/rowland/tuppence/tup/source"

// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .

// Identifier represents a regular identifier in the code (variable name, parameter name, etc.)
type Identifier struct {
	BaseNode
	Name string // The identifier name
}

// NewIdentifier creates a new Identifier node
func NewIdentifier(name string, source *source.Source, startOffset int32, length int32) *Identifier {
	return &Identifier{
		BaseNode: BaseNode{Type: NodeIdentifier, Source: source, StartOffset: startOffset, Length: length},
		Name:     name,
	}
}

// String returns a textual representation of the identifier
func (i *Identifier) String() string {
	return i.Name
}

// type_identifier = uppercase_letter { letter | decimal_digit | "_" } .

// TypeIdentifier represents a type identifier (starts with uppercase)
type TypeIdentifier struct {
	BaseNode
	Name string // The type name
}

// NewTypeIdentifier creates a new TypeIdentifier node
func NewTypeIdentifier(name string, source *source.Source, startOffset int32, length int32) *TypeIdentifier {
	return &TypeIdentifier{
		BaseNode: BaseNode{Type: NodeTypeIdentifier, Source: source, StartOffset: startOffset, Length: length},
		Name:     name,
	}
}

// String returns a textual representation of the type identifier
func (t *TypeIdentifier) String() string {
	return t.Name
}

// function_identifier = lowercase_letter { letter | decimal_digit | "_" } [ "?" | "!" ] .

// FunctionIdentifier represents a function identifier (starts with lowercase, may end with ? or !)
type FunctionIdentifier struct {
	BaseNode
	Name string // The function name
}

// NewFunctionIdentifier creates a new FunctionIdentifier node
func NewFunctionIdentifier(name string, source *source.Source, startOffset int32, length int32) *FunctionIdentifier {
	return &FunctionIdentifier{
		BaseNode: BaseNode{Type: NodeFunctionIdentifier, Source: source, StartOffset: startOffset, Length: length},
		Name:     name,
	}
}

// String returns a textual representation of the function identifier
func (f *FunctionIdentifier) String() string {
	return f.Name
}

type Rename interface {
	// Node
	renameNode()
	Name() string
}

// rename_identifier = identifier [ ":" identifier ] .

// RenameIdentifier represents an identifier with an optional new name for import renaming
type RenameIdentifier struct {
	BaseNode
	Identifier *Identifier
	Original   *Identifier // may be nil if not renamed
}

// NewRenameIdentifier creates a new RenameIdentifier node
func NewRenameIdentifier(identifier *Identifier, original *Identifier) *RenameIdentifier {
	return &RenameIdentifier{
		BaseNode:   BaseNode{Type: NodeRenameIdentifier},
		Identifier: identifier,
		Original:   original,
	}
}

func (r *RenameIdentifier) renameNode() {}

func (r *RenameIdentifier) Name() string {
	return r.Identifier.Name
}

// String returns a textual representation of the rename identifier
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

// NewRenameType creates a new RenameType node
func NewRenameType(identifier *TypeIdentifier, original *TypeIdentifier) *RenameType {
	return &RenameType{
		BaseNode:   BaseNode{Type: NodeRenameType},
		Identifier: identifier,
		Original:   original,
	}
}

func (r *RenameType) renameNode() {}

func (r *RenameType) Name() string {
	return r.Identifier.Name
}

// String returns a textual representation of the rename type
func (r *RenameType) String() string {
	if r.Original != nil {
		return r.Identifier.String() + ": " + r.Original.String()
	}
	return r.Identifier.String()
}
