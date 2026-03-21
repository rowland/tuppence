package ast

import "strings"

// parameter = annotations ( nilable_type
//	                       | type
//	                       | literal
//	                       | union_type
//	                       | union_declaration ) .

type FunctionTypeParameter interface {
	Node
	functionTypeParameterNode()
}

func (n *Parameter) functionTypeParameterNode()        {}
func (n *LabeledParameter) functionTypeParameterNode() {}
func (n *RestParameter) functionTypeParameterNode()    {}
func (n *LabeledRestParameter) functionTypeParameterNode() {
}

type FunctionTypeParameterType interface {
	Node
	functionTypeParameterTypeNode()
}

func (n *NilableType) functionTypeParameterTypeNode()      {}
func (n *DynamicArrayType) functionTypeParameterTypeNode() {}
func (n *FixedSizeArrayType) functionTypeParameterTypeNode() {
}
func (n *FunctionType) functionTypeParameterTypeNode()     {}
func (n *ErrorTuple) functionTypeParameterTypeNode()       {}
func (n *TupleType) functionTypeParameterTypeNode()        {}
func (n *GenericType) functionTypeParameterTypeNode()      {}
func (n *TypeReference) functionTypeParameterTypeNode()    {}
func (n *Identifier) functionTypeParameterTypeNode()       {}
func (n *InlineUnion) functionTypeParameterTypeNode()      {}
func (n *UnionType) functionTypeParameterTypeNode()        {}
func (n *UnionDeclaration) functionTypeParameterTypeNode() {}
func (n *FloatLiteral) functionTypeParameterTypeNode()     {}
func (n *IntegerLiteral) functionTypeParameterTypeNode()   {}
func (n *BooleanLiteral) functionTypeParameterTypeNode()   {}
func (n *StringLiteral) functionTypeParameterTypeNode()    {}
func (n *InterpolatedStringLiteral) functionTypeParameterTypeNode() {
}
func (n *RawStringLiteral) functionTypeParameterTypeNode() {}
func (n *MultiLineStringLiteral) functionTypeParameterTypeNode() {
}
func (n *TupleLiteral) functionTypeParameterTypeNode()          {}
func (n *ArrayLiteral) functionTypeParameterTypeNode()          {}
func (n *SymbolLiteral) functionTypeParameterTypeNode()         {}
func (n *RuneLiteral) functionTypeParameterTypeNode()           {}
func (n *FixedSizeArrayLiteral) functionTypeParameterTypeNode() {}

type Parameter struct {
	BaseNode
	Annotations *Annotations              // Optional annotations
	Type        FunctionTypeParameterType // Parameter type
}

func NewParameter(annotations *Annotations, paramType FunctionTypeParameterType) *Parameter {
	return &Parameter{
		BaseNode:    BaseNode{Type: NodeParameter},
		Annotations: annotations,
		Type:        paramType,
	}
}

func (p *Parameter) String() string {
	var builder strings.Builder

	if p.Annotations != nil {
		builder.WriteString(p.Annotations.String())
	}

	builder.WriteString(p.Type.String())
	return builder.String()
}

// labeled_parameter = annotations identifier ":" ( nilable_type
//	                                              | type
//	                                              | literal
//	                                              | union_type
//	                                              | union_declaration ) .

type LabeledParameter struct {
	BaseNode
	Annotations *Annotations              // Optional annotations
	Identifier  *Identifier               // Parameter name
	Type        FunctionTypeParameterType // Parameter type
}

func NewLabeledParameter(annotations *Annotations, identifier *Identifier, paramType FunctionTypeParameterType) *LabeledParameter {
	return &LabeledParameter{
		BaseNode:    BaseNode{Type: NodeLabeledParameter},
		Annotations: annotations,
		Identifier:  identifier,
		Type:        paramType,
	}
}

func (l *LabeledParameter) String() string {
	var builder strings.Builder

	if l.Annotations != nil {
		builder.WriteString(l.Annotations.String())
	}

	builder.WriteString(l.Identifier.String())
	builder.WriteString(": ")
	builder.WriteString(l.Type.String())
	return builder.String()
}

// rest_parameter = "..." type .

type RestParameter struct {
	BaseNode
	Type TypeArgumentType // Parameter type
}

func NewRestParameter(paramType TypeArgumentType) *RestParameter {
	return &RestParameter{
		BaseNode: BaseNode{Type: NodeRestParameter},
		Type:     paramType,
	}
}

func (r *RestParameter) String() string {
	return "..." + r.Type.String()
}

// labeled_rest_parameter = annotations identifier ":" rest_parameter .

type LabeledRestParameter struct {
	BaseNode
	Annotations *Annotations   // Optional annotations
	Identifier  *Identifier    // Parameter name
	RestType    *RestParameter // Rest parameter type
}

func NewLabeledRestParameter(annotations *Annotations, identifier *Identifier, restType *RestParameter) *LabeledRestParameter {
	return &LabeledRestParameter{
		BaseNode:    BaseNode{Type: NodeLabeledRestParameter},
		Annotations: annotations,
		Identifier:  identifier,
		RestType:    restType,
	}
}

func (l *LabeledRestParameter) String() string {
	var builder strings.Builder

	if l.Annotations != nil {
		builder.WriteString(l.Annotations.String())
	}

	builder.WriteString(l.Identifier.String())
	builder.WriteString(": ")
	builder.WriteString(l.RestType.String())

	return builder.String()
}

// InferredErrorType represents a bare `error` return type whose concrete
// error type or union of error types is inferred from the function body.
type InferredErrorType struct {
	BaseNode
}

func NewInferredErrorType() *InferredErrorType {
	return &InferredErrorType{
		BaseNode: BaseNode{Type: NodeInferredErrorType},
	}
}

func (i *InferredErrorType) String() string {
	return "error"
}

// return_type = union_with_error
//             | union_declaration_with_error
//             | nilable_type
//             | type
//             | "error" .

type ReturnType struct {
	BaseNode
	Type Node // Return type
}

func NewReturnType(returnType Node) *ReturnType {
	return &ReturnType{
		BaseNode: BaseNode{Type: NodeReturnType},
		Type:     returnType,
	}
}

func (r *ReturnType) String() string {
	return r.Type.String()
}

// function_type = ( "fn" | "fx" ) "(" [ labeled_parameters | parameters ] ")" return_type .

type FunctionType struct {
	BaseNode
	HasSideEffects bool                    // True for 'fx', false for 'fn'
	Parameters     []FunctionTypeParameter // Can be Parameter, LabeledParameter, RestParameter, or LabeledRestParameter
	ReturnType     *ReturnType             // Return type
}

func NewFunctionType(hasSideEffects bool, parameters []FunctionTypeParameter, returnType *ReturnType) *FunctionType {
	return &FunctionType{
		BaseNode:       BaseNode{Type: NodeFunctionType},
		HasSideEffects: hasSideEffects,
		Parameters:     parameters,
		ReturnType:     returnType,
	}
}

func (f *FunctionType) String() string {
	var builder strings.Builder

	if f.HasSideEffects {
		builder.WriteString("fx")
	} else {
		builder.WriteString("fn")
	}

	builder.WriteString("(")

	for i, param := range f.Parameters {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(param.String())
	}

	builder.WriteString(")")

	if f.ReturnType != nil {
		builder.WriteString(" ")
		builder.WriteString(f.ReturnType.String())
	}

	return builder.String()
}
