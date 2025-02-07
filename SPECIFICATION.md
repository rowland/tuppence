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

    Car = type (make: String, model: String)
    HttpError = error (code: Int, message: String)

Functions, values, fields and parameters are identified using a lowercase letter or underscore,
followed by any number of letters, decimal digits or underscores.

    i = 1
    point = (x: 5, y: 10)
    my_first_car = Car(make: "Toyota", model: "Corona)

A single underscore identifier serves as a placeholder and may be assigned to, but never accessed.

    x, _ = point

Function identifiers also allow a "?" or a "!" as the final character.

    empty? = fn(a: []Int) { len(a) == 0 }
    fail! = fx(message: String) { ... }

## Modules

A module is a namespaced unit of code that must be compiled together.

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

A module is essentially a tuple with initialization code. Top-level assignments are available to
the functions within the module. This includes values, functions, and types.

    five = 5
    fiver = fn() _ { five }
    Quintet = type(a: Int, b: Int, c: Int, d: Int, e: Int)

Exported symbols, like fields of a type or function parameters, are introduced by an identifer
followed by a ":", e.g.

    greeting: "Hello, world!"
    new_quintet: fn(v: Int) { Quintet(v, v, v, v, v) }
    Trio: type(a: Int, b: Int, c: Int)

Only explicitly exported symbols are accessible outside the module.

A module is imported using the `import` keyword.

    io = import("io")
    hello = fx() { io.print("Hello!") }

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

### Simple Literals

Whether or not a simple literal has been bound to an identifier, it is not coerced to a single precision
or internal representation and carries no semantics.

#### Integers

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

##### Binary

    binary_literal = "0b" ( "0" | "1" ) { "0" | "1" | "_" } .

    0b01001000_01001001 # HI

##### Hexadecimal

    hexadecimal_literal = "0x" hex_digit { hex_digit | "_" } .
    hex_digit = decimal_digit | "a"-"f" | "A"-"F" .
    decimal_digit = "0"-"9" .

    0xDEADBEEF # debug marker magic number

##### Octal

    octal_literal = "0o" octal_digit { octal_digit } .
    octal_digit = "0"-"7" .

    0o660 # read and write permissions

##### Decimal

    decimal_literal = decimal_digit { decimal_digit | "_" } .
    decimal_digit = "0"-"9" .

    1_234_567_890

#### Floats

The literal `1.0` constrains internal representations to floating point types
and byte arrays. If there is ambiguity, such as in the cases of overloaded
and generic functions, it will default to type `Float`.

Float literals larger than `65504.0` will not coerce to `Float16`. Other limits
will apply at higher precisions.

As with integer literals, byte array representations make it possible to implement
arbitrary-precision libraries.

Floats in the source code are represented in by a sequence of one or more decimal digits,
followed by a decimal point and at least one more decimal digit, followed by an optional exponent.
Alternatively, the decimal portion may be omitted if an exponent is present.
Underscores may be included anywhere after the first digit to aid in readability.

    float_literal = decimal_digit { decimal_digit | "_" } "." decimal_digit { decimal_digit | "_" } [ exponent ]
                  | decimal_digit { decimal_digit | "_" } exponent .
    decimal_digit = "0"-"9" .
    exponent = "e" [ "-" | "+" ] decimal_digit { decimal_digit } .

    123_456.789
    6.02214076e23 # one mole
    3e5 # speed of light in km/s

#### Booleans

The boolean literals `true` and `false` can be represented by signed or unsigned
integers of any precision, but will default to `Bool` if there is any ambiguity.

#### Strings

String literals are expected to be a sequence of bytes conforming to UTF-8, but
byte escape sequences are supported, which may result in a binary string which is
not conformant.

A string literal's internal representation is an array of bytes, which may be
used to construct user-defined types other than the standard library `String`, but
`String` is the default.

Again, an identifier may be bound to a string literal without assuming a type
until the point of use.

##### Basic Strings

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

##### Raw Strings

Raw strings dispense with most escape sequences and allow newlines.

    raw_string_literal = "`" { "``" | character - "`" } "`" .

Backticks must be escaped, but any other valid UTF-8 sequence is permitted.

##### Interpolated Strings

An interpolated string is really an expression allowing substrings to be concatenated with interpolated expressions.

    interpolated_string_literal = '"' { byte_escape_sequence 
                                      | unicode_escape_sequence 
                                      | escape_sequence 
                                      | interpolation 
                                      | character - '"' - eol 
                                      } '"' .

    interpolation = "\\(" expression ")"

There must be a `string` function matched to the type of each interpolated value.

##### Multi-line Strings

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

#### Runes

Tuppence uses rune literals to represent individual Unicode code points, similar to character literals in some languages, but with explicit support for Unicode. A rune is stored as a 32-bit integer (Rune), allowing it to represent any valid Unicode scalar value.

A rune literal is written as a single character enclosed in single quotes ('), optionally including escape sequences:

    'A'          # Unicode code point 65
    'Ω'          # Unicode code point 937
    '𝄞'          # Unicode code point 119070
    '\n'         # Newline character (code point 10)
    '\u03A9'     # Explicit Unicode escape for 'Ω'
    '\U0001D11E' # Explicit Unicode escape for '𝄞'
    '\\'         # Backslash character

##### Properties of Rune Literals

1. Equivalent to an Integer

Since Rune is an alias for a 32-bit integer, rune literals can be used in integer operations:

    b = 'A' + 1
    print(b)  # 66 ('B')

2. Supports Escape Sequences

Rune literals accept the same escape sequences as string literals, including:

    \n (newline)
    \t (tab)
    \r (carriage return)
    \' (single quote)
    \\ (backslash)
    \xXX (byte escapes)
    \uXXXX (Unicode escapes)
    \UXXXXXXXX (Unicode escapes)

### Complex Literals

#### Tuples

We will consider values constructed with named types to be expressions, rather than literals.

    p1 = Point(5, 10) # expression of type Point
    p2 = (x: 5, y: 10) # tuple literal with named fields
    p3 = (5, 10) # tuple literal with ordinal fields

The result of an expression will have a type and a memory layout. A tuple literal, whether or not
it has been bound to an identifier, is not coerced to a concrete type until the point of use.
At its point of use, it may be coerced to the expected type, or it may assume an anonymous type.

All anonymous types with matching shapes are unified and may be coerced to a named type with matching shape.

Consider a drawing API with the following types and functions:

    Point = type (x: Float, y: Float)
    draw_circle = fx(center: Point, radius: Float)

We may, of course, construct a value of type `Point` to be passed to `draw_circle`:

    c1 = Point(5, 10)
    draw_circle(c1, 6)

Or we may construct the value right at the call site:

    draw_circle(Point(5, 10), 6)

But we may also bind an identifier to a literal and, because the shape is compatible,
it will be coerced at the call site:

    c2 = (5, 10)
    draw_circle(c2, 6)

Or inline:

    draw_circle((5, 10), 6)

Note that it is not necessary in this case to write floating point literals:

    c2 = (5.0, 10.0)

because the binding allows that the literal could exist in multiple types and precisions.

Similarly, we may return a value from a function with a named return type without ceremony:

    shift_southeast = fn(p: Point) Point {
      (p.x + 10, p.y - 10)
    }

Only if there is ambiguity, such as in the cases of overloaded and generic functions,
will the tuple type default to the composite of the most natural types for each component.

Furthermore, a given binding may be used at multiple call sites with different parameters:

    Point16 = type (Int16, Int16)
    Point32 = type (Int32, Int32)
    PointF = type (Float, Float)
    circle_16 = fx(p: Point16, r: Int16)
    circle_32 = fx(p: Point32, r: Int32)
    circle_f = fx(p: PointF, r: Float)

    c3 = (100, 200)
    circle_16(c3, 15)
    circle_32(c3, 15)
    circle_f(c3, 15)

#### Arrays

Once again, we will consider values constructed with named types to be expressions, rather than literals.

    a1 = Int[1, 3, 3] # expression of type array of Int
    a2 = [1, 2, 3] # array literal

`a1` is of known type and memory layout, but `a2` could match several possibilities.

    sum_16 = fn(a: []Int16) Int16
    sum_32 = fn(a: []Int32) Int32
    sum_64 = fn(a: []Int64) Int64
    sum_f = fn(a: []Float) Float

All four functions could take `a2` as an argument.

Only if there is ambiguity, such as in the cases of overloaded and generic functions,
will the array type default to the composite of the most natural types for each member.

## Function Declaration

### Fully-qualified Names

A function is introduced with an identifier, as described above, but consider a function declared
in module `strconv`:

    atoi = fn(s: String) !Int { ... }

It will implicitly have the name `strconv.atoi`, but within the `strconv` model, or when it has been
imported into another module, it is not necessary to use the module-qualified function name.
However, `strconv.atoi` is not the entirety of the fully-qualified name. The full name includes the
types from the function signature, starting with the return type and continuing with the parameter types.

    atoi[!Int, String] = fn(s: String) !Int { ... }

When introducing an identifier into a namespace, such as a module, no two identifiers may be identical.
However, it is only necessary to explicitly include what would otherwise be implicit until the two
identifiers are no longer ambiguous.

Consider a family of functions with different return types:

    atoi[!Int16] = fn(s: String) !Int16 { ... }
    atoi[!Int32] = fn(s: String) !Int32 { ... }
    atoi[!Int64] = fn(s: String) !Int64 { ... }

The complete identifier for the first function might be `atoi[!Int16,String]`, or the types may also
include the module names where they were declared, but the rules for scope resolution make it
unnecessary to use the most verbose version.

Suppose we only want to convert to type `Int`, but we want for convert from either a `String` or an
array of `Byte`:

    atoi[!Int, String] = fn(s: String) !Int { ... }
    atoi[!Int, []Byte] = fn(b: []Byte) !Int { ... }

When declaring the two functions, they only differ starting with the second type argument, so both
must be included in the function identifier.

Invoking the functions have different ergonomic tradeoffs. In the first scenario, the arguments are
of matching times, so the version of the function with the desired return type must be specified:

    foo = atoi[!Int16]("123")
    bar = atoi[!Int32]("456")

In the second scenario, the type of the argument allows the compiler to resolve which version of
the function to invoke:

    foo = atoi("123")
    bar = atoi([]Byte['4', '5', '6'])

### Return Types

#### Union Types in Function Return Types

Tuppence supports concise union syntax for function return types, allowing functions to return values of multiple possible types.

##### Basic Union Return Type

A function can return a value of one of multiple types using the | operator:

    foo = fn() String | Int { ... }

This means foo can return either a String or an Int.

##### Optional Parentheses for Clarity

Parentheses are *optional but allowed* around unions in return types for readability:

    bar = fn() (String | Int) { ... }  # parentheses allowed

While parentheses *do not change behavior*, they can help visually distinguish unions in complex type signatures.

##### Example With Multiple Types

    baz = fn() []Byte | Ok(Int) | Err(String) | error { ... }  # allowed
    qux = fn() ([]Byte | Ok(Int) | Err(String) | error) { ... }  # parentheses optional

Both function signatures above are valid and equivalent.

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

## Scope Resolution

As per usual, identifiers are resolved from the inside out, from the current scope up to the module scope.

    foo = fn() String { "foo 1" }
    bar = fn() String {
        foo = fn() String { "foo 2" }
        foo()
    }
    bar() # returns "foo 2"

Shadowing identifiers should be avoided, where possible, but this is the basic scope resolution order.

A module may export a set of independent functions, but frequently a module will export a type and a
set of functions for operating on that type. So, when a function is invoked using
Uniform Function Call Syntax, which is to say, like it were a method of a type, the module takes
precedence before local scopes.

Without UFCS:

    (Int, add, sub, mul, div) = import("int")
    div(mul(3, add(2, 2)), 2) # (2 + 2) * 3 / 2 = 6

With UFCS:

    Int = import("int").Int
    2.add(2).mul(3).div(2) # still 6

Note that the operators in `(2 + 2) * 3 / 2` map to the function calls above, so user-defined types
may use familiar syntax by implementing those functions.

If the module where a type is declared does not implement a function, or at least a particular
overloaded version of a function, then a function in the local scope can still apply.

For example, the "string" module may export a type `String` without a `mul` function, but one
could implement:

    mul = fn(s: String, m: Int) String {
        for acc, i = ("", 0); i < m {
            acc + s, i + 1
        }.0
    }

Then, one could invoke:

    "5".mul(5) # "55555"

Or, more naturally,

    "5" * 5 # "55555"

One more scope rule: If a function is invoked as if it is a method of a type:

    List = import("list").List
    List.empty()

Then the scope is confined to the module implementing the type. It does not fall back to the local scope.
This will frequently make it unnecessary to import both a module and a type from the module or else
fully qualify the type.

    list = import("list")

    build_list = fn() list.List {
        l = mut list.empty()
        l << "abc << "def"
    }

## Type Constructors

When a type is declared:

    Complex = type(a: Float, b: Float(0))

A constructor function `new` is provided:

    new[Complex, Float, Float] = fn(a: Float, b: Float(0)) Complex { it }

The constructor parameters match the fields of the tuple type, including default values.

If this type is declared in its own module, then instantiating an instance from another module
would be reasonably convenient:

    Complex = import("complex").Complex
    c = Complex.new(2.2, 4.4)

However, if `Complex` were to share a module with other numeric types, the scope resolution rule
whereby `new` is located in the module declaring `Complex` could need help identifying which
`new` function is being invoked, given that other `new` constructors are likely to also have one
or two `Float` parameters.

    Complex = import("numeric").Complex
    c = Complex.new[Complex](2.2, 4.4)

This is lacking, both ergonomically and aesthetically. Fortunately, there is a shorthand:

    c = Complex(2.2, 4.4)

Other type constructors may be implemented as well, provided the function signature does not match
the default constructor.

    ErrorInComplexFormat = error
    # parse real and imaginary parts from string with format "a+b"
    new[!Complex, String] = fn(s: String) !Complex {
        a, b, ... = s.split("+").map { it.float() }
        switch (a, b) {
            (Float, Float) { Complex(a, b) }
            (Float, _) { Complex(a) }
            else { ErrorInComplexFormat() }
        }
    }

The default constructors will always provide a concrete instance of the type.
Additional constructors may elect to return unions with other types, such as `error` or `Nil`.

## For Loops

A `for` loop may be a little unusual in a functional language, but consider this: In a functional language with tail call optimization, the recursive calls are transformed into a loop under the hood. The functional semantics of immutable values are preserved without creating a dangerous number of stack frames.

Tuppece `for` loops also preserve functional semantics while providing familiar syntax
and, at least in some cases, improved readability.

Tuppence supports two primary forms of `for` loops:
  1. Traditional `for` loops with an `initializer`, `condition`, and optional `step expression`.
  2. `for`...`in` loops for iterating over iterable collections, including arrays, ranges, and user-defined iterables.

### Traditional for Loops

A conventional `for` loop consists of:
  - An initializer (optional)
  - A condition (mandatory)
  - A step expression (optional, placed either in the loop header or at the end of the block)

Example:

    sum = for i = 0; i < 10; i + 1 {
        sum + i
    }

If the step expression is omitted from the header, it must appear as the last expression inside the block:

    sum = for i = 0; i < 10 {
        sum + i
        i + 1 # step expression
    }

The step expression, whether located in the header or at the end of the block, must be compatible with the initializer.

### for ... in loops

The `for`...`in` loop is designed for iterating over iterable collections and ranges. The syntax mirrors conventional `for` loops but replaces the condition with an iterable expression.

    for i, value in numbers {
        print(i, value)
    }

Tuple destructuring allows named access to structured data:

    for (k: key, v: value) in hash_map {
        print(k, v)
    }

#### Step Expressions in for...in Loops

A `for`...`in` loop may optionally include an initializer and a step expression, just like a traditional for loop.

    sum = for acc = 0; n in numbers; acc + n {}

This behaves as:
  - acc starts at 0.
  - n iterates over numbers.
  - The loop accumulates acc + n at each step.
  - The final value of acc is returned.

Similarly, an index can be managed alongside iteration:

    for i = 0; a, b in collection; i + 1 {
        print(i, a, b)
    }

If no initializer is present, the loop evaluates to `nil`.

