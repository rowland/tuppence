## Type Tags

Type tags are structured as an index into a types table, shifted left by 1 bit. The low-order bit serves as a flag to indicate whether the value is a reference (1) or not (0).

This encoding provides an efficient mechanism for distinguishing between value types and references while maintaining compact type representation.

## Tuple Field Layout

Tuples are laid out in memory such that fields are ordered by size in descending order (largest to smallest). For fields of identical size, the order is determined by a second-order sort key:
	•	For labeled tuples, the second-order key is the identifier name.
	•	For unlabeled tuples, the second-order key is the ordinal position of the field.

This layout ensures predictable memory alignment and efficient access.

## Anonymous Types

Anonymous types with identical structures are merged into a single representation. This ensures that equivalent types, even when declared anonymously in different contexts, share the same entry in the types table.

This optimization reduces memory usage and simplifies type comparisons by avoiding redundant type definitions for structurally identical types.

## Error Type Declaration

The shorthand declaration:

```code
Error = error()
```

is equivalent to the expanded form:

```code
@error
@false
Error = type()
```

## Built-in Annotations

`@error` Annotation: Marks the type as representing an error, ensuring it is treated appropriately by the type system and tooling.

`@false` Annotation: Indicates that the type evaluates as false in control flow constructs like `if` or `switch`.

`@true` Annotation: The default for types not directly or indirectly annotated `@false`.

`@bool` Annotation: Signifies that instances of a type are evaluated in control flow constructs like `if` or `switch` based on the result of a user-defined `true?` function. This allows for custom logic to determine the truthiness of a value, enabling flexible and expressive type behaviors. If a type annotated with `@bool` does not have a `true?` function in scope, it results in a compile-time error.

## Function Parameter Packaging

In Tuppence, all formal function parameters declared in the function signature are automatically packaged into a tuple at runtime. This tuple serves as a structured representation of the function’s input arguments. However, additional implicit parameters are also appended to this tuple to provide contextual and environmental information required for the function’s execution.

Implicit Contextual Parameters:
  * Additional parameters are automatically added to the tuple. These include, but are not limited to:
    * Context: Information about the current execution context (e.g., current module, scope).
    * Process Information: Details about the process or thread of execution.
    * Heap/Memory Management: A reference to memory or heap resources accessible to the function.
  * These implicit parameters ensure that functions have the metadata required for efficient execution, debugging, and resource management.

Accessing Implicit Parameters:
  * These implicit parameters are generally not directly accessible in the function body but are used internally by the runtime or standard library. If needed, specific parameters may be exposed through utilities or annotations.

Optimization:
  * If a function does not require certain implicit parameters, the compiler may optimize their inclusion, reducing overhead for simple functions.

