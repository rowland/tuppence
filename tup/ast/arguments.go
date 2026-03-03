package ast

import "strings"

// argument = ( expression | spread_argument ) .
// spread_argument = "..." expression .

type Argument struct {
	BaseNode
	Expr   Expression // The expression being spread
	Spread bool       // True if the argument is a spread argument
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
