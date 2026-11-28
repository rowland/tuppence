package ir

// LiteralKind identifies the source-level numeric literal category.
type LiteralKind uint8

const (
	LiteralInteger LiteralKind = iota + 1
	LiteralFloat
)

// Valid reports whether the literal category is recognized.
func (k LiteralKind) Valid() bool {
	return k >= LiteralInteger && k <= LiteralFloat
}

// LiteralStoragePriority returns storage kinds in the order used when resolving
// an otherwise untyped numeric literal or immutable literal template.
func LiteralStoragePriority(kind LiteralKind) []StorageKind {
	if !kind.Valid() {
		return nil
	}
	if kind == LiteralFloat {
		return []StorageKind{StorageF16, StorageF32, StorageF64, StorageBytes}
	}
	return []StorageKind{
		StorageI8, StorageI16, StorageI32, StorageI64,
		StorageU8, StorageU16, StorageU32, StorageU64,
		StorageF16, StorageF32, StorageF64, StorageBytes,
	}
}

// SelectLiteralStorage chooses the first admissible storage kind in the
// language-defined priority order. The caller filters candidates by range and
// representation before calling this function.
func SelectLiteralStorage(kind LiteralKind, admissible []StorageKind) (StorageKind, bool) {
	if !kind.Valid() {
		return 0, false
	}
	allowed := make(map[StorageKind]struct{}, len(admissible))
	for _, storage := range admissible {
		allowed[storage] = struct{}{}
	}
	for _, storage := range LiteralStoragePriority(kind) {
		if _, ok := allowed[storage]; ok {
			return storage, true
		}
	}
	return 0, false
}

// LiteralDefaultKind identifies the semantic default selected when a literal
// remains ambiguous after candidate discovery.
type LiteralDefaultKind uint8

const (
	LiteralDefaultInt LiteralDefaultKind = iota + 1
	LiteralDefaultFloat
)

// DefaultLiteralKind returns the configured default for an ambiguous literal.
func DefaultLiteralKind(kind LiteralKind, outOfIntegerRange bool) (LiteralDefaultKind, bool) {
	if !kind.Valid() {
		return 0, false
	}
	if kind == LiteralFloat || outOfIntegerRange {
		return LiteralDefaultFloat, true
	}
	return LiteralDefaultInt, true
}
