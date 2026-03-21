package ast

import (
	"fmt"
	"strings"
)

// type_constructor_call = type_reference [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

// TypeConstructorCall represents a type constructor call
type TypeConstructorCall struct {
	BaseNode
	TypeReference *TypeReference     // The type being constructed
	Arguments     *FunctionArguments // The constructor arguments
	FunctionBlock *FunctionBlock     // Optional function block (may be nil)
}

// NewTypeConstructorCall creates a new TypeConstructorCall node
func NewTypeConstructorCall(typeRef *TypeReference, arguments *FunctionArguments, functionBlock *FunctionBlock) *TypeConstructorCall {
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
	if t.Arguments != nil {
		builder.WriteString(t.Arguments.String())
	} else {
		builder.WriteString("()")
	}

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

// member_access_tail = "." ( decimal_literal | identifier ) .

type MemberAccessMember interface {
	Node
	memberAccessMemberNode()
}

func (n *Identifier) memberAccessMemberNode()     {}
func (n *IntegerLiteral) memberAccessMemberNode() {}

// MemberAccess represents a member access expression (e.g., obj.field)
type MemberAccess struct {
	BaseNode
	Object Node               // The receiver expression or type identifier
	Member MemberAccessMember // The selected member
}

// NewMemberAccess creates a new MemberAccess node
func NewMemberAccess(object Node, member MemberAccessMember) *MemberAccess {
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

// indexed_access_tail = "[" index "]" .

// IndexedAccess represents an indexed access expression (e.g., arr[idx])
type IndexedAccess struct {
	BaseNode
	Object Expression // The object being indexed
	Index  Expression // The index expression
}

// NewIndexedAccess creates a new IndexedAccess node
func NewIndexedAccess(object Expression, index Expression) *IndexedAccess {
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

// safe_indexed_access_tail = "[" index "]" "!" .

// SafeIndexedAccess represents a safe indexed access expression (e.g., arr[idx]!)
type SafeIndexedAccess struct {
	BaseNode
	Object Expression // The object being indexed
	Index  Expression // The index expression
}

// NewSafeIndexedAccess creates a new SafeIndexedAccess node
func NewSafeIndexedAccess(object Expression, index Expression) *SafeIndexedAccess {
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

// tuple_update_expression = expression "." labeled_tuple_members .

// TupleUpdateExpression represents a tuple update expression (e.g., obj.(field: value))
type TupleUpdateExpression struct {
	BaseNode
	Object Expression    // The expression being updated
	Update *TupleLiteral // The labeled tuple literal with updated fields
}

// NewTupleUpdateExpression creates a new TupleUpdateExpression node
func NewTupleUpdateExpression(object Expression, update *TupleLiteral) *TupleUpdateExpression {
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
