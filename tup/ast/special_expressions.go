package ast

import "strings"

// typeof_expression = "typeof" "(" expression ")" .

type TypeofExpression struct {
	BaseNode
	Expression Expression // The expression to get the type of
}

func NewTypeofExpression(expression Expression) *TypeofExpression {
	return &TypeofExpression{
		BaseNode:   BaseNode{Type: NodeTypeofExpression},
		Expression: expression,
	}
}

func (t *TypeofExpression) String() string {
	return "typeof(" + t.Expression.String() + ")"
}

// import_expression = "import" "(" string_literal ")" .

type ImportExpression struct {
	BaseNode
	Path *StringLiteral
}

func NewImportExpression(path *StringLiteral) *ImportExpression {
	return &ImportExpression{
		BaseNode: BaseNode{Type: NodeImportExpression},
		Path:     path,
	}
}

func (i *ImportExpression) String() string {
	return "import(" + i.Path.String() + ")"
}

// return_expression = "return" [ expression ] .

type ReturnExpression struct {
	BaseNode
	Expression Expression // Optional expression (may be nil)
}

func NewReturnExpression(expression Expression) *ReturnExpression {
	return &ReturnExpression{
		BaseNode:   BaseNode{Type: NodeReturnExpression},
		Expression: expression,
	}
}

func (r *ReturnExpression) String() string {
	if r.Expression != nil {
		return "return " + r.Expression.String()
	}
	return "return"
}

// break_expression = "break" [ expression ] .

type BreakExpression struct {
	BaseNode
	Expression Expression // Optional expression (may be nil)
}

func NewBreakExpression(expression Expression) *BreakExpression {
	return &BreakExpression{
		BaseNode:   BaseNode{Type: NodeBreakExpression},
		Expression: expression,
	}
}

func (b *BreakExpression) String() string {
	if b.Expression != nil {
		return "break " + b.Expression.String()
	}
	return "break"
}

// continue_expression = "continue" [ expression ] .

type ContinueExpression struct {
	BaseNode
	Expression Node // Optional expression (may be nil)
}

func NewContinueExpression(expression Node) *ContinueExpression {
	return &ContinueExpression{
		BaseNode:   BaseNode{Type: NodeContinueExpression},
		Expression: expression,
	}
}

func (c *ContinueExpression) String() string {
	if c.Expression != nil {
		return "continue " + c.Expression.String()
	}
	return "continue"
}

// TryVariant represents the different variants of try expressions
type TryVariant string

const (
	TryStandard TryVariant = "try"
	TryContinue TryVariant = "try_continue"
	TryBreak    TryVariant = "try_break"
)

// try_expression = "try" expression
//                | "try_continue" [ expression ]
//                | "try_break" [ expression ] .

type TryExpression struct {
	BaseNode
	Variant    TryVariant // The variant (try, try_continue, try_break)
	Expression Node       // The expression being tried (may be nil for try_continue and try_break)
}

func NewTryExpression(variant TryVariant, expression Node) *TryExpression {
	return &TryExpression{
		BaseNode:   BaseNode{Type: NodeTryExpression},
		Variant:    variant,
		Expression: expression,
	}
}

func (t *TryExpression) String() string {
	if t.Expression != nil {
		return string(t.Variant) + " " + t.Expression.String()
	}
	return string(t.Variant)
}

// meta_expression = "$" labeled_tuple .

// MetaExpression represents a compile-time meta expression (e.g., $(key: value))
type MetaExpression struct {
	BaseNode
	KeyValues map[string]Node // The key-value pairs
}

func NewMetaExpression(keyValues map[string]Node) *MetaExpression {
	return &MetaExpression{
		BaseNode:  BaseNode{Type: NodeMetaExpression},
		KeyValues: keyValues,
	}
}

func (m *MetaExpression) String() string {
	var builder strings.Builder
	builder.WriteString("$(")

	i := 0
	for key, value := range m.KeyValues {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(key)
		builder.WriteString(": ")
		builder.WriteString(value.String())
		i++
	}

	builder.WriteString(")")
	return builder.String()
}

// constant = literal
//          | scoped_identifier .

type ConstantValue interface {
	Node
	constantValueNode()
}

func (n *FloatLiteral) constantValueNode()              {}
func (n *IntegerLiteral) constantValueNode()            {}
func (n *BooleanLiteral) constantValueNode()            {}
func (n *StringLiteral) constantValueNode()             {}
func (n *InterpolatedStringLiteral) constantValueNode() {}
func (n *RawStringLiteral) constantValueNode()          {}
func (n *MultiLineStringLiteral) constantValueNode()    {}
func (n *TupleLiteral) constantValueNode()              {}
func (n *ArrayLiteral) constantValueNode()              {}
func (n *SymbolLiteral) constantValueNode()             {}
func (n *RuneLiteral) constantValueNode()               {}
func (n *ScopedIdentifier) constantValueNode()          {}

// Constant represents a constant value reference (e.g., module-level constants)
type Constant struct {
	BaseNode
	Value ConstantValue
}

func NewConstant(value ConstantValue) *Constant {
	return &Constant{
		BaseNode: BaseNode{Type: NodeConstant},
		Value:    value,
	}
}

func (c *Constant) String() string {
	return c.Value.String()
}
