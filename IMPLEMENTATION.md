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

`@cstruct` Annotation: Layout follows C standard for platform, rather than Tuppence-native optimized layouts.

## Function Parameter Packaging

Formal parameters form a logical tuple. The target ABI can split its fields into
registers or stack slots, or materialize a tuple when needed. Runtime context
such as a heap handle can be passed separately when required.

The proposed storage lowering assigns borrowed or consumed modes to logical
inputs and delivers owned results. A consumed heap reference transfers one
strong owner to the callee. See [IR_OWNERSHIP.md](IR_OWNERSHIP.md) for the
counting invariant, duplicate arguments, and the conditions for storage reuse.

## Symbols

A symbol is represented as a tuple (hash, string).
  * hash: The computed hash of the string (MurmurHash or CityHash for balance of speed and collision resistance)
  * string: The original string for debugging or reverse mapping.
  * Symbolic literals found in a module are precomputed at module compile or initialization time.
  * Comparison function `eq?` compares only the hash portion.

