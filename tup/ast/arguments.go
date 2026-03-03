package ast

import "strings"

// argument = ( expression | spread_argument ) .
// spread_argument = "..." expression .

type Argument struct {
	BaseNode
	Expr   Expression
	Spread bool // True if the argument is a spread argument
}

func NewArgument(expr Expression, spread bool) *Argument {
	return &Argument{
		BaseNode: BaseNode{Type: NodeArgument},
		Expr:     expr,
		Spread:   spread,
	}
}

func (a *Argument) String() string {
	if a.Spread {
		return "..." + a.Expr.String()
	}
	return a.Expr.String()
}

// x argument { "," argument } .

type Arguments struct {
	BaseNode
	Args []*Argument
}

func NewArguments(args []*Argument) *Arguments {
	return &Arguments{
		BaseNode: BaseNode{Type: NodeArguments},
		Args:     args,
	}
}

func (a *Arguments) String() string {
	var builder strings.Builder
	for i, arg := range a.Args {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(arg.String())
	}
	return builder.String()
}

// labeled_argument = ( identifier ":" argument ) .

type LabeledArgument struct {
	BaseNode
	Identifier *Identifier // The argument identifier
	Argument   *Argument   // The argument value
}

func NewLabeledArgument(identifier *Identifier, argument *Argument) *LabeledArgument {
	return &LabeledArgument{
		BaseNode:   BaseNode{Type: NodeLabeledArgument},
		Identifier: identifier,
		Argument:   argument,
	}
}

func (l *LabeledArgument) String() string {
	return l.Identifier.String() + ": " + l.Argument.String()
}

// labeled_arguments = labeled_argument { "," ( labeled_argument ) } .

type LabeledArguments struct {
	BaseNode
	Args []*LabeledArgument
}

func NewLabeledArguments(args []*LabeledArgument) *LabeledArguments {
	return &LabeledArguments{
		BaseNode: BaseNode{Type: NodeLabeledArguments},
		Args:     args,
	}
}

func (l *LabeledArguments) String() string {
	var builder strings.Builder
	for i, arg := range l.Args {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(arg.String())
	}
	return builder.String()
}
