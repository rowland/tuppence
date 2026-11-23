package ir

// IDs are meaningful only within the snapshot or function that owns them.
// Zero is reserved as the invalid ID for every IR identity type.
type ModuleID uint32
type ScopeID uint32
type DeclID uint32
type BindingID uint32
type ExprID uint32
type TypeID uint32
type TypeVarID uint32
type FunctionID uint32
type FileID uint32
type OriginID uint32

type ValueID uint32
type BlockID uint32
type TupleSchemaID uint32
type DefaultID uint32
type RuntimeTypeID uint32
type LayoutID uint32
type AllocationID uint32
type CleanupID uint32

func (id ModuleID) Valid() bool      { return id != 0 }
func (id ScopeID) Valid() bool       { return id != 0 }
func (id DeclID) Valid() bool        { return id != 0 }
func (id BindingID) Valid() bool     { return id != 0 }
func (id ExprID) Valid() bool        { return id != 0 }
func (id TypeID) Valid() bool        { return id != 0 }
func (id TypeVarID) Valid() bool     { return id != 0 }
func (id FunctionID) Valid() bool    { return id != 0 }
func (id FileID) Valid() bool        { return id != 0 }
func (id OriginID) Valid() bool      { return id != 0 }
func (id ValueID) Valid() bool       { return id != 0 }
func (id BlockID) Valid() bool       { return id != 0 }
func (id TupleSchemaID) Valid() bool { return id != 0 }
func (id DefaultID) Valid() bool     { return id != 0 }
func (id RuntimeTypeID) Valid() bool { return id != 0 }
func (id LayoutID) Valid() bool      { return id != 0 }
func (id AllocationID) Valid() bool  { return id != 0 }
func (id CleanupID) Valid() bool     { return id != 0 }
