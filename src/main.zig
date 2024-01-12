const std = @import("std");
const debug = std.debug;
const stderr = std.io.getStdErr().writer();
const clap = @import("clap");
const tokenize = @import("tokenizer.zig").tokenize;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    var allocator = gpa.allocator();

    const params = comptime clap.parseParamsComptime(
        \\-h, --help         Display this help and exit.
        \\-v, --version      Output version information and exit.
        \\-i, --input <str>  Input file.
        \\-o, --output <str> Output file.
        \\<string>...
        \\
    );

    var diag = clap.Diagnostic{};
    var res = clap.parse(clap.Help, &params, clap.parsers.default, .{
        .diagnostic = &diag,
        .allocator = allocator,
    }) catch |err| {
        // Report useful error and exit
        diag.report(stderr, err) catch {};
        return err;
    };
    defer res.deinit();

    // Use the parsed arguments
    if (res.args.input) |input| {
        debug.print("Input file: '{s}'\n", .{input});
        const source = try readFile(allocator, input);
        defer allocator.free(source);
        const tokens = try tokenize(allocator, source, input);
        defer allocator.free(tokens);
        debug.print("tokens: {any}\n", .{tokens});
    }
    if (res.args.output) |output| {
        debug.print("Output file: '{s}'\n", .{output});
    }
    for (res.positionals, 0..) |pos, i| {
        debug.print("{}: {s}\n", .{ i, pos });
    }
}

pub fn readFile(allocator: std.mem.Allocator, path: []const u8) ![]u8 {
    const file = try std.fs.cwd().openFile(path, .{});
    defer file.close();

    const file_size = try file.getEndPos();
    const buffer = try allocator.alloc(u8, file_size);
    _ = try file.readAll(buffer);

    return buffer;
}

test "run all tests" {
    std.testing.refAllDeclsRecursive(@import("tokenizer.zig"));
    try std.testing.expectEqual(1, 1);
}
