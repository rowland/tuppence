# Tuppence Language Specification

## Keywords

||||||
| --- | --- | --- | --- | --- |
| fn | fx |
| type  | error | enum | interface |
| switch | for | if | else |
| mut |
| i8 | i16 | i32 | i64 |
| u8 | u16 | u32| u64 |
| f32 | f64 |
| v128 |
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
| += | (a = a + x ) |
| -= | (a = a - x ) |
| *= | (a = a * x ) |
| /= | (a = a / x ) |
| << | (a = append(a, x) ) |
| \|> | (pipe) |
| . | (dereference) |

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
