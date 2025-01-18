# Tuppence Language Specification

## Keywords

||||||
| --- | --- | --- | --- | --- |
| fn | fx |
| type  | error | enum | interface |
| switch | for | if | else |
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
| += | (a = a + x ) |
| -= | (a = a - x ) |
| *= | (a = a * x ) |
| /= | (a = a / x ) |
| <<= | (a = append(a, x) ) |
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
| Float32 | Float64 |
| Byte | Int | Float | Rune |
| Array | String | Range |
