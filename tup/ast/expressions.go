package ast

import (
	"fmt"
	"strings"
)

// expression = try_expression
//            | binary_expression
//            | unary_expression .

type Expression interface {
	Node
	expressionNode()
}

func (n *TryExpression) expressionNode()    {}
func (n *BinaryExpression) expressionNode() {}
func (n *UnaryExpression) expressionNode()  {}

func (n *Prec1Expression) expressionNode() {}
func (n *Prec2Expression) expressionNode() {}
func (n *Prec4Expression) expressionNode() {}
func (n *Prec5Expression) expressionNode() {}
func (n *Prec6Expression) expressionNode() {}

func (n *TypeComparison) expressionNode()        {}
func (n *RelationalComparison) expressionNode()  {}
func (n *Identifier) expressionNode()            {}
func (n *Block) expressionNode()                 {}
func (n *IfExpression) expressionNode()          {}
func (n *ForExpression) expressionNode()         {}
func (n *InlineForExpression) expressionNode()   {}
func (n *ArrayFunctionCall) expressionNode()     {}
func (n *TypeofExpression) expressionNode()      {}
func (n *FunctionCall) expressionNode()          {}
func (n *TypeConstructorCall) expressionNode()   {}
func (n *MemberAccess) expressionNode()          {}
func (n *TupleUpdateExpression) expressionNode() {}
func (n *SafeIndexedAccess) expressionNode()     {}
func (n *IndexedAccess) expressionNode()         {}

// func (n *ImportExpression) expressionNode() {}
// func (n *ReturnExpression) expressionNode()      {}
// func (n *BreakExpression) expressionNode()       {}
// func (n *ContinueExpression) expressionNode()    {}

func (n *FloatLiteral) expressionNode()              {}
func (n *IntegerLiteral) expressionNode()            {}
func (n *BooleanLiteral) expressionNode()            {}
func (n *StringLiteral) expressionNode()             {}
func (n *InterpolatedStringLiteral) expressionNode() {}
func (n *RawStringLiteral) expressionNode()          {}
func (n *MultiLineStringLiteral) expressionNode()    {}
func (n *TupleLiteral) expressionNode()              {}
func (n *ArrayLiteral) expressionNode()              {}
func (n *SymbolLiteral) expressionNode()             {}
func (n *RuneLiteral) expressionNode()               {}
func (n *FixedSizeArrayLiteral) expressionNode()     {}

// prec1_expression = prec2_expression { logical_or_op prec2_expression } .

type Prec1Expression struct {
	BaseNode
	Operands []Expression
}

func NewPrec1Expression(operands []Expression) *Prec1Expression {
	return &Prec1Expression{
		BaseNode: BaseNode{Type: NodePrec1Expression},
		Operands: operands,
	}
}

// String returns a textual representation of the prec1 expression
func (p *Prec1Expression) String() string {
	var builder strings.Builder
	for i, operand := range p.Operands {
		if i > 0 {
			builder.WriteString(" || ")
		}
		builder.WriteString(operand.String())
	}
	return builder.String()
}

// prec2_expression = prec3_expression { logical_and_op prec3_expression } .

type Prec2Expression struct {
	BaseNode
	Operands []Expression
}

// NewPrec2Expression creates a new Prec2Expression node
func NewPrec2Expression(operands []Expression) *Prec2Expression {
	return &Prec2Expression{
		BaseNode: BaseNode{Type: NodePrec2Expression},
		Operands: operands,
	}
}

// String returns a textual representation of the prec2 expression
func (p *Prec2Expression) String() string {
	var builder strings.Builder
	for i, operand := range p.Operands {
		if i > 0 {
			builder.WriteString(" && ")
		}
		builder.WriteString(operand.String())
	}
	return builder.String()
}

// prec3_expression = type_comparison | relational_comparison .

type Prec3Expression interface {
	Node
	prec3ExpressionNode()
}

func (n *TypeComparison) prec3ExpressionNode()       {}
func (n *RelationalComparison) prec3ExpressionNode() {}

// prec4_expression = prec5_expression { add_sub_op prec5_expression } .

type Prec4Expression struct {
	BaseNode
	Left     Expression
	Operator AddSubOp
	Right    Expression
}

func NewPrec4Expression(left Expression, operator AddSubOp, right Expression) *Prec4Expression {
	return &Prec4Expression{
		BaseNode: BaseNode{Type: NodePrec4Expression},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// String returns a textual representation of the prec4 expression
func (p *Prec4Expression) String() string {
	return fmt.Sprintf("%s %s %s", p.Left, p.Operator, p.Right)
}

// prec5_expression = prec6_expression { mul_div_op prec6_expression } .

type Prec5Expression struct {
	BaseNode
	Left     Expression
	Operator MulDivOp
	Right    Expression
}

func NewPrec5Expression(left Expression, operator MulDivOp, right Expression) *Prec5Expression {
	return &Prec5Expression{
		BaseNode: BaseNode{Type: NodePrec5Expression},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (p *Prec5Expression) String() string {
	return fmt.Sprintf("%s %s %s", p.Left, p.Operator, p.Right)
}

// prec6_expression = unary_expression { "^" unary_expression } .

type Prec6Expression struct {
	BaseNode
	Operands []Expression
}

func NewPrec6Expression(operands []Expression) *Prec6Expression {
	return &Prec6Expression{
		BaseNode: BaseNode{Type: NodePrec6Expression},
		Operands: operands,
	}
}

func (p *Prec6Expression) String() string {
	var builder strings.Builder
	for i, operand := range p.Operands {
		if i > 0 {
			builder.WriteString(" ^ ")
		}
		builder.WriteString(operand.String())
	}
	return builder.String()
}

// type_predicate = type_reference | inline_union .

type TypePredicate interface {
	Node
	typePredicateNode()
}

func (n *TypeReference) typePredicateNode() {}
func (n *InlineUnion) typePredicateNode()   {}

// type_comparison = prec4_expression is_op type_predicate .

type TypeComparison struct {
	BaseNode
	Left  Expression
	Right TypePredicate
}

func NewTypeComparison(left Expression, right TypePredicate) *TypeComparison {
	return &TypeComparison{
		BaseNode: BaseNode{Type: NodeTypeComparison},
		Left:     left,
		Right:    right,
	}
}

func (t *TypeComparison) String() string {
	return fmt.Sprintf("%s is %s", t.Left, t.Right)
}

// relational_comparison = prec4_expression { rel_op prec4_expression } .
type RelationalComparison struct {
	BaseNode
	Left     Expression
	Operator RelOp
	Right    Expression
}

func NewRelationalComparison(left Expression, operator RelOp, right Expression) *RelationalComparison {
	return &RelationalComparison{
		BaseNode: BaseNode{Type: NodeRelationalComparison},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (r *RelationalComparison) String() string {
	return fmt.Sprintf("%s %s %s", r.Left, r.Operator, r.Right)
}

// function_call = function_identifier [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

// FunctionCall represents a function call expression
type FunctionCall struct {
	BaseNode
	Function      Node   // The function being called (can be Identifier, MemberAccess, etc.)
	Arguments     []Node // The arguments passed to the function
	FunctionBlock Node   // Optional function block (for higher-order functions, may be nil)
}

// NewFunctionCall creates a new FunctionCall node
func NewFunctionCall(function Node, arguments []Node, functionBlock Node) *FunctionCall {
	return &FunctionCall{
		BaseNode:      BaseNode{Type: NodeFunctionCall},
		Function:      function,
		Arguments:     arguments,
		FunctionBlock: functionBlock,
	}
}

// String returns a textual representation of the function call
func (f *FunctionCall) String() string {
	var builder strings.Builder
	builder.WriteString(f.Function.String())
	builder.WriteString("(")

	for i, arg := range f.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(arg.String())
	}

	builder.WriteString(")")

	if f.FunctionBlock != nil {
		builder.WriteString(" ")
		builder.WriteString(f.FunctionBlock.String())
	}

	return builder.String()
}

// UFCSFunctionCall represents a Uniform Function Call Syntax function call
type UFCSFunctionCall struct {
	BaseNode
	Receiver  Node   // The receiver object
	Function  Node   // The function being called (typically an Identifier)
	Arguments []Node // The arguments passed to the function (excluding the receiver)
}

// NewUFCSFunctionCall creates a new UFCSFunctionCall node
func NewUFCSFunctionCall(receiver Node, function Node, arguments []Node) *UFCSFunctionCall {
	return &UFCSFunctionCall{
		BaseNode:  BaseNode{Type: NodeUFCSFunctionCall},
		Receiver:  receiver,
		Function:  function,
		Arguments: arguments,
	}
}

// String returns a textual representation of the UFCS function call
func (u *UFCSFunctionCall) String() string {
	var builder strings.Builder
	builder.WriteString(u.Receiver.String())
	builder.WriteString(".")
	builder.WriteString(u.Function.String())
	builder.WriteString("(")

	for i, arg := range u.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(arg.String())
	}

	builder.WriteString(")")
	return builder.String()
}

// type_constructor_call = type_reference [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

// TypeConstructorCall represents a type constructor call
type TypeConstructorCall struct {
	BaseNode
	TypeReference Node   // The type being constructed
	Arguments     []Node // The constructor arguments
	FunctionBlock Node   // Optional function block (may be nil)
}

// NewTypeConstructorCall creates a new TypeConstructorCall node
func NewTypeConstructorCall(typeRef Node, arguments []Node, functionBlock Node) *TypeConstructorCall {
	return &TypeConstructorCall{
		BaseNode:      BaseNode{Type: NodeTypeConstructorCall},
		TypeReference: typeRef,
		Arguments:     arguments,
		FunctionBlock: functionBlock,
	}
}

// String returns a textual representation of the type constructor call
func (t *TypeConstructorCall) String() string {
	var builder strings.Builder
	builder.WriteString(t.TypeReference.String())
	builder.WriteString("(")

	for i, arg := range t.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(arg.String())
	}

	builder.WriteString(")")

	if t.FunctionBlock != nil {
		builder.WriteString(" ")
		builder.WriteString(t.FunctionBlock.String())
	}

	return builder.String()
}

// BuiltinFunctionCall represents a call to a built-in function
type BuiltinFunctionCall struct {
	BaseNode
	Name      string // Name of the builtin function
	Arguments []Node // The arguments passed to the function
}

// NewBuiltinFunctionCall creates a new BuiltinFunctionCall node
func NewBuiltinFunctionCall(name string, arguments []Node) *BuiltinFunctionCall {
	return &BuiltinFunctionCall{
		BaseNode:  BaseNode{Type: NodeBuiltinFunctionCall},
		Name:      name,
		Arguments: arguments,
	}
}

// String returns a textual representation of the builtin function call
func (b *BuiltinFunctionCall) String() string {
	var builder strings.Builder
	builder.WriteString(b.Name)
	builder.WriteString("(")

	for i, arg := range b.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(arg.String())
	}

	builder.WriteString(")")
	return builder.String()
}

// array_function_call = "array" "(" type_identifier "," expression ")" .

// ArrayFunctionCall represents a call to the array() function
type ArrayFunctionCall struct {
	BaseNode
	TypeArg Node // The type argument
	SizeArg Node // The size argument (may be nil)
}

// NewArrayFunctionCall creates a new ArrayFunctionCall node
func NewArrayFunctionCall(typeArg Node, sizeArg Node) *ArrayFunctionCall {
	return &ArrayFunctionCall{
		BaseNode: BaseNode{Type: NodeArrayFunctionCall},
		TypeArg:  typeArg,
		SizeArg:  sizeArg,
	}
}

// String returns a textual representation of the array function call
func (a *ArrayFunctionCall) String() string {
	var builder strings.Builder
	builder.WriteString("array(")
	builder.WriteString(a.TypeArg.String())

	if a.SizeArg != nil {
		builder.WriteString(", ")
		builder.WriteString(a.SizeArg.String())
	}

	builder.WriteString(")")
	return builder.String()
}

// member_access = ( expression | type_identifier ) "." ( decimal_literal
// 	                                                 | identifier
// 	                                                 | function_call ) .

// MemberAccess represents a member access expression (e.g., obj.field)
type MemberAccess struct {
	BaseNode
	Object Node // The object being accessed
	Member Node // The member being accessed (can be Identifier, FunctionCall, etc.)
}

// NewMemberAccess creates a new MemberAccess node
func NewMemberAccess(object Node, member Node) *MemberAccess {
	return &MemberAccess{
		BaseNode: BaseNode{Type: NodeMemberAccess},
		Object:   object,
		Member:   member,
	}
}

// String returns a textual representation of the member access
func (m *MemberAccess) String() string {
	return fmt.Sprintf("%s.%s", m.Object, m.Member)
}

// indexed_access = expression "[" index "]" .

// IndexedAccess represents an indexed access expression (e.g., arr[idx])
type IndexedAccess struct {
	BaseNode
	Object Node // The object being indexed
	Index  Node // The index expression
}

// NewIndexedAccess creates a new IndexedAccess node
func NewIndexedAccess(object Node, index Node) *IndexedAccess {
	return &IndexedAccess{
		BaseNode: BaseNode{Type: NodeIndexedAccess},
		Object:   object,
		Index:    index,
	}
}

// String returns a textual representation of the indexed access
func (i *IndexedAccess) String() string {
	return i.Object.String() + "[" + i.Index.String() + "]"
}

// safe_indexed_access = expression "[" index "]" "!" .

// SafeIndexedAccess represents a safe indexed access expression (e.g., arr[idx]!)
type SafeIndexedAccess struct {
	BaseNode
	Object Node // The object being indexed
	Index  Node // The index expression
}

// NewSafeIndexedAccess creates a new SafeIndexedAccess node
func NewSafeIndexedAccess(object Node, index Node) *SafeIndexedAccess {
	return &SafeIndexedAccess{
		BaseNode: BaseNode{Type: NodeSafeIndexedAccess},
		Object:   object,
		Index:    index,
	}
}

// String returns a textual representation of the safe indexed access
func (s *SafeIndexedAccess) String() string {
	return s.Object.String() + "[" + s.Index.String() + "]!"
}

// typeof_expression = "typeof" "(" expression ")" .

// TypeofExpression represents a typeof expression
type TypeofExpression struct {
	BaseNode
	Expression Node // The expression to get the type of
}

// NewTypeofExpression creates a new TypeofExpression node
func NewTypeofExpression(expression Node) *TypeofExpression {
	return &TypeofExpression{
		BaseNode:   BaseNode{Type: NodeTypeofExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the typeof expression
func (t *TypeofExpression) String() string {
	return "typeof(" + t.Expression.String() + ")"
}

// binary_expression = chained_expression .

type BinaryExpression = ChainedExpression

// unary_expression = prefixed_unary_expression
//                  | primary_expression .

type UnaryExpression struct {
	BaseNode
	Operator   UnaryOp
	Expression Expression
}

func NewUnaryExpression(operator UnaryOp, expression Expression) *UnaryExpression {
	return &UnaryExpression{
		BaseNode:   BaseNode{Type: NodeUnaryExpression},
		Operator:   operator,
		Expression: expression,
	}
}

func (u *UnaryExpression) String() string {
	return u.Operator.String() + u.Expression.String()
}

// chained_expression = prec1_expression { "|>" function_call } .

// ChainedExpression represents a chained pipe expression (e.g., a |> b |> c)
type ChainedExpression struct {
	BaseNode
	Initial       Node            // The initial expression
	FunctionCalls []*FunctionCall // The chained expressions (FunctionCall nodes)
}

// NewChainedExpression creates a new ChainedExpression node
func NewChainedExpression(initial Node, functionCalls []*FunctionCall) *ChainedExpression {
	return &ChainedExpression{
		BaseNode:      BaseNode{Type: NodeChainedExpression},
		Initial:       initial,
		FunctionCalls: functionCalls,
	}
}

// String returns a textual representation of the chained expression
func (c *ChainedExpression) String() string {
	var builder strings.Builder
	builder.WriteString(c.Initial.String())

	for _, functionCall := range c.FunctionCalls {
		builder.WriteString(" |> ")
		builder.WriteString(functionCall.String())
	}

	return builder.String()
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

// tuple_update_expression = expression "." tuple_literal .

// TupleUpdateExpression represents a tuple update expression (e.g., obj.(field: value))
type TupleUpdateExpression struct {
	BaseNode
	Object Node // The object being updated
	Update Node // The tuple literal with updated fields
}

// NewTupleUpdateExpression creates a new TupleUpdateExpression node
func NewTupleUpdateExpression(object Node, update Node) *TupleUpdateExpression {
	return &TupleUpdateExpression{
		BaseNode: BaseNode{Type: NodeTupleUpdateExpression},
		Object:   object,
		Update:   update,
	}
}

// String returns a textual representation of the tuple update expression
func (t *TupleUpdateExpression) String() string {
	return t.Object.String() + "." + t.Update.String()
}
