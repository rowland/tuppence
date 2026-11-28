# IR implementation punch list

This is the implementation checklist for the intermediate representation described in
[IR_DESIGN.md](IR_DESIGN.md) and [IR_OWNERSHIP.md](IR_OWNERSHIP.md). It is ordered by
dependencies rather than by document section. A box should be checked only when the
implementation, focused tests, and readable dumps or diagnostics are present.

## Working rules

- Keep semantic IR and executable SSA as separate representations.
- Keep published snapshots immutable; builders are private and mutable until frozen.
- Keep source spans and origins on synthesized entities so diagnostics survive lowering.
- Keep logical member order separate from physical storage order.
- Do not use `any`, nil type pointers, or missing terminators to represent incomplete IR.
- Add invalid-fixture tests for every verifier invariant, not only valid examples.
- Keep parser goldens separate from IR goldens; update `TODO.md` only for covered work.

## Phase 0 — settle the contract before coding

- [ ] Approve the two-stage architecture: semantic snapshots followed by typed monomorphic SSA.
- [ ] Record the initial supported source subset and explicit non-goals for the first executable milestone.
- [ ] Decide the initial rules for literal defaulting and the solver boundary for first-use inference.
- [ ] Decide block declaration visibility and module-initialization ordering/cycle diagnostics.
- [ ] Decide callable compatibility for parameter names, defaults, rest members, and function references.
- [ ] Decide overload ranking for union parameters, generic candidates, structural adaptation, and mixed specificity.
- [ ] Decide expected-result selection, implicit conversion boundaries, and finite-union dispatch limits.
- [ ] Decide initial effect/divergence/panic policy for optimization and cleanup eligibility.
- [ ] Decide the first runtime representation: tag width, canonical payload layout, and supported union forms.
- [ ] Record unresolved ownership policy for cleanup effects, cycles, weak/foreign references, and unwind behavior.

## Phase 1 — semantic foundation

- [ ] Add `tup/ir` identity types with zero reserved as invalid: modules, scopes, declarations,
      bindings, expressions, types, type variables, functions, files, origins, and spans.
- [ ] Add immutable source snapshots and origin/provenance tables based on existing AST offsets.
- [ ] Add diagnostic types with severity, source origin, stable code, and related locations.
- [ ] Add private table builders and validated freeze operations; prevent writable backing slices from escaping.
- [ ] Add scopes, scope entries, declarations, visibility, and declaration collection independent of body checking.
- [ ] Add binding versions for same-scope rebinding, previous-version links, and shared mutable type slots.
- [ ] Add semantic expression alternatives without collapsing unresolved forms into executable instructions.
- [ ] Add indexed, analyzed, and checked snapshot wrappers with validated transitions.
- [ ] Add deterministic semantic formatting/golden dumps.
- [ ] Test same-scope rebinding, right-hand-side lookup, inner shadowing, loop-carried values, and snapshot isolation.
- [ ] Test diagnostics and definition/origin links on incomplete programs for language-server use.

## Phase 2 — types, literals, and a checked subset

- [ ] Add explicit primitive, nominal, anonymous tuple, array, callable, finite-union, open-union,
      and generic-term representations.
- [ ] Add canonical type/tuple-schema interning without erasing nominal identity or recursive declarations.
- [ ] Add tuple schemas for named types, anonymous tuples, formal parameters, captures, defaults, and rest members.
- [ ] Add literal templates that preserve exact numeric syntax and source provenance.
- [ ] Add constraints for equality, literal admissibility, union inclusion, structural adaptation,
      callable compatibility, and required operations.
- [ ] Add dependency-graph/worklist handling for recursive signatures and inference variables.
- [ ] Add direct pure/effectful call plans and explicit return constraints for the initial subset.
- [ ] Add purity and cleanup-effect checks for the initial subset.
- [ ] Test literal instantiation, mutable-binding fixed types, exact numeric handling, closure capture identities,
      purity violations, and recursive-inference diagnostics.

## Phase 3 — name resolution and elaboration

- [ ] Add candidate discovery with lexical locations, receiver home modules, explicit qualification, and source origins.
- [ ] Implement exact, generic, then `any` candidate priority within each location group.
- [ ] Preserve rejection reasons and distinguish preferred-location ambiguity from outer-scope fallback.
- [ ] Implement tuple construction plans for labels, positional arguments, spreads, defaults, rest packing,
      structural adaptation, and explicit conversions.
- [ ] Add generic requirements with definition-site origins and specialization keys.
- [ ] Add recursive-specialization guards and useful diagnostics for unbounded specialization growth.
- [ ] Add UFCS, operator, constructor, indirect-call, and callable-field resolution rules.
- [ ] Add finite-union dispatch without treating `any` or `error` as finite unions.
- [ ] Add partial-application plans that preserve remaining defaults and slot mappings.
- [ ] Add closure plans with explicit captured value identities and home environments.
- [ ] Add accepted, ambiguous, rejected, and source-order/effect-order fixtures for each resolution rule.

## Phase 4 — typed monomorphic SSA

- [ ] Add function-local value/block IDs and immutable SSA function/program snapshots.
- [ ] Add typed basic blocks with parameters, instructions, terminators, and origins.
- [ ] Add initial instructions: constants, direct/indirect calls, tuple make/get/update, array make/get/update,
      union inject/tag/project, closure make, and primitive operations.
- [ ] Add jump, branch, switch, return, and unreachable terminators with edge arguments.
- [ ] Add compiler-only effect tokens for effectful calls and effect-token block parameters at joins.
- [ ] Lower the initial checked subset into SSA, including explicit block results and loop-carried values.
- [ ] Lower `break`, `continue`, `return`, and `try` exits with their cleanup obligations represented.
- [ ] Add deterministic SSA formatting/golden dumps.
- [ ] Implement `VerifySSA`: ownership of IDs, definitions, dominance, targets, types, edge arity/types,
      complete terminators, union projection proofs, effect chains, purity, and origins.
- [ ] Add deliberately invalid SSA fixtures for each verifier failure class.

## Phase 5 — baseline C execution

- [ ] Choose the smallest executable runtime subset and define its primitive arithmetic/failure semantics.
- [ ] Add runtime type metadata and deterministic reachable-type IDs.
- [ ] Add baseline layouts with logical-to-physical field mappings and alignment checks.
- [ ] Add ABI planning separate from semantic tuple types and overload resolution.
- [ ] Lower block arguments with parallel-copy semantics and temporaries where needed.
- [ ] Emit sequenced C for effectful calls; do not depend on C argument evaluation order.
- [ ] Emit a minimal label/goto backend and compiler-owned runtime.
- [ ] Compile and run representative generated programs; compare values and effect traces.
- [ ] Add sanitizer-backed tests where supported.

## Phase 6 — ownership, storage, and cleanup

- [ ] Add ownership/storage summaries: parameter modes, result ownership, placement, escapes, and cleanup effects.
- [ ] Add storage plans, layout descriptors, ownership shapes, allocation identities, and cleanup identities.
- [ ] Implement conservative borrowed-input/owned-result lowering before consumed-input inference.
- [ ] Implement explicit copy/share, borrow, move, drop, and reuse operations.
- [ ] Lower cleanup on all scope exits, branches, loops, returns, breaks, continues, and error paths.
- [ ] Verify counts never underflow and every ownership obligation is discharged exactly once.
- [ ] Add ownership traces for stack aggregates, heap aggregates, duplicate arguments, nested aliases,
      closures, partial applications, returned values, and early exits.
- [ ] Add tuple-update tests covering field exchange, shared parents, unique parents, and shared children.
- [ ] Add array append/swap/quicksort tests covering growth, old-buffer element sources, aliases, and clone counts.
- [ ] Add union inline/boxed representation and layout-conversion tests.
- [ ] Keep resource cleanup ordering conservative until the final cleanup-effect policy is accepted.

## Phase 7 — optimization and representation expansion

- [ ] Add constant folding using Tuppence arithmetic rules.
- [ ] Add CFG simplification and dead pure-value removal under the selected divergence policy.
- [ ] Add local CSE keyed by resolved callee identity and type arguments.
- [ ] Add ownership-aware retain/drop cancellation and proven last-use moves.
- [ ] Add escape analysis and safe storage reuse after ownership verification is reliable.
- [ ] Add richer layout alternatives only after measuring the baseline: split tags, separate tags, and optional SoA.
- [ ] Add inlining/global value numbering only after provenance and ownership rewrites are verified.
- [ ] Keep unsupported cycles, weak references, foreign references, and threading explicit in diagnostics or design.

## Cross-cutting acceptance criteria

- [ ] Every stage publishes a stated invariant and has a verifier or validated constructor.
- [ ] Every transformation preserves source provenance, or records a derived origin.
- [ ] Every pass has focused tests plus readable semantic/SSA golden output where applicable.
- [ ] Every invalid graph or unresolved semantic state produces structured diagnostics rather than a panic.
- [ ] `make test` passes for every implementation milestone.
- [ ] Documentation and specification are reconciled when an open language rule becomes accepted.
- [ ] `TODO.md` is updated only after the corresponding behavior is implemented and covered.
