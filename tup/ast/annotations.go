package ast

import "strings"

// annotation = namespaced_annotation | simple_annotation .

type Annotation interface {
	Node
	annotationNode()
}

func (n *SimpleAnnotation) annotationNode()     {}
func (n *NamespacedAnnotation) annotationNode() {}

// simple_annotation = "@" identifier eol .

type SimpleAnnotation struct {
	BaseNode
	Identifier string // The name of the annotation
}

func NewSimpleAnnotation(identifier string) *SimpleAnnotation {
	return &SimpleAnnotation{
		BaseNode:   BaseNode{Type: NodeAnnotation},
		Identifier: identifier,
	}
}

// String returns a textual representation of the annotation
func (a *SimpleAnnotation) String() string {
	return "@" + a.Identifier
}

// namespaced_annotation = "@" namespace ":" identifier annotation_value eol .

type NamespacedAnnotation struct {
	BaseNode
	Namespace  string          // The namespace of the annotation
	Identifier string          // The name of the annotation
	Value      AnnotationValue // Optional argument value
}

func NewNamespacedAnnotation(namespace, name string, value AnnotationValue) *NamespacedAnnotation {
	return &NamespacedAnnotation{
		BaseNode:   BaseNode{Type: NodeAnnotation},
		Namespace:  namespace,
		Identifier: name,
		Value:      value,
	}
}

// String returns a textual representation of the annotation
func (a *NamespacedAnnotation) String() string {
	return "@" + a.Namespace + ":" + a.Identifier + " " + a.Value.String()
}

// annotation_value = string_literal | ["-"] number | boolean_literal | type_reference .

type AnnotationValue interface {
	Node
	annotationValueNode()
}

// NewAnnotationValue creates a new AnnotationValue node
func (n *StringLiteral) annotationValueNode()  {}
func (n *IntegerLiteral) annotationValueNode() {}
func (n *FloatLiteral) annotationValueNode()   {}
func (n *BooleanLiteral) annotationValueNode() {}
func (n *TypeReference) annotationValueNode()  {}

// Annotations represents a list of annotations
type Annotations struct {
	BaseNode
	Annotations []Annotation // List of annotations
}

// NewAnnotations creates a new Annotations node
func NewAnnotations(annotations []Annotation) *Annotations {
	return &Annotations{
		BaseNode:    BaseNode{Type: NodeAnnotations},
		Annotations: annotations,
	}
}

// String returns a textual representation of the annotations
func (a *Annotations) String() string {
	var builder strings.Builder
	for _, annotation := range a.Annotations {
		builder.WriteString(annotation.String())
		builder.WriteString("\n")
	}
	return builder.String()
}
