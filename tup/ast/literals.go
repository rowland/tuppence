package ast

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/source"
)

// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .

// Base type for all literals
type Literal interface {
	Node
	literalNode()
}

func (n *BooleanLiteral) literalNode()            {}
func (n *StringLiteral) literalNode()             {}
func (n *InterpolatedStringLiteral) literalNode() {}
func (n *RawStringLiteral) literalNode()          {}
func (n *MultiLineStringLiteral) literalNode()    {}
func (n *TupleLiteral) literalNode()              {}
func (n *ArrayLiteral) literalNode()              {}
func (n *SymbolLiteral) literalNode()             {}
func (n *RuneLiteral) literalNode()               {}
func (n *FloatLiteral) literalNode()              {}
func (n *FixedSizeArrayLiteral) literalNode()     {}

// float_literal = decimal_digit { decimal_digit | "_" } "." decimal_digit { decimal_digit | "_" } [ exponent ]
//               | decimal_digit { decimal_digit | "_" } exponent .

// FloatLiteral represents a floating point literal in the code
type FloatLiteral struct {
	BaseNode
	Value      string  // Original text representation
	FloatValue float64 // The floating point value
}

// NewFloatLiteral creates a new FloatLiteral node
func NewFloatLiteral(value string, floatValue float64, source *source.Source, startOffset int32, length int32) *FloatLiteral {
	return &FloatLiteral{
		BaseNode:   BaseNode{Type: NodeFloatLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:      value,
		FloatValue: floatValue,
	}
}

// integer_literal = binary_literal
//                 | hexadecimal_literal
//                 | octal_literal
//                 | decimal_literal .

// Base type for integer literals
type IntegerLiteral interface {
	Literal
	integerLiteralNode()
}

func (n *BinaryLiteral) integerLiteralNode()      {}
func (n *HexadecimalLiteral) integerLiteralNode() {}
func (n *OctalLiteral) integerLiteralNode()       {}
func (n *DecimalLiteral) integerLiteralNode()     {}

// binary_literal = "0b" ( "0" | "1" ) { "0" | "1" | "_" } .

type BinaryLiteral struct {
	BaseNode
	Value        string // Original text representation
	IntegerValue int64
}

// NewBinaryLiteral creates a new BinaryLiteral node
func NewBinaryLiteral(value string, integerValue int64, source *source.Source, startOffset int32, length int32) *BinaryLiteral {
	return &BinaryLiteral{
		BaseNode:     BaseNode{Type: NodeIntegerLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		IntegerValue: integerValue,
	}
}

// hexadecimal_literal = "0x" hex_digit { hex_digit | "_" } .
// hex_digit = decimal_digit | "a"-"f" | "A"-"F" .

type HexadecimalLiteral struct {
	BaseNode
	Value        string // Original text representation
	IntegerValue int64
}

// NewHexadecimalLiteral creates a new HexadecimalLiteral node
func NewHexadecimalLiteral(value string, source *source.Source, startOffset int32, length int32) *HexadecimalLiteral {
	integerValue, err := strconv.ParseInt(value, 16, 64)
	if err != nil {
		return nil
	}
	return &HexadecimalLiteral{
		BaseNode:     BaseNode{Type: NodeIntegerLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		IntegerValue: integerValue,
	}
}

// octal_literal = "0o" octal_digit { octal_digit } .
// octal_digit = "0"-"7" .

type OctalLiteral struct {
	BaseNode
	Value        string // Original text representation
	IntegerValue int64
}

// NewOctalLiteral creates a new OctalLiteral node
func NewOctalLiteral(value string, source *source.Source, startOffset int32, length int32) *OctalLiteral {
	integerValue, err := strconv.ParseInt(value, 8, 64)
	if err != nil {
		return nil
	}
	return &OctalLiteral{
		BaseNode:     BaseNode{Type: NodeIntegerLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		IntegerValue: integerValue,
	}
}

// decimal_literal = decimal_digit { decimal_digit | "_" } .
// decimal_digit = "0"-"9" .

type DecimalLiteral struct {
	BaseNode
	Value        string // Original text representation
	IntegerValue int64
}

// NewDecimalLiteral creates a new DecimalLiteral node
func NewDecimalLiteral(value string, source *source.Source, startOffset int32, length int32) *DecimalLiteral {
	integerValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil
	}
	return &DecimalLiteral{
		BaseNode:     BaseNode{Type: NodeIntegerLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		IntegerValue: integerValue,
	}
}

// boolean_literal = "true" | "false" .

// BooleanLiteral represents a boolean literal (true or false)
type BooleanLiteral struct {
	BaseNode
	Value        string // Original text representation
	BooleanValue bool
}

// NewBooleanLiteral creates a new BooleanLiteral node
func NewBooleanLiteral(value string, source *source.Source, startOffset int32, length int32) *BooleanLiteral {
	return &BooleanLiteral{
		BaseNode:     BaseNode{Type: NodeBooleanLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:        value,
		BooleanValue: value == "true",
	}
}

// string_literal = '"' { byte_escape_sequence | unicode_escape_sequence | escape_sequence | character - '"' - eol } '"' .

// StringLiteral represents a string literal in the code
type StringLiteral struct {
	BaseNode
	Value       string // Original text representation
	StringValue string // The string value (without quotes)
}

// NewStringLiteral creates a new StringLiteral node
func NewStringLiteral(value string, source *source.Source, startOffset int32, length int32) *StringLiteral {
	return &StringLiteral{
		BaseNode:    BaseNode{Type: NodeStringLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:       value,
		StringValue: value,
	}
}

// raw_string_literal = "`" { "``" | character - "`" } "`" .

// RawStringLiteral represents a raw string literal enclosed in backticks
type RawStringLiteral struct {
	BaseNode
	Value       string // Original text representation
	StringValue string // The string value (without backticks)
}

// NewRawStringLiteral creates a new RawStringLiteral node
func NewRawStringLiteral(value string, source *source.Source, startOffset int32, length int32) *RawStringLiteral {
	return &RawStringLiteral{
		BaseNode:    BaseNode{Type: NodeRawStringLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:       value,
		StringValue: value,
	}
}

// interpolation = "\\(" expression ")" .

// Interpolation represents an interpolated expression within a string
type Interpolation struct {
	BaseNode
	Expression Expression // The expression to be interpolated
}

// NewInterpolation creates a new Interpolation node
func NewInterpolation(expr Expression, source *source.Source) *Interpolation {
	return &Interpolation{
		BaseNode:   BaseNode{Type: NodeInterpolation, Source: source, StartOffset: int32(expr.Pos().Offset), Length: int32(expr.End().Offset - expr.Pos().Offset)},
		Expression: expr,
	}
}

// interpolated_string_literal = '"' { byte_escape_sequence | unicode_escape_sequence | escape_sequence | interpolation | character - '"' - eol } '"' .

// InterpolatedStringLiteral represents a string literal with interpolated expressions
type InterpolatedStringLiteral struct {
	BaseNode
	Parts []Node // Mix of StringLiteral and Interpolation nodes
}

// NewInterpolatedStringLiteral creates a new InterpolatedStringLiteral node
func NewInterpolatedStringLiteral(parts []Node, source *source.Source) *InterpolatedStringLiteral {
	return &InterpolatedStringLiteral{
		BaseNode: BaseNode{Type: NodeInterpolatedStringLiteral, Source: source, StartOffset: int32(parts[0].Pos().Offset), Length: int32(parts[len(parts)-1].End().Offset - parts[0].Pos().Offset)},
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
func (i *InterpolatedStringLiteral) LiteralValue() any {
	// This is a complex operation that would be resolved at runtime
	// For now, return a simplified string representation
	return i.String()
}

// multi_line_string_literal = "```" [ function_call_context ] eol { indented_line } indented_closing .

// MultiLineStringLiteral represents a multi-line string literal with optional processor
type MultiLineStringLiteral struct {
	BaseNode
	Lines     []string      // The actual string content (by lines)
	Processor *FunctionCall // Optional processor function (may be nil)
}

// NewMultiLineStringLiteral creates a new MultiLineStringLiteral node
func NewMultiLineStringLiteral(lines []string, processor *FunctionCall) *MultiLineStringLiteral {
	return &MultiLineStringLiteral{
		BaseNode:  BaseNode{Type: NodeMultiLineStringLiteral},
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
func (m *MultiLineStringLiteral) LiteralValue() any {
	return strings.Join(m.Lines, "\n")
}

// rune_literal = "'" ( byte_escape_sequence | unicode_escape_sequence | escape_sequence | character - eol ) "'" .

// RuneLiteral represents a rune literal in the code
type RuneLiteral struct {
	BaseNode
	Value     string // Original text representation
	RuneValue rune   // The rune value
}

// NewRuneLiteral creates a new RuneLiteral node
func NewRuneLiteral(value string, source *source.Source, startOffset int32, length int32) *RuneLiteral {
	runeValue, _, _, _ := strconv.UnquoteChar(value, '\'')
	return &RuneLiteral{
		BaseNode:  BaseNode{Type: NodeRuneLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:     value,
		RuneValue: runeValue,
	}
}

// symbol_literal = ":" identifier .

// SymbolLiteral represents a symbol literal in the code (e.g., :name)
type SymbolLiteral struct {
	BaseNode
	Value string // Original text representation
}

// NewSymbolLiteral creates a new SymbolLiteral node
func NewSymbolLiteral(value string, source *source.Source, startOffset int32, length int32) *SymbolLiteral {
	return &SymbolLiteral{
		BaseNode: BaseNode{Type: NodeSymbolLiteral, Source: source, StartOffset: startOffset, Length: length},
		Value:    value,
	}
}

// tuple_member = expression .

// TupleMember represents a member of a tuple
type TupleMember = Expression

// labeled_tuple_member = identifier ":" tuple_member .

// LabeledTupleMember represents a labeled member of a tuple
type LabeledTupleMember struct {
	BaseNode
	Label *Identifier // The label
	Value TupleMember // The value expression
}

// tuple_literal = "(" [ labeled_tuple_members | tuple_members ] ")" .

// TupleLiteral represents a tuple literal in the code
type TupleLiteral struct {
	BaseNode
	Labeled bool   // Whether the tuple is labeled
	Members []Node // Mix of TupleMember and LabeledTupleMember nodes
}

// NewTupleLiteral creates a new TupleLiteral node
func NewTupleLiteral(labeled bool, members []Node) *TupleLiteral {
	return &TupleLiteral{
		BaseNode: BaseNode{Type: NodeTupleLiteral},
		Labeled:  labeled,
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
func (t *TupleLiteral) LiteralValue() any {
	// This is a complex value that would be represented differently at runtime
	return t.String()
}

// array_literal = "[" [ array_members | array_literal ] "]"
//               | type_identifier "[" [ array_members | array_literal ] "]" .

// ArrayLiteral represents an array literal in the code
type ArrayLiteral struct {
	BaseNode
	Elements      []Node          // The array elements
	TypeSpecifier *TypeIdentifier // Optional type specifier (may be nil)
}

// NewArrayLiteral creates a new ArrayLiteral node
func NewArrayLiteral(elements []Node, typeSpecifier *TypeIdentifier) *ArrayLiteral {
	return &ArrayLiteral{
		BaseNode:      BaseNode{Type: NodeArrayLiteral},
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
func (a *ArrayLiteral) LiteralValue() any {
	// This is a complex value that would be represented differently at runtime
	return a.String()
}

// fixed_size_array_literal = fixed_size_array "[" array_members "]" .

// FixedSizeArrayLiteral represents a fixed-size array literal in the code
type FixedSizeArrayLiteral struct {
	BaseNode
	ArrayType *ArrayType // The array type
	Elements  []Node     // The array elements
}

// NewFixedSizeArrayLiteral creates a new FixedSizeArrayLiteral node
func NewFixedSizeArrayLiteral(arrayType *ArrayType, elements []Node) *FixedSizeArrayLiteral {
	return &FixedSizeArrayLiteral{
		BaseNode:  BaseNode{Type: NodeFixedSizeArrayLiteral},
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
func (f *FixedSizeArrayLiteral) LiteralValue() any {
	// This is a complex value that would be represented differently at runtime
	return f.String()
}
