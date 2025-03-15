package ast

// Identifier node types
const (
	NodeIdentifier         NodeType = "Identifier"
	NodeTypeIdentifier     NodeType = "TypeIdentifier"
	NodeFunctionIdentifier NodeType = "FunctionIdentifier"
)

// Identifier represents a regular identifier in the code (variable name, parameter name, etc.)
type Identifier struct {
	BaseNode
	Name string // The identifier name
}

// NewIdentifier creates a new Identifier node
func NewIdentifier(name string) *Identifier {
	return &Identifier{
		BaseNode: BaseNode{NodeType: NodeIdentifier},
		Name:     name,
	}
}

// String returns a textual representation of the identifier
func (i *Identifier) String() string {
	return i.Name
}

// TypeIdentifier represents a type identifier (starts with uppercase)
type TypeIdentifier struct {
	BaseNode
	Name      string // The type name
	Namespace string // Optional namespace (e.g., module name)
}

// NewTypeIdentifier creates a new TypeIdentifier node
func NewTypeIdentifier(name string, namespace string) *TypeIdentifier {
	return &TypeIdentifier{
		BaseNode:  BaseNode{NodeType: NodeTypeIdentifier},
		Name:      name,
		Namespace: namespace,
	}
}

// String returns a textual representation of the type identifier
func (t *TypeIdentifier) String() string {
	if t.Namespace != "" {
		return t.Namespace + "." + t.Name
	}
	return t.Name
}

// QualifiedName returns the fully qualified name including namespace
func (t *TypeIdentifier) QualifiedName() string {
	if t.Namespace != "" {
		return t.Namespace + "." + t.Name
	}
	return t.Name
}

// FunctionIdentifier represents a function identifier (starts with lowercase, may end with ? or !)
type FunctionIdentifier struct {
	BaseNode
	Name               string // The function name
	Namespace          string // Optional namespace (e.g., module name or type name)
	IsPredicateOrPanic bool   // True if the function name ends with ? or !
}

// NewFunctionIdentifier creates a new FunctionIdentifier node
func NewFunctionIdentifier(name string, namespace string) *FunctionIdentifier {
	isPredOrPanic := false
	// Check if the function name ends with ? or !
	if len(name) > 0 {
		lastChar := name[len(name)-1]
		isPredOrPanic = lastChar == '?' || lastChar == '!'
	}

	return &FunctionIdentifier{
		BaseNode:           BaseNode{NodeType: NodeFunctionIdentifier},
		Name:               name,
		Namespace:          namespace,
		IsPredicateOrPanic: isPredOrPanic,
	}
}

// String returns a textual representation of the function identifier
func (f *FunctionIdentifier) String() string {
	if f.Namespace != "" {
		return f.Namespace + "." + f.Name
	}
	return f.Name
}

// QualifiedName returns the fully qualified name including namespace
func (f *FunctionIdentifier) QualifiedName() string {
	if f.Namespace != "" {
		return f.Namespace + "." + f.Name
	}
	return f.Name
}

// RenameIdentifier represents an identifier with an optional new name for import renaming
type RenameIdentifier struct {
	BaseNode
	Original *Identifier // Original identifier
	Renamed  *Identifier // New identifier (may be nil if not renamed)
}

// NewRenameIdentifier creates a new RenameIdentifier node
func NewRenameIdentifier(original *Identifier, renamed *Identifier) *RenameIdentifier {
	return &RenameIdentifier{
		BaseNode: BaseNode{NodeType: "RenameIdentifier"},
		Original: original,
		Renamed:  renamed,
	}
}

// String returns a textual representation of the rename identifier
func (r *RenameIdentifier) String() string {
	if r.Renamed != nil {
		return r.Original.String() + ": " + r.Renamed.String()
	}
	return r.Original.String()
}

// Children returns the child nodes
func (r *RenameIdentifier) Children() []Node {
	if r.Renamed != nil {
		return []Node{r.Original, r.Renamed}
	}
	return []Node{r.Original}
}

// RenameType represents a type identifier with an optional new name for import renaming
type RenameType struct {
	BaseNode
	Original *TypeIdentifier // Original type identifier
	Renamed  *TypeIdentifier // New type identifier (may be nil if not renamed)
}

// NewRenameType creates a new RenameType node
func NewRenameType(original *TypeIdentifier, renamed *TypeIdentifier) *RenameType {
	return &RenameType{
		BaseNode: BaseNode{NodeType: "RenameType"},
		Original: original,
		Renamed:  renamed,
	}
}

// String returns a textual representation of the rename type
func (r *RenameType) String() string {
	if r.Renamed != nil {
		return r.Original.String() + ": " + r.Renamed.String()
	}
	return r.Original.String()
}

// Children returns the child nodes
func (r *RenameType) Children() []Node {
	if r.Renamed != nil {
		return []Node{r.Original, r.Renamed}
	}
	return []Node{r.Original}
}
