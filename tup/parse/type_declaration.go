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
		return nil, tokens, ErrNoMatch
	}

	var rhs ast.TypeDeclarationRHS
	if rhs, remainder, err = TypeDeclarationRHS(remainder); err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
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
//                      | type_tuple
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

	if tupleType, remainder, err := TypeTuple(tokens); err == nil {
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

// fallible_type = "!" union_member .

func FallibleType(tokens []tok.Token) (*ast.FallibleType, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokOpNot {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	member, remainder, err := UnionMember(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("union member", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewFallibleType(member), remainder, nil
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
				return nil, tokens, ErrNoMatch
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

// type_tuple = "type" tuple_type .

func TypeTuple(tokens []tok.Token) (*ast.TypeTuple, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwType {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	tupleType, remainder, err := TupleType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("tuple type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewTypeTuple(tupleType), remainder, nil
}

// error_tuple .

func ErrorTuple(tokens []tok.Token) (*ast.ErrorTuple, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwError {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	tupleType, remainder, err := TupleType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("tuple type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewErrorTuple(tupleType), remainder, nil
}

// dynamic_array = "[" "]" (type_reference | array_type) .

func DynamicArray(tokens []tok.Token) (*ast.DynamicArrayType, []tok.Token, error) {
	remainder, found := OpenBracket(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	elementType, remainder, err := ArrayElementType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("array element type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewDynamicArrayType(elementType), remainder, nil
}

// fixed_size_array .

func FixedSizeArray(tokens []tok.Token) (*ast.FixedSizeArrayType, []tok.Token, error) {
	remainder, found := OpenBracket(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	size, remainder, err := Size(remainder)
	if err == ErrNoMatch {
		if peek(remainder).Type == tok.TokCloseBracket {
			return nil, tokens, ErrNoMatch
		}
		return nil, remainder, errorExpecting("array size", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
	}

	elementType, remainder, err := ArrayElementType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("array element type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewFixedSizeArrayType(elementType, size), remainder, nil
}

// union_type = "any"
//            | union_member "|" union_member { "|" union_member } .

func UnionType(tokens []tok.Token) (*ast.UnionType, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if identifier, remainder2, err := Identifier(remainder); err == nil && identifier.Name == "any" {
		return ast.NewUnionType(nil), remainder2, nil
	}

	first, remainder, err := UnionMember(tokens)
	if err != nil {
		return nil, remainder, err
	}

	remainder2, found := Pipe(remainder)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	second, remainder, err := UnionMember(remainder2)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("union member", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	members := []ast.UnionMemberType{first, second}
	for {
		var found bool
		if remainder, found = Pipe(remainder); !found {
			break
		}

		member, remainder2, err := UnionMember(remainder)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("union member", remainder2)
		} else if err != nil {
			return nil, remainder2, err
		}

		members = append(members, member)
		remainder = remainder2
	}

	return ast.NewUnionType(members), remainder, nil
}

// inline_union = "(" union_type ")" .

func InlineUnion(tokens []tok.Token) (*ast.InlineUnion, []tok.Token, error) {
	remainder, found := OpenParen(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	unionType, remainder, err := UnionType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("union type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewInlineUnion(unionType), remainder, nil
}

// union_with_error = ( "!" union_member )
//                  | ( union_member { "|" union_member } "|" "error" )
//                  | ( "(" union_member { "|" union_member } "|" "error" ")" ) .

func UnionWithError(tokens []tok.Token) (*ast.UnionWithError, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type == tok.TokOpNot {
		remainder = remainder[1:]

		member, remainder, err := UnionMember(remainder)
		if err == ErrNoMatch {
			return nil, remainder, errorExpecting("union member", remainder)
		} else if err != nil {
			return nil, remainder, err
		}

		return ast.NewUnionWithError([]ast.UnionMemberType{member}, true), remainder, nil
	}

	if remainder, found := OpenParen(tokens); found {
		unionWithError, remainder, err := unionMembersWithError(remainder)
		if err != nil {
			return nil, remainder, err
		}

		if remainder, found = CloseParen(remainder); !found {
			return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
		}

		return unionWithError, remainder, nil
	}

	return unionMembersWithError(tokens)
}

func unionMembersWithError(tokens []tok.Token) (*ast.UnionWithError, []tok.Token, error) {
	first, remainder, err := UnionMember(tokens)
	if err != nil {
		return nil, remainder, err
	}

	members := []ast.UnionMemberType{first}
	for {
		remainder2, found := Pipe(remainder)
		if !found {
			return nil, tokens, ErrNoMatch
		}

		remainder = remainder2
		remainder2 = skipTrivia(remainder)
		if peek(remainder2).Type == tok.TokKwError {
			return ast.NewUnionWithError(members, false), remainder2[1:], nil
		}

		member, remainder2, err := UnionMember(remainder)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("union member or error", remainder2)
		} else if err != nil {
			return nil, remainder2, err
		}

		members = append(members, member)
		remainder = remainder2
	}
}

// generic_type = type_reference type_argument_list .

func GenericType(tokens []tok.Token) (*ast.GenericType, []tok.Token, error) {
	typeReference, remainder, err := TypeReference(tokens)
	if err != nil {
		return nil, remainder, err
	}

	typeArguments, remainder, err := TypeArgumentList(remainder)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewGenericType(typeReference, typeArguments), remainder, nil
}

// type_argument = type .

func TypeArgument(tokens []tok.Token) (*ast.TypeArgument, []tok.Token, error) {
	if typeNode, remainder, err := Type(tokens); err == nil {
		return ast.NewTypeArgument(typeNode), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// type_argument_list = "[" type_argument { "," type_argument } "]" .

func TypeArgumentList(tokens []tok.Token) (*ast.TypeArgumentList, []tok.Token, error) {
	remainder, found := OpenBracket(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	first, remainder, err := TypeArgument(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("type argument", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	arguments := []*ast.TypeArgument{first}
	for {
		var found bool
		if remainder, found = Comma(remainder); !found {
			break
		}

		argument, remainder2, err := TypeArgument(remainder)
		if err == ErrNoMatch {
			return nil, remainder, errorExpecting("type argument", remainder)
		} else if err != nil {
			return nil, remainder2, err
		}

		arguments = append(arguments, argument)
		remainder = remainder2
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
	}

	return ast.NewTypeArgumentList(arguments), remainder, nil
}

// return_type = union_with_error
//             | union_declaration_with_error
//             | nilable_type
//             | "error"
//             | type .

func ReturnType(tokens []tok.Token) (*ast.ReturnType, []tok.Token, error) {
	if unionWithError, remainder, err := UnionWithError(tokens); err == nil {
		return ast.NewReturnType(unionWithError), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if unionDeclarationWithError, remainder, err := UnionDeclarationWithError(tokens); err == nil {
		return ast.NewReturnType(unionDeclarationWithError), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if nilableType, remainder, err := NilableType(tokens); err == nil {
		return ast.NewReturnType(nilableType), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if inferredErrorType, remainder, err := InferredErrorType(tokens); err == nil {
		return ast.NewReturnType(inferredErrorType), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if typeNode, remainder, err := Type(tokens); err == nil {
		return ast.NewReturnType(typeNode), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// inferred_error_type = "error" .

func InferredErrorType(tokens []tok.Token) (*ast.InferredErrorType, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwError {
		return nil, tokens, ErrNoMatch
	}

	return ast.NewInferredErrorType(), remainder[1:], nil
}

// union_declaration = "union" "(" eol union_members ")" .

func UnionDeclaration(tokens []tok.Token) (*ast.UnionDeclaration, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwUnion {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	if remainder, found = EOL(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokEOL, remainder)
	}

	members, remainder, err := UnionMembers(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("union members", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewUnionDeclaration(members), remainder, nil
}

// union_declaration_with_error = "union" "(" eol
//                              union_member_declaration eol
//                              { union_member_declaration eol }
//                              "error" eol
//                              ")" .

func UnionDeclarationWithError(tokens []tok.Token) (*ast.UnionDeclarationWithError, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwUnion {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	if remainder, found = EOL(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokEOL, remainder)
	}

	members, remainder, err := UnionMembers(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("union members", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	remainder = skipComments(remainder)
	if peek(remainder).Type != tok.TokKwError {
		return nil, remainder, errorExpecting("error", remainder)
	}
	remainder = remainder[1:]

	if remainder, found = EOL(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokEOL, remainder)
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewUnionDeclarationWithError(members), remainder, nil
}

// enum_declaration = "enum" "(" eol enum_members ")" .

func EnumDeclaration(tokens []tok.Token) (*ast.EnumDeclaration, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwEnum {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	if remainder, found = EOL(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokEOL, remainder)
	}

	members, remainder, err := EnumMembers(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("enum members", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewEnumDeclaration(members), remainder, nil
}

// contract_declaration = "contract" "(" eol contract_members ")" .

func ContractDeclaration(tokens []tok.Token) (*ast.ContractDeclaration, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwContract {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	if remainder, found = EOL(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokEOL, remainder)
	}

	members, remainder, err := ContractMembers(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("contract members", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewContractDeclaration(members), remainder, nil
}

// enum_members = enum_member_declaration { eol enum_member_declaration } eol .

func EnumMembers(tokens []tok.Token) (*ast.EnumMembers, []tok.Token, error) {
	first, remainder, err := EnumMemberDeclaration(tokens)
	if err != nil {
		return nil, remainder, err
	}

	members := []*ast.EnumMember{first}
	for {
		remainder2, found := EOL(remainder)
		if !found {
			return nil, remainder, errorExpectingTokenType(tok.TokEOL, remainder)
		}

		next, remainder3, err := EnumMemberDeclaration(remainder2)
		if err == ErrNoMatch {
			return ast.NewEnumMembers(members), remainder2, nil
		} else if err != nil {
			return nil, remainder3, err
		}

		members = append(members, next)
		remainder = remainder3
	}
}

// enum_member_declaration = annotations identifier [ "=" integer_literal ] .

func EnumMemberDeclaration(tokens []tok.Token) (*ast.EnumMember, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	name, remainder, err := Identifier(remainder)
	if err != nil {
		return nil, remainder, err
	}

	var value *ast.IntegerLiteral
	if remainder2, found := AssignOp(remainder); found {
		value, remainder2, err = IntegerLiteral(remainder2)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("integer literal", remainder2)
		} else if err != nil {
			return nil, remainder2, err
		}
		remainder = remainder2
	}

	return ast.NewEnumMember(annotations, name, value), remainder, nil
}

// contract_members = contract_member { eol contract_member } eol .

func ContractMembers(tokens []tok.Token) (*ast.ContractMembers, []tok.Token, error) {
	first, remainder, err := ContractMember(tokens)
	if err != nil {
		return nil, remainder, err
	}

	members := []ast.ContractMemberNode{first}
	for {
		remainder2, found := EOL(remainder)
		if !found {
			return nil, remainder, errorExpectingTokenType(tok.TokEOL, remainder)
		}

		next, remainder3, err := ContractMember(remainder2)
		if err == ErrNoMatch {
			return ast.NewContractMembers(members), remainder2, nil
		} else if err != nil {
			return nil, remainder3, err
		}

		members = append(members, next)
		remainder = remainder3
	}
}

// contract_member = contract_function | contract_field .

func ContractMember(tokens []tok.Token) (ast.ContractMemberNode, []tok.Token, error) {
	if contractFunction, remainder, err := ContractFunction(tokens); err == nil {
		return contractFunction, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if contractField, remainder, err := ContractField(tokens); err == nil {
		return contractField, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// contract_function = function_declaration_lhs "=" function_type .

func ContractFunction(tokens []tok.Token) (*ast.ContractFunction, []tok.Token, error) {
	lhs, remainder, err := FunctionDeclarationLHS(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = AssignOp(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	functionType, remainder, err := FunctionType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("function type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewContractFunction(lhs, functionType), remainder, nil
}

// contract_field = identifier [ "[" type_parameter "]" ] ":" ( nilable_type | type ) .

func ContractField(tokens []tok.Token) (*ast.ContractField, []tok.Token, error) {
	name, remainder, err := Identifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var typeParameter *ast.TypeParameter
	if remainder2, found := OpenBracket(remainder); found {
		typeParameter, remainder2, err = TypeParameter(remainder2)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("type parameter", remainder2)
		} else if err != nil {
			return nil, remainder2, err
		}

		if remainder2, found = CloseBracket(remainder2); !found {
			return nil, remainder2, errorExpectingTokenType(tok.TokCloseBracket, remainder2)
		}

		remainder = remainder2
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	fieldType, remainder, err := ContractFieldType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("contract field type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewContractField(name, typeParameter, fieldType), remainder, nil
}

// union_member_declaration = annotations named_tuple
//                          | union_member_no_annotations .

func UnionMemberDeclaration(tokens []tok.Token) (*ast.UnionMemberDeclaration, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	// annotations is optional in the grammar, so introduced named_tuple members
	// may appear with or without annotations. We give that form precedence over
	// union_member_no_annotations so Ok() and Err(a) are treated as introduced
	// members, not as failed existing-type members.
	if namedTuple, remainder, err := NamedTuple(remainder); err == nil {
		return ast.NewUnionMemberDeclaration(annotations.Annotations, namedTuple), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if len(annotations.Annotations) > 0 {
		return nil, remainder, errorExpecting("named tuple", remainder)
	}

	member, remainder, err := UnionMemberNoAnnotations(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewUnionMemberDeclaration(nil, member), remainder, nil
}

// union_members = union_member_declaration { eol union_member_declaration } eol .

func UnionMembers(tokens []tok.Token) (ast.UnionMembers, []tok.Token, error) {
	first, remainder, err := UnionMemberDeclaration(tokens)
	if err != nil {
		return nil, remainder, err
	}

	members := ast.UnionMembers{first}
	for {
		remainder2, found := EOL(remainder)
		if !found {
			return nil, remainder, errorExpectingTokenType(tok.TokEOL, remainder)
		}

		next, remainder3, err := UnionMemberDeclaration(remainder2)
		if err == ErrNoMatch {
			return members, remainder2, nil
		} else if err != nil {
			return nil, remainder3, err
		}

		members = append(members, next)
		remainder = remainder3
	}
}

// union_member_no_annotations = generic_type
//                             | dynamic_array
//                             | fixed_size_array
//                             | type_reference .

func UnionMemberNoAnnotations(tokens []tok.Token) (ast.UnionDeclarationMemberType, []tok.Token, error) {
	if genericType, remainder, err := GenericType(tokens); err == nil {
		return genericType, remainder, nil
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

	if typeReference, remainder, err := TypeReference(tokens); err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// union_member = named_tuple
//              | generic_type
//              | dynamic_array
//              | fixed_size_array
//              | local_type_reference
//              | contract_declaration .
//
// The parser currently implements the already-supported subset: local type
// references, array types, and named tuples.

func UnionMember(tokens []tok.Token) (ast.UnionMemberType, []tok.Token, error) {
	if namedTuple, remainder, err := NamedTuple(tokens); err == nil {
		return namedTuple, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if genericType, remainder, err := GenericType(tokens); err == nil {
		return genericType, remainder, nil
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

	// local_type_reference = type_reference | identifier .
	// We spell those alternatives out here instead of calling
	// LocalTypeReference(...) so UnionMember(...) can return the narrower
	// UnionMemberType interface directly.
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

	if contractDeclaration, remainder, err := ContractDeclaration(tokens); err == nil {
		return contractDeclaration, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// named_tuple = type_identifier tuple_type .

func NamedTuple(tokens []tok.Token) (*ast.NamedTuple, []tok.Token, error) {
	typeIdentifier, remainder, err := TypeIdentifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	tupleType, remainder, err := TupleType(remainder)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewNamedTuple(typeIdentifier, tupleType), remainder, nil
}

// array_type = fixed_size_array | dynamic_array .

func ArrayType(tokens []tok.Token) (ast.ArrayElementType, []tok.Token, error) {
	if fixedSizeArray, remainder, err := FixedSizeArray(tokens); err == nil {
		return fixedSizeArray, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if dynamicArray, remainder, err := DynamicArray(tokens); err == nil {
		return dynamicArray, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

func ArrayElementType(tokens []tok.Token) (ast.ArrayElementType, []tok.Token, error) {
	if typeReference, remainder, err := TypeReference(tokens); err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if arrayType, remainder, err := ArrayType(tokens); err == nil {
		return arrayType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// tuple_type = "(" [ labeled_tuple_type_members | tuple_type_members ] ")" .

func TupleType(tokens []tok.Token) (*ast.TupleType, []tok.Token, error) {
	remainder, found := OpenParen(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	remainder = skipTrivia(remainder)
	if remainder, found = CloseParen(remainder); found {
		return ast.NewTupleType(nil), remainder, nil
	}

	var members []ast.TupleTypeMemberNode
	var err error
	if members, remainder, err = LabeledTupleTypeMembers(remainder); err == nil {
		if remainder, found = CloseParen(remainder); !found {
			return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
		}
		return ast.NewTupleType(members), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if members, remainder, err = TupleTypeMembers(remainder); err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewTupleType(members), remainder, nil
}

// labeled_tuple_type_member = annotations identifier ":" tuple_type_member .

func LabeledTupleTypeMember(tokens []tok.Token) (*ast.LabeledTupleTypeMember, []tok.Token, error) {
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

	member, remainder, err := TupleTypeMember(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewLabeledTupleTypeMember(annotations, identifier, member.Type), remainder, nil
}

// labeled_tuple_type_members = labeled_tuple_type_member { "," labeled_tuple_type_member } .

func LabeledTupleTypeMembers(tokens []tok.Token) ([]ast.TupleTypeMemberNode, []tok.Token, error) {
	first, remainder, err := LabeledTupleTypeMember(tokens)
	if err != nil {
		return nil, remainder, err
	}

	members := []ast.TupleTypeMemberNode{first}
	for {
		var found bool
		if remainder, found = Comma(remainder); !found {
			break
		}

		next, remainder2, err := LabeledTupleTypeMember(remainder)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("labeled tuple type member", remainder2)
		} else if err != nil {
			return nil, remainder2, err
		}

		members = append(members, next)
		remainder = remainder2
	}

	return members, remainder, nil
}

// tuple_type_member = annotations ( nilable_type
//                                 | type
//                                 | union_type
//                                 | union_declaration
//                                 | literal ) .
//
// The parser currently implements the subset that is already in use elsewhere:
// nilable types, tuple types, local type references, and literals.

func TupleTypeMember(tokens []tok.Token) (*ast.TupleTypeMember, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var memberType ast.FunctionTypeParameterType
	if memberType, remainder, err = tupleTypeMemberType(remainder); err != nil {
		return nil, remainder, err
	}

	return ast.NewTupleTypeMember(annotations, memberType), remainder, nil
}

// tuple_type_members = tuple_type_member { "," tuple_type_member } .

func TupleTypeMembers(tokens []tok.Token) ([]ast.TupleTypeMemberNode, []tok.Token, error) {
	first, remainder, err := TupleTypeMember(tokens)
	if err != nil {
		return nil, remainder, err
	}

	members := []ast.TupleTypeMemberNode{first}
	for {
		var found bool
		if remainder, found = Comma(remainder); !found {
			break
		}

		next, remainder2, err := TupleTypeMember(remainder)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("tuple type member", remainder2)
		} else if err != nil {
			return nil, remainder2, err
		}

		members = append(members, next)
		remainder = remainder2
	}

	return members, remainder, nil
}

func tupleTypeMemberType(tokens []tok.Token) (ast.FunctionTypeParameterType, []tok.Token, error) {
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

	if tupleType, remainder, err := TupleType(tokens); err == nil {
		return tupleType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if localTypeReference, remainder, err := LocalTypeReference(tokens); err == nil {
		memberType, ok := any(localTypeReference).(ast.FunctionTypeParameterType)
		if !ok {
			return nil, remainder, errorExpecting("tuple type member type", remainder)
		}
		return memberType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if literal, remainder, err := Literal(tokens); err == nil {
		memberType, ok := any(literal).(ast.FunctionTypeParameterType)
		if !ok {
			return nil, remainder, errorExpecting("tuple type member type", remainder)
		}
		return memberType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

func ContractFieldType(tokens []tok.Token) (ast.ContractFieldType, []tok.Token, error) {
	if nilableType, remainder, err := NilableType(tokens); err == nil {
		return nilableType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if fieldType, remainder, err := Type(tokens); err == nil {
		return fieldType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}
