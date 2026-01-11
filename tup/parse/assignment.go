package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// assignment = assignment_lhs "=" [ "mut" ] expression .

func Assignment(tokens []tok.Token) (*ast.Assignment, []tok.Token, error) {
	left, remainder, err := AssignmentLHS(tokens)
	if err != nil {
		return nil, nil, err
	}

	if peek(remainder).Type != tok.TokOpAssign {
		return nil, nil, errorExpecting("=", remainder)
	}
	remainder = remainder[1:]

	mut := false
	if peek(remainder).Type == tok.TokKwMut {
		mut = true
		remainder = remainder[1:]
	}

	right, remainder, err := Expression(remainder)
	if err != nil {
		return nil, nil, err
	}

	return ast.NewAssignment(left, mut, right), remainder, nil
}

// assignment_lhs = ordinal_assignment_lhs
//                | "(" labeled_assignment_lhs ")" .

func AssignmentLHS(tokens []tok.Token) (item ast.AssignmentLHS, remainder []tok.Token, err error) {
	var errors []error
	ordinalLHS, remainder, err := ordinalAssignmentLHS(tokens)
	if err != nil {
		errors = append(errors, err)
	} else if ordinalLHS != nil {
		return ordinalLHS, remainder, nil
	}

	labeledLHS, remainder, err := labeledAssignmentLHS(tokens)
	if err != nil {
		errors = append(errors, err)
	} else if labeledLHS != nil {
		return labeledLHS, remainder, nil
	}

	return nil, nil, errorExpectingOneOf("assignment lhs", tokens, errors)
}

// ordinal_assignment_lhs = identifier { "," identifier } [ "," rest_operator ] .

func ordinalAssignmentLHS(tokens []tok.Token) (item *ast.OrdinalAssignmentLHS, remainder []tok.Token, err error) {
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

	return ast.NewOrdinalAssignmentLHS(identifiers, restOperator), remainder, nil
}

func commaIdentifier(tokens []tok.Token) (item *ast.Identifier, remainder []tok.Token, err error) {
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

func commaRestOperator(tokens []tok.Token) (item *ast.RestOperator, remainder []tok.Token, err error) {
	if peek(tokens).Type != tok.TokComma {
		return nil, nil, errorExpecting(",", tokens)
	}
	remainder = tokens[1:]
	restOperator, remainder, err := RestOperator(remainder)
	if restOperator == nil || err != nil {
		return nil, nil, err
	}
	return restOperator, remainder, nil
}

// rest_operator = "..." [ identifier ] .

func RestOperator(tokens []tok.Token) (item *ast.RestOperator, remainder []tok.Token, err error) {
	if peek(tokens).Type != tok.TokOpRest {
		return nil, nil, errorExpecting("...", tokens)
	}
	remainder = tokens[1:]
	identifier, remainder, err := Identifier(remainder)
	if identifier == nil || err != nil {
		return nil, nil, err
	}
	return ast.NewRestOperator(identifier), remainder, nil
}

// labeled_assignment_lhs = ( rename_identifier | rename_type ) { "," ( rename_identifier | rename_type ) } .

func labeledAssignmentLHS(tokens []tok.Token) (item *ast.LabeledAssignmentLHS, remainder []tok.Token, err error) {
	var renames []ast.Rename

	rename, remainder, err := Rename(tokens)
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

	return ast.NewLabeledAssignmentLHS(renames), remainder, nil
}

// rename_identifier | rename_type

func Rename(tokens []tok.Token) (item ast.Rename, remainder []tok.Token, err error) {
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

func commaRename(tokens []tok.Token) (item ast.Rename, remainder []tok.Token, err error) {
	if peek(tokens).Type != tok.TokComma {
		return nil, nil, errorExpecting(",", tokens)
	}
	remainder = tokens[1:]
	rename, remainder, err := Rename(remainder)
	if rename == nil || err != nil {
		return nil, nil, err
	}
	return rename, remainder, nil
}
