package ast

// size = decimal_literal | identifier .

type Size interface {
	Node
	sizeNode()
}

func (n *IntegerLiteral) sizeNode() {}
func (n *Identifier) sizeNode()     {}

type ArrayElementType interface {
	Node
	arrayElementTypeNode()
}

func (n *TypeReference) arrayElementTypeNode()      {}
func (n *DynamicArrayType) arrayElementTypeNode()   {}
func (n *FixedSizeArrayType) arrayElementTypeNode() {}

// array_type = fixed_size_array | dynamic_array .

type ArrayType struct {
	BaseNode
	ElementType ArrayElementType
}

func (a *ArrayType) String() string {
	return "[" + "]" + a.ElementType.String()
}

// dynamic_array = "[" "]" (type_reference | array_type) .

type DynamicArrayType struct {
	ArrayType
}

func NewDynamicArrayType(elementType ArrayElementType) *DynamicArrayType {
	return &DynamicArrayType{
		ArrayType: ArrayType{
			BaseNode:    BaseNode{Type: NodeArrayType},
			ElementType: elementType,
		},
	}
}

func (d *DynamicArrayType) String() string {
	return "[]" + d.ElementType.String()
}

// fixed_size_array = "[" size "]" (type_reference | array_type) .

type FixedSizeArrayType struct {
	ArrayType
	Size Size
}

func NewFixedSizeArrayType(elementType ArrayElementType, size Size) *FixedSizeArrayType {
	return &FixedSizeArrayType{
		ArrayType: ArrayType{
			BaseNode:    BaseNode{Type: NodeArrayType},
			ElementType: elementType,
		},
		Size: size,
	}
}

func (f *FixedSizeArrayType) String() string {
	return "[" + f.Size.String() + "]" + f.ElementType.String()
}
