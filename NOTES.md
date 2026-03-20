Functions have an implicit return type, which must be a subset of any formal return type, if one is specified.

Auditing use of "private" fields (having names starting with an underscore) and unused exported functions should be relegated to a vet utility.

Unlabeled fields may be labeled _0, _1, _2, etc.

Augmenting a tuple:
    bar = (...foo, baz)
    bar = Bar(...foo, baz)

Types are @true by default, unless annotated @false or @bool.

Generic functions have an implicit constraint.
We may be able to express constraints explicitly.
Implicit contraints may be more minimal.
The purpose of these constraints is better error messages at generic boundaries.

Operator =~ maps to matches?[T, U] = fn(a: T, b: U) bool

For indicating whether a type is true or false, the following annotations can be used.

    @true
    Ok = type()
    @false
    Err = type(String)
    Result = Ok | Err

    Result[a] = union(
      @true
      Ok()
      @false
      Err(a)
    )

:foo == Symbol("foo")

``` triple backtick for multi-line strings

	•	" => handled by tokenizer as the start of a TokenStringLiteral (instead of a separate token)
	•	` => handled by tokenizer as the start of a TokenRawStringLiteral
