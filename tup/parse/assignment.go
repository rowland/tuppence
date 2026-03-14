package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// assignment = assignment_lhs "=" [ "mut" ] expression .

func Assignment(tokens []tok.Token) (assignment *ast.Assignment, remainder []tok.Token, err error) {
	// fmt.Println("Assignment", tok.Types(tokens))
	var left ast.AssignmentLHS
	if left, remainder, err = AssignmentLHS(tokens); err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = AssignOp(remainder); !found {
		return nil, remainder, errorExpecting("=", remainder)
	}

	mut := ast.Immutable
	if peek(remainder).Type == tok.TokKwMut {
		mut = ast.Mutable
		remainder = remainder[1:]
	}

	var right ast.Expression
	if right, remainder, err = Expression(remainder); err != nil {
		return nil, remainder, err
	}

	return ast.NewAssignment(left, mut, right), remainder, nil
}

// assignment_lhs = labeled_assignment_lhs
//                | ordinal_assignment_lhs  .

func AssignmentLHS(tokens []tok.Token) (lhs ast.AssignmentLHS, remainder []tok.Token, err error) {
	var labeledLHS *ast.LabeledAssignmentLHS
	if labeledLHS, remainder, err = labeledAssignmentLHS(tokens); err == nil {
		return labeledLHS, remainder, nil
	}

	var ordinalLHS *ast.OrdinalAssignmentLHS
	if ordinalLHS, remainder, err = ordinalAssignmentLHS(tokens); err == nil {
		return ordinalLHS, remainder, nil
	}

	return nil, tokens, ErrNoMatch
}

// ordinal_assignment_lhs = identifier { "," identifier } [ "," rest_operator ] .

func ordinalAssignmentLHS(tokens []tok.Token) (lhs *ast.OrdinalAssignmentLHS, remainder []tok.Token, err error) {
	var identifiers []*ast.Identifier

	identifier, remainder, err := Identifier(tokens)
	if identifier == nil || err != nil {
		return nil, remainder, err
	}

	identifiers = append(identifiers, identifier)

	identifier, remainder2, err := commaIdentifier(remainder)
	for err == nil {
		identifiers = append(identifiers, identifier)
		remainder = remainder2
		identifier, remainder2, err = commaIdentifier(remainder2)
	}

	restOperator, remainder3, err := commaRestOperator(remainder)
	if err == nil {
		remainder = remainder3
	}
	if _, _, err := commaRestOperator(remainder); err == nil {
		return nil, remainder, errorNotExpecting(remainder)
	}

	return ast.NewOrdinalAssignmentLHS(identifiers, restOperator), remainder, nil
}

func commaIdentifier(tokens []tok.Token) (ident *ast.Identifier, remainder []tok.Token, err error) {
	if peek(tokens).Type != tok.TokComma {
		return nil, nil, errorExpecting(",", tokens)
	}
	remainder = tokens[1:]
	identifier, remainder, err := Identifier(remainder)
	if identifier == nil || err != nil {
		return nil, nil, err
	}
	return identifier, remainder, nil
}

// "," rest_operator
// rest_operator = "..." [ identifier ] .

func commaRestOperator(tokens []tok.Token) (op *ast.RestOperator, remainder []tok.Token, err error) {
	var found bool
	if remainder, found = Comma(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	restOperator, remainder, err := RestOperator(remainder)
	if restOperator == nil || err != nil {
		return nil, nil, err
	}
	return restOperator, remainder, nil
}

// rest_operator = "..." [ identifier ] .

func RestOperator(tokens []tok.Token) (op *ast.RestOperator, remainder []tok.Token, err error) {
	var found bool
	if remainder, found = RestOp(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	// returns nil if no identifier found, but we already matched the "..."
	identifier, remainder, _ := Identifier(remainder)

	return ast.NewRestOperator(identifier), remainder, nil
}

// labeled_assignment_lhs = "(" ( rename_identifier | rename_type ) { "," ( rename_identifier | rename_type ) } ")" .

func labeledAssignmentLHS(tokens []tok.Token) (lhs *ast.LabeledAssignmentLHS, remainder []tok.Token, err error) {
	var renames []ast.Rename

	var found bool
	if remainder, found = OpenParen(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	rename, remainder, err := Rename(remainder)
	if rename == nil || err != nil {
		return nil, nil, err
	}
	renames = append(renames, rename)

	rename, remainder2, err := commaRename(remainder)
	for err == nil {
		renames = append(renames, rename)
		remainder = remainder2
		rename, remainder2, err = commaRename(remainder2)
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, ErrNoMatch
	}

	return ast.NewLabeledAssignmentLHS(renames), remainder, nil
}

// rename_identifier | rename_type

func Rename(tokens []tok.Token) (ren ast.Rename, remainder []tok.Token, err error) {
	var errors []error

	id, remainder, err := RenameIdentifier(tokens)
	if err != nil {
		errors = append(errors, err)
	} else if id != nil {
		return id, remainder, nil
	}

	typeId, remainder, err := RenameType(tokens)
	if err != nil {
		errors = append(errors, err)
	} else if typeId != nil {
		return typeId, remainder, nil
	}

	return nil, nil, errorExpectingOneOf("rename_identifier or rename_type", tokens, errors)
}

// "," ( rename_identifier | rename_type )

func commaRename(tokens []tok.Token) (ren ast.Rename, remainder []tok.Token, err error) {
	var found bool
	if remainder, found = Comma(tokens); !found {
		return nil, tokens, ErrNoMatch
	}
	var rename ast.Rename
	rename, remainder, err = Rename(remainder)
	if err != nil {
		return nil, remainder, err
	}
	return rename, remainder, nil
}
