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
| import | array | cap | len |
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
| []! | safe_index |
| << | shl(a, x) |
| >> | shr(a, x) |
| <<= | (a = shl(a, x) ) |
| >>= | (a = shr(a, x) ) |
| += | (a = a + x ) |
| -= | (a = a - x ) |
| *= | (a = a * x ) |
| /= | (a = a / x ) |
| \|> | (pipe) |
| . | (dereference) |
| ! | not (Bool) |
| ~ | not (Int) |
| - | neg (Int) |

## Internal Types

|||||
| --- | --- | --- | --- |
| I8 | I16 | I32 | I64 |
| U8 | U16 | U32| U64 |
| F16 | F32 | F64 |
| V128 | PTR |

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

Tuppence supports two types of Unicode escape sequences:

  - `\uXXXX` - 4-digit hexadecimal escape for Basic Multilingual Plane (BMP) code points.
  - `\UXXXXXXXX` - 8-digit hexadecimal escape for full Unicode range (outside BMP).

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

### Type-qualified Declarations

Tuppence allows defining values and functions directly associated with a type using type-qualified declarations. These declarations are always written in the form:

    Type.identifier = expression

or:

    Type.identifier = fn(...) { ... }

This enables encapsulation of constants and functions within a type’s namespace.

Type-qualified declarations *must be in the same scope as the type they belong to*.
This ensures consistency and prevents modifying a type from an unrelated context.

If a type is *declared at the module level*, its *type-qualified members* must also be at the module level.

    Foo = type(x: Int, y: Int)

    Foo.default = Foo(x: 0, y: 0) # Allowed: Foo.default is set in the same scope

    Foo.describe = fn(self: Foo) String {
        "Foo(\(self.x), \(self.y))"
    }

If a type is *declared inside a function, type-qualified members must also be declared inside the function*.

    bar = fn() {
        Baz = type(a: String)

        Baz.example = Baz(a: "Hello") # Allowed: Declared in the same function scope
    }

Declaring type-qualified values outside of the scope where the type is defined is not allowed.

    Foo = type(x: Int, y: Int)

    bar = fn() {
        Foo.default = Foo(x: 0, y: 0)  # Error: Foo was declared at module scope, but this is inside a function.
    }

## Type Constructors

When a type is declared in Tuppence, a default constructor is automatically provided.
Constructors are now declared within the type’s namespace, ensuring clarity and preventing name conflicts.

When defining a *new type*, a *default constructor* is automatically created with 
*parameters matching the fields of the type*, including *default values*:

    Complex = type(a: Float, b: Float(0))
    Complex.new = fn(a: Float, b: Float(0)) Complex { it } # compiler supplied default

If Complex is part of a large module containing multiple types, its constructor remains unambiguous
because it is declared within its type’s namespace:

    Complex = import("numeric").Complex
    c = Complex(2.2, 4.4)

### Defining Custom Constructors

You can define additional constructors for a type within its namespace by declaring them as type-qualified functions.

    ErrorInComplexFormat = error
    # parse real and imaginary parts from string with format "a+b"
    Complex.from_string = fn(s: String) !Complex {
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

## Operator Details

Operators in Tuppence are *syntactic sugar* for function calls, allowing *overloading and customization*.
Most operators *delegate* to functions that can be redefined, except for certain core operators that must
be handled by the compiler (e.g., `=`, `.` for dereferencing, and logical short-circuiting operators `&&` and `||`).

### Assignment (=)

Assignment binds a value to an identifier. Unlike other operators, `=` *is not overloadable*.

    x = 10
    y = x + 5  # y is now 15

### Equality (==)

Equality comparison is handled by the function `eq?`.

    eq?(3, 3)   # true
    eq?("abc", "def")  # false

This function can be overridden for custom types.

### Pattern Matching (=~)

The `matches?` function determines if a value matches a pattern.

    matches?("hello", "\w+") # true (regex match)

### Comparison (<, >, <=, >=, <=>)

Comparison operators map to `lt?`, `gt?`, `lte?`, `gte?`, and `compare_to`. The `compare_to` function returns:

    -1 if left-hand side is smaller,
    0 if both are equal,
    1 if left-hand side is larger.

### Arithmetic (+, -, *, /, %)

Arithmetic operators are function calls:

    add(5, 3)   # 8
    sub(10, 4)  # 6
    mul(6, 7)   # 42
    div(8, 2)   # 4
    mod(10, 3)  # 1

The checked arithmetic versions (`?+`, `?-`, etc.) return an error instead of overflowing.

    checked_add(Int32.max, 1)  # error

### Exponentiation (^)

The `^` operator in Tuppence is used for exponentiation (raising a number to a power).
It corresponds to the `pow` function.

    result = base ^ exponent

this is equivalent to:

    result = pow(base, exponent)

Examples:

    x = 2 ^ 3   # x = 8
    y = 5 ^ 0   # y = 1
    z = 9 ^ 0.5 # z = 3.0 (square root of 9)

Notes:

  - The exponent can be an integer or a floating-point number.
  - Negative exponents return fractional results (2 ^ -1 is 0.5).
  - For integer bases with negative exponents, the result is a floating-point number.

### Bitwise Operators (&, |, ^)

These operators perform bitwise logic:

    and(0b1010, 0b1100)  # 0b1000
    or(0b1010, 0b1100)   # 0b1110

### Bit Shifting (<<, >>)

Bitwise shifting operators shift bits left or right. These call shl(a, x) and shr(a, x), which must be implemented for numeric types.

    shl(0b0001, 2)  # 0b0100 (shift left)
    shr(0b0100, 1)  # 0b0010 (shift right)

Compound assignment variants (<<= and >>=) update in place:

    x = mut 1
    x <<= 3  # x is now 8

### Indexing ([])

Indexing calls `index(array, key)`. This can be overloaded for custom collections.

    index(my_array, 2)  # Get 3rd element

### Chaining (|>)

The pipe operator passes results between functions:

    "hello" |> upper() |> print()  # "HELLO"

Equivalent to:

    print(upper("hello"))

### Logical Operators (&&, ||)

Logical AND (`&&`) and OR (`||`) *short-circuit*, meaning evaluation stops early.

    false && expensive_function()  # `expensive_function` is never called
    true || expensive_function()   # `expensive_function` is never called

### Dereferencing (.)

The dot (`.`) operator accesses *tuple fields, module members, and function calls*.

    point = (x: 3, y: 5)
    print(point.x)  # 3

### Logical and Bitwise Negation (! and ~)

#### Boolean Negation (!)

The ! operator negates a boolean value and is overloadable via the `not` function.
It is commonly used in logical expressions.

    not(true)   # false
    not(false)  # true

Using the operator:

    !true   # false
    !false  # true

For user-defined types, the `!` operator invokes the `not` function if the type is annotated with `@bool`:

    UserDefined = type(x: Int)

    not[Bool, UserDefined] = fn(u: UserDefined) Bool {
        u.x == 0
    }

    a = UserDefined(0)
    b = UserDefined(42)

    !a   # true
    !b   # false

#### Bitwise Negation (~)

The `~` operator inverts all bits of an integer (bitwise NOT).
This operation is also overloadable via `not`, keeping naming consistent with other bitwise operations.

    not[Int8] = fn(x: Int8) Int8 {
        x.xor(-1)  # Flips all bits
    }

    not[Int16] = fn(x: Int16) Int16 {
        x.xor(-1)
    }

    not[Int32] = fn(x: Int32) Int32 {
        x.xor(-1)
    }

Using the operator:

    ~0b00001111  # 0b11110000
    ~42          # -43 (in two’s complement)

For user-defined types, the ~ operator invokes the not function:

    Flags = type(bits: Int8)

    not[Flags] = fn(f: Flags) Flags {
        Flags(~f.bits)
    }

    f = Flags(0b10101010)
    ~f  # Flags(0b01010101)

## Dynamic Array Instantiation

Tuppence allows creating dynamic arrays with a specified capacity. This ensures memory is preallocated, avoiding frequent resizing.

    arr = array(Int, 10)  # Creates a dynamic array of Int with capacity 10

  - The first argument is the type of elements.
  - The second argument is the initial capacity (optional, defaults to 0 if omitted).

The built-in function `array` is joined by `cap` and `len`.

    numbers = mut array(Int, 5)   # Allocates space for 5 integers but is empty initially
    numbers <<= 42                # Appends a value
    print(len(numbers))           # 1 (not 5, since length starts at 0)
    print(cap(numbers))           # 5

## Fixed-Size Array Initialization

Fixed-size arrays are preallocated with a specific size and do not grow dynamically.
The initialization block determines how values are assigned.

### Examples

1. Specify a Single Default Value

Use one value to fill the entire array.

    Zeroes = [5]Int
    z = Zeroes[0]     # [0, 0, 0, 0, 0]

All elements are initialized to zero.

2. Specify All Values

Provide explicit values.

    Colors = [3]String
    c = Colors["red", "green", "blue"]  # ["red", "green", "blue"]

Must match exactly the array’s size.

3. Index-Based Initialization

Initialize using the index.

    Indices = [8]Int
    idx = Indices { it }  # [0, 1, 2, 3, 4, 5, 6, 7]

Each element gets its index value.

4. Multi-Dimensional Initialization

    Table = [3, 3]Int
    t = Table { |x, y| (x + 1) * (y + 1) }

Creates:

    [
      [1, 2, 3],
      [2, 4, 6],
      [3, 6, 9]
    ]

## Tuppence Type System

### Integer Types

Tuppence provides both signed and unsigned integer types of various bit widths. The naming convention follows `IntN` for signed integers and `UIntN` for unsigned integers, where `N` represents the bit width:

| Signed Type | Unsigned Type | Alias |
|------------|--------------|-------|
| `Int8`  | `UInt8`  | `Byte` (alias for `UInt8`) |
| `Int16` | `UInt16` | - |
| `Int32` | `UInt32` | `Rune` (alias for `Int32`) |
| `Int64` | `UInt64` | - |

- **Signed integers** (`Int8`, `Int16`, `Int32`, `Int64`) support both positive and negative values.
- **Unsigned integers** (`UInt8`, `UInt16`, `UInt32`, `UInt64`) only support non-negative values (0 and above).
- **`Byte` is an alias for `UInt8`**, making it clear when working with raw byte-level data.

### Character Representation
- **`Rune` is an alias for `Int32`**, representing a Unicode code point.
- Tuppence does not have a dedicated `char` type; instead, `Rune` provides full Unicode support.
- Character literals use single quotes and are interpreted as their corresponding `Rune` integer values:
  
  ```tuppence
  a = 'A'  # Equivalent to 65
  euro = '€'  # Equivalent to 8364
  ```

### Type System Properties
- All types in Tuppence are **immutable**.
- Integer types are fixed-width, except for Int and UInt, which adjust to the architecture’s natural word size.

- **Type aliases** (such as `Byte` and `Rune`) provide better readability but do not introduce new underlying types.

### Operations
- **Signed integers support arithmetic and bitwise operations** (`+`, `-`, `*`, `/`, `&`, `|`, `^`, `<<`, `>>`).
- **Unsigned integers support the same operations** but **disallow negation (`-x`)**, ensuring safe unsigned arithmetic.
- `UIntN` to `IntN` conversions must be explicit, and truncation rules apply when converting to a smaller bit width.

### Example Usage
```tuppence
x = Int32(100)
b = Byte(255)
r = Rune('λ')

sum = UInt64(x) + UInt64(b)
print(sum.string())  # "355"
```

## **Contract Annotations**

In Tuppence, **contracts** define a set of required functions that a type must implement. To formally associate
a type with a contract, we use **annotations**.

Tuppence supports **contract annotations** using the `@type:implements` directive, which allows a type to
explicitly declare that it conforms to a contract. The compiler will enforce this by checking whether all 
required functions are correctly implemented.

### **Syntax**

```tuppence
@type:implements module.ContractName
TypeName = type(...)
```
- `module.ContractName` must be a **valid contract**.
- `TypeName` must implement all required functions.

### **Example: Declaring an Unsigned Integer Type**

Tuppence includes `UnsignedInt[a]`, a contract for unsigned integer types.
We can declare a custom `UInt24` type and specify that it conforms to `UnsignedInt`.

```tuppence
@type:implements core.UnsignedInt
UInt24 = type(lo: UInt8, mid: UInt8, hi: UInt8)
```

This means:
- `UInt24` **must implement** all functions required by `core.UnsignedInt`.
- If it **fails to do so**, the compiler will produce an **error**.

### **Example: Declaring a Numeric Type**

A `BigDecimal` type can declare that it satisfies the `math.Numeric` contract.

```tuppence
@type:implements math.Numeric
BigDecimal = type(value: []Byte)
```

`BigDecimal` **must implement** functions like `add`, `sub`, `mul`, and `div`.  
If it does **not** provide these functions, compilation will **fail**.

### **Compiler Behavior**

When encountering `@type:implements`, the compiler will:

1. **Resolve the referenced contract** (`core.UnsignedInt`, `math.Numeric`, etc.).
2. **Verify that all required functions** are implemented for `TypeName`.
3. **Raise an error** if:
   - The referenced contract does not exist.
   - `TypeName` does not provide all required functions.

## Type Introspection with typeof

`typeof` retrieves a uniform type descriptor for any value or type.

    typeof(expression)

Returns a type descriptor, which can be used for type comparisons and introspection.

If `typeof` is applied to a literal, the result is resolved at compile-time.

To check if a value is an instance of a type:

    if typeof(x) == typeof(Int) {
        print("x is an integer")
    }

## Inline `for` Loops

Inline `for` loops in Tuppence allow compile-time iteration over fixed-size structures, such as tuples. These loops are fully unrolled at compile time and enable operations like tuple transformation, filtering, or field-based computations.

### Syntax

An `inline for` loop follows this structure:

```tuppence
new_tuple = inline for acc = (); name, value in some_tuple {
    # Compile-time transformation logic
    (...acc, name: modified_value)
}
```

- The loop iterates over `some_tuple`, producing named tuples `(name: Symbol, value: T)`, where `T` is the type of the corresponding tuple field.
- The `acc` variable accumulates the transformed values as a new tuple.
- Each iteration appends a new field to `acc`, potentially modifying field names and values.
- The loop is unrolled at compile time, ensuring efficient and predictable behavior.

### Example: Transforming Tuple Fields

```tuppence
ABC = type(a: Int, b: String, c: Float)
abc = ABC(1, "Hello", 5.5)

def = inline for acc = (); name, value in abc {
    switch name {
        :a { (...acc, d: value + 1) }
        :b { (...acc, e: value + " World") }
        :c { (...acc, f: value + 4.5) }
    }
}
```

Result:
```tuppence
def == (d: 2, e: "Hello World", f: 10.0)
```

### Compile-Time Constraints

- The tuple fields are iterated in declaration order.
- The `name` must be a known symbol at compile time; dynamic field generation is not permitted.
- The `value` is statically known per iteration, allowing type specialization.

### Use Cases

- Renaming or restructuring tuples
- Compile-time tuple field validation or filtering
- Extracting metadata or generating derived values

## Tuppence Meta Functions and Use Cases

Tuppence introduces `$()` as a **compile-time meta-function mechanism**, allowing key-value pairs in **labeled tuple syntax**. The compiler resolves these expressions at compile-time, embedding values or executing safe operations. The parsing phase treats `$()` as a labeled tuple, while the semantic phase determines validity and execution.

### General Syntax

```tuppence
$(key: value, key2: value2)
```

- **`key`**: Identifies the meta-function (e.g., `file`, `hash`, `env`).
- **`value`**: Can be a **string, expression, or compile-time evaluable function**.
- **Order of keys does not matter**, as names provide disambiguation.

### Use Cases

#### 1. File Embedding

Embed file contents as a string at compile-time:

```tuppence
data = $(file: "config.json")
```

Equivalent to:

```tuppence
data = "{\"key\": \"value\"}"
```

Embed a file as a function that returns the content:

```tuppence
get_config = $(embed: "config.json")
```

Equivalent to:

```tuppence
get_config = fn() String { "{\"key\": \"value\"}" }
```

#### 2. Environment Variables

Retrieve environment variables at compile time:

```tuppence
home_dir = $(env: "HOME")
```

Equivalent to:

```tuppence
home_dir = "/Users/alice"
```

#### 3. Compile-Time Hashing

Compute a hash at compile-time:

```tuppence
hash_value = $(hash: "hello", algorithm: "sha256")
```

Equivalent to:

```tuppence
hash_value = "2cf24dba5fb0a30e..."
```

#### 4. Including Other Tuppence Files

Include another file at compile-time:

```tuppence
$(include: "config.tup")
```

### Error Handling and Validation

- If a required file is missing: **Compiler error**.
- If an expression is not compile-time evaluable: **Compiler error**.
- If an unknown key is used: **Compiler error**.

## sizeof build-in function

Tuppence does **not** treat `sizeof` as a reserved keyword. Instead, it is a built-in function that follows normal scope resolution:

- **Shadowing Allowed**: Users can define a local function named `sizeof`.
- **Scope Resolution**:
  - If `sizeof` exists in a local scope, it takes precedence.
  - If used as `x.sizeof()`, it resolves based on `x`'s module.
- **No Special Grammar Rule**: `sizeof` behaves like `cap` and `len`, remaining a standard function.

Example:

```tuppence
x = (a: 1, b: 2)
size = sizeof(x)  # Resolves to built-in sizeof function

sizeof = fn(x) 0  # Shadowing allowed
size2 = sizeof(x)  # Calls the locally defined sizeof
```

## Indexing in Tuppence

Tuppence provides two forms of indexing for accessing elements within arrays and other collections:

- Standard Indexing
- Safe Indexing

### Standard Indexing ([])

Indexing an array or collection with `[]` calls `index(self, key)`, which should return either an element or an error.

Example:

```tuppence
values = [5]Int[10, 20, 30, 40, 50]

x = values[2]    # Ok: Returns 30
y = values[10]   # ERROR: Returns ErrIndexOutOfBounds
```

The return type is a union of the element type and `ErrIndexOutOfBounds`.

```tuppence
value = values[10]  # Type: Int | ErrIndexOutOfBounds
```

### Custom Overloading

Library authors can define `index()` for their own data structures.

```tuppence
Person = type (
    id: Int
    name: String
    age: Int
)

People = type (
    len: Int
    # other fields organizing Person records
)

index = fn(self: People, key: Int) Person | ErrIndexOutOfBounds {
    if key < 0 || key >= self.len { return ErrIndexOutOfBounds }
    # code returning a Person value
}
```

### Safe Indexing ([]!)

Indexing an array or collection with `[]!` guarantees a valid result by calling `safe_index(array, key)`.

If `safe_index` is not defined, an error is reported.

Safe Indexing Rules

Safe indexing is allowed when the compiler can prove the index is within bounds, such as:

1. Compile-time constant indices

```tuppence
values = [5]Int[10, 20, 30, 40, 50]

x = values[2]!   # Ok: Returns 30
y = values[10]!  # Compile-time ERROR: Index out of bounds
```

2. Looping over a known range

```tuppence
for i in 0..4 {
    print(values[i]!)  # Ok: Guaranteed in range
}
```

3. Comparing index to constants

```tuppence
for i = 0; i < 100; i + 1 {
    if i < len(values) && i >= -len(values) {
        print(values[i]!)
    }
}
```

### Custom Overloading

Until the compiler can ascertain the expected size of a custom data structure, the `[]!` operator cannot be supported.
