package ast

import (
	"strings"
)

type ContractMember struct {
	BaseNode
	Member Node // Either ContractField or ContractFunction
}

func NewContractMember(member Node) *ContractMember {
	return &ContractMember{
		BaseNode: BaseNode{Type: NodeContractMember},
		Member:   member,
	}
}

func (c *ContractMember) String() string {
	return c.Member.String()
}

type ContractMembers struct {
	BaseNode
	Members []*ContractMember // The contract members
}

func NewContractMembers(members []*ContractMember) *ContractMembers {
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

type ContractField struct {
	BaseNode
	Name        *Identifier  // The field name
	Type        Node         // The field type
	Annotations []Annotation // Field annotations
	Docs        string       // Documentation comments
}

func NewContractField(name *Identifier, fieldType Node, annotations []Annotation, docs string) *ContractField {
	return &ContractField{
		BaseNode:    BaseNode{Type: NodeContractField},
		Name:        name,
		Type:        fieldType,
		Annotations: annotations,
		Docs:        docs,
	}
}

func (c *ContractField) String() string {
	var builder strings.Builder
	for _, annotation := range c.Annotations {
		builder.WriteString(annotation.String())
		builder.WriteString(" ")
	}
	builder.WriteString(c.Name.String())
	builder.WriteString(": ")
	builder.WriteString(c.Type.String())
	builder.WriteString(";")
	return builder.String()
}

type ContractFunction struct {
	BaseNode
	Name        *FunctionIdentifier // The function name
	Parameters  []Node              // Function parameters
	ReturnType  Node                // The return type
	Annotations []Annotation        // Function annotations
}

func NewContractFunction(name *FunctionIdentifier, parameters []Node, returnType Node, annotations []Annotation) *ContractFunction {
	return &ContractFunction{
		BaseNode:    BaseNode{Type: NodeContractFunction},
		Name:        name,
		Parameters:  parameters,
		ReturnType:  returnType,
		Annotations: annotations,
	}
}

func (c *ContractFunction) String() string {
	var builder strings.Builder
	for _, annotation := range c.Annotations {
		builder.WriteString(annotation.String())
		builder.WriteString(" ")
	}
	builder.WriteString("fn ")
	builder.WriteString(c.Name.String())
	builder.WriteString("(")
	for i, param := range c.Parameters {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(param.String())
	}
	builder.WriteString(")")
	if c.ReturnType != nil {
		builder.WriteString(" -> ")
		builder.WriteString(c.ReturnType.String())
	}
	builder.WriteString(";")
	return builder.String()
}

type ContractImplementsAnnotation struct {
	BaseNode
	Contract *TypeIdentifier // The contract being implemented
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

type ContractDeclaration struct {
	BaseNode
	Name        *TypeIdentifier                 // The contract name
	TypeParams  []*GenericTypeParam             // Type parameters if generic
	Members     *ContractMembers                // The contract members
	Implements  []*ContractImplementsAnnotation // Implemented contracts
	Annotations []Annotation                    // Contract annotations
}

func NewContractDeclaration(name *TypeIdentifier, typeParams []*GenericTypeParam, members *ContractMembers, implements []*ContractImplementsAnnotation, annotations []Annotation) *ContractDeclaration {
	return &ContractDeclaration{
		BaseNode:    BaseNode{Type: NodeContractDeclaration},
		Name:        name,
		TypeParams:  typeParams,
		Members:     members,
		Implements:  implements,
		Annotations: annotations,
	}
}

func (c *ContractDeclaration) String() string {
	var builder strings.Builder
	for _, annotation := range c.Annotations {
		builder.WriteString(annotation.String())
		builder.WriteString("\n")
	}

	for _, impl := range c.Implements {
		builder.WriteString(impl.String())
		builder.WriteString("\n")
	}

	builder.WriteString("contract ")
	builder.WriteString(c.Name.String())

	if len(c.TypeParams) > 0 {
		builder.WriteString("<")
		for i, param := range c.TypeParams {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(param.String())
		}
		builder.WriteString(">")
	}

	builder.WriteString(" {\n")
	if c.Members != nil {
		builder.WriteString(c.Members.String())
	}
	builder.WriteString("\n}")

	return builder.String()
}
