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

        interpolated_string_literal,
        octal_literal,
        raw_string_literal,
        string_literal,
    };

    pub fn next(self: *Tokenizer) Token {
        if (self.pending_invalid_token) |token| {
            self.pending_invalid_token = null;
            return token;
        }
        var state: State = .start;
        var start = self.index;
        var typ: TokenType = .eof;
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
                    'a'...'z', 'A'...'Z', '_' => {
                        typ = .identifier;
                        state = .identifier;
                    },

                    '"' => {
                        state = .string_literal;
                        typ = .string_literal;
                    },
                    '0'...'9' => {
                        // state = .int;
                        // typ = .number_literal;
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
                    else => {
                        break;
                    },
                },
                .op_div => switch (c) {
                    '=' => {
                        typ = .op_div_equal;
                        self.index += 1;
                        break;
                    },
                    else => {
                        break;
                    },
                },
                .op_minus => switch (c) {
                    '=' => {
                        typ = .op_minus_equal;
                        self.index += 1;
                        break;
                    },
                    else => {
                        break;
                    },
                },
                .op_mod => switch (c) {
                    '=' => {
                        typ = .op_mod_equal;
                        self.index += 1;
                        break;
                    },
                    else => {
                        break;
                    },
                },
                .op_mul => switch (c) {
                    '=' => {
                        typ = .op_mul_equal;
                        self.index += 1;
                        break;
                    },
                    else => {
                        break;
                    },
                },
                .op_not => switch (c) {
                    '=' => {
                        typ = .op_not_equal;
                        self.index += 1;
                        break;
                    },
                    else => {
                        break;
                    },
                },
                .op_plus => switch (c) {
                    '=' => {
                        typ = .op_plus_equal;
                        self.index += 1;
                        break;
                    },
                    else => {
                        break;
                    },
                },
                .op_pow => switch (c) {
                    '=' => {
                        typ = .op_pow_equal;
                        self.index += 1;
                        break;
                    },
                    else => {
                        break;
                    },
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
                    else => {
                        break;
                    },
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
                    else => {
                        break;
                    },
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
                    else => {
                        break;
                    },
                },
                .op_logical_and => switch (c) {
                    '=' => {
                        typ = .op_logical_and_equal;
                        self.index += 1;
                        break;
                    },
                    else => {
                        break;
                    },
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
                    else => {
                        break;
                    },
                },
                .op_logical_or => switch (c) {
                    '=' => {
                        typ = .op_logical_or_equal;
                        self.index += 1;
                    },
                    else => {
                        break;
                    },
                },
                .op_equal => switch (c) {
                    '=' => {
                        typ = .op_equal_equal;
                        self.index += 1;
                    },
                    else => {
                        break;
                    },
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
                .interpolated_string_literal => {},
                .octal_literal => {},
                .raw_string_literal => {},
                .string_literal => {},
            }
        }
        return Token{
            .typ = typ,
            .line = self.line,
            .column = std.math.cast(u32, self.index - self.bol + 1) orelse 0,
            .value = self.source[start..self.index],
            .filename = self.filename,
        };
    }
};

fn testTokenize(source: []const u8, expected_token_tags: []const TokenType) !void {
    var tokenizer = Tokenizer.init(source, "test.zig");
    for (expected_token_tags) |expected_token_tag| {
        const token = tokenizer.next();
        try std.testing.expectEqual(expected_token_tag, token.typ);
    }
    const last_token = tokenizer.next();
    try std.testing.expectEqual(TokenType.eof, last_token.typ);
    try std.testing.expectEqual(source.len, last_token.column - 1);
}

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
