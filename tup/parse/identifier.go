package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .

func Identifier(tokens []tok.Token) (ident *ast.Identifier, remainder []tok.Token, err error) {
	// fmt.Println("Identifier", tokens)
	remainder = skipComments(tokens)
	t := peek(remainder)
	if t.Type != tok.TokID && t.Type != tok.TokKwIt {
		return nil, tokens, ErrNoMatch
	}
	return ast.NewIdentifier(t.Value(), t.File, t.Offset, t.Length), remainder[1:], nil
}

// type_identifier = uppercase_letter { letter | decimal_digit | "_" } .

func TypeIdentifier(tokens []tok.Token) (typeIdent *ast.TypeIdentifier, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	t := peek(remainder)
	if t.Type != tok.TokTypeID {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokTypeID, remainder)
	}

	return ast.NewTypeIdentifier(t.Value(), t.File, t.Offset, t.Length), remainder[1:], nil
}

// rename_identifier = identifier [ ":" identifier ] .

func RenameIdentifier(tokens []tok.Token) (ident *ast.RenameIdentifier, remainder []tok.Token, err error) {
	var identifier *ast.Identifier
	if identifier, remainder, err = Identifier(tokens); err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	remainder2, found := Colon(remainder)
	if !found {
		// no colon found, so no renaming
		return ast.NewRenameIdentifier(identifier, nil), remainder2, nil
	}

	// colon found, so original identifier expected
	original, remainder3, err3 := Identifier(remainder2)
	if err3 != nil {
		return nil, remainder3, err3
	}

	return ast.NewRenameIdentifier(identifier, original), remainder3, nil
}

// rename_type = type_identifier [ ":" type_identifier ] .

func RenameType(tokens []tok.Token) (typeIdent *ast.RenameType, remainder []tok.Token, err error) {
	var typeIdentifier *ast.TypeIdentifier
	if typeIdentifier, remainder, err = TypeIdentifier(tokens); err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	remainder2, found := Colon(remainder)
	if !found {
		// no colon found, so no renaming
		return ast.NewRenameType(typeIdentifier, nil), remainder, nil
	}

	// colon found, so original type identifier expected
	original, remainder3, err3 := TypeIdentifier(remainder2)
	if err3 != nil {
		return nil, remainder3, err3
	}

	return ast.NewRenameType(typeIdentifier, original), remainder3, nil
}

// type_reference = [ identifier { "." identifier } "." ] type_identifier .

func TypeReference(tokens []tok.Token) (typeRef *ast.TypeReference, remainder []tok.Token, err error) {
	// fmt.Println("TypeReference", tokens)
	remainder = skipComments(tokens)

	var identifiers []*ast.Identifier
	for {
		var identifier *ast.Identifier
		var remainder2 []tok.Token
		identifier, remainder2, err = Identifier(remainder)

		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder2, err
		}

		remainder = remainder2
		identifiers = append(identifiers, identifier)
		var found bool
		if remainder, found = Dot(remainder); !found {
			return nil, tokens, ErrNoMatch
		}
	}

	var typeIdentifier *ast.TypeIdentifier
	if typeIdentifier, remainder, err = TypeIdentifier(remainder); err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, nil, err
	}

	var src *source.Source
	var startOffset, length int32
	if len(identifiers) > 0 {
		src = identifiers[0].Source
		startOffset = identifiers[0].StartOffset
		length = typeIdentifier.StartOffset + typeIdentifier.Length - identifiers[0].StartOffset
	} else {
		src = typeIdentifier.Source
		startOffset = typeIdentifier.StartOffset
		length = typeIdentifier.Length
	}
	return ast.NewTypeReference(identifiers, typeIdentifier, src, startOffset, length), remainder, nil
}

// function_identifier = lowercase_letter { letter | decimal_digit | "_" } [ "?" | "!" ] .

func FunctionIdentifier(tokens []tok.Token) (fundIdent *ast.FunctionIdentifier, remainder []tok.Token, err error) {
	// fmt.Println("FunctionIdentifier", tok.Types(tokens))
	remainder = skipComments(tokens)

	t := peek(remainder)
	if t.Type != tok.TokFuncID && t.Type != tok.TokID {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokFuncID, remainder)
	}

	return ast.NewFunctionIdentifier(t.Value(), t.File, t.Offset, t.Length), remainder[1:], nil
}
