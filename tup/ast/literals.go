package ast

import (
	"strconv"
	"strings"
)

// Base type for all literals
type Literal interface {
	Node
	LiteralValue() interface{}
}

// IntegerLiteral represents an integer literal in the code
type IntegerLiteral struct {
	BaseNode
	Value      string // Original text representation
	Base       int    // 2, 8, 10, or 16
	IsNegative bool
}

// NewIntegerLiteral creates a new IntegerLiteral node
func NewIntegerLiteral(value string, base int, isNegative bool) *IntegerLiteral {
	return &IntegerLiteral{
		BaseNode:   BaseNode{NodeType: NodeIntegerLiteral},
		Value:      value,
		Base:       base,
		IsNegative: isNegative,
	}
}

// String returns a textual representation of the integer literal
func (i *IntegerLiteral) String() string {
	prefix := ""
	if i.IsNegative {
		prefix = "-"
	}

	switch i.Base {
	case 2:
		return prefix + "0b" + i.Value
	case 8:
		return prefix + "0o" + i.Value
	case 16:
		return prefix + "0x" + i.Value
	default:
		return prefix + i.Value
	}
}

// LiteralValue returns the Go value of the literal
func (i *IntegerLiteral) LiteralValue() interface{} {
	// Remove underscores for parsing
	val := strings.ReplaceAll(i.Value, "_", "")

	// Parse according to base
	n, _ := strconv.ParseInt(val, i.Base, 64)
	if i.IsNegative {
		n = -n
	}
	return n
}

// FloatLiteral represents a floating point literal in the code
type FloatLiteral struct {
	BaseNode
	Value       string // Original text representation
	HasExponent bool
	IsNegative  bool
}

// NewFloatLiteral creates a new FloatLiteral node
func NewFloatLiteral(value string, hasExponent bool, isNegative bool) *FloatLiteral {
	return &FloatLiteral{
		BaseNode:    BaseNode{NodeType: NodeFloatLiteral},
		Value:       value,
		HasExponent: hasExponent,
		IsNegative:  isNegative,
	}
}

// String returns a textual representation of the float literal
func (f *FloatLiteral) String() string {
	if f.IsNegative {
		return "-" + f.Value
	}
	return f.Value
}

// LiteralValue returns the Go value of the literal
func (f *FloatLiteral) LiteralValue() interface{} {
	// Remove underscores for parsing
	val := strings.ReplaceAll(f.Value, "_", "")
	n, _ := strconv.ParseFloat(val, 64)
	if f.IsNegative {
		n = -n
	}
	return n
}

// BooleanLiteral represents a boolean literal (true or false)
type BooleanLiteral struct {
	BaseNode
	Value bool
}

// NewBooleanLiteral creates a new BooleanLiteral node
func NewBooleanLiteral(value bool) *BooleanLiteral {
	return &BooleanLiteral{
		BaseNode: BaseNode{NodeType: NodeBooleanLiteral},
		Value:    value,
	}
}

// String returns a textual representation of the boolean literal
func (b *BooleanLiteral) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// LiteralValue returns the Go value of the literal
func (b *BooleanLiteral) LiteralValue() interface{} {
	return b.Value
}

// StringLiteral represents a string literal in the code
type StringLiteral struct {
	BaseNode
	Value string // The string value (without quotes)
	Raw   string // Original text representation (with quotes)
}

// NewStringLiteral creates a new StringLiteral node
func NewStringLiteral(value string, raw string) *StringLiteral {
	return &StringLiteral{
		BaseNode: BaseNode{NodeType: NodeStringLiteral},
		Value:    value,
		Raw:      raw,
	}
}

// String returns a textual representation of the string literal
func (s *StringLiteral) String() string {
	return s.Raw
}

// LiteralValue returns the Go value of the literal
func (s *StringLiteral) LiteralValue() interface{} {
	return s.Value
}

// RawStringLiteral represents a raw string literal enclosed in backticks
type RawStringLiteral struct {
	BaseNode
	Value string // The string value (without backticks)
}

// NewRawStringLiteral creates a new RawStringLiteral node
func NewRawStringLiteral(value string) *RawStringLiteral {
	return &RawStringLiteral{
		BaseNode: BaseNode{NodeType: NodeRawStringLiteral},
		Value:    value,
	}
}

// String returns a textual representation of the raw string literal
func (r *RawStringLiteral) String() string {
	return "`" + r.Value + "`"
}

// LiteralValue returns the Go value of the literal
func (r *RawStringLiteral) LiteralValue() interface{} {
	return r.Value
}

// Interpolation represents an interpolated expression within a string
type Interpolation struct {
	BaseNode
	Expression Node // The expression to be interpolated
}

// NewInterpolation creates a new Interpolation node
func NewInterpolation(expr Node) *Interpolation {
	return &Interpolation{
		BaseNode:   BaseNode{NodeType: NodeInterpolation},
		Expression: expr,
	}
}

// String returns a textual representation of the interpolation
func (i *Interpolation) String() string {
	return "\\(" + i.Expression.String() + ")"
}

// Children returns the child nodes
func (i *Interpolation) Children() []Node {
	return []Node{i.Expression}
}

// InterpolatedStringLiteral represents a string literal with interpolated expressions
type InterpolatedStringLiteral struct {
	BaseNode
	Parts []Node // Mix of StringLiteral and Interpolation nodes
}

// NewInterpolatedStringLiteral creates a new InterpolatedStringLiteral node
func NewInterpolatedStringLiteral(parts []Node) *InterpolatedStringLiteral {
	return &InterpolatedStringLiteral{
		BaseNode: BaseNode{NodeType: NodeInterpolatedStringLiteral},
		Parts:    parts,
	}
}

// String returns a textual representation of the interpolated string literal
func (i *InterpolatedStringLiteral) String() string {
	var builder strings.Builder
	builder.WriteString("\"")
	for _, part := range i.Parts {
		builder.WriteString(part.String())
	}
	builder.WriteString("\"")
	return builder.String()
}

// LiteralValue returns the Go value of the literal (concatenated parts)
func (i *InterpolatedStringLiteral) LiteralValue() interface{} {
	// This is a complex operation that would be resolved at runtime
	// For now, return a simplified string representation
	return i.String()
}

// Children returns the child nodes
func (i *InterpolatedStringLiteral) Children() []Node {
	return i.Parts
}

// MultiLineStringLiteral represents a multi-line string literal with optional processor
type MultiLineStringLiteral struct {
	BaseNode
	Lines     []string      // The actual string content (by lines)
	Processor *FunctionCall // Optional processor function (may be nil)
}

// NewMultiLineStringLiteral creates a new MultiLineStringLiteral node
func NewMultiLineStringLiteral(lines []string, processor *FunctionCall) *MultiLineStringLiteral {
	return &MultiLineStringLiteral{
		BaseNode:  BaseNode{NodeType: NodeMultiLineStringLiteral},
		Lines:     lines,
		Processor: processor,
	}
}

// String returns a textual representation of the multi-line string literal
func (m *MultiLineStringLiteral) String() string {
	var builder strings.Builder
	builder.WriteString("```")
	if m.Processor != nil {
		builder.WriteString(m.Processor.String())
	}
	builder.WriteString("\n")
	for _, line := range m.Lines {
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("```")
	return builder.String()
}

// LiteralValue returns the Go value of the literal
func (m *MultiLineStringLiteral) LiteralValue() interface{} {
	return strings.Join(m.Lines, "\n")
}

// Children returns the child nodes
func (m *MultiLineStringLiteral) Children() []Node {
	if m.Processor != nil {
		return []Node{m.Processor}
	}
	return nil
}

// RuneLiteral represents a rune literal in the code
type RuneLiteral struct {
	BaseNode
	Value rune   // The rune value
	Raw   string // Original text representation
}

// NewRuneLiteral creates a new RuneLiteral node
func NewRuneLiteral(value rune, raw string) *RuneLiteral {
	return &RuneLiteral{
		BaseNode: BaseNode{NodeType: NodeRuneLiteral},
		Value:    value,
		Raw:      raw,
	}
}

// String returns a textual representation of the rune literal
func (r *RuneLiteral) String() string {
	return r.Raw
}

// LiteralValue returns the Go value of the literal
func (r *RuneLiteral) LiteralValue() interface{} {
	return r.Value
}

// SymbolLiteral represents a symbol literal in the code (e.g., :name)
type SymbolLiteral struct {
	BaseNode
	Value string // Symbol name without the colon
}

// NewSymbolLiteral creates a new SymbolLiteral node
func NewSymbolLiteral(value string) *SymbolLiteral {
	return &SymbolLiteral{
		BaseNode: BaseNode{NodeType: NodeSymbolLiteral},
		Value:    value,
	}
}

// String returns a textual representation of the symbol literal
func (s *SymbolLiteral) String() string {
	return ":" + s.Value
}

// LiteralValue returns the Go value of the literal
func (s *SymbolLiteral) LiteralValue() interface{} {
	return s.Value
}

// TupleMember represents a member of a tuple
type TupleMember struct {
	BaseNode
	Value Node // The value expression
}

// LabeledTupleMember represents a labeled member of a tuple
type LabeledTupleMember struct {
	BaseNode
	Label *Identifier // The label
	Value Node        // The value expression
}

// String returns a textual representation of the labeled tuple member
func (l *LabeledTupleMember) String() string {
	return l.Label.String() + ": " + l.Value.String()
}

// Children returns the child nodes
func (l *LabeledTupleMember) Children() []Node {
	return []Node{l.Label, l.Value}
}

// TupleLiteral represents a tuple literal in the code
type TupleLiteral struct {
	BaseNode
	Members []Node // Mix of TupleMember and LabeledTupleMember nodes
}

// NewTupleLiteral creates a new TupleLiteral node
func NewTupleLiteral(members []Node) *TupleLiteral {
	return &TupleLiteral{
		BaseNode: BaseNode{NodeType: NodeTupleLiteral},
		Members:  members,
	}
}

// String returns a textual representation of the tuple literal
func (t *TupleLiteral) String() string {
	var builder strings.Builder
	builder.WriteString("(")
	for i, member := range t.Members {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(member.String())
	}
	builder.WriteString(")")
	return builder.String()
}

// LiteralValue returns the Go value of the literal
func (t *TupleLiteral) LiteralValue() interface{} {
	// This is a complex value that would be represented differently at runtime
	return t.String()
}

// Children returns the child nodes
func (t *TupleLiteral) Children() []Node {
	return t.Members
}

// ArrayLiteral represents an array literal in the code
type ArrayLiteral struct {
	BaseNode
	Elements      []Node          // The array elements
	TypeSpecifier *TypeIdentifier // Optional type specifier (may be nil)
}

// NewArrayLiteral creates a new ArrayLiteral node
func NewArrayLiteral(elements []Node, typeSpecifier *TypeIdentifier) *ArrayLiteral {
	return &ArrayLiteral{
		BaseNode:      BaseNode{NodeType: NodeArrayLiteral},
		Elements:      elements,
		TypeSpecifier: typeSpecifier,
	}
}

// String returns a textual representation of the array literal
func (a *ArrayLiteral) String() string {
	var builder strings.Builder
	if a.TypeSpecifier != nil {
		builder.WriteString(a.TypeSpecifier.String())
	}
	builder.WriteString("[")
	for i, elem := range a.Elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(elem.String())
	}
	builder.WriteString("]")
	return builder.String()
}

// LiteralValue returns the Go value of the literal
func (a *ArrayLiteral) LiteralValue() interface{} {
	// This is a complex value that would be represented differently at runtime
	return a.String()
}

// Children returns the child nodes
func (a *ArrayLiteral) Children() []Node {
	var children []Node
	if a.TypeSpecifier != nil {
		children = append(children, a.TypeSpecifier)
	}
	children = append(children, a.Elements...)
	return children
}

// FixedSizeArrayLiteral represents a fixed-size array literal in the code
type FixedSizeArrayLiteral struct {
	BaseNode
	ArrayType *ArrayType // The array type
	Elements  []Node     // The array elements
}

// NewFixedSizeArrayLiteral creates a new FixedSizeArrayLiteral node
func NewFixedSizeArrayLiteral(arrayType *ArrayType, elements []Node) *FixedSizeArrayLiteral {
	return &FixedSizeArrayLiteral{
		BaseNode:  BaseNode{NodeType: NodeFixedSizeArrayLiteral},
		ArrayType: arrayType,
		Elements:  elements,
	}
}

// String returns a textual representation of the fixed-size array literal
func (f *FixedSizeArrayLiteral) String() string {
	var builder strings.Builder
	builder.WriteString(f.ArrayType.String())
	builder.WriteString("[")
	for i, elem := range f.Elements {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(elem.String())
	}
	builder.WriteString("]")
	return builder.String()
}

// LiteralValue returns the Go value of the literal
func (f *FixedSizeArrayLiteral) LiteralValue() interface{} {
	// This is a complex value that would be represented differently at runtime
	return f.String()
}

// Children returns the child nodes
func (f *FixedSizeArrayLiteral) Children() []Node {
	children := []Node{f.ArrayType}
	children = append(children, f.Elements...)
	return children
}
