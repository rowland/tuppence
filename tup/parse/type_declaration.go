package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// type_declaration = type_declaration_lhs "=" type_declaration_rhs .

func TypeDeclaration(tokens []tok.Token) (*ast.TypeDeclaration, []tok.Token, error) {
	var lhs *ast.TypeDeclarationLHS
	remainder := tokens
	var err error

	if lhs, remainder, err = TypeDeclarationLHS(remainder); err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = AssignOp(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpAssign, remainder)
	}

	var rhs ast.TypeDeclarationRHS
	if rhs, remainder, err = TypeDeclarationRHS(remainder); err == ErrNoMatch {
		return nil, remainder, errorExpecting("type declaration right-hand side", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewTypeDeclaration(lhs, rhs), remainder, nil
}

// type_declaration_lhs = annotations type_identifier [ type_parameters ] .

func TypeDeclarationLHS(tokens []tok.Token) (*ast.TypeDeclarationLHS, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var name *ast.TypeIdentifier
	if name, remainder, err = TypeIdentifier(remainder); err != nil {
		return nil, remainder, err
	}

	var typeParameters *ast.TypeParameters
	if typeParameters, remainder, err = TypeParameters(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	return ast.NewTypeDeclarationLHS(annotations.Annotations, name, typeParameters), remainder, nil
}

// type_declaration_rhs = nilable_type
//                      | "type" tuple_type
//                      | error_tuple
//                      | dynamic_array
//                      | fixed_size_array
//                      | union_type
//                      | union_declaration
//                      | enum_declaration
//                      | contract_declaration
//                      | type_reference .

func TypeDeclarationRHS(tokens []tok.Token) (ast.TypeDeclarationRHS, []tok.Token, error) {
	if nilableType, remainder, err := NilableType(tokens); err == nil {
		return nilableType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if tupleType, remainder, err := TypeTupleDeclarationRHS(tokens); err == nil {
		return tupleType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if errorTuple, remainder, err := ErrorTuple(tokens); err == nil {
		return errorTuple, remainder, nil
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

	if unionType, remainder, err := UnionType(tokens); err == nil {
		return unionType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if unionDeclaration, remainder, err := UnionDeclaration(tokens); err == nil {
		return unionDeclaration, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if enumDeclaration, remainder, err := EnumDeclaration(tokens); err == nil {
		return enumDeclaration, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if contractDeclaration, remainder, err := ContractDeclaration(tokens); err == nil {
		return contractDeclaration, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if typeReference, remainder, err := TypeReference(tokens); err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// nilable_type = "?" local_type_reference .

func NilableType(tokens []tok.Token) (*ast.NilableType, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokQuestionMark {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	localTypeReference, remainder, err := LocalTypeReference(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("local type reference", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewNilableType(localTypeReference), remainder, nil
}

// type_parameter = identifier .

func TypeParameter(tokens []tok.Token) (*ast.TypeParameter, []tok.Token, error) {
	identifier, remainder, err := Identifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewTypeParameter(identifier), remainder, nil
}

// type_parameters = "[" type_parameter { "," type_parameter } "]" .

func TypeParameters(tokens []tok.Token) (*ast.TypeParameters, []tok.Token, error) {
	remainder, found := OpenBracket(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	var parameters []*ast.TypeParameter
	for {
		parameter, remainder2, err := TypeParameter(remainder)
		if err == ErrNoMatch {
			if len(parameters) == 0 {
				return nil, remainder, errorExpecting("type parameter", remainder)
			}
			break
		} else if err != nil {
			return nil, remainder2, err
		}

		parameters = append(parameters, parameter)
		remainder = remainder2

		if remainder2, found = Comma(remainder); found {
			remainder = remainder2
			continue
		}
		break
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
	}

	return ast.NewTypeParameters(parameters), remainder, nil
}

// "type" tuple_type .

func TypeTupleDeclarationRHS(tokens []tok.Token) (*ast.TupleType, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}

// error_tuple .

func ErrorTuple(tokens []tok.Token) (ast.TypeDeclarationRHS, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}

// dynamic_array .

func DynamicArray(tokens []tok.Token) (ast.TypeDeclarationRHS, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}

// fixed_size_array .

func FixedSizeArray(tokens []tok.Token) (ast.TypeDeclarationRHS, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}

// union_type .

func UnionType(tokens []tok.Token) (ast.TypeDeclarationRHS, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}

// union_declaration .

func UnionDeclaration(tokens []tok.Token) (ast.TypeDeclarationRHS, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}

// enum_declaration .

func EnumDeclaration(tokens []tok.Token) (ast.TypeDeclarationRHS, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}

// contract_declaration .

func ContractDeclaration(tokens []tok.Token) (ast.TypeDeclarationRHS, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}
