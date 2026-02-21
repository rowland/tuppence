package parse

import (
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// annotation = namespaced_annotation | simple_annotation .

func Annotation(tokens []tok.Token) (item ast.Annotation, remainder []tok.Token, err error) {
	namespacedAnnotation, remainder, err := NamespacedAnnotation(tokens)
	if err == nil {
		return namespacedAnnotation, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	simpleAnnotation, remainder, err := SimpleAnnotation(tokens)
	if err == nil {
		return simpleAnnotation, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// namespaced_annotation = "@" namespace ":" identifier annotation_value eol .

func NamespacedAnnotation(tokens []tok.Token) (item *ast.NamespacedAnnotation, remainder []tok.Token, err error) {
	remainder, err = At(tokens)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	namespace, remainder, err := Namespace(remainder)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	remainder, err = Colon(remainder)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	identifier, remainder, err := Identifier(remainder)
	if err == ErrNoMatch {
		return nil, tokens, errorExpecting("identifier", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	annotationValue, remainder, err := AnnotationValue(remainder)
	if err != nil {
		return nil, remainder, errorExpecting("annotation_value", remainder)
	}

	remainder, err = EOL(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewNamespacedAnnotation(namespace, identifier.Name, annotationValue), remainder, nil
}

// simple_annotation = "@" identifier eol .

func SimpleAnnotation(tokens []tok.Token) (item *ast.SimpleAnnotation, remainder []tok.Token, err error) {
	remainder, err = At(tokens)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	identifier, remainder, err := Identifier(remainder)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	remainder, err = EOL(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewSimpleAnnotation(identifier.Name), remainder, nil
}

// namespace = letter { letter | decimal_digit | "_" } .

func Namespace(tokens []tok.Token) (namespace string, remainder []tok.Token, err error) {
	// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .
	identifier, remainder, err := Identifier(tokens)
	if err == ErrNoMatch {
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

func AnnotationValue(tokens []tok.Token) (item ast.AnnotationValue, remainder []tok.Token, err error) {
	stringLiteral, remainder, err := StringLiteral(tokens)
	if err == nil {
		return stringLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	integerLiteral, remainder, err := IntegerLiteral(tokens)
	if err == nil {
		return integerLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	floatLiteral, remainder, err := FloatLiteral(tokens)
	if err == nil {
		return floatLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	booleanLiteral, remainder, err := BooleanLiteral(tokens)
	if err == nil {
		return booleanLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	typeReference, remainder, err := TypeReference(tokens)
	if err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}
