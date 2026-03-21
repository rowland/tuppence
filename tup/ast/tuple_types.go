package ast

import "strings"

// error_tuple = "error" tuple_type .

type ErrorTuple struct {
	BaseNode
	TupleType *TupleType
}

func NewErrorTuple(tupleType *TupleType) *ErrorTuple {
	return &ErrorTuple{
		BaseNode:  BaseNode{Type: NodeErrorTuple},
		TupleType: tupleType,
	}
}

func (e *ErrorTuple) String() string {
	return "error" + e.TupleType.String()
}

// type_tuple = "type" tuple_type .

type TypeTuple struct {
	BaseNode
	TupleType *TupleType
}

func NewTypeTuple(tupleType *TupleType) *TypeTuple {
	return &TypeTuple{
		BaseNode:  BaseNode{Type: NodeTypeTuple},
		TupleType: tupleType,
	}
}

func (t *TypeTuple) String() string {
	return "type" + t.TupleType.String()
}

type TupleTypeMemberNode interface {
	Node
	tupleTypeMemberNode()
}

func (n *TupleTypeMember) tupleTypeMemberNode()        {}
func (n *LabeledTupleTypeMember) tupleTypeMemberNode() {}

// tuple_type_member = annotations ( nilable_type
// 	                               | type
// 	                               | union_type
// 	                               | union_declaration
// 	                               | literal ) .

type TupleTypeMember struct {
	BaseNode
	Annotations *Annotations // Optional annotations
	Type        Node         // Member type
}

func NewTupleTypeMember(annotations *Annotations, memberType Node) *TupleTypeMember {
	return &TupleTypeMember{
		BaseNode:    BaseNode{Type: NodeTupleTypeMember},
		Annotations: annotations,
		Type:        memberType,
	}
}

func (t *TupleTypeMember) String() string {
	var builder strings.Builder

	if t.Annotations != nil {
		builder.WriteString(t.Annotations.String())
	}

	builder.WriteString(t.Type.String())
	return builder.String()
}

// labeled_tuple_type_member = annotations identifier ":" tuple_type_member .

type LabeledTupleTypeMember struct {
	BaseNode
	Annotations *Annotations // Optional annotations
	Identifier  *Identifier  // Field name
	Type        Node         // Field type
}

func NewLabeledTupleTypeMember(annotations *Annotations, identifier *Identifier, memberType Node) *LabeledTupleTypeMember {
	return &LabeledTupleTypeMember{
		BaseNode:    BaseNode{Type: NodeLabeledTupleTypeMember},
		Annotations: annotations,
		Identifier:  identifier,
		Type:        memberType,
	}
}

func (l *LabeledTupleTypeMember) String() string {
	var builder strings.Builder

	if l.Annotations != nil {
		builder.WriteString(l.Annotations.String())
	}

	builder.WriteString(l.Identifier.String())
	builder.WriteString(": ")
	builder.WriteString(l.Type.String())
	return builder.String()
}

// tuple_type = "(" [ labeled_tuple_type_members | tuple_type_members ] ")" .

type TupleType struct {
	BaseNode
	Members []TupleTypeMemberNode
}

func NewTupleType(members []TupleTypeMemberNode) *TupleType {
	return &TupleType{
		BaseNode: BaseNode{Type: NodeTupleType},
		Members:  members,
	}
}

func (t *TupleType) String() string {
	if len(t.Members) == 0 {
		return "()"
	}

	multiline := len(t.Members) > 1
	if !multiline {
		for _, member := range t.Members {
			if strings.Contains(member.String(), "\n") {
				multiline = true
				break
			}
		}
	}

	var builder strings.Builder
	builder.WriteString("(")
	if multiline {
		builder.WriteString("\n")
		for _, member := range t.Members {
			builder.WriteString(indentString(member.String()))
			builder.WriteString(",\n")
		}
		builder.WriteString(")")
		return builder.String()
	}

	for i, member := range t.Members {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(member.String())
	}
	builder.WriteString(")")
	return builder.String()
}

// named_tuple = type_identifier tuple_type .

type NamedTuple struct {
	BaseNode
	TypeIdentifier *TypeIdentifier
	TupleType      *TupleType
}

func NewNamedTuple(typeIdentifier *TypeIdentifier, tupleType *TupleType) *NamedTuple {
	return &NamedTuple{
		BaseNode:       BaseNode{Type: NodeNamedTuple},
		TypeIdentifier: typeIdentifier,
		TupleType:      tupleType,
	}
}

func (n *NamedTuple) String() string {
	return n.TypeIdentifier.String() + n.TupleType.String()
}
