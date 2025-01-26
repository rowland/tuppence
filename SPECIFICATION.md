# Tuppence Language Specification

## Keywords

|||||||
| --- | --- | --- | --- | --- | --- |
| fn | fx | return |
| type  | error | enum | union | contract |
| if | else | switch |
| for | in | break | continue |
| try | try_break | try_continue |
| mut |
| import |
| it |
| typeof |
| true | false |

## Operators

| Operator | Function |
| --- | --- |
| = | (assignment) |
| == | eq? |
| =~ | matches? |
| < | lt? |
| > | gt? |
| <= | lte? |
| >= | gte? |
| <=> | compare_to |
| + | add |
| - | sub |
| * | mul |
| / | div |
| ?+ | checked_add |
| ?- | checked_sub |
| ?* | checked_mul |
| ?/ | checked_div |
| ?% | checked_mod |
| \| | or |
| \|\| | (logical) |
| & | and |
| && | (logical) |
| % | mod |
| ^ | pow |
| [] | index |
| << | append(a, x) |
| <<= | (a = append(a, x) ) |
| += | (a = a + x ) |
| -= | (a = a - x ) |
| *= | (a = a * x ) |
| /= | (a = a / x ) |
| \|> | (pipe) |
| . | (dereference) |
| ! | not |

## Internal Types

|||||
| --- | --- | --- | --- |
| I8 | I16 | I32 | I64 |
| U8 | U16 | U32| U64 |
| F32 | F64 |
| V128 |

## Standard Types

|||||||
| --- | --- | --- | --- | --- | --- |
| Nil |
| Bool |
| Int8 | Int16 | Int32 | Int64 |
| UInt8 | UInt16 | UInt32 | UInt64 |
| Float16 | Float32 | Float64 |
| Byte | Int | Float | Rune |
| String | Range |

## Core Elements

Tuppence is comprised of functions, tuple types and arrays.

### Functions

Functions may be divided into those that are "pure" and those having side effects.

Pure functions are introduced with the `fn` keyword:

    sqr = fn(i: Int) Int { i * i }

Functions having side effects are introduced with the `fx` keyword:

    log = fx(s: String) { ... }

Functions parameters may have default values:

    hello = fn(entity: "World") String { "Hello, " + entity + "!" }
    hello_world = hello() # "Hello, World!"
    hello_john = hello("John") # "Hello, John!"

### Tuple Types

A tuple type is a heterogenous collection of values, addressable by ordinal position and, optionally, by name.

Tuple types may be introduced with the `type` keyword:

    Person = type (name: String, age: Int)

or using the `error` keyword:

    HttpError = error (code: Int, message: String)

Structurally, the two are the same, but they have different defaults and semantics.

A field may have a default value, which it will be created with if the field is not specified:

  Account = type (name: String, balance: Money(0))
  acct1 = Account("Acme") # name == "Acme", balance == 0
  acct2 = Account("Bob", Money(1000)) # name == "Bob", balance == 1000

### Array Types

An array is a homogenoous collection of values addressable by ordinal position.

An array type may be of fixed size:

    IPv4Address = [4]Byte

or of dynamic size:

    Vec = []Float16

### Values

Tuple types and arrays are instantiated using their constructors:

    person1 = Person(name: "John", age: 18)
    person2 = Person("Jane", age: 25)
    people = Person[person1, person2]

## Identifiers

Types are identified using an uppercase letter followed by any number of letters, decimal digits or underscores.

Functions, values, fields and parameters are identified using a lowercase letter or underscore,
followed by any number of letters, decimal digits or underscores.

A single underscore identifier serves as a placeholder and may be assigned to, but never accessed.

Function identifiers also allow a "?" or a "!" as the final character.

## Modules

A module name is derived from the file name up to the first character that is not a valid identifier.
Multiple files with the same identifier prefix are combined for compilation purposes.

    math.tup
    math.complex.tup
    math.vectors.tup
    math.matrices.tup
    math.trig.tup

Source files containing platform-specific and/or architecture-specific code are named with a suffix
starting with `--` and continuing with OS and ARCH:

    math.vectors--amd64.tup
    math.vectors--arm64.tup
    filesystem--darwin.tup
    filesystem--linux.tup
    inference--darwin-arm64.tup

## Expressions

Tuppence expressions are constructs that evaluate to a value and include literals, function calls,
control flow structures, and operations.

    5
    "five"
    f(x)
    1 + 1
    x * y

## Statements

Tuppence statements are the building blocks of program logic, encompassing assignments, declarations,
conditionals, loops, and expressions.

## Blocks

Except for module initialization, all sequential statements are found within blocks delimited by braces:

    {
      a = 5
      b = 3
      c = 8
      x = 4
      a * x^2 + b * x + c
    }

Identifiers introduced within a block shadow any found within the surrounding scope, out to the module level,
but go out of scope at the end of the block:

    {
      a = 5
      b = {
        a = a * 2
        a + 1
      }
    }

A block evaluates to the value of the final expression.

## Comments

Tuppence uses `#` to introduce comments, which continue until the end of the line.

    # quadratic is a polynomial function of degree 2.
    quadratic = fn(
      a: Int, # leading coefficient
      b: Int, # linear coefficient
      c: Int, # constant term
      x: Int, # independent variable
    ) Int {
      a * x^2 + b * x + c
    }

## Assignments

Named types, functions and values are always introduced through assignment.

### Types

    Name = type (String)
    FullName = type (first: Name, last: Name)
    IndexOutOfBounds = error ()

### Functions

    sqr = fn(n: Float64) Float64 { n * n }
    log = fx(s: String) { ... }

### Values

    n = 5
    v = sqr(n)

### Destructuring

The left side of the `=` may contain multiple identifiers if the expression on the right side
can be destructured into multiple parts:

    a, b = (1, 2)

The placeholder `_` may be used to ignore a single part of the expression:

    c, d, _ = (3, 4, 5)

The "rest" operator `...` may be used to gather zero to many parts of the expression:

    e, f, ...g = (6, 7, 8, 9) # e == 6, f == 7, g == (8, 9)

The "rest" operator without an identifier can be used to ignore zero to many parts of the expression:

    e, f, ... = (6, 7, 8, 9) # e == 6, f == 7

Tuples with named fields can be destructured by name, rather than by ordinal position,
by enclosing the identifiers on the left side in parenthesis:

    Person = type (name: String, age: Int)
    p = Person("John", 18)
    (age, name) = p # age == 18, name == "John"

It would be an error if an identifier on the left did not match a named part of the expression.
However, named parts may be bound to new identifiers by specifying both new and old names:

    (a: age, n: name) = p # a == 18, n == "John"

## Literals

Literals may be divided into the simple, like booleans, numbers and strings, and the complex,
like tuples and arrays, which are composed of simpler types.

## Simple Literals

Whether or not a simple literal has been bound to an identifier, it is not coerced to a single precision
or internal representation and carries no semantics.

### Integers

Take, for example, the literal `1`:

It could be represented by an 8-, 16-, 32- or 64-bit signed integer,
corresponding to internal types I8, I16, I32 or I64.

Or, it could be represented by an 8-, 16-, 32- or 64-bit unsigned integer,
corresponding to internal types U8, U16, U32 or U64.

Or, it could be represented by a 16-, 32- or 64-bit IEEE float,
corresponding to internal types F16, F32 or F64.

It could even be represented by an array of bytes.

It is only at the point of use where there is sufficient information to determine
which variant is called for. If there is ambiguity, such as in the cases of overloaded
and generic functions, it will default to type`Int`, unless it is out of range,
in which case it will default to type `Float`.

The literal `-1` cannot be represented by the unsigned internal types, so it canot be
used in a context like `Bool`. The literal `1000` cannot be represented by an
8-bit internal type, so it cannot be used in a context like `Bool` either.

The literal `18446744073709551616` (2^64) cannot be represented by even an unsigned
64-bit integer, so possible representations are constrained to 32-bit and larger
floats and, of course, byte arrays.

The option of representing numeric literals as byte arrays means user-defined types,
like `BigInt`, can be defined without being forced to use string literals, which
do not have the correctness guarantees the language otherwise provides.

Integers in the source code may be represented in binary, hexadecimal, octal or decimal format.
Underscores may be included anywhere after any prefix and the first digit to aid in readability.

#### Binary

    binary_literal = "0b" ( "0" | "1" ) { "0" | "1" | "_" } .

    0b01001000_01001001 # HI

#### Hexadecimal

    hexadecimal_literal = "0x" hex_digit { hex_digit | "_" } .
    hex_digit = decimal_digit | "a"-"f" | "A"-"F" .
    decimal_digit = "0"-"9" .

    0xDEADBEEF # debug marker magic number

#### Octal

    octal_literal = "0o" octal_digit { octal_digit } .
    octal_digit = "0"-"7" .

    0o660 # read and write permissions

#### Decimal

    decimal_literal = decimal_digit { decimal_digit | "_" } .
    decimal_digit = "0"-"9" .

    1_234_567_890

### Floats

The literal `1.0` constrains internal representations to floating point types
and byte arrays. If there is ambiguity, such as in the cases of overloaded
and generic functions, it will default to type `Float`.

Float literals larger than `65504.0` will not coerce to `Float16`. Other limits
will apply at higher precisions.

As with integer literals, byte array representations make it possible to implement
arbitrary-precision libraries.

Floats in the source code are represented in by a sequence of one or more decimal digits,
followd by a decimal point and at least one more decimal digit, followed by an optional exponent.
Alternatively, the decimal portion may be omitted if an exponent is present.
Underscores may be included anywhere after the first digit to aid in readability.

    float_literal = decimal_digit { decimal_digit | "_" } "." decimal_digit { decimal_digit | "_" } [ exponent ]
                  | decimal_digit { decimal_digit | "_" } exponent .
    decimal_digit = "0"-"9" .
    exponent = "e" [ "-" | "+" ] decimal_digit { decimal_digit } .

    123_456.789
    6.02214076e23 # one mole
    3e5 # speed of light in km/s

### Booleans

The boolean literals `true` and `false` can be represented by signed or unsigned
integers of any precision, but will default to `Bool` if there is any ambiguity.

### Strings

String literals are expected to be a sequence of bytes conforming to UTF-8, but
byte escape sequences are supported, which may result in a binary string which is
not conformant.

A string literal's internal representation is an array of bytes, which may be
used to construct user-defined types other than the standard library `String`, but
`String` is the default.

Again, an identifier may be bound to a string literal without assuming a type
until the point of use.

#### Basic Strings

The most basic syntax for a string is enclosed between double quotes and allows for a range of escape sequences.

    string_literal = '"' { byte_escape_sequence | unicode_escape_sequence | escape_sequence | character - '"' - eol } '"' .
    byte_escape_sequence = "\\" "x" hex_digit hex_digit
    hex_digit = decimal_digit | "a"-"f" | "A"-"F" .
    decimal_digit = "0"-"9" .
    unicode_escape_sequence = "\\" "u" hex_digit hex_digit hex_digit hex_digit
                            | "\\" "U" hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit .
    escape_sequence = ( "\\n" | "\\t" | "\\\"" | "\\'" | "\\\\" | "\\r" | "\\b" | "\\f" | "\\v" | "\\0" | "\\`" ) .
    eol = ( "\r\n" | "\r" | "\n" ) .

    "Hello, World!"
    "Hello, \"World\"--if that is your real name?"

    "First line\r\nSecond line"
    "First line\rSecond line"
    "First line\nSecond line"

    "Enclosing \xe2\x80\x9cWorld\xe2\x80\x9d in left and right double quotes."

    "Unicode snowman: \u2603"
    "Unicode rocket: \U0001F680"

    "Newline: \\n, Tab: \\t, Double Quote: \\\""
    "Backslash: \\\\, Carriage Return: \\r"
    "Bell: \\b, Form Feed: \\f, Vertical Tab: \\v"
    "Null: \\0, Backtick: \\`"

#### Raw Strings

Raw strings dispense with most escape sequences and allow newlines.

    raw_string_literal = "`" { "``" | character - "`" } "`" .

Backticks must be escaped, but any other valid UTF-8 sequence is permitted.

#### Interpolated Strings

An interpolated string is really an expression allowing substrings to be concatenated with interpolated expressions.

    interpolated_string_literal = '"' { byte_escape_sequence 
                                      | unicode_escape_sequence 
                                      | escape_sequence 
                                      | interpolation 
                                      | character - '"' - eol 
                                      } '"' .

    interpolation = "\\(" expression ")"

There must be a `string` function matched to the type of each interpolated value.

#### Multi-line Strings

A multi-line string combines features from raw strings and interpolated strings, allowing newlines,
escape sequences and interpolation. Leading whitespace from the first line will be removed and
matching leading whitespace from each successive line will also be removed.

    text = ```
      This is a multi-line string.
      It spans multiple lines
      without needing escape characters.
    ```

If a processor function is specified following the initial "```", it will be invoked on the value:

    foo = ```json
      { "a": "b" }
    ```

If the processor function is basically the identity function, the notation could serve as a hint for syntax highlighters.

Or, the processor function can accept parameters and do real work.

    context = (name: "Alice")
    bar = ```mustache(context)
      Hello, {{name}}
    ```


## Function Invocation

Tuppence supports both standard function-calling syntax:

    foo = bar(baz)

and Uniform Function Call Syntax (UFCS):

    foo = baz.bar()

where the "receiver" becomes the first argument.

Functions may be invoked using ordinal arguments:

    foo = bar(baz, boom)

or using named arguments:

    foo = bar(baz: baz, boom: boom)

or using ordinal arguments followed by named arguments:

    foo = bar(baz, boom: boom)

Using a combination of ordinal and named arguments is convenient when some parameters have default arguments.
