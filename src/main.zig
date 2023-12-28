const std = @import("std");
const clap = @import("clap");

const debug = std.debug;

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
        diag.report(std.io.getStdErr().writer(), err) catch {};
        return err;
    };
    defer res.deinit();

    // Use the parsed arguments
    if (res.args.input) |input|
        debug.print("Input file: '{s}'\n", .{input});
    if (res.args.output) |output|
        debug.print("Output file: '{s}'\n", .{output});
    for (res.positionals, 0..) |pos, i|
        debug.print("{}: {s}\n", .{ i, pos });
}
