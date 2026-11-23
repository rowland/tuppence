# Tuppence IR: ownership and storage proposal

Design review draft — 2026-11-22. Companion to [IR_DESIGN.md](IR_DESIGN.md).
The language's value semantics and allocation flexibility are agreed; the proposed
v1 lowering uses ordinary strong counts and explicit transfers at consuming calls.
This is a compiler design, pending implementation.

## 1. Separate values, storage, and ownership

A semantic value is immutable. Storage may be reused when doing so cannot change
any observable alias. A value's type does not determine whether it lives on the
stack, heap, in registers, or in split fields.

An SSA value is a compiler identity, not a runtime reference-count contribution.
Several SSA names may denote one moved ownership obligation, or several actual
references to the same heap allocation. Conversely, an unused SSA name may require
no storage. Reference counts must follow ownership/storage analysis, not the
number of SSA uses or variable names.

For a stack tuple containing heap references, copying fields into a new owning
container adds references to those heap allocations. Sharing a heap tuple adds a
reference to the tuple allocation, not recursively to all of its members. Destroying
the tuple releases its owned member references, potentially cascading destruction.

Propose an ownership vocabulary independent of the eventual counter encoding:

| Operation | Meaning |
| --- | --- |
| Copy/share | Create an additional owning reference; retain as needed. |
| Borrow | Temporarily access a value while an owner protects the old value throughout the access. |
| Move | Transfer an ownership obligation without duplicating it. |
| Drop | End an ownership obligation, eventually triggering required destruction. |
| Reuse | Transfer suitable storage from a consumed value to a replacement value. |

Borrow here is an internal compiler concept, not a proposed source-language
borrow checker. Reuse needs alias, lifetime, representation, and effect proofs;
SSA alone supplies none of those proofs automatically.

## 2. Proposed analysis summaries

```go
type ParameterMode uint8 // Borrowed or consumed in this compiled calling convention.
type Placement uint8     // Register, stack, heap, split aggregate.
type LayoutID uint32
type AllocationID uint32
type CleanupID uint32

type FunctionOwnership struct {
    Parameters []ParameterMode
    Result ResultOwnership
    Escapes EscapeSummary
    CleanupEffects CleanupEffectSummary
}

type StoragePlan struct {
    // Private maps keyed by function/value/allocation identity.
    // Maps logical values to layouts and ownership obligations.
}
```

Start conservatively: borrowed parameters remain alive during a call; returning a
borrowed argument requires creating/transferring ownership explicitly. A later
consuming convention can expose last-use reuse across calls. Summaries for indirect
calls must hold for every possible target or use a conservative common convention.

Closures own their captured references when those captures outlive the caller.
Semantic parameter tuples and partial applications' supplied-value tuples follow
the same ownership rules as other tuples. Ownership summaries may describe modes
per logical member without replacing the tuple-based callable signature with a
flat type list. Completing a partial application must preserve captured values
for any later invocation unless analysis proves that callable is consumed; it
cannot move captures out unconditionally on the first call.
Escape analysis may keep an environment on the stack when safe. Do not create a
heap environment pointing into a dead stack frame. Whole-program compilation
helps analysis but does not prove every indirect call or recursive lifetime safe.

## 3. Recommended v1 protocol

Use ordinary strong counts for managed heap allocations. Construction delivers one
owned result at count one. Sharing adds one, moving preserves the count, and
dropping subtracts one. The drop that reaches zero runs the layout's destructor
and reclaims the allocation. Returning or forwarding a moved value transfers its
cleanup obligation too. The sender has no scope-exit cleanup for that obligation.

A count of one answers “how many owning edges reach this allocation?” It cannot
answer “does the caller still need the old value?” The call convention answers
that question independently. Modes are a contract of a compiled entry point.
Source callers pass ordinary immutable values.

| Parameter mode | Caller obligation | Callee permission |
| --- | --- | --- |
| Borrowed | Keep a protecting owner and the old contents valid through the call. | Read and reborrow. Share into a new owner before escaping or returning it. |
| Consumed | Deliver one independent owner. Move a dead input or share a still-live input. | Forward, return, drop, or attempt reuse with that owner. |

Start with borrowed inputs and owned results. Use consumed inputs for selected
update operations and, later, functions that forward to them. This choice is
made when compiling the entry point. A caller's last use alone cannot upgrade a
borrowed parameter. Apply the same contracts to heap-owning fields passed inline.
An indirect call uses the convention of its callable type or an adapter.

### Who selects the parameter mode?

Library authors write ordinary Tuppence. Parameter modes belong to the compiled
entry point and require no source annotation. For the first implementation,
give ordinary functions borrowed inputs and owned results. Give each internal
value operation a compiler-defined contract. For example, an array append can
consume the array and borrow the element, acquiring an owner for the stored
element when its representation requires one. Identify the operation through
resolved compiler identity. A user-defined overload called `append` has its
own body and follows the ordinary rules.

The existing `to_array` in `lib/list.tup` illustrates the boundary. It traverses
its input list and builds an array using `acc.append(current.head)`. The list
parameter can stay borrowed throughout. The loop owns its local accumulator, so
its final use on the backedge can transfer that owner to array append and carry
the returned owner into the next iteration. The algorithm is entirely Tuppence.
Only the backend needs to know how to grow or reuse the backing allocation.

Later, infer a consuming convention for a simple wrapper whose final action
forwards its parameter to a known consuming slot. Begin with direct calls and
a single unconditional forward. Keep the borrowed default for ambiguous cases.
A missed inference costs a retain and can force a copy. Changing a mode requires
lowering both the body and its callers against the same contract and preserving
cleanup effects on every exit.

Expose only operations with immutable value semantics to Tuppence code,
including the standard library. Allocation, reference-count inspection, and
writes to reused storage belong to backend lowering or its private helpers.
Even an ordinary `fn` can append efficiently under this rule. Its result and
all surviving aliases must match the behavior of constructing a fresh array.

### What crosses the C17 call boundary

Record each parameter's mode in the compiler's table of compiled entry points.
Call lowering reads that table when emitting a call. Body lowering uses the same
table to decide which inputs the callee must discharge on each exit. Both modes
can pass the same ordinary C pointer. The generated call carries no mode flag
and the allocation header needs no borrow count. A move changes which code is
responsible for cleanup and can emit zero instructions.

For example, suppose `Buf` is a writable, directly owned buffer of integers,
and the index has already been checked. These are two compiled entry points:

```c
int buf_at(const Buf *b, size_t i) {       // b is borrowed
    return b->data[i];
}

Buf *buf_set(Buf *b, size_t i, int n) {   // b is consumed, result is owned
    if (b->refs != 1) {
        Buf *copy = buf_clone(b);         // fresh backing buffer, refs = 1
        buf_drop(b);                      // discharge the incoming owner
        b = copy;
    }
    b->data[i] = n;
    return b;                             // transfer the result owner
}
```

The C `const` qualifier documents read access in this example. The compiler's
mode table supplies the ownership contract even when the C types are identical.
For an owned input `x`, a final `return buf_set(x, i, n);` transfers it directly.
If the caller later reads `x`, it emits `buf_retain(x);` before the same call
and keeps responsibility for dropping its original owner. A borrowed input
likewise needs a retain to supply a consuming slot, even on its final local use.
The callee body is identical for all three calls.

### Lowering a call

Evaluate arguments in language order, keeping earlier results alive across later
argument evaluation. Account for the whole call, including the callable's
environment, before transferring ownership. Each consumed occurrence needs an
owner. Borrowed occurrences may use a common anchor that survives the call.
A move is legal only if its source obligation is available on this path and no
outstanding borrow depends on it. Otherwise share for the consumed occurrence.
This rule works when different expressions happen to yield the same allocation.

For example, start with a managed array in `x`, with one owner and no borrows.
Here `update` consumes its array and returns an owned result.

| Call and subsequent use | Counts and consequences |
| --- | --- |
| `update(x)`, then read `x` | Share for the call. At count two, update builds the changed value in fresh storage and drops its input. The caller reads the original. |
| Return `update(x)` | Move into the call at count one. Update can recycle the buffer and return its owner. Neither caller nor callee drops the forwarded owner. |
| `f(x, x)`, both borrowed | Keep the caller's owner through the call. Both parameters read the original at count one. Neither has permission to recycle it. |
| `f(x, x)`, both consumed, last caller use | Share once and move once. At entry there are two independent owners. Updating the first while the second is live takes the copying path. If the first is dropped after its final read, the second may then recycle at count one. |
| `f(x, x)`, borrowed and consumed, last caller use | Keep the original owner as the borrow anchor through return and share once for the consumed slot. The consumed slot sees at least two owners while the borrow is active. Parameter order makes no difference. |

For a nested example, let `f(a, b)` consume both inputs, compute `r = update(a)`,
then read `b` and return `r`. In `f(x, x)`, moving `a` into the nested call
leaves its count at two. Dropping that input after copying leaves `b` at one
with its original contents. If `f` instead finishes reading and drops `b`
before forwarding `a`, the nested call sees one and may reuse. Neither function
needs to recognize the duplicate pointers. An extra owner elsewhere simply keeps
the count higher.

### Reuse at an update

A consuming update must have a writable managed allocation of compatible layout,
size and alignment, and must finish any reads or interior borrows that need the
old bytes before overwriting them. Check for one strong owner of the **backing
allocation** or prove exclusivity statically. Other live views, including slices
and captured values, must have their own counted protection or a borrow tied to
an owner that cannot move, drop, or update until the borrow ends. A borrowed
view of a parent also protects the reachable fields it reads. Static literal
storage and untracked foreign storage are ineligible for this runtime fast path.
The baseline assumes single-threaded execution and no untracked observers.

On success, transfer the input's storage obligation to the replacement. On
failure, construct the replacement independently and discharge the consumed
input. Both paths return an owned value. A register or stack local can likewise
be overwritten after its last observable use without consulting a counter.

### Why defer two-part release?

The earlier proposal decremented before final use (`release1`), allowed reuse
at zero, then conditionally destroyed storage later (`release2`). Its motivation
is sound: a caller on its final use should be able to give storage to its callee.
Moving an owner at count one conveys exactly that permission without changing
the count. A borrow from a still-interested caller conveys no such permission.

Under early decrement, two arguments can both be in their final use while either
can still read the original. Delaying the overlapping decrements through the
last read is safe if the corresponding owners stay counted and finalization
happens exactly once. Allowing a live zero instead requires tracking which use
protects it, which other uses overlap, and who cancels pending cleanup when a
nested call returns, recycles, or frees it. Keep this as an alternative only if
a measured workload justifies the extra protocol.

### Eliminating ARC traffic

Prefer moves during lowering. A proven last-use handoff across several consuming
calls can require no counter traffic. A temporary share followed by a drop can
also become a borrow if an independent anchor remains valid for the entire
interval. Eliminate a retain/drop pair only when lifetime, old-value observations,
and required destruction effects are preserved on every path. In particular a
uniqueness test inside the interval can change from shared to unique when a
retain disappears. Treat such tests and unknown calls as barriers unless the
optimizer separately proves the newly enabled reuse safe.

### Minimal array operations for append and sort

Assume the backend already supports array construction, length and indexed reads.
Two additional private lowering operations suffice for these algorithms:

| Operation | Contract |
| --- | --- |
| Prepare writable storage | Consume an array and a checked required capacity. Deliver exclusive writable storage containing the same logical elements and length. Reuse eligible backing storage when it fits, grow or relocate a unique buffer when necessary, or copy the logical elements of a shared array or view. Acquire ownership of copied elements according to their layout. |
| Write one slot | Transfer a prepared element into the destination. Replacing an initialized slot discharges its previous owner. Initializing the first unused slot has no previous owner and advances length only after initialization. |

Preparation can serve one checked update or a sequence of updates while the
compiler can establish exclusive access. Writable storage stays private to lowering. Source-level replacement requires
an index within the current length. Append checks length arithmetic, requests
space for one more element, and initializes the tail. Use geometric growth for
repeated single-element appends. Stabilize the incoming element before preparing
storage if it refers into the old buffer, since growth can relocate that buffer.

Swap can be ordinary Tuppence: validate both indices, read both old elements,
then thread the array through two replacements. Equal valid indices can return
the input directly. Lowering can later fuse this sequence into one preparation
and two moves of complete slots (including union tags), with no element ARC
traffic. Sort can express its comparisons and swaps entirely in Tuppence.
Even without fusion a second replacement can reuse the result of the first.
For efficient sorting, the loop must forward its owned accumulator into each
swap. A borrowed wrapper around each swap would otherwise force repeated copies.
Inlining it or inferring a consuming mode for its array parameter enables reuse.
End comparator borrows before a write and reestablish exclusivity if any new
alias can observe the current array.

### One preparation for a whole quicksort

Target at most one shallow clone of the array backing storage for a complete
sequential quicksort, including all recursive partitions. Sort consumes an
array and returns an owned array. After the trivial-length check, prepare it
once with required capacity equal to its current length. This takes the
existing writable buffer if exclusive or detaches a private working copy.
The count of a shared input can stay above one throughout the sort.

Use a private compiled worker for the ordinary Tuppence range algorithm.
Its entry contract requires an already-exclusive buffer. The worker passes
that buffer through partition and both recursive calls in sequence, along with
index bounds. Each call returns the working owner to its caller. The caller
needs only bounds and the partition index across these transfers. Passing two
owning slices into the recursion would defeat this argument. Copy the pivot
value (sharing its heap members as needed) or track its index correctly across
swaps. An interior pointer into a slot cannot stand for a saved pivot value.

During storage lowering, verify the worker body and its specialized callees:
each control-flow edge carries one working owner, temporary reads finish before
overlapping writes, and no operation publishes a surviving alias of that buffer.
A comparator may borrow elements for a comparison. A closure that captures the
original input keeps that original alive, which entry preparation accounts for.
The compiler can then emit slot writes throughout the worker without repeating
the uniqueness test. This stronger internal contract can be established for a
known library worker before implementing general inference. Arbitrary user
functions can always use ordinary checked updates. If an algorithm deliberately
saves an intermediate array snapshot, its next write must preserve that snapshot
and can require another clone.

This makes flat arrays the intended common representation for short-lived
sequences with one evolving version. Persistent collections serve workloads
that keep many large versions alive and benefit from sharing unchanged paths.
The expected coverage is a workload hypothesis to measure.

### Tuple update

Lower `foo.(bar: baz)` as a value-producing `TupleUpdate`. Evaluate the base
and replacements once in language order. Secure the replacement values before
any write that could invalidate their sources. For example,
`p.(left: p.right, right: p.left)` reads both fields of the original tuple.

- For a register or split tuple, construct the result from the replacement and
  the unchanged fields. A dead stack slot may serve as the result destination.
- For a heap tuple, reuse its allocation only with a transferable owner,
  exclusive writable storage and a compatible result layout. Install the
  prepared replacements and discharge displaced field owners exactly once.
- Otherwise construct a fresh outer tuple. Copy unchanged fields with their
  layout's ownership rules and initialize the changed fields directly from the
  replacements. This implements a shallow clone followed by update while
  avoiding a retain/drop pair for each overwritten field. Keep the base valid
  until all reads finish, then discharge any consumed base owner.

A heap-backed child stays shared when its field is copied. Updating the child
requires its own uniqueness decision, even when the outer tuple is unique.

### Control flow

Last use is path-sensitive. Place ownership transfers on branch edges and loop
backedges as necessary. An unexecuted branch must neither retain nor drop its
values. Scope exits include `return`, `break`, `continue`, and error propagation;
each must carry the correct cleanup obligations. A later panic/unwind design will
need its own policy rather than inheriting one accidentally from generated C.

Loop iterations must not accumulate deferred cleanup until function exit unless
that is an intentional lifetime rule. Escaping loop results transfer obligations
out of the iteration; temporary values are finalized on the relevant iteration
exit.

## 4. Destruction effects and limits

Resource wrappers may close a file or perform other foreign cleanup on final
release. The provisional language rule is that `fn` must not introduce effectful
cleanup. This needs more than checking explicit calls: dropping the final copy of
a passed resource wrapper or replacing such a field can also cause cleanup.

Propose a type/ownership cleanup summary with a conservative "may have effects"
case for open unions and unknown foreign resources. A borrowed parameter can be
inspected by pure code without destroying its caller-owned resource. Operations
that would own and destroy effectful resources need a defined restriction or
effectful lowering; do not silently attach an effect to a pure function.

Open: exact cleanup timing, order among multiple resources, early destruction,
and whether ownership is allowed to delay cleanup across effects. Until settled,
preserve scope-end cleanup sites and their recorded order conservatively. ARC
cancellation must not move resource closure across observable operations.

Reference-count cycles, weak references, foreign references, and threading are
also open. Immutable construction restricts some cycles but does not justify
assuming recursive closures or foreign objects can never form them. v1 should
explicitly restrict unsupported cycles or define a collection mechanism. Atomic
counts alone would not make opportunistic mutation thread-safe.

Validate optimized traces against immutable copying with balanced ownership and
the same resource-cleanup policy. Compare observed values and cleanup events.
The two executions may legitimately allocate different numbers of buffers.

## 5. Type table and layout descriptors

```go
type RuntimeTypeID uint32        // Assigned after reachable concrete types are known.

type PhysicalField struct {
    LogicalSlot uint32
    Offset uint64
    Layout LayoutID
}

type Layout struct {
    SemanticType TypeID
    Size, Alignment uint64
    Fields []PhysicalField      // Physical order, with logical correspondence.
    Ownership OwnershipShape
}

type RuntimeTypeEntry struct {
    ID RuntimeTypeID
    SemanticType TypeID
    CanonicalLayout LayoutID
    Metadata TypeMetadata
}
```

Type metadata includes logical fields, names, annotations, and concrete type
arguments as needed. Layout-specific drop/copy helpers can be generated functions;
these are compiler runtime operations, not source-language method tables.

An important constraint: a TypeID plus an inline/heap bit does not distinguish
arbitrarily many byte layouts for the same type. Proposed baseline: when values
cross a boundary known only through their tag, use one canonical payload layout
per concrete type. Specialized layouts are allowed where the enclosing array,
function ABI, or static descriptor identifies them, with explicit conversion at
the boundary. An alternative is an allocation header carrying a layout descriptor;
that is a further design choice, not information already supplied by the type ID.

Assign runtime IDs deterministically after reachability/specialization. Reserve
only internal primitive IDs. Standard types receive ordinary table entries.
The provisional encoding `(TypeID << 1) | heapFlag` needs an explicit width,
overflow limit, and masking convention before C implementation.

## 6. Logical and physical member order

Preserve declaration order for positional construction, destructuring, ordinal
access, and reflection. Physical placement sorts by descending chosen storage size;
ties use member-name order for named tuples and logical slot for ordinal tuples.
Choose a deterministic name comparison rule (proposed: bytewise UTF-8 ordering).

Honor target alignment and add padding as needed. Size order alone does not prove
globally minimal padding on every target. A reference field is pointer-sized even
if the referenced value is large. The existing `@cstruct` exception follows the
target's C layout rather than this sorting rule.

Recursive layouts require indirection on every infinitely expanding cycle. This
does not force the entire enclosing value onto the heap. Compute layout dependency
components and choose representation breaks before calculating finite sizes.

## 7. Union and array representation

The tag logically belongs to a union slot, not exclusively to the referenced heap
object. A local may carry tag/payload in separate temporaries; a tuple field may
store them together. A heap alternative places the pointer in the payload, with
the representation bit telling access/cleanup how to interpret it. Inline payloads
may themselves contain owned heap references.

For mixed unions, reserve space for selected inline alternatives and a pointer,
boxing alternatives that do not fit the chosen policy. Union injection/projection
remains explicit so storage lowering can copy, box, unbox, or transfer as needed.
Borrowed interior views must not outlive the containing union storage.

Array layout alternatives:

| Layout | Storage | Status |
| --- | --- | --- |
| Interleaved union | Repeated tag/payload slots | Proposed simple baseline. |
| Separate tags | Tag region plus fixed-stride payload region | Deferred optimization; can share one backing allocation. |
| Struct of arrays | Separate columns for logical fields | Stretch goal. |

Separate tags can reduce per-element alignment padding and support tag-only scans;
they do not remove unused space inside fixed-stride payload slots. Measure complete
element access as well as scans before selecting defaults. A more elaborate
per-alternative payload arrangement is a separate future option.

Array descriptors, not semantic element types, identify the chosen layout.
Logical `array.get`/`field.get` can lower to column access without assembling an
entire element. When a complete value escapes, materialize the required canonical
representation with correct ownership transfers/copies. No stable element address
is promised. Slices must preserve the backing layout and lifetime information.

## 8. Ownership acceptance tests

Before enabling reuse, require tests covering:

- Stack aggregates owning heap members; heap aggregates sharing heap members.
- Sharing versus moving; final use through ordinary and indirect calls.
- Duplicate arguments, nested aliases, and overlapping borrows.
- Branch-specific last uses and loop-carried ownership.
- Closure captures surviving their creator, including partial applications.
- Returned/stored/reused allocations with invalidated old cleanup obligations.
- Replacing a member with a value derived from another member of the same object.
- Array append with spare capacity, growth, and an element sourced from the old
  buffer. Swap equal and distinct indices, including heap-owning elements.
- Instrument backing-buffer clones across a complete multi-partition quicksort.
  Require zero for a transferred, unique writable buffer and at most one with
  a live alias. For shared, unsorted input verify exactly one clone, sorted
  output and original order through the alias. Count across all worker calls.
- Tuple updates on shared and unique parents, including a two-field exchange
  and an independently shared child.
- Inline/boxed union alternatives and conversions between layouts.
- Early exits and deterministic resource-cleanup traces under the selected policy.

Check that every allocation is valid at access, every ownership obligation is
discharged exactly once, counts never underflow, and reuse preserves all remaining
aliases. Value equality alone is insufficient evidence: tests must also expose
leaks, stale cleanup, and premature destruction. Use generated C sanitizer runs
where supported once the backend exists.
