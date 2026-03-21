package ast

import "github.com/rowland/tuppence/tup/source"

// type_reference = [ identifier { "." identifier } "." ] type_identifier .

type TypeReference struct {
	BaseNode
	Identifiers    []*Identifier   // The identifiers in the type reference
	TypeIdentifier *TypeIdentifier // The type identifier
}

func NewTypeReference(identifiers []*Identifier, typeIdentifier *TypeIdentifier, source *source.Source, startOffset int32, length int32) *TypeReference {
	return &TypeReference{
		BaseNode:       BaseNode{Type: NodeTypeReference, Source: source, StartOffset: startOffset, Length: length},
		TypeIdentifier: typeIdentifier,
		Identifiers:    identifiers,
	}
}

func (t *TypeReference) String() string {
	return t.TypeIdentifier.String()
}

// local_type_reference = type_reference | identifier .

type LocalTypeReference interface {
	Node
	localTypeReferenceNode()
}

func (n *TypeReference) localTypeReferenceNode() {}
func (n *Identifier) localTypeReferenceNode()    {}

// nilable_type = "?" local_type_reference .

type NilableType struct {
	BaseNode
	InnerType Node // The type that is made nilable
}

func NewNilableType(innerType Node) *NilableType {
	return &NilableType{
		BaseNode:  BaseNode{Type: NodeNilableType},
		InnerType: innerType,
	}
}

func (n *NilableType) String() string {
	return "?" + n.InnerType.String()
}
