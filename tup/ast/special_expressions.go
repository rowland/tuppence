package ast

import "strings"

// typeof_expression = "typeof" "(" expression ")" .

// TypeofExpression represents a typeof expression
type TypeofExpression struct {
	BaseNode
	Expression Expression // The expression to get the type of
}

// NewTypeofExpression creates a new TypeofExpression node
func NewTypeofExpression(expression Expression) *TypeofExpression {
	return &TypeofExpression{
		BaseNode:   BaseNode{Type: NodeTypeofExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the typeof expression
func (t *TypeofExpression) String() string {
	return "typeof(" + t.Expression.String() + ")"
}

// import_expression = "import" "(" string_literal ")" .

// ImportExpression represents an import expression
type ImportExpression struct {
	BaseNode
	Path *StringLiteral
}

// NewImportExpression creates a new ImportExpression node
func NewImportExpression(path *StringLiteral) *ImportExpression {
	return &ImportExpression{
		BaseNode: BaseNode{Type: NodeImportExpression},
		Path:     path,
	}
}

// String returns a textual representation of the import expression
func (i *ImportExpression) String() string {
	return "import(" + i.Path.String() + ")"
}

// return_expression = "return" expression .

// ReturnExpression represents a return expression
type ReturnExpression struct {
	BaseNode
	Expression Node // The expression being returned
}

// NewReturnExpression creates a new ReturnExpression node
func NewReturnExpression(expression Node) *ReturnExpression {
	return &ReturnExpression{
		BaseNode:   BaseNode{Type: NodeReturnExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the return expression
func (r *ReturnExpression) String() string {
	return "return " + r.Expression.String()
}

// break_expression = "break" [ expression ] .

// BreakExpression represents a break expression
type BreakExpression struct {
	BaseNode
	Expression Node // Optional expression (may be nil)
}

// NewBreakExpression creates a new BreakExpression node
func NewBreakExpression(expression Node) *BreakExpression {
	return &BreakExpression{
		BaseNode:   BaseNode{Type: NodeBreakExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the break expression
func (b *BreakExpression) String() string {
	if b.Expression != nil {
		return "break " + b.Expression.String()
	}
	return "break"
}

// continue_expression = "continue" [ expression ] .

// ContinueExpression represents a continue expression
type ContinueExpression struct {
	BaseNode
	Expression Node // Optional expression (may be nil)
}

// NewContinueExpression creates a new ContinueExpression node
func NewContinueExpression(expression Node) *ContinueExpression {
	return &ContinueExpression{
		BaseNode:   BaseNode{Type: NodeContinueExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the continue expression
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

// TryExpression represents a try expression
type TryExpression struct {
	BaseNode
	Variant    TryVariant // The variant (try, try_continue, try_break)
	Expression Node       // The expression being tried (may be nil for try_continue and try_break)
}

// NewTryExpression creates a new TryExpression node
func NewTryExpression(variant TryVariant, expression Node) *TryExpression {
	return &TryExpression{
		BaseNode:   BaseNode{Type: NodeTryExpression},
		Variant:    variant,
		Expression: expression,
	}
}

// String returns a textual representation of the try expression
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

// NewMetaExpression creates a new MetaExpression node
func NewMetaExpression(keyValues map[string]Node) *MetaExpression {
	return &MetaExpression{
		BaseNode:  BaseNode{Type: NodeMetaExpression},
		KeyValues: keyValues,
	}
}

// String returns a textual representation of the meta expression
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
//          | scoped_identifier
//          | identifier .

// Constant represents a constant value reference (e.g., module-level constants)
type Constant struct {
	BaseNode
	Identifier Node // The identifier of the constant
}

// NewConstant creates a new Constant node
func NewConstant(identifier Node) *Constant {
	return &Constant{
		BaseNode:   BaseNode{Type: NodeConstant},
		Identifier: identifier,
	}
}

// String returns a textual representation of the constant
func (c *Constant) String() string {
	return c.Identifier.String()
}
