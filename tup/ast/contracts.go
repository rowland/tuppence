package ast

import (
	"strings"
)

// ContractMember represents a member of a contract
type ContractMember struct {
	BaseNode
	Member Node // Either ContractField or ContractFunction
}

// NewContractMember creates a new ContractMember node
func NewContractMember(member Node) *ContractMember {
	return &ContractMember{
		BaseNode: BaseNode{NodeType: NodeContractMember},
		Member:   member,
	}
}

// String returns a textual representation of the contract member
func (c *ContractMember) String() string {
	return c.Member.String()
}

// Children returns the child nodes
func (c *ContractMember) Children() []Node {
	return []Node{c.Member}
}

// ContractMembers represents a collection of contract members
type ContractMembers struct {
	BaseNode
	Members []*ContractMember // The contract members
}

// NewContractMembers creates a new ContractMembers node
func NewContractMembers(members []*ContractMember) *ContractMembers {
	return &ContractMembers{
		BaseNode: BaseNode{NodeType: NodeContractMembers},
		Members:  members,
	}
}

// String returns a textual representation of the contract members
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

// Children returns the child nodes
func (c *ContractMembers) Children() []Node {
	children := make([]Node, len(c.Members))
	for i, member := range c.Members {
		children[i] = member
	}
	return children
}

// ContractField represents a field in a contract
type ContractField struct {
	BaseNode
	Name        *Identifier   // The field name
	Type        Node          // The field type
	Annotations []*Annotation // Field annotations
	Docs        string        // Documentation comments
}

// NewContractField creates a new ContractField node
func NewContractField(name *Identifier, fieldType Node, annotations []*Annotation, docs string) *ContractField {
	return &ContractField{
		BaseNode:    BaseNode{NodeType: NodeContractField},
		Name:        name,
		Type:        fieldType,
		Annotations: annotations,
		Docs:        docs,
	}
}

// String returns a textual representation of the contract field
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

// Children returns the child nodes
func (c *ContractField) Children() []Node {
	children := make([]Node, 0, len(c.Annotations)+2)
	for _, annotation := range c.Annotations {
		children = append(children, annotation)
	}
	children = append(children, c.Name, c.Type)
	return children
}

// ContractFunction represents a function in a contract
type ContractFunction struct {
	BaseNode
	Name        *FunctionIdentifier // The function name
	Parameters  []Node              // Function parameters
	ReturnType  Node                // The return type
	Annotations []*Annotation       // Function annotations
	Docs        string              // Documentation comments
}

// NewContractFunction creates a new ContractFunction node
func NewContractFunction(name *FunctionIdentifier, parameters []Node, returnType Node, annotations []*Annotation, docs string) *ContractFunction {
	return &ContractFunction{
		BaseNode:    BaseNode{NodeType: NodeContractFunction},
		Name:        name,
		Parameters:  parameters,
		ReturnType:  returnType,
		Annotations: annotations,
		Docs:        docs,
	}
}

// String returns a textual representation of the contract function
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

// Children returns the child nodes
func (c *ContractFunction) Children() []Node {
	children := make([]Node, 0, len(c.Annotations)+len(c.Parameters)+2)
	for _, annotation := range c.Annotations {
		children = append(children, annotation)
	}
	children = append(children, c.Name)
	children = append(children, c.Parameters...)
	if c.ReturnType != nil {
		children = append(children, c.ReturnType)
	}
	return children
}

// ContractImplementsAnnotation represents an @implements annotation for a contract
type ContractImplementsAnnotation struct {
	BaseNode
	Contract *TypeIdentifier // The contract being implemented
}

// NewContractImplementsAnnotation creates a new ContractImplementsAnnotation node
func NewContractImplementsAnnotation(contract *TypeIdentifier) *ContractImplementsAnnotation {
	return &ContractImplementsAnnotation{
		BaseNode: BaseNode{NodeType: NodeContractImplementsAnnotation},
		Contract: contract,
	}
}

// String returns a textual representation of the contract implements annotation
func (c *ContractImplementsAnnotation) String() string {
	return "@implements(" + c.Contract.String() + ")"
}

// Children returns the child nodes
func (c *ContractImplementsAnnotation) Children() []Node {
	return []Node{c.Contract}
}

// ContractDeclaration represents a contract declaration
type ContractDeclaration struct {
	BaseNode
	Name        *TypeIdentifier                 // The contract name
	TypeParams  []*GenericTypeParam             // Type parameters if generic
	Members     *ContractMembers                // The contract members
	Implements  []*ContractImplementsAnnotation // Implemented contracts
	Annotations []*Annotation                   // Contract annotations
	Docs        string                          // Documentation comments
}

// NewContractDeclaration creates a new ContractDeclaration node
func NewContractDeclaration(name *TypeIdentifier, typeParams []*GenericTypeParam, members *ContractMembers, implements []*ContractImplementsAnnotation, annotations []*Annotation, docs string) *ContractDeclaration {
	return &ContractDeclaration{
		BaseNode:    BaseNode{NodeType: NodeContractDeclaration},
		Name:        name,
		TypeParams:  typeParams,
		Members:     members,
		Implements:  implements,
		Annotations: annotations,
		Docs:        docs,
	}
}

// String returns a textual representation of the contract declaration
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

// Children returns the child nodes
func (c *ContractDeclaration) Children() []Node {
	children := make([]Node, 0, len(c.Annotations)+len(c.TypeParams)+len(c.Implements)+2)
	for _, annotation := range c.Annotations {
		children = append(children, annotation)
	}

	children = append(children, c.Name)

	for _, param := range c.TypeParams {
		children = append(children, param)
	}

	for _, impl := range c.Implements {
		children = append(children, impl)
	}

	if c.Members != nil {
		children = append(children, c.Members)
	}

	return children
}
