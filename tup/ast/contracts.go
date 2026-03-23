package ast

import "strings"

type ContractMemberNode interface {
	Node
	contractMemberNode()
}

func (n *ContractFunction) contractMemberNode() {}
func (n *ContractField) contractMemberNode()    {}

type ContractFieldType interface {
	Node
	contractFieldTypeNode()
}

func (n *NilableType) contractFieldTypeNode()   {}
func (n *TypeReference) contractFieldTypeNode() {}
func (n *Identifier) contractFieldTypeNode()    {}
func (n *DynamicArrayType) contractFieldTypeNode() {
}
func (n *FixedSizeArrayType) contractFieldTypeNode() {}
func (n *FunctionType) contractFieldTypeNode()       {}
func (n *ErrorTuple) contractFieldTypeNode()         {}
func (n *TupleType) contractFieldTypeNode()          {}
func (n *GenericType) contractFieldTypeNode()        {}
func (n *InlineUnion) contractFieldTypeNode()        {}

// contract_function = function_declaration_lhs "=" function_type .

type ContractFunction struct {
	BaseNode
	LHS  *FunctionDeclarationLHS
	Type *FunctionType
}

func NewContractFunction(lhs *FunctionDeclarationLHS, functionType *FunctionType) *ContractFunction {
	return &ContractFunction{
		BaseNode: BaseNode{Type: NodeContractFunction},
		LHS:      lhs,
		Type:     functionType,
	}
}

func (c *ContractFunction) String() string {
	var builder strings.Builder
	builder.WriteString(c.LHS.String())
	builder.WriteString(" = ")
	builder.WriteString(c.Type.String())
	return builder.String()
}

// contract_field = identifier [ "[" type_parameter "]" ] ":" ( nilable_type | type ) .

type ContractField struct {
	BaseNode
	Name          *Identifier
	TypeParameter *TypeParameter
	Type          ContractFieldType
}

func NewContractField(name *Identifier, typeParameter *TypeParameter, fieldType ContractFieldType) *ContractField {
	return &ContractField{
		BaseNode:      BaseNode{Type: NodeContractField},
		Name:          name,
		TypeParameter: typeParameter,
		Type:          fieldType,
	}
}

func (c *ContractField) String() string {
	var builder strings.Builder
	builder.WriteString(c.Name.String())
	if c.TypeParameter != nil {
		builder.WriteString("[")
		builder.WriteString(c.TypeParameter.String())
		builder.WriteString("]")
	}
	builder.WriteString(": ")
	builder.WriteString(c.Type.String())
	return builder.String()
}

// contract_members = contract_member { eol contract_member } eol .

type ContractMembers struct {
	BaseNode
	Members []ContractMemberNode
}

func NewContractMembers(members []ContractMemberNode) *ContractMembers {
	return &ContractMembers{
		BaseNode: BaseNode{Type: NodeContractMembers},
		Members:  members,
	}
}

func (c *ContractMembers) String() string {
	var builder strings.Builder
	for i, member := range c.Members {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(member.String())
	}
	return builder.String()
}

// contract_declaration = "contract" "(" eol contract_members ")" .

type ContractDeclaration struct {
	BaseNode
	Members *ContractMembers
}

func NewContractDeclaration(members *ContractMembers) *ContractDeclaration {
	return &ContractDeclaration{
		BaseNode: BaseNode{Type: NodeContractDeclaration},
		Members:  members,
	}
}

func (c *ContractDeclaration) String() string {
	var builder strings.Builder
	builder.WriteString("contract(\n")
	if c.Members != nil {
		builder.WriteString(indentString(c.Members.String()))
		builder.WriteString("\n")
	}
	builder.WriteString(")")
	return builder.String()
}

type ContractImplementsAnnotation struct {
	BaseNode
	Contract *TypeIdentifier
}

func NewContractImplementsAnnotation(contract *TypeIdentifier) *ContractImplementsAnnotation {
	return &ContractImplementsAnnotation{
		BaseNode: BaseNode{Type: NodeContractImplementsAnnotation},
		Contract: contract,
	}
}

func (c *ContractImplementsAnnotation) String() string {
	return "@implements(" + c.Contract.String() + ")"
}
