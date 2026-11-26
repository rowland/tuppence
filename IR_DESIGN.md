# Tuppence intermediate representation proposal

Draft 1, revised 2026-09-08 for tuple-based parameters and construction. For discussion; this document does not authorize or implement
language, parser, runtime, or backend changes.

This proposal records the design discussion and recommends an implementation
structure. **Agreed** means established in that discussion. **Proposed** means a
recommendation awaiting review. **Open** means a decision is still needed. Existing
specification text is not silently superseded: discrepancies are called out below
and should be reconciled when the corresponding design is accepted.

The companion [ownership and layout proposal](IR_OWNERSHIP.md) examines allocation,
union representations, and the proposed strong-count and call-transfer protocol in more detail.

## 1. Recommendation

Introduce a `tup/ir` package containing two related representations:

1. A semantic representation that retains scopes, declarations, structured
   expressions, generic requirements, and inference constraints. It has explicit
   value identities early, but can represent unresolved information.
2. A typed, monomorphic SSA representation with basic blocks, block parameters,
   explicit calls, and explicit control flow. This is the input to optimization
   and eventually ownership and storage lowering.

Do not require every semantic-analysis intermediate state to be an executable SSA
program. In particular, unresolved overload sets, reusable literals, and generic
requirements are meaningful compiler objects, but not executable instructions.

Use immutable published snapshots and private mutable builders. Keep semantic
type identity, runtime type-table indices, and physical storage layouts separate.
Target C first; keep the semantic and SSA layers independent of C syntax and ABI.

This is a proposal for related stages, not a request to implement every stage at
once. The first coding milestone should establish IDs, semantic tables, binding
versions, and diagnostics; executable SSA follows once a small typed subset works.

## 2. Established constraints

| Area | Agreed direction |
| --- | --- |
| Backend | Whole-program compilation to C in v1. Independently linkable Tuppence libraries are a non-goal. LLVM/WASM are later possibilities. |
| Value semantics | No mutation visible through aliases; no language-level allocation identity. |
| `mut` | Sugar for successive bindings of one fixed type, with rebinding allowed only in the exact declaring scope. Inner bindings shadow outer names. |
| Blocks | Values leave inner scopes through block results or explicit control-flow exits, never by assigning an outer binding. |
| Closures | Capture snapshots of values, not mutable binding cells. |
| Call effects | `fn` may call `fn` by name or reference. It may store an `fx` callable but cannot invoke it. `fx` may invoke either. |
| Parameters | Formal parameters form a tuple definition, including names and defaults. Actual arguments form a tuple construction. Anonymous tuples also have generated types, constructors, and enumerable metadata. |
| Effect ordering | Invocations of `fx` inside `fx` preserve their source order. Pure calls may be deduplicated and reordered subject to dependencies and the eventual failure/divergence rules. |
| Literals | An immutable binding whose value is a literal remains a reusable literal template: each use may instantiate it under the same rules as the literal written at that use. Integer and float literals/templates begin with their category's admissible storage-type superset, including byte arrays. Mutable bindings, function-call results, and values with an established concrete type do not behave this way. |
| Literal defaults | If one non-generic candidate remains applicable, literal/template storage types are tried in the specified priority order through the candidate's type constructors. If multiple candidates remain, the literal/template assumes the configured default, such as `Int` or `Float` when out of range. |
| Conversions | A known concrete value is not implicitly converted merely because another parameter type would fit. If a suitable existing, non-lossy cast function is in scope, the compiler may insert a call to it; the compiler does not synthesize the cast function. Anonymous structures may adapt to an expected named type; the reverse is not implicit. |
| Lookup | UFCS searches the receiver type's home module before lexical scopes. Ordinary calls search inward to outward. Within a location, candidates are considered in priority tiers: exact concrete matches, generic matches, then `any` matches. A lower tier cannot displace an applicable higher tier. |
| Generic environment | Specializations retain the generic declaration's home module and lexical environment. Caller-local functions are not injected into that environment. |
| Contracts | Compile-time requirements, including requirements inferred from operations on generic parameters. No runtime interface or method tables. |
| Unions | Finite unions may dispatch by generated switches when all alternatives have applicable calls. `any` and `error` are open unions; switches over them require an `else`. |
| Return inference | Infer when unambiguous; require annotations when recursion or other dependencies make inference indeterminate. |
| Types | Only internal primitives have hard-coded type identities. Standard types are language-defined, although the compiler may know their roles. |
| Placement | One semantic type may have stack or heap storage in different contexts. Heap references in stack values still carry ownership obligations. |
| Layout | Logical member order is distinct from physical order. Sort physical members by decreasing storage size, then name or logical slot, respecting alignment. |
| Runtime tags | Program type IDs identify concrete union alternatives; a representation flag may distinguish inline payloads from heap pointers. Encoding remains provisional. |
| Future layout | Separate array tag regions and struct-of-arrays are stretch goals. Logical array elements need not have contiguous storage. |

## 3. Pipeline and stage guarantees

```text
AST + source snapshots
  -> declaration collection and semantic indexing
  -> name/candidate discovery + constraint solving + generic requirements
  -> checked semantic snapshot
  -> reachable specialization and compile-time elaboration
  -> typed monomorphic SSA
  -> value/control-flow optimization
  -> closure, escape, ownership, and representation planning
  -> explicit ownership + storage operations
  -> ownership-aware cleanup optimization and verification
  -> C lowering, runtime metadata, and C emission
```

This is a dependency outline, not necessarily one linear pass per line. Name
resolution and typing require a worklist. Specializing a generic body may discover
more specializations. Layout and escape decisions may require bounded iteration.
Each stage must nevertheless publish a result with a stated invariant.

Effectful-cleanup eligibility must be checked before value optimization, even
though the actual retain/release instructions arrive later. Otherwise dead-value
removal could discard a resource's required cleanup before ownership lowering
ever sees it. Carry abstract cleanup obligations/effect summaries through early
stages; only their concrete implementation is deferred.

| Published result | What it guarantees |
| --- | --- |
| Indexed semantic snapshot | Source identities, scope graph, declarations, and expression identities exist; unresolved references are explicit. |
| Analyzed semantic snapshot | Available resolutions and types are recorded with diagnostics; pending/error states remain usable by the language server. |
| Checked program | All reachable code is valid, except generic operations represented by checked requirements to discharge at specialization. |
| Monomorphic program | Every reachable executable function has concrete parameter/result types and resolved callees; no runtime generics or unresolved literals. |
| SSA program | Definitions dominate uses, edge arguments match block parameters, effects are ordered, and every block has a terminator. |
| Storage program | ABI, layouts, allocation classes, ownership transfers, and cleanup paths are explicit and verified. |

Do not encode stage completeness using an `any` type or a nil type pointer.
Unknown, erroneous, generic, open-union, and concrete types are different concepts.

### Package direction

Proposed initial organization, all under the new `tup/ir` sub-package:

```text
ir/
  ids.go, source.go, diagnostics.go
  scopes.go, declarations.go, bindings.go
  types.go, constants.go, constraints.go, calls.go
  semantic.go, builder.go, verify_semantic.go, format_semantic.go
  ssa.go, instructions.go, terminators.go
  verify_ssa.go, format_ssa.go
```

Add files as their milestones arrive; do not create empty shells. Semantic passes
can initially live in this package with unexported implementation details. A later
`sema` package may own orchestration when a useful boundary emerges. C emission
belongs in a separate backend package. `ast` and `parse` must not import `ir`.

## 4. Identities, snapshots, and builders

The following Go sketches illustrate shapes, not a compilable API. Supporting
types and closed alternatives are described in prose rather than fully expanded.
Exported fields in these sketches show content for review. The implementation
should keep published mutable containers private as described below.

```go
type ModuleID uint32
type ScopeID uint32
type DeclID uint32
type BindingID uint32
type ExprID uint32
type TypeID uint32       // Compiler identity; not the emitted runtime tag.
type TypeVarID uint32
type FunctionID uint32
type FileID uint32
type OriginID uint32

type Span struct {
    File FileID
    Start, End uint32    // Half-open byte offsets in a source snapshot.
}

type Origin struct {
    Primary Span
    DerivedFrom []OriginID
}
```

Reserve zero for invalid IDs. IDs are meaningful in their owning snapshot; SSA
value/block IDs additionally belong to one function. A pass preserving an entity
preserves its ID where practical; cloning gets a new ID with provenance. IDs need
not survive source edits. Verification rejects out-of-range and cross-owner IDs.

Published tables use private slices/maps with query methods. Go does not enforce
deep immutability: constructors copy or take exclusive ownership of input slices,
and getters must not expose writable backing arrays. Builders may append and patch
private state, then freeze it once. They cannot mutate a published result afterward.

Start with copying small tables and sharing unchanged function bodies. Persistent
map libraries and cross-edit caching are not prerequisites. Pass-local caches,
worklists, union-find structures, and predecessor/use indexes may be mutable and
are not part of the published language model.

Keep byte spans derived from the AST's existing source offsets. An IR origin table
can preserve both a generic definition and its specialization/call site, or an
operator and its elaborated call. Do not rely on AST pretty-print output for
semantic identity. Retain source snapshots for diagnostics; never reread a changed
file to interpret an old span.

## 5. Semantic representation

### Scopes and binding versions

```go
type Scope struct {
    ID ScopeID
    Parent ScopeID
    Module ModuleID
    Kind ScopeKind              // Module, function, block, loop, etc.
    Entries []ScopeEntry         // Name, declaration, visibility/source position.
}

type Binding struct {
    ID BindingID
    Scope ScopeID
    Name string
    Previous BindingID          // Previous same-scope mutable version, if any.
    TypeSlot TypeVarID          // Shared constraint slot for mutable versions.
    Value ExprID
    Mutable bool
    Origin OriginID
}
```

Every assignment creates a binding identity. The same-scope name environment is
updated only after the right-hand side is resolved. Thus `x = x + 1` reads the old
version. In an inner scope, the right-hand side can read an outer `x`, but the new
binding belongs to the inner scope. It never changes the parent's environment.

An immutable binding that denotes a reusable literal uses a literal template,
rather than forcing a runtime SSA value immediately. All versions of a mutable
binding share one type constraint slot even when their value identities differ.

Loop-carried values are explicit loop inputs/results, not exceptions permitting
inner assignment to outer state. Liveness diagnostics operate on source binding
versions before optimization; an optimizer removing a use must not manufacture a
new source-level "never used" error. The exact unused-binding exemptions remain
open, including exported values and deliberately discarded effects.

### Declarations and visibility

Collect function/type declarations independently of checking bodies. Proposed
block policy: functions/types are visible throughout their scope, ordinary values
only after initialization. This policy is not yet a settled language rule.

Provisionally allow forward references to module values. Record both initializer
dependencies and source order. A module name being resolvable does not imply its
value is initialized. An initializer dependency graph must distinguish creating a
closure from executing its body; a capture may still require an initialized value.

Module cycles, cross-file order, and effectful initialization scheduling are open.
Do not silently topologically reorder effectful initializers. Do not replace this
problem with uninitialized C globals. Publish the graph before choosing its policy.

### Semantic expressions and analysis state

Use closed, typed alternatives for expressions: literal template, binding use,
aggregate construction, member selection, call site, block, conditional, loop,
match, closure, and nonlocal exit. Preserve structured scopes and argument source
order while analysis is incomplete.

```go
type TypeState interface { typeState() }
// PendingType{Variable}, KnownType{Type}, GenericType{Term}, ErrorType{Diagnostic}

type Resolution interface { resolution() }
// PendingCall{Candidates, Constraints}, ResolvedCall{Plan}, InvalidCall{Diagnostic}

type AnalysisSnapshot struct {
    // Private indexed syntax/semantic tables shared with earlier snapshots.
    types       map[ExprID]TypeState
    resolutions map[ExprID]Resolution
    diagnostics []Diagnostic
}
```

Use separate wrapper types such as `IndexedProgram`, `AnalyzedProgram`, and
`CheckedProgram`, with private constructors and validated transitions. A checked
result cannot be obtained by setting a public boolean. The language server may
query candidate sets, definition links, inferred types, and diagnostics without
requiring successful C compilation.

## 6. Types, literals, and constraints

### Type model

Represent primitive, nominal, anonymous tuple, array, callable, finite union,
open union (`any`/`error`), and generic type terms explicitly. A contract is a
constraint object, not an existential runtime object with a method table.

```go
type Field struct {
    Name string                 // Empty for ordinal fields.
    LogicalSlot uint32
    Type TypeID
    Origin OriginID
}

type TupleSchemaID uint32
type DefaultID uint32

type TupleMember struct {
    Field Field
    Default DefaultID           // Invalid when required; expression plus environment.
    Mode MemberMode             // Ordinary or rest parameter where permitted.
}

type TupleSchema struct {
    Type TypeID                 // Logical tuple type after type resolution.
    Members []TupleMember       // Logical order, including names and defaults.
    Constructor DeclID          // Compiler-generated constructor.
    Origin OriginID
}

type CallableType struct {
    Parameters TupleSchemaID    // Formal parameter tuple, not just member types.
    Result TypeID
    Effect EffectKind           // Pure or Effectful.
}

type NominalType struct {
    Declaration DeclID
    HomeModule ModuleID
    Arguments []TypeID
    Shape TypeID
}
```

Intern anonymous structures and canonical union sets. Nominal identity is not
erased just because two declarations have equal layouts. Aliases refer to the
aliased identity. Recursive nominal types use declaration identities before their
shapes are complete; interning must not require infinitely unfolding them.

### Tuple schemas and construction

Use the same tuple-definition and construction machinery for named types,
anonymous tuple literals, and function parameter tuples. A schema carries member
names, logical slots, types, default expressions, and constructor provenance.
During inference its member types may be type terms; the sketch above shows its
resolved form. Defaults retain their declaration environment and source origins.

Separate a tuple's structural type identity from its construction schema: two
schemas can produce the same structural value type while supplying different
defaults. This permits anonymous type interning without merging distinct default
behavior. Whether parameter names/default availability affect callable type
compatibility still needs a rule; retaining schemas prevents their accidental
erasure. In particular, calls through function references must have a defined way
to select defaults, not assume a directly named declaration is always available.

An anonymous tuple literal elaborates to a generated constructor for an interned
anonymous tuple type. Its metadata enumerates logical members just as for a named
tuple, supporting serialization and introspection. Compiler-generated parameter
and capture tuples use the same machinery; policy for exposing internal capture
names remains separate from preserving their metadata. Retain observable metadata
even if optimization eliminates the allocation or inlines the constructor.

Construction is a semantic operation, not a requirement for an out-of-line machine
call. Generated constructors may lower directly to `tuple.make`, constants, or
separate fields. A zero-argument call constructs the empty parameter tuple.
Constructor selection and actual user-defined constructor behavior must remain
explicit; only generated structural construction may be erased without examining
an arbitrary user-defined constructor body.

Union subset inclusion is distinct from equality and from conversion. A declared
result union may accept a subset; SSA edges may need explicit union injection to
make their block argument types equal. Named-versus-structural adaptation is also
distinct from numeric conversion and should appear in resolved call plans.

"Concrete" here does not require a singleton static type: a fully resolved finite
union, `any`, or `error` is a valid known static type. Each runtime member has a
concrete identity, while a generic type variable or pending inference variable
does not belong in executable SSA.

Internal primitives have reserved compiler/runtime identities. Standard types are
ordinary declarations referenced through a proposed `BuiltinRoles` configuration
(default integer, default float, boolean, string, rune, and required operations).
Do not infer privileged status solely from a declaration's spelling.

### Literal handling

Preserve exact numeric syntax or an exact coefficient/exponent representation;
never parse every number through host `float64` before selecting a target. Keep
literal admissibility, precision/rounding policy, and constructor selection
separate. A literal such as `12.34e+500` must not overflow merely because the
compiler's implementation language uses a smaller built-in floating type.

For an immutable direct literal binding:

```text
x = 1
takes32(x)  -> instantiate literal 1 as Int32
takes64(x)  -> instantiate literal 1 as Int64
```

The binding `x` above is therefore matched like the literal `1` written directly in
those call sites. This rule applies only while the binding remains an immutable,
direct literal template. A mutable binding, a function result, or a binding whose
concrete type has already been fixed is an ordinary known-typed value and cannot
silently change type at a later use.

A literal template begins as the superset of admissible storage types for its
category: integer and float templates may use the applicable signed, unsigned, or
floating storage types and an array of bytes, subject to range and representation
rules. If one non-generic function candidate remains applicable, resolve the
candidate's parameter through its type constructor for each admissible storage type
in the specified storage-type priority order. If multiple function candidates
remain, default the template to `Int`, or `Float` when the value is out of integer
range. Generic and `any` consumers use their configured defaults.

For a function result, independent instantiation does not apply. For a mutable
binding, every version must agree on one type. A uniquely determined concrete
consumer provides a constraint; a generic/`any` consumer selects the configured
default. No allocation is implied by binding a literal.

Open: the deterministic meaning of "first use" across control flow, how later
constraints interact with a default already selected, whether aliases of literal
templates remain reusable, and whether compound constant expressions qualify.
Do not silently replace first-use inference with unrestricted whole-function
backward inference. Proposed initial restriction: only direct immutable literals
are reusable, and defaulting occurs at an explicit solver boundary with provenance.

### Constraints and recursion

Use explicit constraints for equality, result-union inclusion, literal
admissibility, structural adaptation, callable compatibility, and required calls.
Candidate-local substitutions are speculative; rejecting a candidate must not
leave its inferred types in the shared snapshot.

Construct a dependency graph for signatures and inference variables. Process
acyclic dependencies directly and recursive components with a worklist. If a
recursive component cannot determine a result unambiguously, report the cycle and
an annotation location. An annotation on one function may suffice to resolve the
others. Do not infer an arbitrary union simply to terminate a cycle.

## 7. Call resolution and generic elaboration

```go
type CallSite struct {
    Origin OriginID
    Form CallForm               // Ordinary, UFCS, operator, constructor, indirect.
    Name string
    LexicalScope ScopeID
    Receiver ExprID             // Present only for receiver-based forms.
    Selectors []TypeTerm
    Arguments TupleLiteralPlan  // Shared literal syntax model: labels, spread, order.
    ExpectedResult TypeState
}

type CallPlan struct {
    Target Callee               // Direct, specialized, indirect, or dispatch plan.
    Arguments TupleConstructionPlan // Construct the target's formal parameter tuple.
    Result TypeID
    Effect EffectKind
    Evidence []RequirementEvidence
}
```

`TupleConstructionPlan` identifies its schema/constructor and maps source arguments
to logical members. Its member plans have alternatives for supplied value, instantiated literal,
default expression, rest packing, structural adaptation, and a cast-function call
when a suitable existing non-lossy cast is in scope. The compiler may insert that
call, but candidate discovery must not invent a cast function. Literal storage-type
constructors are tried in the configured priority order. This does not permit
arbitrary implicit numeric conversions.

Proposed resolution algorithm:

1. Identify the call form and candidate location groups. UFCS starts with the
   concrete receiver's home module, followed by the declaration's lexical scopes.
   Ordinary calls use lexical scopes. Explicit qualification restricts discovery.
2. Within one location, check arity, labels, selectors, literal/storage-type
   constraints, generic substitutions, contract requirements, and cast-function
   availability. A known concrete argument matches an exact parameter when its
   type is already equal. It may also use an existing suitable non-lossy cast
   function, but that cast is an inserted call rather than an exact match.
3. Apply exact/concrete, generic-specialization, then `any` priority. An immutable
   literal template is eligible for contextual storage-type instantiation in the
   exact tier. If one non-generic candidate remains, try its type constructor over
   the admissible storage types in priority order. If multiple candidates remain,
   default the literal template to `Int`, or `Float` when out of integer range, and
   then consider existing cast functions. Proceed to an outer location only when
   no candidate is applicable. Ambiguity in a preferred location is a diagnostic,
   not a reason to fall back.
4. Select argument adaptations, type constructors, cast-function calls, and defaults
   only after choosing a viable candidate. The compiler may insert a call to an
   existing suitable cast function, but must not synthesize the function itself.
   Keep source evaluation order separate from formal parameter packaging.
5. For finite union arguments, construct an exhaustive dispatch plan if necessary.
   Resolve each concrete combination in its appropriate environment, and form the
   result union. Avoid duplicating evaluation of arguments before dispatch.

Open ranking details: union-accepting parameters, multiple generic candidates,
structural adaptations, default/rest preferences, and mixed-specificity signatures
such as `f(Int, any)` versus `f(any, Int)`. Do not rank these by declaration order.
Ordinary-scope hiding by a same-name but inapplicable declaration also needs final
confirmation. Type-qualified members and callable-valued fields need an explicit
precedence relative to UFCS. Operator lookup is provisionally receiver-first,
consistent with the specification examples, but should be confirmed.

Multiple finite union arguments may yield a Cartesian product of cases. A compiler
limit must produce a useful diagnostic or retain a correct compact dispatch tree;
it must not silently change overload choice. `any` and `error` never become finite
just because the current executable happens to contain finitely many types. Their
dispatch needs an applicable fallback, or explicit source narrowing with `else`.

Return selectors disambiguate calls whose inputs are insufficient. Whether an
expected result without an explicit selector can select such an overload remains
open. Store expected-result information without making it a hidden ranking rule.

### Generic requirements and specialization

Represent `a + b` in a generic body as a required operation with type terms and
definition-site origin. At specialization, substitute concrete types and resolve
the operation using the generic's captured lexical environment and applicable
receiver-module priority. Never use the caller's local overloads as implicit
dictionary entries.

A specialization key is `(generic declaration, canonical type arguments)`;
compile-time value arguments can be added if the language adopts them. Insert an
in-progress entry before processing a recursive specialization. Diagnose unbounded
specialization growth rather than hanging the compiler. Concrete runtime code has
no generic parameters or contract dispatch tables.

Explicit conformance annotation versus structural satisfaction, overlapping
contract requirements, and contract composition remain open. Evidence objects
should retain the chosen declarations and field projections for diagnostics.

### Defaults, partial application, and closures

Default expressions belong to the declaration's environment and execute per full
call when omitted. Their precise order and effect restrictions remain open.
Ordinary invocation constructs the formal parameter tuple, supplying default
values for omitted members through the generated constructor's construction plan.
Named and positional arguments, spreading, structural adaptation, and rest packing
use the shared tuple construction rules; function-specific rest constraints are
validated on the formal schema.

`foo(1, *)` produces a callable even if the remaining parameter has a default;
partial application must not eagerly fill that parameter. Evaluation of supplied
arguments occurs once when constructing the partial application, provisionally;
remaining defaults are deferred to completion of the call. Their precise ordering
and interaction with subsequent partial applications still need specification.

Propose a `PartialApplicationPlan` containing the original callee, a generated
supplied-value tuple schema, a remaining-parameter schema, and a mapping from both
sets of logical slots into the original parameter tuple. The remaining schema
preserves unfilled members' names, types, and defaults in their original environment.
The resulting callable accepts that remaining tuple. Its generated body combines
captured values with newly supplied/defaulted values to construct the original
argument tuple exactly once before invocation.

For `foo = fn(a: Int, b: Int = 2) ...` (schematic default notation), partial
application of `a` produces a captured tuple `(a: Int)` and a remaining schema
`(b: Int = 2)`. The latter retains a default; it does not fill `b` at partial
application time. This notation illustrates schemas, not a proposed syntax change.
Repeated partial application composes the slot mappings. Equivalent generated
shapes may share types; schemas with distinct defaults/environments must not be
conflated. Neither tuple necessarily requires a heap allocation.

Closures explicitly list captured value identities. Early SSA may retain
`closure.make` and `call.indirect`; closure conversion later creates a generated
environment tuple plus code reference. A captured `fx` callable is data until
invoked. Partial applications are a structured instance of this closure model.

### Compile-time elaboration

`inline for`, type-dependent tuple expansion, `typeof`, and meta expressions need
an explicit boundary between compile-time facts and runtime computations. Expand
an inline tuple loop using logical member order after the tuple shape is known;
the member values themselves may still require runtime evaluation. Generic bodies
may defer this expansion until specialization.

Fold `typeof` when static information determines the answer; otherwise lower it
to the union tag and runtime metadata as appropriate. Do not treat an open union's
static descriptor as the concrete dynamic member descriptor by accident.

Meta file/environment inputs must be recorded as compilation dependencies and
their evaluated results frozen in the snapshot. Their evaluation policy and any
compile-time interpreter are separate decisions. Unsupported meta operations must
produce diagnostics before SSA lowering rather than become opaque runtime calls.

## 8. Typed SSA

Proposed convention: use block parameters instead of separate phi instructions.
Every incoming edge supplies arguments for the destination block's parameters.
This makes joins, loop-carried values, and ownership transfers visible on edges.

```go
type ValueID uint32             // Function-local.
type BlockID uint32

type SSAValue struct {
    ID ValueID
    Type TypeID
    Origin OriginID
}

type Instruction interface { instruction() }
// Const, Call, IndirectCall, TupleMake, FieldGet, TupleUpdate,
// ArrayMake, ArrayGet, ArrayUpdate, UnionInject, UnionTag, UnionProject,
// ClosureMake, PrimitiveOp, ... each with typed fields.

type Terminator interface { terminator() }
// Jump{Target, Args}, Branch{Condition, Then, Else},
// Switch{Tag, Cases, Default}, Return{Value, Effect}, Unreachable{Reason}.

type BasicBlock struct {
    ID BlockID
    Parameters []SSAValue
    Instructions []Instruction
    Terminator Terminator
    Origin OriginID
}

type Function struct {
    ID FunctionID
    Declaration DeclID
    HomeModule ModuleID
    Signature TypeID
    Entry BlockID
    Blocks []BasicBlock
}
```

Representative concrete instruction alternatives:

```go
type FieldGet struct {
    Result SSAValue
    Aggregate ValueID
    LogicalSlot uint32
}

type UnionInject struct {
    Result SSAValue             // A resolved union type.
    Member ValueID              // A value of one accepted member type.
}

type DirectCall struct {
    Result SSAValue
    Target FunctionID
    Arguments ValueID           // One constructed semantic parameter tuple.
    EffectIn, EffectOut ValueID // Both invalid for pure calls.
}
```

Effect tokens have an internal IR-only type distinct from source/runtime types.
A stronger Go API could split pure/effectful call alternatives to remove invalid
token combinations; either way, verification enforces the distinction.
Indirect calls use the same tuple convention. Function entry initially receives
the parameter tuple and projects logical members for bindings. Block parameters
remain compiler control-flow arguments, not source formal-parameter schemas.
Scalar replacement may split the tuple later while preserving its semantic type
and metadata. The textual SSA examples below abbreviate tuple construction and
projection for readability.

Prefer concrete instruction alternatives over a universal opcode record with many
optional fields. Each instruction exposes defined values and operands through
shared query helpers. Derived use/predecessor indexes belong to analyses, avoiding
two competing sources of truth in the published graph.

Values denote logical language values. They are not automatically registers,
heap objects, stack slots, or owned references. A tuple is initially one logical
value even if later split into fields. Field IDs/slots are logical; byte offsets
appear only during storage lowering.

### Example: rebinding through a block result

```tuppence
x = mut seed()
x = if condition() { replacement() } else { x }
consume(x)
```

Assuming pure calls and a common concrete result type `T`:

```text
entry:
  %x0 : T = call seed()
  %c : Bool = call condition()
  branch %c, yes(), no()
yes:
  %y : T = call replacement()
  jump join(%y)
no:
  jump join(%x0)
join(%x1 : T):
  %r = call consume(%x1)
  return %r
```

No inner assignment updates `%x0`. `%x1` is the new same-scope binding's value.
The branch edge explicitly conveys which value becomes the block result.

Loops similarly have header parameters for accumulators, backedges carrying the
next values, and exit parameters for results. `break`, `continue`, `return`, and
their `try` variants become explicit edges. Their cleanup obligations must follow
the actual exited scopes. Short-circuit operators and pattern guards remain control
flow; they are not eager calls evaluating both branches.

### Effects

Propose a compiler-only SSA effect token. Every effectful call consumes the current
token and produces its successor; branch joins merge tokens through block
parameters. The token is not a Tuppence value or an emitted C argument.

```text
%e1, %a = call.fx read(%e0, path)
%b = call.fn transform(%a)
%e2, %unit = call.fx write(%e1, %b)
return %unit, %e2
```

Foreign operations and effectful resource cleanup participate in this ordering.
Pure instructions have ordinary data dependencies. Declaring an effectful foreign
operation pure is a promise to the optimizer: hidden effects are not protected
against CSE or reordering.

Purity alone is not a complete rule for speculation: nontermination and eventual
panic semantics still need a decision. Proposed conservative initial optimizer:
do not hoist a possibly diverging/panicking call onto a path that did not execute
it. Ordinary error union values are results, not implicit exceptional control flow.

## 9. Pass operations and verification

Suggested operations (names are illustrative):

| Operation | Result |
| --- | --- |
| `Index(astModules, sources)` | Indexed semantic snapshot and collection diagnostics. |
| `Analyze(indexed, builtinRoles)` | Resolutions, constraints, inferred requirements, diagnostics. |
| `Check(analyzed, roots)` | Checked reachable program, or diagnostics preventing lowering. |
| `Specialize(checked, roots)` | Concrete reachable functions and types. |
| `LowerSSA(monomorphic)` | Typed blocks, values, calls, and effects. |
| `VerifySemantic` / `VerifySSA` | Structured invariant failures with origin information. |
| `RewriteFunction(fn, replacements)` | New immutable function plus ID/provenance mapping. |
| `OptimizeValues(ssa)` | New SSA program; invalidated analyses are recomputed. |
| `PlanStorage(ssa, target)` | Escape/ownership summaries and representation descriptors. |
| `LowerOwnershipAndStorage(plan)` | Explicit transfers, cleanup, layouts, and ABI operations. |
| `EmitC(storageProgram)` | C source and runtime type metadata. |

`RewriteFunction` must update operands, edge arguments, and origins consistently;
it cannot merely replace one table entry while leaving dangling references.
Effect rewrites and ownership rewrites need their own legality checks.

Builder operations should be small and explicit: `NewBlock(parameterTypes)`,
`Emit(block, instruction)`, `SetTerminator(block, terminator)`, and `Freeze()`.
Reserve IDs for forward edges without publishing half-built blocks. `Freeze`
rejects missing terminators, unresolved reservations, and invalid graph references.
Read-only operations include `Definition(value)`, `Operands(instruction)`,
`Successors(block)`, `TypeOf(value)`, and `OriginOf(entity)`. Dominance, liveness,
and use lists are computed analyses, cached only for the matching immutable body.

Verify SSA at construction and after transforming passes in tests/debug builds:

- Unique definitions, correct function ownership, dominance, and valid targets.
- Correct instruction types and call signatures; no pending types or candidates.
- Edge argument arity/types agree with block parameters; all blocks terminate.
- A union projection is justified by a matching discriminant path or static proof.
- Effects form valid per-path chains and cannot be silently dropped or duplicated.
- Pure functions do not invoke effectful operations.
- Source provenance exists for diagnostics, including synthesized operations.

Keep ownership verification separate and stronger after ownership lowering. A
well-typed SSA value is not by itself proof of a valid retain/release schedule.

Initial optimizations: constant folding with Tuppence arithmetic rules, dead pure
value removal under the chosen divergence policy, CFG simplification, and local
CSE. Add global value numbering and inlining later. CSE uses resolved callee
identity and type arguments, never just a source name. Once ownership is explicit,
rewrites must preserve transfers and cleanup ordering.

## 10. C boundary

Emit sequenced statements for effectful calls; do not rely on C argument evaluation
order to implement Tuppence ordering. Lower block arguments to edge assignments
with parallel-copy semantics, using temporaries when sources/destinations overlap.
An initial label-and-goto emitter is acceptable; structured C is an output-quality
improvement, not a semantic requirement.

Lower language arithmetic to defined C operations/helpers. Do not inherit C signed
overflow, invalid shifts, aliasing violations, or unaligned accesses as accidental
language behavior. Resolve Tuppence's arithmetic/failure rules before the relevant
primitive milestone.

Keep semantic function arguments separate from ABI packaging. The implementation
notes describe parameter tuples and implicit context parameters; neither forces
an actual allocation nor requires all context arguments in every function.
An ABI plan maps the semantic tuple to C parameters, a struct, or another selected
representation, and separately records implicit context inputs such as heap/runtime
state. These implementation inputs are not user-supplied tuple members and do not
participate in overload resolution or partial application's remaining parameters.
Compute context requirements transitively through calls, using a compatible
convention for indirect targets. Whether a context is passed at invocation or
captured in an environment needs an explicit lifetime rule. Compile-time effect
tokens are separate from these actual runtime context values.

C interop annotations, foreign ownership, callbacks, panic behavior, and platform
resource APIs need a separate ABI design before implementation. They are not
prerequisites for the initial IR package. A minimal C executable can use a small
compiler-owned runtime without establishing the public foreign interface.

## 11. Implementation sequence and acceptance criteria

1. **Semantic foundation.** Add IDs, spans, immutable tables/builders, scopes,
   declarations, binding versions, deterministic dumps, and diagnostics. Test
   same-scope rebinding, inner shadowing, right-hand-side lookup, and snapshot
   isolation. No backend or ARC implementation.
2. **A small checked subset.** Add primitive/nominal/structural types, literal
   constraints, direct pure/effectful calls, and explicit return constraints.
   Test literal instantiation, fixed mutable types, closure capture identities,
   purity violations, and recursive inference diagnostics.
3. **Resolution and elaboration.** Add overload location groups, generic
   requirements/specialization, UFCS, finite union dispatch, defaults, and partial
   application. Use curated fixtures for ambiguous and rejected calls as well as
   accepted ones. Resolve the relevant open ranking rules first.
4. **SSA subset.** Add block parameters, effects, conditionals, loops, union
   operations, and early exits. Verify every fixture and include deliberately
   invalid graphs to test dominance, type, effect, and edge checks.
5. **First C execution.** Select a small concrete runtime subset, fixed baseline
   layouts, and defined primitive semantics. Compile/run generated C and compare
   results/effect traces against expected behavior. Add closure and aggregate
   runtime support incrementally.
6. **Ownership and reuse.** Validate the protocol in the companion document using
   executable ownership traces before enabling reuse. Test aliases, nested heap
   members, stack containers, returns, loops, union boxing, and effectful cleanup.
7. **Optimization and layout expansion.** Add ARC cancellation, reuse, richer
   escape analysis, and measured layout alternatives. SoA remains optional.

Each milestone should add focused Go tests and readable semantic/SSA golden dumps.
Parser goldens change only when source syntax/parser behavior changes; IR goldens
belong in their own testdata directory. Update `TODO.md` only for implemented,
covered work. Run `make test` for implementation milestones; grammar regeneration
is unnecessary for a document-only proposal.

## 12. Review queue

Review these in small rounds rather than deciding the entire queue at once.

| Priority | Decision | Current recommendation/status |
| --- | --- | --- |
| First | Two semantic/SSA representations; block parameters; immutable snapshots | Proposed architecture for approval. |
| First | Literal defaulting and first-use boundaries | Agreed categories, detailed solver policy open. |
| First | Block visibility and module initialization | Proposed visibility; module execution order/cycles open. |
| Next | Complete overload ranking and contract conformance | Preserve unresolved cases; no declaration-order tiebreak. |
| Next | Expected-result overload selection and implicit conversions | Open; avoid automatic concrete conversions initially. |
| Next | Tuple schemas and callable compatibility | Preserve names/defaults; define compatibility and default selection through function references. |
| Next | Effect ordering, defaults, divergence, and panic | Effectful calls ordered; speculation details open. |
| Runtime | Strong counts, borrowed/consumed calls, cleanup timing, cycles | Proposed baseline. Cleanup policy and cycle handling remain open. |
| Runtime | Literal construction roles and C ABI annotations | Design required before corresponding backend support. |
| Later | Tag width, layout heuristics, split tags, SoA | Keep abstract; measure after a baseline implementation. |

Documentation reconciliation after acceptance: add the exact `mut` scope/type
rules and literal categories to the specification; formalize open unions and
overload ranking; clarify purity in examples that call printing from `fn`; keep
logical/physical order and context-parameter notes aligned with chosen layouts.
The existing AST README mentions traversal APIs not present in the current node
interface, so implementation should follow actual AST definitions rather than
depending on that README's illustrative API.
