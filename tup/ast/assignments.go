package ast

import (
	"strings"

	"github.com/rowland/tuppence/tup/source"
)

// assignment = assignment_lhs "=" [ "mut" ] expression .

type Assignment struct {
	BaseNode
	Left  AssignmentLHS // Target of the assignment (typically an identifier)
	Mut   bool          // True if the assignment is mutable
	Right Expression    // Value being assigned
}

// NewAssignment creates a new Assignment node
func NewAssignment(left AssignmentLHS, mut bool, right Expression) *Assignment {
	return &Assignment{
		BaseNode: BaseNode{Type: NodeAssignment},
		Left:     left,
		Mut:      mut,
		Right:    right,
	}
}

// String returns a textual representation of the assignment
func (a *Assignment) String() string {
	var builder strings.Builder
	builder.WriteString(a.Left.String())
	builder.WriteString(" = ")
	if a.Mut {
		builder.WriteString("mut ")
	}
	builder.WriteString(a.Right.String())
	return builder.String()
}

// assignment_lhs = ordinal_assignment_lhs
//                | "(" labeled_assignment_lhs ")" .

type AssignmentLHS interface {
	Node
	assignmentLHSNode()
}

// ordinal_assignment_lhs = identifier { "," identifier } [ "," rest_operator ] .

type OrdinalAssignmentLHS struct {
	BaseNode
	Identifiers  []*Identifier
	RestOperator *RestOperator
}

func NewOrdinalAssignmentLHS(identifiers []*Identifier, restOperator *RestOperator) *OrdinalAssignmentLHS {
	source := identifiers[0].Source
	startOffset := identifiers[0].StartOffset
	length := int32(0)
	if restOperator != nil {
		length = restOperator.StartOffset + restOperator.Length - startOffset
	} else {
		length = identifiers[len(identifiers)-1].StartOffset + identifiers[len(identifiers)-1].Length - startOffset
	}
	return &OrdinalAssignmentLHS{
		BaseNode: BaseNode{
			Type:        NodeOrdinalAssignmentLHS,
			Source:      source,
			StartOffset: startOffset,
			Length:      length,
		},
		Identifiers:  identifiers,
		RestOperator: restOperator,
	}
}

func (o *OrdinalAssignmentLHS) String() string {
	var builder strings.Builder
	for i, identifier := range o.Identifiers {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(identifier.String())
	}
	return builder.String()
}

func (o *OrdinalAssignmentLHS) assignmentLHSNode() {}

// labeled_assignment_lhs = ( rename_identifier | rename_type ) { "," ( rename_identifier | rename_type ) } .

type LabeledAssignmentLHS struct {
	BaseNode
	Renames []Rename
}

func NewLabeledAssignmentLHS(renames []Rename) *LabeledAssignmentLHS {
	var source *source.Source
	var startOffset int32
	var length int32

	switch rename := renames[0].(type) {
	case *RenameIdentifier:
		source = rename.Identifier.Source
		startOffset = rename.Identifier.StartOffset
	case *RenameType:
		source = rename.Identifier.Source
		startOffset = rename.Identifier.StartOffset
	}

	switch rename := renames[len(renames)-1].(type) {
	case *RenameIdentifier:
		if rename.Original != nil {
			length = rename.Original.StartOffset + rename.Original.Length - rename.Original.StartOffset
		} else {
			length = rename.Identifier.StartOffset + rename.Identifier.Length - rename.Identifier.StartOffset
		}
	case *RenameType:
		if rename.Original != nil {
			length = rename.Original.StartOffset + rename.Original.Length - rename.Original.StartOffset
		} else {
			length = rename.Identifier.StartOffset + rename.Identifier.Length - rename.Identifier.StartOffset
		}
	}

	return &LabeledAssignmentLHS{
		BaseNode: BaseNode{
			Type:        NodeLabeledAssignmentLHS,
			Source:      source,
			StartOffset: startOffset,
			Length:      length,
		},
		Renames: renames,
	}
}

func (l *LabeledAssignmentLHS) String() string {
	var builder strings.Builder
	for i, rename := range l.Renames {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(rename.Name())
	}
	return builder.String()
}

func (l *LabeledAssignmentLHS) assignmentLHSNode() {}
