package parse

import (
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// annotation = namespaced_annotation | simple_annotation .

func Annotation(tokens []tok.Token) (annot ast.Annotation, remainder []tok.Token, err error) {
	var namespacedAnnotation *ast.NamespacedAnnotation
	if namespacedAnnotation, remainder, err = NamespacedAnnotation(tokens); err == nil {
		return namespacedAnnotation, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var simpleAnnotation *ast.SimpleAnnotation
	if simpleAnnotation, remainder, err = SimpleAnnotation(tokens); err == nil {
		return simpleAnnotation, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// namespaced_annotation = "@" namespace ":" identifier annotation_value eol .

func NamespacedAnnotation(tokens []tok.Token) (annot *ast.NamespacedAnnotation, remainder []tok.Token, err error) {
	var found bool
	if remainder, found = At(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	var namespace string
	if namespace, remainder, err = Namespace(remainder); err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var identifier *ast.Identifier
	if identifier, remainder, err = Identifier(remainder); err == ErrNoMatch {
		return nil, tokens, errorExpecting("identifier", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	var annotationValue ast.AnnotationValue
	if annotationValue, remainder, err = AnnotationValue(remainder); err != nil {
		return nil, remainder, errorExpecting("annotation_value", remainder)
	}

	if remainder, found = EOL(remainder); !found {
		return nil, remainder, ErrNoMatch
	}

	return ast.NewNamespacedAnnotation(namespace, identifier.Name, annotationValue), remainder, nil
}

// simple_annotation = "@" identifier eol .

func SimpleAnnotation(tokens []tok.Token) (annot *ast.SimpleAnnotation, remainder []tok.Token, err error) {
	var found bool
	if remainder, found = At(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	var identifier *ast.Identifier
	if identifier, remainder, err = Identifier(remainder); err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder, found = EOL(remainder); !found {
		return nil, remainder, ErrNoMatch
	}

	return ast.NewSimpleAnnotation(identifier.Name), remainder, nil
}

// namespace = letter { letter | decimal_digit | "_" } .

func Namespace(tokens []tok.Token) (namespace string, remainder []tok.Token, err error) {
	// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .
	var identifier *ast.Identifier
	if identifier, remainder, err = Identifier(tokens); err == ErrNoMatch {
		return "", tokens, ErrNoMatch
	} else if err != nil {
		return "", remainder, err
	}
	if strings.HasPrefix(identifier.Name, "_") {
		return "", tokens, ErrNoMatch
	}
	return identifier.Name, remainder, nil
}

// annotation_value = string_literal | ["-"] number | boolean_literal | type_reference .

func AnnotationValue(tokens []tok.Token) (value ast.AnnotationValue, remainder []tok.Token, err error) {
	var stringLiteral *ast.StringLiteral
	if stringLiteral, remainder, err = StringLiteral(tokens); err == nil {
		return stringLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var integerLiteral *ast.IntegerLiteral
	if integerLiteral, remainder, err = IntegerLiteral(tokens); err == nil {
		return integerLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var floatLiteral *ast.FloatLiteral
	if floatLiteral, remainder, err = FloatLiteral(tokens); err == nil {
		return floatLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var booleanLiteral *ast.BooleanLiteral
	if booleanLiteral, remainder, err = BooleanLiteral(tokens); err == nil {
		return booleanLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var typeReference *ast.TypeReference
	if typeReference, remainder, err = TypeReference(tokens); err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}
