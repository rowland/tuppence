# tuppence
The Tuppence Programming Language

- [Specification](SPECIFICATION.md)
- [FAQ](FAQ.md)

## Repository Guide

- [SPECIFICATION.md](SPECIFICATION.md) is the human-readable language guide.
- [tup.ebnf](tup.ebnf) is the formal grammar.
- [tup.ebnf.html](tup.ebnf.html) is the generated HTML view of the grammar.
- [examples](examples) contains valid and intentionally invalid examples; pay attention to the comments.
- [tup](tup) contains the tokenizer, parser, AST, and related tests.

## Common Tasks

Update the generated grammar HTML:

```sh
make grammar
```

Run the Go test suite:

```sh
make test
```

Refresh the curated parser golden outputs:

```sh
make goldens
```

The parser goldens live under:

- [tup/parse/testdata/top_level/input](tup/parse/testdata/top_level/input)
- [tup/parse/testdata/top_level/output](tup/parse/testdata/top_level/output)

Each `.tup` input file contains curated top-level items labeled with `# case name` comments. The
matching output file contains the current AST serialization for the same cases. These fixtures are
intended to be curated, high-value examples of what the parser can currently parse; they are not
meant to be exhaustive.
