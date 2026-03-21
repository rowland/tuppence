package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_type = ( "fn" | "fx" ) "(" [ labeled_parameters | parameters ] ")" return_type .

func FunctionType(tokens []tok.Token) (*ast.FunctionType, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	hasSideEffects := false
	switch peek(remainder).Type {
	case tok.TokKwFn:
		remainder = remainder[1:]
	case tok.TokKwFx:
		hasSideEffects = true
		remainder = remainder[1:]
	default:
		return nil, tokens, ErrNoMatch
	}

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	var parameters []ast.FunctionTypeParameter
	if remainder2, found := CloseParen(remainder); found {
		remainder = remainder2
	} else {
		if parameters, remainder, found = functionTypeParameters(remainder); !found {
			return nil, remainder, errorExpecting("function parameters", remainder)
		}

		if remainder, found = CloseParen(remainder); !found {
			return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
		}
	}

	returnType, remainder, err := ReturnType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("return type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewFunctionType(hasSideEffects, parameters, returnType), remainder, nil
}

func functionTypeParameters(tokens []tok.Token) ([]ast.FunctionTypeParameter, []tok.Token, bool) {
	if parameters, remainder, err := LabeledParameters(tokens); err == nil {
		return parameters, remainder, true
	}

	if parameters, remainder, err := Parameters(tokens); err == nil {
		return parameters, remainder, true
	}

	return nil, tokens, false
}

// labeled_parameters = ( labeled_parameter | labeled_rest_parameter ) { "," ( labeled_parameter | labeled_rest_parameter ) } [ "," ] .

func LabeledParameters(tokens []tok.Token) ([]ast.FunctionTypeParameter, []tok.Token, error) {
	first, remainder, err := labeledFunctionTypeParameter(tokens)
	if err != nil {
		return nil, remainder, err
	}

	parameters := []ast.FunctionTypeParameter{first}
	for {
		remainder2, found := Comma(remainder)
		if !found {
			break
		}

		next, remainder3, err := labeledFunctionTypeParameter(remainder2)
		if err == ErrNoMatch {
			return parameters, remainder2, nil
		} else if err != nil {
			return nil, remainder3, err
		}

		parameters = append(parameters, next)
		remainder = remainder3
	}

	return parameters, remainder, nil
}

func labeledFunctionTypeParameter(tokens []tok.Token) (ast.FunctionTypeParameter, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	identifier, remainder, err := Identifier(remainder)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	if _, _, err := RestOperator(remainder); err == nil {
		restParameter, remainder, err := RestParameter(remainder)
		if err != nil {
			return nil, remainder, err
		}
		return ast.NewLabeledRestParameter(annotations, identifier, restParameter), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	paramType, remainder, err := functionTypeParameterValue(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("parameter type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewLabeledParameter(annotations, identifier, paramType), remainder, nil
}

// parameters = ( parameter | rest_parameter ) { "," ( parameter | rest_parameter ) } [ "," ] .

func Parameters(tokens []tok.Token) ([]ast.FunctionTypeParameter, []tok.Token, error) {
	first, remainder, err := positionalFunctionTypeParameter(tokens)
	if err != nil {
		return nil, remainder, err
	}

	parameters := []ast.FunctionTypeParameter{first}
	for {
		remainder2, found := Comma(remainder)
		if !found {
			break
		}

		next, remainder3, err := positionalFunctionTypeParameter(remainder2)
		if err == ErrNoMatch {
			return parameters, remainder2, nil
		} else if err != nil {
			return nil, remainder3, err
		}

		parameters = append(parameters, next)
		remainder = remainder3
	}

	return parameters, remainder, nil
}

func positionalFunctionTypeParameter(tokens []tok.Token) (ast.FunctionTypeParameter, []tok.Token, error) {
	if parameter, remainder, err := RestParameter(tokens); err == nil {
		return parameter, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if parameter, remainder, err := Parameter(tokens); err == nil {
		return parameter, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// parameter = annotations ( nilable_type
//                         | type
//                         | literal
//                         | union_type
//                         | union_declaration ) .

func Parameter(tokens []tok.Token) (*ast.Parameter, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	paramType, remainder, err := functionTypeParameterValue(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewParameter(annotations, paramType), remainder, nil
}

// labeled_parameter = annotations identifier ":" ( nilable_type
//                                              | type
//                                              | literal
//                                              | union_type
//                                              | union_declaration ) .

func LabeledParameter(tokens []tok.Token) (*ast.LabeledParameter, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	identifier, remainder, err := Identifier(remainder)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	paramType, remainder, err := functionTypeParameterValue(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("parameter type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewLabeledParameter(annotations, identifier, paramType), remainder, nil
}

// rest_parameter = "..." type .

func RestParameter(tokens []tok.Token) (*ast.RestParameter, []tok.Token, error) {
	_, remainder, err := RestOperator(tokens)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	paramType, remainder, err := functionTypeType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewRestParameter(paramType), remainder, nil
}

// labeled_rest_parameter = annotations identifier ":" rest_parameter .

func LabeledRestParameter(tokens []tok.Token) (*ast.LabeledRestParameter, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	identifier, remainder, err := Identifier(remainder)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	restParameter, remainder, err := RestParameter(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("rest parameter", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewLabeledRestParameter(annotations, identifier, restParameter), remainder, nil
}

func functionTypeParameterValue(tokens []tok.Token) (ast.FunctionTypeParameterType, []tok.Token, error) {
	if nilableType, remainder, err := NilableType(tokens); err == nil {
		return nilableType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if functionType, remainder, err := FunctionType(tokens); err == nil {
		return functionType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if dynamicArray, remainder, err := DynamicArray(tokens); err == nil {
		return dynamicArray, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if fixedSizeArray, remainder, err := FixedSizeArray(tokens); err == nil {
		return fixedSizeArray, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if errorTuple, remainder, err := ErrorTuple(tokens); err == nil {
		return errorTuple, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if tupleType, remainder, err := TupleType(tokens); err == nil {
		return tupleType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if genericType, remainder, err := GenericType(tokens); err == nil {
		return genericType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if inlineUnion, remainder, err := InlineUnion(tokens); err == nil {
		return inlineUnion, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if unionDeclaration, remainder, err := UnionDeclaration(tokens); err == nil {
		return unionDeclaration, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if unionType, remainder, err := UnionType(tokens); err == nil {
		return unionType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if literal, remainder, err := Literal(tokens); err == nil {
		paramType, ok := literal.(ast.FunctionTypeParameterType)
		if !ok {
			return nil, remainder, errorExpecting("parameter type", remainder)
		}
		return paramType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	// local_type_reference = type_reference | identifier .
	if typeReference, remainder, err := TypeReference(tokens); err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if identifier, remainder, err := Identifier(tokens); err == nil {
		return identifier, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

func functionTypeType(tokens []tok.Token) (ast.TypeArgumentType, []tok.Token, error) {
	if functionType, remainder, err := FunctionType(tokens); err == nil {
		return functionType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if dynamicArray, remainder, err := DynamicArray(tokens); err == nil {
		return dynamicArray, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if fixedSizeArray, remainder, err := FixedSizeArray(tokens); err == nil {
		return fixedSizeArray, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if errorTuple, remainder, err := ErrorTuple(tokens); err == nil {
		return errorTuple, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if tupleType, remainder, err := TupleType(tokens); err == nil {
		return tupleType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if genericType, remainder, err := GenericType(tokens); err == nil {
		return genericType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if inlineUnion, remainder, err := InlineUnion(tokens); err == nil {
		return inlineUnion, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	// local_type_reference = type_reference | identifier .
	if typeReference, remainder, err := TypeReference(tokens); err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if identifier, remainder, err := Identifier(tokens); err == nil {
		return identifier, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}
