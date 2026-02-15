package parse

import (
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// annotation = namespaced_annotation | simple_annotation .

func Annotation(tokens []tok.Token) (item ast.Annotation, remainder []tok.Token, err error) {
	var errors []error
	namespacedAnnotation, remainder, err := NamespacedAnnotation(tokens)
	if err == nil {
		return namespacedAnnotation, remainder, nil
	}
	errors = append(errors, err)

	simpleAnnotation, remainder, err := SimpleAnnotation(tokens)
	if err == nil {
		return simpleAnnotation, remainder, nil
	}
	errors = append(errors, err)

	return nil, nil, errorExpectingOneOf("annotation", tokens, errors)
}

// namespaced_annotation = "@" namespace ":" identifier annotation_value eol .

func NamespacedAnnotation(tokens []tok.Token) (item *ast.NamespacedAnnotation, remainder []tok.Token, err error) {
	remainder, err = At(tokens)
	if err != nil {
		return nil, nil, err
	}

	namespace, remainder, err := Namespace(remainder)
	if err != nil {
		return nil, nil, err
	}

	remainder, err = Colon(remainder)
	if err != nil {
		return nil, nil, err
	}

	identifier, remainder, err := Identifier(remainder)
	if err != nil {
		return nil, nil, err
	}

	annotationValue, remainder, err := AnnotationValue(remainder)
	if err != nil {
		return nil, nil, err
	}

	remainder, err = EOL(remainder)
	if err != nil {
		return nil, nil, err
	}

	return ast.NewNamespacedAnnotation(namespace, identifier.Name, annotationValue), remainder, nil
}

// simple_annotation = "@" identifier eol .

func SimpleAnnotation(tokens []tok.Token) (item *ast.SimpleAnnotation, remainder []tok.Token, err error) {
	remainder, err = At(tokens)
	if err != nil {
		return nil, nil, err
	}

	identifier, remainder, err := Identifier(remainder)
	if err != nil {
		return nil, nil, err
	}

	remainder, err = EOL(remainder)
	if err != nil {
		return nil, nil, err
	}

	return ast.NewSimpleAnnotation(identifier.Name), remainder, nil
}

// namespace = letter { letter | decimal_digit | "_" } .

func Namespace(tokens []tok.Token) (namespace string, remainder []tok.Token, err error) {
	// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .
	identifier, remainder, err := Identifier(tokens)
	if err != nil {
		return "", nil, err
	}
	if strings.HasPrefix(identifier.Name, "_") {
		return "", nil, errorExpecting("namespace", tokens)
	}
	return identifier.Name, remainder, nil
}

// annotation_value = string_literal | ["-"] number | boolean_literal | type_reference .

func AnnotationValue(tokens []tok.Token) (item ast.AnnotationValue, remainder []tok.Token, err error) {
	var errors []error

	stringLiteral, remainder, err := StringLiteral(tokens)
	if err == nil {
		return stringLiteral, remainder, nil
	}
	errors = append(errors, err)

	integerLiteral, remainder, err := IntegerLiteral(tokens)
	if err == nil {
		return integerLiteral, remainder, nil
	}
	errors = append(errors, err)

	floatLiteral, remainder, err := FloatLiteral(tokens)
	if err == nil {
		return floatLiteral, remainder, nil
	}
	errors = append(errors, err)

	booleanLiteral, remainder, err := BooleanLiteral(tokens)
	if err == nil {
		return booleanLiteral, remainder, nil
	}
	errors = append(errors, err)

	typeReference, remainder, err := TypeReference(tokens)
	if err == nil {
		return typeReference, remainder, nil
	}
	errors = append(errors, err)

	return nil, nil, errorExpectingOneOf("annotation_value", tokens, errors)
}
