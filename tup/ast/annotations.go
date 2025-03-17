package ast

// AnnotationValue represents a value in an annotation
type AnnotationValue struct {
	BaseNode
	Value Node // The value (could be a literal, identifier, or expression)
}

// NewAnnotationValue creates a new AnnotationValue node
func NewAnnotationValue(value Node) *AnnotationValue {
	return &AnnotationValue{
		BaseNode: BaseNode{NodeType: NodeAnnotationValue},
		Value:    value,
	}
}

// String returns a textual representation of the annotation value
func (a *AnnotationValue) String() string {
	return a.Value.String()
}

// Children returns the child nodes
func (a *AnnotationValue) Children() []Node {
	return []Node{a.Value}
}

// Annotations represents a list of annotations
type Annotations struct {
	BaseNode
	Annotations []Node // List of annotations
}

// NewAnnotations creates a new Annotations node
func NewAnnotations(annotations []Node) *Annotations {
	return &Annotations{
		BaseNode:    BaseNode{NodeType: NodeAnnotations},
		Annotations: annotations,
	}
}

// String returns a textual representation of the annotations
func (a *Annotations) String() string {
	result := ""
	for i, annotation := range a.Annotations {
		if i > 0 {
			result += " "
		}
		result += annotation.String()
	}
	return result
}

// Children returns the child nodes
func (a *Annotations) Children() []Node {
	return a.Annotations
}
