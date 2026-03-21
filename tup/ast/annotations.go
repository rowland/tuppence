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

func (a *NamespacedAnnotation) String() string {
	return "@" + a.Namespace + ":" + a.Identifier + " " + a.Value.String()
}

// annotation_value = string_literal | ["-"] number | boolean_literal | type_reference .

type AnnotationValue interface {
	Node
	annotationValueNode()
}

func (n *StringLiteral) annotationValueNode()  {}
func (n *IntegerLiteral) annotationValueNode() {}
func (n *FloatLiteral) annotationValueNode()   {}
func (n *BooleanLiteral) annotationValueNode() {}
func (n *TypeReference) annotationValueNode()  {}

type Annotations struct {
	BaseNode
	Annotations []Annotation // List of annotations
}

func NewAnnotations(annotations []Annotation) *Annotations {
	return &Annotations{
		BaseNode:    BaseNode{Type: NodeAnnotations},
		Annotations: annotations,
	}
}

func (a *Annotations) String() string {
	var builder strings.Builder
	for _, annotation := range a.Annotations {
		builder.WriteString(annotation.String())
		builder.WriteString("\n")
	}
	return builder.String()
}
