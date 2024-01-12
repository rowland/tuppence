pub const TokenType = enum {
    // Symbols
    at, // @
    close_brace, // }
    close_bracket, // ]
    close_paren, // )
    colon, // :
    comma, // ,
    dot, // .
    open_brace, // {
    open_bracket, // [
    open_paren, // (
    question_mark, // ?
    semi_colon, // ;

    // Operators
    op_checked_add, // ?+
    op_checked_div, // ?/
    op_checked_mod, // ?%
    op_checked_mul, // ?*
    op_checked_sub, // ?-
    op_div, // /
    op_minus, // -
    op_mod, // %
    op_mul, // *
    op_not, // !
    op_plus, // +
    op_pow, // ^
    op_shift_left, // <<
    op_shift_right, // >>

    // Bitwise Operators
    op_bitwise_and, // &
    op_bitwise_or, // |

    // Relational Operators
    op_equal_equal, // ==
    op_greater_equal, // >=
    op_greater_than, // >
    op_less_equal, // <=
    op_less_than, // <
    op_not_equal, // !=

    // Logical Operators
    op_logical_and, // &&
    op_logical_or, // ||

    // Assignment
    op_bitwise_and_equal, // &=
    op_bitwise_or_equal, // |=
    op_div_equal, // /=
    op_equal, // =
    op_logical_and_equal, // &&=
    op_logical_or_equal, // ||=
    op_minus_equal, // -=
    op_mod_equal, // %=
    op_mul_equal, // *=
    op_plus_equal, // +=
    op_pow_equal, // ^=

    // Identifiers
    identifier,
    type_identifier,

    // Keywords
    keyword_break, // break
    keyword_else, // else
    keyword_enum, // enum
    keyword_error, // error
    keyword_fn, // fn
    keyword_fx, // fx
    keyword_for, // for
    keyword_if, // if
    keyword_import, // import
    keyword_mut, // mut
    keyword_return, // return
    keyword_switch, // switch
    keyword_try, // try
    keyword_type, // type
    keyword_typeof, // typeof

    // Literals
    binary_literal,
    boolean_literal,
    character_literal,
    decimal_literal,
    float_literal,
    hexadecimal_literal,
    interpolated_string_literal,
    octal_literal,
    raw_string_literal,
    string_literal,

    // Comments
    comment,

    // Special tokens
    eol,
    eof,
    invalid,
};
