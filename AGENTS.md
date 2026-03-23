# AGENTS.md

This file is a short working guide for agents making changes in this repository.

## Repository Layout

- `SPECIFICATION.md`
  Human-readable language specification. Start here for intended syntax and examples.

- `tup.ebnf`
  Formal grammar. Keep this in sync with parser behavior whenever precedence or accepted syntax changes.

- `tup.ebnf.html`
  Generated HTML rendering of the grammar. Regenerate it with `make grammar` after changing `tup.ebnf`.

- `examples/`
  Curated language examples. These are not throwaway samples; they help define intended behavior and style.

- `TODO.md`
  Active implementation checklist and punchlist. Update it when rules are truly implemented and covered.

- `IMPLEMENTATION.md`
  Backend/runtime notes and design context. Useful, but usually secondary to the spec/grammar when doing parser work.

- `tup/`
  Go module containing the tokenizer, parser, AST, and tests.

## Project Structure Inside `tup/`

- `tup/ast/`
  AST node definitions and interface tests.

- `tup/parse/`
  Parser implementation, direct parser tests, and top-level golden tests.

- `tup/tok/`
  Tokenizer.

Most language work happens in these three directories.

## Normal Workflow

When changing language syntax or parser behavior:

1. Start from `SPECIFICATION.md`, `tup.ebnf`, existing parser code, and nearby tests.
2. Keep grammar, parser, AST, and tests aligned.
3. Add focused direct tests for the specific rule being implemented.
4. Add or update curated top-level golden coverage when the new syntax is user-facing.
5. Update examples and spec text if the intended language surface changed.
6. Check off the relevant items in `TODO.md` only when the rule is genuinely implemented.

## Testing and Regeneration

From the repo root:

- `make grammar`
  Regenerates `tup.ebnf.html` from `tup.ebnf`.

- `make test`
  Runs the Go test suite in the `tup` module.

- `make goldens`
  Regenerates top-level parser goldens.

- `make test-goldens`
  Verifies top-level parser goldens without rewriting them.

Focused commands are often useful from `tup/`, for example:

- `go test ./parse`
- `go test ./ast`
- `go test ./tok ./parse`
- `UPDATE_TOP_LEVEL_GOLDENS=1 go test ./parse -run TestTopLevelGoldenFixtures`

If `go test` hits cache-permission issues in an automated environment, retry with:

- `GOCACHE=/tmp/codex-gocache go test ./...`

## Goldens

Top-level golden fixtures live in:

- `tup/parse/testdata/top_level/input`
- `tup/parse/testdata/top_level/output`

Guidelines:

- Prefer a small number of representative, readable examples.
- Keep examples intentionally curated rather than exhaustive.
- Use comments to explain the purpose of a golden case when helpful.
- When a new syntax shape is important, make sure it appears in the goldens.

## Conventions for Parser and AST Work

- Prefer grammar-shaped AST nodes and interfaces over catch-all `Node` fields where practical.
- If a parser rule corresponds to a real grammar rule, add or keep a rule comment above the parse function.
- When a grammar rule is represented by an interface, use interface embedding to reflect grammar inclusion relationships when possible.
- Treat parser-phase validation and later semantic validation separately:
  the parser should accept valid syntax without trying to enforce type-system or exhaustiveness rules.

## Practical Advice

- Use `rg` and `rg --files` for code search.
- Look for existing tests before inventing new patterns.
- Keep changes local and mechanical when possible.
- Do not treat generated or legacy AST shapes as automatically authoritative; prefer the current grammar and spec.
- If parser behavior changes, check whether comments, `tup.ebnf`, examples, and goldens also need to change.

## Good First Places to Look

For any language feature, the most useful files are usually:

- `SPECIFICATION.md`
- `tup.ebnf`
- the relevant files in `tup/parse`
- the relevant files in `tup/ast`
- nearby tests in `tup/parse/*_test.go`
- `tup/parse/testdata/top_level`

If those disagree, prefer the current intended language design, then bring the implementation and docs back into sync.
