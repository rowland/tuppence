const std = @import("std");
const Token = @import("token.zig").Token;
const Tokenizer = @import("tokenizer.zig").Tokenizer;

pub fn tokenize(allocator: std.mem.Allocator, source: [:0]const u8, filename: []const u8) ![]Token {
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
