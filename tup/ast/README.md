# Tuppence Abstract Syntax Tree (AST) Package

The `ast` package defines the Abstract Syntax Tree (AST) data structures for the Tuppence programming language. This package provides a comprehensive set of types that represent the syntactic structure of Tuppence code.

## Overview

The AST is a hierarchical representation of a Tuppence program, where each node in the tree corresponds to a construct in the source code. The AST serves as an intermediate representation between the parser and later stages of the compiler such as type checking, optimization, and code generation.

## Key Components

### Base Node Interface

All AST nodes implement the `Node` interface defined in `node.go`, which provides common functionality:

```go
type Node interface {
	// Pos returns the position of the first character belonging to the node
	Pos() Position
	// End returns the position of the first character immediately after the node
	End() Position
	// Type returns the type of the node
	Type() NodeType
	// String returns a textual representation of the node for debugging
	String() string
	// Children returns all the child nodes
	Children() []Node
}
```

### Node Categories

The AST nodes are organized into several logical categories:

1. **Expressions** - Represent computations that yield values (e.g., `BinaryExpression`, `FunctionCall`)
2. **Literals** - Represent constant values (e.g., `IntegerLiteral`, `StringLiteral`)
3. **Identifiers** - Represent named entities (e.g., `Identifier`, `TypeIdentifier`)
4. **Types** - Represent type-related constructs (e.g., `ArrayType`, `FunctionType`)
5. **Declarations** - Represent declarations of variables, functions, and types (e.g., `FunctionDeclaration`, `TypeDeclaration`)
6. **Statements** - Represent actions to be performed (e.g., `Assignment`, `ReturnExpression`)
7. **Control Flow** - Represent flow control constructs (e.g., `IfExpression`, `ForExpression`)
8. **Pattern Matching** - Represent pattern matching constructs (e.g., `MatchExpression`, `PatternMatch`)
9. **Operators** - Represent operator constructs (e.g., `AddSubOp`, `RelOp`)
10. **Miscellaneous** - Other constructs (e.g., `Comment`, `Module`)

### File Organization

The AST nodes are organized into several files based on their categories:

- `node.go` - Base node interface and common structures
- `identifiers.go` - Identifier-related nodes
- `literals.go` - Literal value nodes
- `types.go` - Type-related nodes
- `expressions.go` - Expression nodes
- `declarations.go` - Declaration nodes
- `operators.go` - Operator nodes
- `controlflow.go` - Control flow nodes
- `pattern_matching.go` - Pattern matching nodes
- `primary_expressions.go` - Primary expression nodes
- `contracts.go` - Contract and union declaration nodes
- `annotations.go` - Annotation-related nodes
- `ranges.go` - Range-related nodes
- `tuples.go` - Tuple-related nodes
- `exports.go` - Export-related nodes
- `interpolations.go` - String interpolation nodes
- `misc.go` - Miscellaneous nodes

## Usage

### Creating AST Nodes

Each AST node type has a corresponding constructor function that facilitates creating new instances:

```go
// Creating an identifier
id := ast.NewIdentifier("myVariable")

// Creating a binary expression
left := ast.NewIntegerLiteral(42, 10, false)
right := ast.NewIntegerLiteral(18, 10, false)
expr := ast.NewBinaryExpression(left, ast.NewAddOp(), right)
```

### Traversing the AST

The AST can be traversed using the `Children()` method, which returns the immediate child nodes of a given node:

```go
func traverse(node ast.Node, depth int) {
    indent := strings.Repeat("  ", depth)
    fmt.Printf("%s%s\n", indent, node.String())
    
    for _, child := range node.Children() {
        traverse(child, depth+1)
    }
}
```

### Position Information

Each node includes position information that indicates its location in the source code:

```go
pos := node.GetPos()
fmt.Printf("Line: %d, Column: %d\n", pos.Line, pos.Column)
```

## Implementation Notes

- All AST nodes embed `BaseNode` which provides common functionality.
- The AST is designed to be immutable after construction to prevent unintended modifications.
- String methods are provided for debugging and visualization purposes.
- Type-specific methods are available on relevant node types to simplify common operations.

## Future Considerations

- Adding visitor pattern support for more flexible traversal strategies.
- Enhancing position information with source file details.
- Adding serialization/deserialization support for persisting ASTs.
- Adding utilities for AST manipulation and transformation.

## References

- `ast.md` - Contains the complete punch list of all implemented AST nodes.
- The Tuppence language specification document. 