const std = @import("std");
const TokenType = @import("token_type.zig").TokenType;

pub const Token = struct {
    typ: TokenType,
    line: u32,
    column: u32,
    value: []const u8,
    filename: []const u8,

    pub const renamed = std.ComptimeStringMap(TokenType, .{
        .{ "break", .keyword_break },
        .{ "else", .keyword_else },
        .{ "enum", .keyword_enum },
        .{ "error", .keyword_error },
        .{ "false", .boolean_literal },
        .{ "fn", .keyword_fn },
        .{ "fx", .keyword_fx },
        .{ "for", .keyword_for },
        .{ "if", .keyword_if },
        .{ "import", .keyword_import },
        .{ "mut", .keyword_mut },
        .{ "return", .keyword_return },
        .{ "switch", .keyword_switch },
        .{ "true", .boolean_literal },
        .{ "try", .keyword_try },
        .{ "type", .keyword_type },
        .{ "typeof", .keyword_typeof },
    });

    pub fn getReserved(bytes: []const u8) ?TokenType {
        return renamed.get(bytes);
    }
};
