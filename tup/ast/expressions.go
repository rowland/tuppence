package ast

import (
	"strings"
)

// Base type for all expressions
type Expression interface {
	Node
}

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
		BaseNode:      BaseNode{NodeType: NodeFunctionCall},
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

// Children returns the child nodes
func (f *FunctionCall) Children() []Node {
	children := make([]Node, 0, len(f.Arguments)+2)
	children = append(children, f.Function)
	children = append(children, f.Arguments...)
	if f.FunctionBlock != nil {
		children = append(children, f.FunctionBlock)
	}
	return children
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
		BaseNode:  BaseNode{NodeType: NodeUFCSFunctionCall},
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

// Children returns the child nodes
func (u *UFCSFunctionCall) Children() []Node {
	children := make([]Node, 0, len(u.Arguments)+2)
	children = append(children, u.Receiver, u.Function)
	children = append(children, u.Arguments...)
	return children
}

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
		BaseNode:      BaseNode{NodeType: NodeTypeConstructorCall},
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

// Children returns the child nodes
func (t *TypeConstructorCall) Children() []Node {
	children := make([]Node, 0, len(t.Arguments)+2)
	children = append(children, t.TypeReference)
	children = append(children, t.Arguments...)
	if t.FunctionBlock != nil {
		children = append(children, t.FunctionBlock)
	}
	return children
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
		BaseNode:  BaseNode{NodeType: NodeBuiltinFunctionCall},
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

// Children returns the child nodes
func (b *BuiltinFunctionCall) Children() []Node {
	return b.Arguments
}

// ArrayFunctionCall represents a call to the array() function
type ArrayFunctionCall struct {
	BaseNode
	TypeArg Node // The type argument
	SizeArg Node // The size argument (may be nil)
}

// NewArrayFunctionCall creates a new ArrayFunctionCall node
func NewArrayFunctionCall(typeArg Node, sizeArg Node) *ArrayFunctionCall {
	return &ArrayFunctionCall{
		BaseNode: BaseNode{NodeType: NodeArrayFunctionCall},
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

// Children returns the child nodes
func (a *ArrayFunctionCall) Children() []Node {
	if a.SizeArg != nil {
		return []Node{a.TypeArg, a.SizeArg}
	}
	return []Node{a.TypeArg}
}

// MemberAccess represents a member access expression (e.g., obj.field)
type MemberAccess struct {
	BaseNode
	Object Node // The object being accessed
	Member Node // The member being accessed (can be Identifier, FunctionCall, etc.)
}

// NewMemberAccess creates a new MemberAccess node
func NewMemberAccess(object Node, member Node) *MemberAccess {
	return &MemberAccess{
		BaseNode: BaseNode{NodeType: NodeMemberAccess},
		Object:   object,
		Member:   member,
	}
}

// String returns a textual representation of the member access
func (m *MemberAccess) String() string {
	return m.Object.String() + "." + m.Member.String()
}

// Children returns the child nodes
func (m *MemberAccess) Children() []Node {
	return []Node{m.Object, m.Member}
}

// IndexedAccess represents an indexed access expression (e.g., arr[idx])
type IndexedAccess struct {
	BaseNode
	Object Node // The object being indexed
	Index  Node // The index expression
}

// NewIndexedAccess creates a new IndexedAccess node
func NewIndexedAccess(object Node, index Node) *IndexedAccess {
	return &IndexedAccess{
		BaseNode: BaseNode{NodeType: NodeIndexedAccess},
		Object:   object,
		Index:    index,
	}
}

// String returns a textual representation of the indexed access
func (i *IndexedAccess) String() string {
	return i.Object.String() + "[" + i.Index.String() + "]"
}

// Children returns the child nodes
func (i *IndexedAccess) Children() []Node {
	return []Node{i.Object, i.Index}
}

// SafeIndexedAccess represents a safe indexed access expression (e.g., arr[idx]!)
type SafeIndexedAccess struct {
	BaseNode
	Object Node // The object being indexed
	Index  Node // The index expression
}

// NewSafeIndexedAccess creates a new SafeIndexedAccess node
func NewSafeIndexedAccess(object Node, index Node) *SafeIndexedAccess {
	return &SafeIndexedAccess{
		BaseNode: BaseNode{NodeType: NodeSafeIndexedAccess},
		Object:   object,
		Index:    index,
	}
}

// String returns a textual representation of the safe indexed access
func (s *SafeIndexedAccess) String() string {
	return s.Object.String() + "[" + s.Index.String() + "]!"
}

// Children returns the child nodes
func (s *SafeIndexedAccess) Children() []Node {
	return []Node{s.Object, s.Index}
}

// TypeofExpression represents a typeof expression
type TypeofExpression struct {
	BaseNode
	Expression Node // The expression to get the type of
}

// NewTypeofExpression creates a new TypeofExpression node
func NewTypeofExpression(expression Node) *TypeofExpression {
	return &TypeofExpression{
		BaseNode:   BaseNode{NodeType: NodeTypeofExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the typeof expression
func (t *TypeofExpression) String() string {
	return "typeof(" + t.Expression.String() + ")"
}

// Children returns the child nodes
func (t *TypeofExpression) Children() []Node {
	return []Node{t.Expression}
}

// Operator represents an operator in an expression
type Operator string

// BinaryExpression represents a binary expression (e.g., a + b)
type BinaryExpression struct {
	BaseNode
	Left     Node     // The left operand
	Operator Operator // The operator
	Right    Node     // The right operand
}

// NewBinaryExpression creates a new BinaryExpression node
func NewBinaryExpression(left Node, operator Operator, right Node) *BinaryExpression {
	return &BinaryExpression{
		BaseNode: BaseNode{NodeType: NodeBinaryExpression},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// String returns a textual representation of the binary expression
func (b *BinaryExpression) String() string {
	return b.Left.String() + " " + string(b.Operator) + " " + b.Right.String()
}

// Children returns the child nodes
func (b *BinaryExpression) Children() []Node {
	return []Node{b.Left, b.Right}
}

// UnaryExpression represents a unary expression (e.g., -x, !y)
type UnaryExpression struct {
	BaseNode
	Operator   Operator // The operator
	Expression Node     // The operand
}

// NewUnaryExpression creates a new UnaryExpression node
func NewUnaryExpression(operator Operator, expression Node) *UnaryExpression {
	return &UnaryExpression{
		BaseNode:   BaseNode{NodeType: NodeUnaryExpression},
		Operator:   operator,
		Expression: expression,
	}
}

// String returns a textual representation of the unary expression
func (u *UnaryExpression) String() string {
	return string(u.Operator) + u.Expression.String()
}

// Children returns the child nodes
func (u *UnaryExpression) Children() []Node {
	return []Node{u.Expression}
}

// ChainedExpression represents a chained pipe expression (e.g., a |> b |> c)
type ChainedExpression struct {
	BaseNode
	Initial     Node   // The initial expression
	Expressions []Node // The chained expressions (FunctionCall nodes)
}

// NewChainedExpression creates a new ChainedExpression node
func NewChainedExpression(initial Node, expressions []Node) *ChainedExpression {
	return &ChainedExpression{
		BaseNode:    BaseNode{NodeType: NodeChainedExpression},
		Initial:     initial,
		Expressions: expressions,
	}
}

// String returns a textual representation of the chained expression
func (c *ChainedExpression) String() string {
	var builder strings.Builder
	builder.WriteString(c.Initial.String())

	for _, expr := range c.Expressions {
		builder.WriteString(" |> ")
		builder.WriteString(expr.String())
	}

	return builder.String()
}

// Children returns the child nodes
func (c *ChainedExpression) Children() []Node {
	children := make([]Node, 0, len(c.Expressions)+1)
	children = append(children, c.Initial)
	children = append(children, c.Expressions...)
	return children
}

// ReturnExpression represents a return expression
type ReturnExpression struct {
	BaseNode
	Expression Node // The expression being returned
}

// NewReturnExpression creates a new ReturnExpression node
func NewReturnExpression(expression Node) *ReturnExpression {
	return &ReturnExpression{
		BaseNode:   BaseNode{NodeType: NodeReturnExpression},
		Expression: expression,
	}
}

// String returns a textual representation of the return expression
func (r *ReturnExpression) String() string {
	return "return " + r.Expression.String()
}

// Children returns the child nodes
func (r *ReturnExpression) Children() []Node {
	return []Node{r.Expression}
}

// BreakExpression represents a break expression
type BreakExpression struct {
	BaseNode
	Expression Node // Optional expression (may be nil)
}

// NewBreakExpression creates a new BreakExpression node
func NewBreakExpression(expression Node) *BreakExpression {
	return &BreakExpression{
		BaseNode:   BaseNode{NodeType: NodeBreakExpression},
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

// Children returns the child nodes
func (b *BreakExpression) Children() []Node {
	if b.Expression != nil {
		return []Node{b.Expression}
	}
	return nil
}

// ContinueExpression represents a continue expression
type ContinueExpression struct {
	BaseNode
	Expression Node // Optional expression (may be nil)
}

// NewContinueExpression creates a new ContinueExpression node
func NewContinueExpression(expression Node) *ContinueExpression {
	return &ContinueExpression{
		BaseNode:   BaseNode{NodeType: NodeContinueExpression},
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

// Children returns the child nodes
func (c *ContinueExpression) Children() []Node {
	if c.Expression != nil {
		return []Node{c.Expression}
	}
	return nil
}

// TryVariant represents the different variants of try expressions
type TryVariant string

const (
	TryStandard TryVariant = "try"
	TryContinue TryVariant = "try_continue"
	TryBreak    TryVariant = "try_break"
)

// TryExpression represents a try expression
type TryExpression struct {
	BaseNode
	Variant    TryVariant // The variant (try, try_continue, try_break)
	Expression Node       // The expression being tried (may be nil for try_continue and try_break)
}

// NewTryExpression creates a new TryExpression node
func NewTryExpression(variant TryVariant, expression Node) *TryExpression {
	return &TryExpression{
		BaseNode:   BaseNode{NodeType: NodeTryExpression},
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

// Children returns the child nodes
func (t *TryExpression) Children() []Node {
	if t.Expression != nil {
		return []Node{t.Expression}
	}
	return nil
}

// TypeComparison represents a type comparison expression (e.g., x is Type)
type TypeComparison struct {
	BaseNode
	Expression Node // The expression being tested
	Type       Node // The type being tested against
}

// NewTypeComparison creates a new TypeComparison node
func NewTypeComparison(expression Node, typeNode Node) *TypeComparison {
	return &TypeComparison{
		BaseNode:   BaseNode{NodeType: NodeTypeComparison},
		Expression: expression,
		Type:       typeNode,
	}
}

// String returns a textual representation of the type comparison
func (t *TypeComparison) String() string {
	return t.Expression.String() + " is " + t.Type.String()
}

// Children returns the child nodes
func (t *TypeComparison) Children() []Node {
	return []Node{t.Expression, t.Type}
}

// RelationalComparison represents a relational comparison expression (e.g., a < b)
type RelationalComparison struct {
	BaseNode
	Left     Node     // The left operand
	Operator Operator // The comparison operator
	Right    Node     // The right operand
}

// NewRelationalComparison creates a new RelationalComparison node
func NewRelationalComparison(left Node, operator Operator, right Node) *RelationalComparison {
	return &RelationalComparison{
		BaseNode: BaseNode{NodeType: NodeRelationalComparison},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// String returns a textual representation of the relational comparison
func (r *RelationalComparison) String() string {
	return r.Left.String() + " " + string(r.Operator) + " " + r.Right.String()
}

// Children returns the child nodes
func (r *RelationalComparison) Children() []Node {
	return []Node{r.Left, r.Right}
}

// MetaExpression represents a compile-time meta expression (e.g., $(key: value))
type MetaExpression struct {
	BaseNode
	KeyValues map[string]Node // The key-value pairs
}

// NewMetaExpression creates a new MetaExpression node
func NewMetaExpression(keyValues map[string]Node) *MetaExpression {
	return &MetaExpression{
		BaseNode:  BaseNode{NodeType: NodeMetaExpression},
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

// Children returns the child nodes
func (m *MetaExpression) Children() []Node {
	children := make([]Node, 0, len(m.KeyValues))
	for _, value := range m.KeyValues {
		children = append(children, value)
	}
	return children
}

// Constant represents a constant value reference (e.g., module-level constants)
type Constant struct {
	BaseNode
	Identifier Node // The identifier of the constant
}

// NewConstant creates a new Constant node
func NewConstant(identifier Node) *Constant {
	return &Constant{
		BaseNode:   BaseNode{NodeType: NodeConstant},
		Identifier: identifier,
	}
}

// String returns a textual representation of the constant
func (c *Constant) String() string {
	return c.Identifier.String()
}

// Children returns the child nodes
func (c *Constant) Children() []Node {
	return []Node{c.Identifier}
}

// TupleUpdateExpression represents a tuple update expression (e.g., obj.(field: value))
type TupleUpdateExpression struct {
	BaseNode
	Object Node // The object being updated
	Update Node // The tuple literal with updated fields
}

// NewTupleUpdateExpression creates a new TupleUpdateExpression node
func NewTupleUpdateExpression(object Node, update Node) *TupleUpdateExpression {
	return &TupleUpdateExpression{
		BaseNode: BaseNode{NodeType: NodeTupleUpdateExpression},
		Object:   object,
		Update:   update,
	}
}

// String returns a textual representation of the tuple update expression
func (t *TupleUpdateExpression) String() string {
	return t.Object.String() + "." + t.Update.String()
}

// Children returns the child nodes
func (t *TupleUpdateExpression) Children() []Node {
	return []Node{t.Object, t.Update}
}
