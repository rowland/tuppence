const std = @import("std");
const Token = @import("token.zig").Token;
const TokenType = @import("token_type.zig").TokenType;

pub fn tokenize(allocator: std.mem.Allocator, source: []const u8, filename: []const u8) ![]Token {
    var tokens = std.ArrayList(Token).init(allocator);
    defer tokens.deinit();

    var tokenizer = Tokenizer.init(source, filename);
    while (true) {
        const token = tokenizer.next();
        try tokens.append(token);
        if (token.typ == .eof) {
            break;
        }
    }
    return tokens.toOwnedSlice();
}

pub const Tokenizer = struct {
    source: []const u8,
    filename: []const u8,
    index: usize,
    line: u32,
    bol: usize,
    pending_invalid_token: ?Token,

    pub fn init(source: []const u8, filename: []const u8) Tokenizer {
        // Skip the UTF-8 BOM if present
        const index: usize = if (std.mem.startsWith(u8, source, "\xEF\xBB\xBF")) 3 else 0;
        return Tokenizer{
            .source = source,
            .filename = filename,
            .line = 1,
            .bol = 0,
            .index = index,
            .pending_invalid_token = null,
        };
    }

    const State = enum {
        start,
        question_mark,
        op_div,
        op_minus,
        op_mod,
        op_mul,
        op_not,
        op_plus,
        op_pow,
        op_less_than,
        op_greater_than,
        op_bitwise_and,
        op_logical_and,
        op_bitwise_or,
        op_logical_or,
        op_equal,
        identifier,
        number,
        int,
        int_dot,
        float,
        exponent,
        exponent_sign,
        exponent_int,
        binary,
        hexadecimal,
        octal,
        raw_string_literal,
        raw_string_literal_end,

        interpolated_string_literal,
        string_literal,
    };

    fn peek(self: *Tokenizer) ?u8 {
        if (self.index + 1 < self.source.len) {
            return self.source[self.index + 1];
        } else {
            return null;
        }
    }

    pub fn next(self: *Tokenizer) Token {
        if (self.pending_invalid_token) |token| {
            self.pending_invalid_token = null;
            return token;
        }
        var state: State = .start;
        var start = self.index;
        var typ: TokenType = .eof;
        var invalid = false;
        // var seen_escape_digits: usize = undefined;
        // _ = seen_escape_digits;
        // var remaining_code_units: usize = undefined;
        // _ = remaining_code_units;
        while (self.index <= self.source.len) : (self.index += 1) {
            var c: u8 = undefined;
            if (self.index < self.source.len) {
                c = self.source[self.index];
            } else c = 0;
            switch (state) {
                .start => switch (c) {
                    0 => {
                        if (self.index != self.source.len) {
                            typ = .invalid;
                        }
                        break;
                    },
                    ' ', '\t', '\r' => {
                        start = self.index + 1;
                    },
                    '\n' => {
                        self.line += 1;
                        start = self.index + 1;
                        self.bol = start;
                    },
                    '@' => {
                        typ = .at;
                        self.index += 1;
                        break;
                    },
                    '}' => {
                        typ = .close_brace;
                        self.index += 1;
                        break;
                    },
                    ']' => {
                        typ = .close_bracket;
                        self.index += 1;
                        break;
                    },
                    ')' => {
                        typ = .close_paren;
                        self.index += 1;
                        break;
                    },
                    ':' => {
                        typ = .colon;
                        self.index += 1;
                        break;
                    },
                    ',' => {
                        typ = .comma;
                        self.index += 1;
                        break;
                    },
                    '.' => {
                        typ = .dot;
                        self.index += 1;
                        break;
                    },
                    '{' => {
                        typ = .open_brace;
                        self.index += 1;
                        break;
                    },
                    '[' => {
                        typ = .open_bracket;
                        self.index += 1;
                        break;
                    },
                    '(' => {
                        typ = .open_paren;
                        self.index += 1;
                        break;
                    },
                    '?' => {
                        typ = .question_mark;
                        state = .question_mark;
                    },
                    ';' => {
                        typ = .semi_colon;
                        self.index += 1;
                        break;
                    },
                    '/' => {
                        typ = .op_div;
                        state = .op_div;
                    },
                    '-' => {
                        typ = .op_minus;
                        state = .op_minus;
                    },
                    '%' => {
                        typ = .op_mod;
                        state = .op_mod;
                    },
                    '*' => {
                        typ = .op_mul;
                        state = .op_mul;
                    },
                    '!' => {
                        typ = .op_not;
                        state = .op_not;
                    },
                    '+' => {
                        typ = .op_plus;
                        state = .op_plus;
                    },
                    '^' => {
                        typ = .op_pow;
                        state = .op_pow;
                    },
                    '<' => {
                        typ = .op_less_than;
                        state = .op_less_than;
                    },
                    '>' => {
                        typ = .op_greater_than;
                        state = .op_greater_than;
                    },
                    '&' => {
                        typ = .op_bitwise_and;
                        state = .op_bitwise_and;
                    },
                    '|' => {
                        typ = .op_bitwise_or;
                        state = .op_bitwise_or;
                    },
                    '=' => {
                        typ = .op_equal;
                        state = .op_equal;
                    },
                    'A'...'Z', 'a'...'z', '_' => {
                        typ = .identifier;
                        state = .identifier;
                    },
                    '0' => {
                        state = .number;
                        typ = .decimal_literal;
                    },
                    '1'...'9' => {
                        state = .int;
                        typ = .decimal_literal;
                    },
                    '`' => {
                        state = .raw_string_literal;
                        typ = .raw_string_literal;
                    },

                    '"' => {
                        state = .string_literal;
                        typ = .string_literal;
                    },
                    else => {
                        typ = .invalid;
                        self.index += 1;
                        break;
                    },
                },
                .question_mark => switch (c) {
                    '+' => {
                        typ = .op_checked_add;
                        self.index += 1;
                        break;
                    },
                    '/' => {
                        typ = .op_checked_div;
                        self.index += 1;
                        break;
                    },
                    '%' => {
                        typ = .op_checked_mod;
                        self.index += 1;
                        break;
                    },
                    '*' => {
                        typ = .op_checked_mul;
                        self.index += 1;
                        break;
                    },
                    '-' => {
                        typ = .op_checked_sub;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_div => switch (c) {
                    '=' => {
                        typ = .op_div_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_minus => switch (c) {
                    '=' => {
                        typ = .op_minus_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_mod => switch (c) {
                    '=' => {
                        typ = .op_mod_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_mul => switch (c) {
                    '=' => {
                        typ = .op_mul_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_not => switch (c) {
                    '=' => {
                        typ = .op_not_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_plus => switch (c) {
                    '=' => {
                        typ = .op_plus_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_pow => switch (c) {
                    '=' => {
                        typ = .op_pow_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_less_than => switch (c) {
                    '<' => {
                        typ = .op_shift_left;
                        self.index += 1;
                        break;
                    },
                    '=' => {
                        typ = .op_less_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_greater_than => switch (c) {
                    '>' => {
                        typ = .op_shift_right;
                        self.index += 1;
                        break;
                    },
                    '=' => {
                        typ = .op_greater_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_bitwise_and => switch (c) {
                    '=' => {
                        typ = .op_bitwise_and_equal;
                        self.index += 1;
                        break;
                    },
                    '&' => {
                        typ = .op_logical_and;
                        state = .op_logical_and;
                    },
                    else => break,
                },
                .op_logical_and => switch (c) {
                    '=' => {
                        typ = .op_logical_and_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_bitwise_or => switch (c) {
                    '|' => {
                        typ = .op_logical_or;
                        state = .op_logical_or;
                    },
                    '=' => {
                        typ = .op_bitwise_or_equal;
                        self.index += 1;
                        break;
                    },
                    else => break,
                },
                .op_logical_or => switch (c) {
                    '=' => {
                        typ = .op_logical_or_equal;
                        self.index += 1;
                    },
                    else => break,
                },
                .op_equal => switch (c) {
                    '=' => {
                        typ = .op_equal_equal;
                        self.index += 1;
                    },
                    else => break,
                },
                .identifier => switch (c) {
                    'a'...'z', 'A'...'Z', '_', '0'...'9' => {},
                    else => {
                        const identifier = self.source[start..self.index];
                        if (Token.getReserved(identifier)) |reserved| {
                            typ = reserved;
                        } else if (std.ascii.isUpper(identifier[0])) {
                            typ = .type_identifier;
                        }
                        break;
                    },
                },
                .number => switch (c) {
                    '0'...'9', '_' => {
                        state = .int;
                    },
                    '.' => {
                        typ = .float_literal;
                        state = .int_dot;
                    },
                    'b' => {
                        typ = .binary_literal;
                        state = .binary;
                    },
                    'o' => {
                        typ = .octal_literal;
                        state = .octal;
                    },
                    'x' => {
                        typ = .hexadecimal_literal;
                        state = .hexadecimal;
                    },
                    'A'...'Z', 'a', 'c'...'n', 'p'...'w', 'y'...'z' => {
                        invalid = true;
                    },
                    else => break,
                },
                .int => switch (c) {
                    '0'...'9', '_' => {},
                    '.' => {
                        typ = .float_literal;
                        state = .int_dot;
                    },
                    'e' => {
                        typ = .float_literal;
                        state = .exponent;
                    },
                    'A'...'Z', 'a'...'d', 'f'...'z' => {
                        invalid = true;
                    },
                    else => break,
                },
                .int_dot => switch (c) {
                    '0'...'9' => {
                        state = .float;
                    },
                    'A'...'Z', '_', 'a'...'z' => {
                        state = .float;
                        invalid = true;
                    },
                    else => {
                        invalid = true;
                        break;
                    },
                },
                .float => switch (c) {
                    '0'...'9', '_' => {},
                    'e' => {
                        state = .exponent;
                    },
                    'A'...'Z', 'a'...'d', 'f'...'z' => {
                        invalid = true;
                    },
                    else => break,
                },
                .exponent => switch (c) {
                    '+', '-', '0'...'9' => {
                        state = .exponent_sign;
                    },
                    'A'...'Z', '_', 'a'...'z' => {
                        invalid = true;
                    },
                    else => {
                        invalid = true;
                        break;
                    },
                },
                .exponent_sign => switch (c) {
                    '0'...'9' => {
                        state = .exponent_int;
                    },
                    'A'...'Z', '_', 'a'...'z' => {
                        invalid = true;
                    },
                    else => {
                        invalid = true;
                        break;
                    },
                },
                .exponent_int => switch (c) {
                    '0'...'9' => {},
                    'A'...'Z', '_', 'a'...'z' => {
                        invalid = true;
                    },
                    else => {
                        invalid = true;
                        break;
                    },
                },
                .binary => switch (c) {
                    '0'...'1', '_' => {},
                    '.', '2'...'9', 'A'...'Z', 'a'...'z' => {
                        invalid = true;
                    },
                    else => break,
                },
                .hexadecimal => switch (c) {
                    '0'...'9', 'A'...'F', '_', 'a'...'f' => {},
                    '.', 'G'...'Z', 'g'...'z' => {
                        invalid = true;
                    },
                    else => break,
                },
                .octal => switch (c) {
                    '0'...'7', '_' => {},
                    '.', '8'...'9', 'A'...'Z', 'a'...'z' => {
                        invalid = true;
                    },
                    else => break,
                },
                .raw_string_literal => switch (c) {
                    0 => {
                        invalid = true;
                        break;
                    },
                    '`' => {
                        state = .raw_string_literal_end;
                    },
                    else => {},
                },
                .raw_string_literal_end => switch (c) {
                    '`' => {
                        state = .raw_string_literal;
                    },
                    else => break,
                },

                .interpolated_string_literal => {},
                .string_literal => {},
            }
        }
        return Token{
            .typ = typ,
            .invalid = invalid,
            .line = self.line,
            .column = std.math.cast(u32, self.index - self.bol + 1) orelse 0,
            .value = self.source[start..self.index],
            .filename = self.filename,
        };
    }
};

test "symbols" {
    try testTokenize(
        "@ } ] ) : , . { [ ( ? ;",
        &.{
            .at, // @
            .close_brace, // }
            .close_bracket, // ]
            .close_paren, // )
            .colon, // :
            .comma, // ,
            .dot, // .
            .open_brace, // {
            .open_bracket, // [
            .open_paren, // (
            .question_mark, // ?
            .semi_colon, // ;
        },
    );
}

test "operators" {
    try testTokenize(
        "?+ ?/ ?% ?* ?- / - % * ! + ^ << >>",
        &.{
            .op_checked_add, // ?+
            .op_checked_div, // ?/
            .op_checked_mod, // ?%
            .op_checked_mul, // ?*
            .op_checked_sub, // ?-
            .op_div, // /
            .op_minus, // -
            .op_mod, // %
            .op_mul, // *
            .op_not, // !
            .op_plus, // +
            .op_pow, // ^
            .op_shift_left, // <<
            .op_shift_right, // >>
        },
    );
}

test "bitwise operators" {
    try testTokenize(
        "& |",
        &.{
            .op_bitwise_and, // &
            .op_bitwise_or, // |
        },
    );
}

test "relational operators" {
    try testTokenize(
        "== >= > <= < !=",
        &.{
            .op_equal_equal, // ==
            .op_greater_equal, // >=
            .op_greater_than, // >
            .op_less_equal, // <=
            .op_less_than, // <
            .op_not_equal, // !=
        },
    );
}

test "logical operators" {
    try testTokenize(
        "&& ||",
        &.{
            .op_logical_and, // &&
            .op_logical_or, // ||
        },
    );
}

test "assignment" {
    try testTokenize(
        "&= |= /= = &&= ||= -= %= *= += ^=",
        &.{
            .op_bitwise_and_equal, // &=
            .op_bitwise_or_equal, // |=
            .op_div_equal, // /=
            .op_equal, // =
            .op_logical_and_equal, // &&=
            .op_logical_or_equal, // ||=
            .op_minus_equal, // -=
            .op_mod_equal, // %=
            .op_mul_equal, // *=
            .op_plus_equal, // +=
            .op_pow_equal, // ^=
        },
    );
}

test "identifiers" {
    try testTokenize(
        "abc Def",
        &.{
            .identifier,
            .type_identifier,
        },
    );
}

test "keywords" {
    try testTokenize(
        "break else enum error fn for fx if import mut return switch try type typeof",
        &.{
            .keyword_break,
            .keyword_else,
            .keyword_enum,
            .keyword_error,
            .keyword_fn,
            .keyword_for,
            .keyword_fx,
            .keyword_if,
            .keyword_import,
            .keyword_mut,
            .keyword_return,
            .keyword_switch,
            .keyword_try,
            .keyword_type,
            .keyword_typeof,
        },
    );
}

test "binary literals" {
    try testTokenize("0b", &.{.binary_literal});
    try testTokenize("0b0", &.{.binary_literal});
    try testTokenize("0b1", &.{.binary_literal});
    try testTokenize("0b10101100", &.{.binary_literal});

    try testTokenizeInvalid("0b2", &.{.binary_literal}, true);
    try testTokenizeInvalid("0b3", &.{.binary_literal}, true);
    try testTokenizeInvalid("0b4", &.{.binary_literal}, true);
    try testTokenizeInvalid("0b5", &.{.binary_literal}, true);
    try testTokenizeInvalid("0b6", &.{.binary_literal}, true);
    try testTokenizeInvalid("0b7", &.{.binary_literal}, true);
    try testTokenizeInvalid("0b8", &.{.binary_literal}, true);
    try testTokenizeInvalid("0b9", &.{.binary_literal}, true);
    try testTokenizeInvalid("0ba", &.{.binary_literal}, true);
    try testTokenizeInvalid("0bb", &.{.binary_literal}, true);
    try testTokenizeInvalid("0bc", &.{.binary_literal}, true);
    try testTokenizeInvalid("0bd", &.{.binary_literal}, true);
    try testTokenizeInvalid("0be", &.{.binary_literal}, true);
    try testTokenizeInvalid("0bf", &.{.binary_literal}, true);
    try testTokenizeInvalid("0bz", &.{.binary_literal}, true);

    try testTokenizeInvalid("0b1.", &.{.binary_literal}, true);
    try testTokenizeInvalid("0b1.0", &.{.binary_literal}, true);

    try testTokenizeInvalid("0B0", &.{.decimal_literal}, true);
    try testTokenizeInvalid("0b_", &.{.binary_literal}, false);
    try testTokenizeInvalid("0b_0", &.{.binary_literal}, false);
    try testTokenizeInvalid("0b1_", &.{.binary_literal}, false);
    try testTokenizeInvalid("0b0__1", &.{.binary_literal}, false);
    try testTokenizeInvalid("0b0_1_", &.{.binary_literal}, false);
    try testTokenizeInvalid("0b1e", &.{.binary_literal}, true);
}

test "boolean literals" {
    try testTokenize("false", &.{.boolean_literal});
    try testTokenize("true", &.{.boolean_literal});
}

test "decimal literals" {
    try testTokenize("0", &.{.decimal_literal});
    try testTokenize("1", &.{.decimal_literal});
    try testTokenize("2", &.{.decimal_literal});
    try testTokenize("3", &.{.decimal_literal});
    try testTokenize("4", &.{.decimal_literal});
    try testTokenize("5", &.{.decimal_literal});
    try testTokenize("6", &.{.decimal_literal});
    try testTokenize("7", &.{.decimal_literal});
    try testTokenize("8", &.{.decimal_literal});
    try testTokenize("9", &.{.decimal_literal});

    try testTokenize("0_0", &.{.decimal_literal});
    try testTokenize("0001", &.{.decimal_literal});
    try testTokenize("01234567890", &.{.decimal_literal});
    try testTokenize("012_345_6789_0", &.{.decimal_literal});
    try testTokenize("0_1_2_3_4_5_6_7_8_9_0", &.{.decimal_literal});
}

test "octal literals" {
    try testTokenize("0o0", &.{.octal_literal});
    try testTokenize("0o1", &.{.octal_literal});
    try testTokenize("0o2", &.{.octal_literal});
    try testTokenize("0o3", &.{.octal_literal});
    try testTokenize("0o4", &.{.octal_literal});
    try testTokenize("0o5", &.{.octal_literal});
    try testTokenize("0o6", &.{.octal_literal});
    try testTokenize("0o7", &.{.octal_literal});

    try testTokenize("0o01234567", &.{.octal_literal});
    try testTokenize("0o0123_4567", &.{.octal_literal});
    try testTokenize("0o01_23_45_67", &.{.octal_literal});
    try testTokenize("0o0_1_2_3_4_5_6_7", &.{.octal_literal});

    try testTokenizeInvalid("0o7.", &.{.octal_literal}, true);
    try testTokenizeInvalid("0o7.0", &.{.octal_literal}, true);

    try testTokenizeInvalid("0O0", &.{.decimal_literal}, true);
    try testTokenizeInvalid("0o_", &.{.octal_literal}, false);
    try testTokenizeInvalid("0o_0", &.{.octal_literal}, false);
    try testTokenizeInvalid("0o1_", &.{.octal_literal}, false);
    try testTokenizeInvalid("0o0__1", &.{.octal_literal}, false);
    try testTokenizeInvalid("0o0_1_", &.{.octal_literal}, false);
    try testTokenizeInvalid("0o1e", &.{.octal_literal}, true);
    try testTokenizeInvalid("0o1e0", &.{.octal_literal}, true);
    try testTokenizeInvalid("0o_,", &.{ .octal_literal, .comma }, false);
}

test "hexadecimal literals" {
    try testTokenize("0x0", &.{.hexadecimal_literal});
    try testTokenize("0x1", &.{.hexadecimal_literal});
    try testTokenize("0x2", &.{.hexadecimal_literal});
    try testTokenize("0x3", &.{.hexadecimal_literal});
    try testTokenize("0x4", &.{.hexadecimal_literal});
    try testTokenize("0x5", &.{.hexadecimal_literal});
    try testTokenize("0x6", &.{.hexadecimal_literal});
    try testTokenize("0x7", &.{.hexadecimal_literal});
    try testTokenize("0x8", &.{.hexadecimal_literal});
    try testTokenize("0x9", &.{.hexadecimal_literal});
    try testTokenize("0xa", &.{.hexadecimal_literal});
    try testTokenize("0xb", &.{.hexadecimal_literal});
    try testTokenize("0xc", &.{.hexadecimal_literal});
    try testTokenize("0xd", &.{.hexadecimal_literal});
    try testTokenize("0xe", &.{.hexadecimal_literal});
    try testTokenize("0xf", &.{.hexadecimal_literal});
    try testTokenize("0xA", &.{.hexadecimal_literal});
    try testTokenize("0xB", &.{.hexadecimal_literal});
    try testTokenize("0xC", &.{.hexadecimal_literal});
    try testTokenize("0xD", &.{.hexadecimal_literal});
    try testTokenize("0xE", &.{.hexadecimal_literal});
    try testTokenize("0xF", &.{.hexadecimal_literal});

    try testTokenize("0x0000", &.{.hexadecimal_literal});
    try testTokenize("0xAA", &.{.hexadecimal_literal});
    try testTokenize("0xFFFF", &.{.hexadecimal_literal});

    try testTokenize("0x0123456789ABCDEF", &.{.hexadecimal_literal});
    try testTokenize("0x0123_4567_89AB_CDEF", &.{.hexadecimal_literal});
    try testTokenize("0x01_23_45_67_89AB_CDE_F", &.{.hexadecimal_literal});
    try testTokenize("0x0_1_2_3_4_5_6_7_8_9_A_B_C_D_E_F", &.{.hexadecimal_literal});

    try testTokenizeInvalid("0X0", &.{.decimal_literal}, true);
    try testTokenizeInvalid("0x_", &.{.hexadecimal_literal}, false);
    try testTokenizeInvalid("0x_1", &.{.hexadecimal_literal}, false);
    try testTokenizeInvalid("0x1_", &.{.hexadecimal_literal}, false);
    try testTokenizeInvalid("0x0__1", &.{.hexadecimal_literal}, false);
    try testTokenizeInvalid("0x0_1_", &.{.hexadecimal_literal}, false);
    try testTokenizeInvalid("0x_,", &.{ .hexadecimal_literal, .comma }, false);

    try testTokenizeInvalid("0x1.0", &.{.hexadecimal_literal}, true);
    try testTokenizeInvalid("0xF.0", &.{.hexadecimal_literal}, true);
    try testTokenizeInvalid("0xF.F", &.{.hexadecimal_literal}, true);

    try testTokenizeInvalid("0x1.", &.{.hexadecimal_literal}, true);
    try testTokenizeInvalid("0xF.", &.{.hexadecimal_literal}, true);
}

test "raw string literals" {
    try testTokenizeValue("`abc`", .raw_string_literal, "`abc`", false);
    try testTokenizeValue("`abc``def`", .raw_string_literal, "`abc``def`", false);
    try testTokenizeValue("`abc``", .raw_string_literal, "`abc``", true);
}

fn testTokenizeInvalid(source: []const u8, expected_token_typs: []const TokenType, invalid: bool) !void {
    var tokenizer = Tokenizer.init(source, "test.zig");
    for (expected_token_typs) |expected_token_tag| {
        const token = tokenizer.next();
        try std.testing.expectEqual(expected_token_tag, token.typ);
        try std.testing.expectEqual(invalid, token.invalid);
    }
    const last_token = tokenizer.next();
    try std.testing.expectEqual(TokenType.eof, last_token.typ);
    try std.testing.expectEqual(source.len + 1, last_token.column);
}

fn testTokenize(source: []const u8, expected_token_typs: []const TokenType) !void {
    var tokenizer = Tokenizer.init(source, "test.zig");
    for (expected_token_typs) |expected_token_tag| {
        const token = tokenizer.next();
        try std.testing.expectEqual(expected_token_tag, token.typ);
        try std.testing.expect(!token.invalid);
    }
    const last_token = tokenizer.next();
    try std.testing.expectEqual(TokenType.eof, last_token.typ);
    try std.testing.expectEqual(source.len + 1, last_token.column);
}

fn testTokenizeValue(source: []const u8, token_typ: TokenType, value: []const u8, invalid: bool) !void {
    var tokenizer = Tokenizer.init(source, "test.zig");
    const token = tokenizer.next();
    try std.testing.expectEqual(token_typ, token.typ);
    try std.testing.expectEqual(invalid, token.invalid);
    try std.testing.expectEqualStrings(value, token.value);
    const last_token = tokenizer.next();
    try std.testing.expectEqual(TokenType.eof, last_token.typ);
    try std.testing.expectEqual(source.len + 1, last_token.column);
}
