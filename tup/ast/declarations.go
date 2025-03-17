package ast

// NamespaceDeclaration represents a namespace declaration
type NamespaceDeclaration struct {
	BaseNode
	Name     string // The name of the namespace
	Elements []Node // The elements in the namespace
}

// NewNamespaceDeclaration creates a new NamespaceDeclaration node
func NewNamespaceDeclaration(name string, elements []Node) *NamespaceDeclaration {
	return &NamespaceDeclaration{
		BaseNode: BaseNode{NodeType: NodeNamespaceDeclaration},
		Name:     name,
		Elements: elements,
	}
}

// String returns a textual representation of the namespace declaration
func (n *NamespaceDeclaration) String() string {
	result := "namespace " + n.Name + " {\n"
	for _, elem := range n.Elements {
		result += "  " + elem.String() + "\n"
	}
	result += "}"
	return result
}

// Children returns the child nodes
func (n *NamespaceDeclaration) Children() []Node {
	return n.Elements
}

// Annotation represents an annotation applied to a declaration
type Annotation struct {
	BaseNode
	Name  string // The name of the annotation
	Value Node   // Optional argument value
}

// NewAnnotation creates a new Annotation node
func NewAnnotation(name string, value Node) *Annotation {
	return &Annotation{
		BaseNode: BaseNode{NodeType: NodeAnnotation},
		Name:     name,
		Value:    value,
	}
}

// String returns a textual representation of the annotation
func (a *Annotation) String() string {
	if a.Value != nil {
		return "@" + a.Name + "(" + a.Value.String() + ")"
	}
	return "@" + a.Name
}

// Children returns the child nodes
func (a *Annotation) Children() []Node {
	if a.Value != nil {
		return []Node{a.Value}
	}
	return nil
}

// GenericTypeParam represents a type parameter for generic types and functions
type GenericTypeParam struct {
	BaseNode
	Name       string // The name of the type parameter
	Constraint Node   // Optional type constraint
}

// NewGenericTypeParam creates a new GenericTypeParam node
func NewGenericTypeParam(name string, constraint Node) *GenericTypeParam {
	return &GenericTypeParam{
		BaseNode:   BaseNode{NodeType: NodeGenericTypeParam},
		Name:       name,
		Constraint: constraint,
	}
}

// String returns a textual representation of the type parameter
func (t *GenericTypeParam) String() string {
	if t.Constraint != nil {
		return t.Name + ": " + t.Constraint.String()
	}
	return t.Name
}

// Children returns the child nodes
func (t *GenericTypeParam) Children() []Node {
	if t.Constraint != nil {
		return []Node{t.Constraint}
	}
	return nil
}

// TypeDeclaration represents a type declaration
type TypeDeclaration struct {
	BaseNode
	Name           *TypeIdentifier     // The name of the type
	TypeParameters []*GenericTypeParam // Type parameters if this is a generic type
	Type           Node                // The type definition
	Annotations    []*Annotation       // Annotations applied to the type
	Docs           string              // Documentation comments
}

// NewTypeDeclaration creates a new TypeDeclaration node
func NewTypeDeclaration(name *TypeIdentifier, typeParams []*GenericTypeParam, typeNode Node, annotations []*Annotation, docs string) *TypeDeclaration {
	return &TypeDeclaration{
		BaseNode:       BaseNode{NodeType: NodeTypeDeclaration},
		Name:           name,
		TypeParameters: typeParams,
		Type:           typeNode,
		Annotations:    annotations,
		Docs:           docs,
	}
}

// String returns a textual representation of the type declaration
func (d *TypeDeclaration) String() string {
	result := ""
	for _, a := range d.Annotations {
		result += a.String() + "\n"
	}
	result += "type " + d.Name.String()

	if len(d.TypeParameters) > 0 {
		result += "<"
		for i, param := range d.TypeParameters {
			if i > 0 {
				result += ", "
			}
			result += param.String()
		}
		result += ">"
	}

	result += " = " + d.Type.String()
	return result
}

// Children returns the child nodes
func (d *TypeDeclaration) Children() []Node {
	children := []Node{d.Name, d.Type}

	for _, ann := range d.Annotations {
		children = append(children, ann)
	}

	for _, param := range d.TypeParameters {
		children = append(children, param)
	}

	return children
}

// FunctionDeclaration represents a function declaration
type FunctionDeclaration struct {
	BaseNode
	Name           *FunctionIdentifier // The name of the function
	TypeParameters []*GenericTypeParam // Type parameters for generic functions
	Parameters     []Node              // Function parameters
	ReturnType     Node                // The return type
	Body           Node                // The function body
	Annotations    []*Annotation       // Annotations applied to the function
	Docs           string              // Documentation comments
	IsLocal        bool                // Whether this is a local function
}

// NewFunctionDeclaration creates a new FunctionDeclaration node
func NewFunctionDeclaration(name *FunctionIdentifier, typeParams []*GenericTypeParam, params []Node, returnType Node, body Node, annotations []*Annotation, docs string, isLocal bool) *FunctionDeclaration {
	return &FunctionDeclaration{
		BaseNode:       BaseNode{NodeType: NodeFunctionDeclaration},
		Name:           name,
		TypeParameters: typeParams,
		Parameters:     params,
		ReturnType:     returnType,
		Body:           body,
		Annotations:    annotations,
		Docs:           docs,
		IsLocal:        isLocal,
	}
}

// String returns a textual representation of the function declaration
func (d *FunctionDeclaration) String() string {
	result := ""
	for _, a := range d.Annotations {
		result += a.String() + "\n"
	}

	if d.IsLocal {
		result += "local "
	}

	result += "fn " + d.Name.String()

	if len(d.TypeParameters) > 0 {
		result += "<"
		for i, param := range d.TypeParameters {
			if i > 0 {
				result += ", "
			}
			result += param.String()
		}
		result += ">"
	}

	result += "("
	for i, param := range d.Parameters {
		if i > 0 {
			result += ", "
		}
		result += param.String()
	}
	result += ")"

	if d.ReturnType != nil {
		result += " -> " + d.ReturnType.String()
	}

	result += " " + d.Body.String()
	return result
}

// Children returns the child nodes
func (d *FunctionDeclaration) Children() []Node {
	children := []Node{d.Name}

	for _, param := range d.TypeParameters {
		children = append(children, param)
	}

	children = append(children, d.Parameters...)

	if d.ReturnType != nil {
		children = append(children, d.ReturnType)
	}

	children = append(children, d.Body)

	for _, ann := range d.Annotations {
		children = append(children, ann)
	}

	return children
}

// VariableDeclaration represents a variable declaration
type VariableDeclaration struct {
	BaseNode
	Name        Node   // The name of the variable
	Type        Node   // Optional type annotation
	Initializer Node   // Optional initializer expression
	IsMutable   bool   // Whether the variable is mutable
	IsLocal     bool   // Whether the variable is local
	Docs        string // Documentation comments
}

// NewVariableDeclaration creates a new VariableDeclaration node
func NewVariableDeclaration(name Node, typ Node, initializer Node, isMutable bool, isLocal bool, docs string) *VariableDeclaration {
	return &VariableDeclaration{
		BaseNode:    BaseNode{NodeType: NodeVariableDeclaration},
		Name:        name,
		Type:        typ,
		Initializer: initializer,
		IsMutable:   isMutable,
		IsLocal:     isLocal,
		Docs:        docs,
	}
}

// String returns a textual representation of the variable declaration
func (d *VariableDeclaration) String() string {
	result := ""

	if d.IsLocal {
		result += "local "
	}

	if d.IsMutable {
		result += "mut "
	} else {
		result += "let "
	}

	result += d.Name.String()

	if d.Type != nil {
		result += ": " + d.Type.String()
	}

	if d.Initializer != nil {
		result += " = " + d.Initializer.String()
	}

	return result
}

// Children returns the child nodes
func (d *VariableDeclaration) Children() []Node {
	var children []Node

	children = append(children, d.Name)

	if d.Type != nil {
		children = append(children, d.Type)
	}

	if d.Initializer != nil {
		children = append(children, d.Initializer)
	}

	return children
}

// ConstDeclaration represents a constant declaration
type ConstDeclaration struct {
	BaseNode
	Name    Node   // The name of the constant
	Type    Node   // Optional type annotation
	Value   Node   // The constant value
	IsLocal bool   // Whether the constant is local
	Docs    string // Documentation comments
}

// NewConstDeclaration creates a new ConstDeclaration node
func NewConstDeclaration(name Node, typ Node, value Node, isLocal bool, docs string) *ConstDeclaration {
	return &ConstDeclaration{
		BaseNode: BaseNode{NodeType: NodeConstDeclaration},
		Name:     name,
		Type:     typ,
		Value:    value,
		IsLocal:  isLocal,
		Docs:     docs,
	}
}

// String returns a textual representation of the constant declaration
func (d *ConstDeclaration) String() string {
	result := ""

	if d.IsLocal {
		result += "local "
	}

	result += "const " + d.Name.String()

	if d.Type != nil {
		result += ": " + d.Type.String()
	}

	result += " = " + d.Value.String()

	return result
}

// Children returns the child nodes
func (d *ConstDeclaration) Children() []Node {
	var children []Node

	children = append(children, d.Name)

	if d.Type != nil {
		children = append(children, d.Type)
	}

	children = append(children, d.Value)

	return children
}

// EnumMember represents a member of an enum
type EnumMember struct {
	BaseNode
	Name  string // The name of the enum member
	Value Node   // Optional explicit value
	Docs  string // Documentation comments
}

// NewEnumMember creates a new EnumMember node
func NewEnumMember(name string, value Node, docs string) *EnumMember {
	return &EnumMember{
		BaseNode: BaseNode{NodeType: NodeEnumMember},
		Name:     name,
		Value:    value,
		Docs:     docs,
	}
}

// String returns a textual representation of the enum member
func (m *EnumMember) String() string {
	if m.Value != nil {
		return m.Name + " = " + m.Value.String()
	}
	return m.Name
}

// Children returns the child nodes
func (m *EnumMember) Children() []Node {
	if m.Value != nil {
		return []Node{m.Value}
	}
	return nil
}

// EnumDeclaration represents an enum declaration
type EnumDeclaration struct {
	BaseNode
	Name        *TypeIdentifier // The name of the enum
	Members     []*EnumMember   // The enum members
	Annotations []*Annotation   // Annotations applied to the enum
	Docs        string          // Documentation comments
}

// NewEnumDeclaration creates a new EnumDeclaration node
func NewEnumDeclaration(name *TypeIdentifier, members []*EnumMember, annotations []*Annotation, docs string) *EnumDeclaration {
	return &EnumDeclaration{
		BaseNode:    BaseNode{NodeType: NodeEnumDeclaration},
		Name:        name,
		Members:     members,
		Annotations: annotations,
		Docs:        docs,
	}
}

// String returns a textual representation of the enum declaration
func (d *EnumDeclaration) String() string {
	result := ""

	for _, a := range d.Annotations {
		result += a.String() + "\n"
	}

	result += "enum " + d.Name.String() + " {\n"

	for _, member := range d.Members {
		result += "  " + member.String() + ",\n"
	}

	result += "}"
	return result
}

// Children returns the child nodes
func (d *EnumDeclaration) Children() []Node {
	children := []Node{d.Name}

	for _, member := range d.Members {
		children = append(children, member)
	}

	for _, ann := range d.Annotations {
		children = append(children, ann)
	}

	return children
}

// GlobalDeclaration represents a global variable declaration
type GlobalDeclaration struct {
	BaseNode
	Name        Node   // The name of the global
	Type        Node   // Type annotation
	Initializer Node   // Optional initializer expression
	IsMutable   bool   // Whether the global is mutable
	Docs        string // Documentation comments
}

// NewGlobalDeclaration creates a new GlobalDeclaration node
func NewGlobalDeclaration(name Node, typ Node, initializer Node, isMutable bool, docs string) *GlobalDeclaration {
	return &GlobalDeclaration{
		BaseNode:    BaseNode{NodeType: NodeGlobalDeclaration},
		Name:        name,
		Type:        typ,
		Initializer: initializer,
		IsMutable:   isMutable,
		Docs:        docs,
	}
}

// String returns a textual representation of the global declaration
func (d *GlobalDeclaration) String() string {
	result := "global "

	if d.IsMutable {
		result += "mut "
	}

	result += d.Name.String() + ": " + d.Type.String()

	if d.Initializer != nil {
		result += " = " + d.Initializer.String()
	}

	return result
}

// Children returns the child nodes
func (d *GlobalDeclaration) Children() []Node {
	var children []Node

	children = append(children, d.Name, d.Type)

	if d.Initializer != nil {
		children = append(children, d.Initializer)
	}

	return children
}

// ImportDeclaration represents an import declaration
type ImportDeclaration struct {
	BaseNode
	Path     string // The path to the imported module
	Items    []Node // Imported items (can be Identifier or RenameIdentifier)
	IsExport bool   // Whether this import is re-exported
}

// NewImportDeclaration creates a new ImportDeclaration node
func NewImportDeclaration(path string, items []Node, isExport bool) *ImportDeclaration {
	return &ImportDeclaration{
		BaseNode: BaseNode{NodeType: NodeImportDeclaration},
		Path:     path,
		Items:    items,
		IsExport: isExport,
	}
}

// String returns a textual representation of the import declaration
func (d *ImportDeclaration) String() string {
	result := ""

	if d.IsExport {
		result += "export "
	}

	result += "import {"

	for i, item := range d.Items {
		if i > 0 {
			result += ", "
		}
		result += item.String()
	}

	result += "} from \"" + d.Path + "\""
	return result
}

// Children returns the child nodes
func (d *ImportDeclaration) Children() []Node {
	return d.Items
}

// ErrorDeclaration represents an error type declaration
type ErrorDeclaration struct {
	BaseNode
	Name        *TypeIdentifier // The name of the error
	Fields      []Node          // Fields for the error type
	Annotations []*Annotation   // Annotations applied to the error
	Docs        string          // Documentation comments
}

// NewErrorDeclaration creates a new ErrorDeclaration node
func NewErrorDeclaration(name *TypeIdentifier, fields []Node, annotations []*Annotation, docs string) *ErrorDeclaration {
	return &ErrorDeclaration{
		BaseNode:    BaseNode{NodeType: NodeErrorDeclaration},
		Name:        name,
		Fields:      fields,
		Annotations: annotations,
		Docs:        docs,
	}
}

// String returns a textual representation of the error declaration
func (d *ErrorDeclaration) String() string {
	result := ""

	for _, a := range d.Annotations {
		result += a.String() + "\n"
	}

	result += "error " + d.Name.String()

	if len(d.Fields) > 0 {
		result += " {\n"
		for _, field := range d.Fields {
			result += "  " + field.String() + ",\n"
		}
		result += "}"
	}

	return result
}

// Children returns the child nodes
func (d *ErrorDeclaration) Children() []Node {
	children := []Node{d.Name}

	children = append(children, d.Fields...)

	for _, ann := range d.Annotations {
		children = append(children, ann)
	}

	return children
}

// TypeAlias represents a type alias declaration
type TypeAlias struct {
	BaseNode
	Name           *TypeIdentifier     // The name of the alias
	TypeParameters []*GenericTypeParam // Type parameters if this is a generic alias
	Target         Node                // The target type
	Annotations    []*Annotation       // Annotations applied to the alias
	Docs           string              // Documentation comments
}

// NewTypeAlias creates a new TypeAlias node
func NewTypeAlias(name *TypeIdentifier, typeParams []*GenericTypeParam, target Node, annotations []*Annotation, docs string) *TypeAlias {
	return &TypeAlias{
		BaseNode:       BaseNode{NodeType: NodeTypeAlias},
		Name:           name,
		TypeParameters: typeParams,
		Target:         target,
		Annotations:    annotations,
		Docs:           docs,
	}
}

// String returns a textual representation of the type alias
func (a *TypeAlias) String() string {
	result := ""

	for _, ann := range a.Annotations {
		result += ann.String() + "\n"
	}

	result += "type " + a.Name.String()

	if len(a.TypeParameters) > 0 {
		result += "<"
		for i, param := range a.TypeParameters {
			if i > 0 {
				result += ", "
			}
			result += param.String()
		}
		result += ">"
	}

	result += " = " + a.Target.String()
	return result
}

// Children returns the child nodes
func (a *TypeAlias) Children() []Node {
	children := []Node{a.Name, a.Target}

	for _, param := range a.TypeParameters {
		children = append(children, param)
	}

	for _, ann := range a.Annotations {
		children = append(children, ann)
	}

	return children
}
