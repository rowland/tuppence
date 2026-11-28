package ir

import "github.com/rowland/tuppence/tup/ast"

// CollectedTypeDeclarations contains IR declarations and scope entries derived
// from parsed type declarations. Type bodies are resolved separately.
type CollectedTypeDeclarations struct {
	Declarations DeclarationTable
	Entries      []ScopeEntry
	Types        TypeTable
	References   []TypeReference
	Structured   []CollectedStructuredType
	Diagnostics  []Diagnostic
}

// CollectedStructuredType is a structured type whose element references are
// resolved before its semantic Type record is installed.
type CollectedStructuredType struct {
	ID                TypeID
	Declaration       DeclID
	Kind              TypeKind
	ElementReferences []TypeReferenceID
	ElementTypes      []TypeID
	ElementOrder      []CollectedStructuredElement
	FixedLength       bool
	Length            uint64
	Origin            Origin
}

// CollectedStructuredElement preserves the source order of a structured type's members.
type CollectedStructuredElement struct {
	Reference TypeReferenceID
	Type      TypeID
}

// InstallCollectedTypeDeclarations adds collected declarations to a module's
// initial scope atomically.
func (p *IndexedProgram) InstallCollectedTypeDeclarations(module ModuleID, collected CollectedTypeDeclarations) bool {
	if p == nil || !module.Valid() || len(collected.Diagnostics) != 0 {
		return false
	}
	var scope *Scope
	for i := range p.Scopes {
		candidate := &p.Scopes[i]
		if candidate.Module == module && candidate.Kind == ScopeModule && candidate.Initial {
			scope = candidate
			break
		}
	}
	if scope == nil || len(collected.Declarations) != len(collected.Entries) {
		return false
	}
	seenNames := make(map[string]struct{}, len(scope.Entries)+len(collected.Entries))
	for _, entry := range scope.Entries {
		seenNames[entry.Name] = struct{}{}
	}
	for i, declaration := range collected.Declarations {
		entry := collected.Entries[i]
		if declaration.Scope != scope.ID || entry.Declaration != declaration.ID || entry.Name != declaration.Name || declaration.ID == 0 {
			return false
		}
		if _, exists := seenNames[entry.Name]; exists {
			return false
		}
		if _, exists := p.Declarations.Find(declaration.ID); exists {
			return false
		}
		if typ, exists := p.Types.Find(TypeID(declaration.ID)); exists || typ.ID.Valid() {
			return false
		}
		seenNames[entry.Name] = struct{}{}
	}
	resolved, diagnostics := p.ResolveTypeReferences(collected.References)
	resolvedByID := make(map[TypeReferenceID]TypeID, len(resolved))
	for _, reference := range resolved {
		resolvedByID[reference.Reference.ID] = reference.Type
	}
	aliases := make([]TypeAlias, 0, len(resolved)+len(collected.Structured))
	for _, reference := range resolved {
		declaration := DeclID(reference.Reference.ID)
		for _, collectedDeclaration := range collected.Declarations {
			if collectedDeclaration.ID == declaration {
				aliases = append(aliases, TypeAlias{Declaration: declaration, Target: reference.Type, Origin: reference.Reference.Origin})
				break
			}
		}
	}
	structuredTypes := make(TypeTable, 0, len(collected.Structured))
	for _, structured := range collected.Structured {
		elements := make([]TypeID, 0, len(structured.ElementOrder))
		if len(structured.ElementOrder) > 0 {
			for _, member := range structured.ElementOrder {
				if member.Reference != 0 {
					element, ok := resolvedByID[member.Reference]
					if !ok {
						elements = nil
						break
					}
					elements = append(elements, element)
					continue
				}
				_, ok := p.Types.Find(member.Type)
				if !ok {
					for _, local := range structuredTypes {
						if local.ID == member.Type {
							ok = true
							break
						}
					}
				}
				if !ok {
					elements = nil
					break
				}
				elements = append(elements, member.Type)
			}
		} else {
			for _, referenceID := range structured.ElementReferences {
				element, ok := resolvedByID[referenceID]
				if !ok {
					elements = nil
					break
				}
				elements = append(elements, element)
			}
			for _, element := range structured.ElementTypes {
				_, exists := p.Types.Find(element)
				if !exists {
					for _, local := range structuredTypes {
						if local.ID == element {
							exists = true
							break
						}
					}
				}
				if !exists {
					elements = nil
					break
				}
				elements = append(elements, element)
			}
		}
		if len(elements) == 0 || (len(structured.ElementOrder) == 0 && len(elements) != len(structured.ElementReferences)+len(structured.ElementTypes)) || (len(structured.ElementOrder) > 0 && len(elements) != len(structured.ElementOrder)) {
			continue
		}
		typ := Type{ID: structured.ID, Kind: structured.Kind, FixedLength: structured.FixedLength, Length: structured.Length, Origin: structured.Origin}
		if structured.Kind == TypeArray {
			typ.Element = elements[0]
		} else {
			typ.Elements = elements
		}
		structuredTypes = append(structuredTypes, typ)
		aliases = append(aliases, TypeAlias{Declaration: structured.Declaration, Target: structured.ID, Origin: structured.Origin})
	}
	for _, typ := range collected.Types {
		if typ.Declaration == 0 || typ.HomeModule != module || !typ.Valid() {
			return false
		}
	}
	for _, typ := range structuredTypes {
		if _, exists := p.Types.Find(typ.ID); exists {
			return false
		}
	}
	p.Declarations = append(p.Declarations, collected.Declarations...)
	scope.Entries = append(scope.Entries, collected.Entries...)
	p.Types = append(p.Types, collected.Types...)
	p.Types = append(p.Types, structuredTypes...)
	p.TypeReferences = append(p.TypeReferences, resolved...)
	p.Diagnostics = append(p.Diagnostics, diagnostics...)
	p.Aliases = append(p.Aliases, aliases...)
	return true
}

// CollectTypeDeclarations converts parsed type declaration names into semantic
// declarations without interpreting their right-hand-side type bodies.
func CollectTypeDeclarations(module ModuleID, scope ScopeID, file FileID, first DeclID, declarations []*ast.TypeDeclaration) CollectedTypeDeclarations {
	result := CollectedTypeDeclarations{}
	next := uint32(first)
	nextType := TypeID(0x71000000)
	nextReference := TypeReferenceID(0x72000000)
	addReference := func(id TypeReferenceID, reference *ast.TypeReference) {
		origin := Origin{Primary: Span{File: file, Start: uint32(reference.StartOffset), End: uint32(reference.StartOffset + reference.Length)}}
		result.References = append(result.References, TypeReference{ID: id, Scope: scope, Name: reference.String(), Origin: origin})
	}
	structuredFields := func(elements []CollectedStructuredElement) ([]TypeReferenceID, []TypeID) {
		references := make([]TypeReferenceID, 0, len(elements))
		types := make([]TypeID, 0, len(elements))
		for _, element := range elements {
			if element.Reference != 0 {
				references = append(references, element.Reference)
			} else {
				types = append(types, element.Type)
			}
		}
		return references, types
	}
	refID := func(kind TypeKind, depth, index int, declaration DeclID) TypeReferenceID {
		if depth == 0 {
			if kind == TypeArray {
				return TypeReferenceID(0x80000000) | TypeReferenceID(declaration)
			}
			if kind == TypeUnion {
				return TypeReferenceID(0x20000000) | TypeReferenceID(declaration)<<8 | TypeReferenceID(index+1)
			}
			if kind == TypeTuple {
				return TypeReferenceID(0x40000000) | TypeReferenceID(declaration)<<8 | TypeReferenceID(index+1)
			}
		}
		if depth == 1 {
			if kind == TypeTuple {
				return TypeReferenceID(0x53000000) | TypeReferenceID(declaration)<<8 | TypeReferenceID(index+1)
			}
			if kind == TypeUnion || kind == TypeArray {
				return TypeReferenceID(0x51000000) | TypeReferenceID(declaration)<<8 | TypeReferenceID(index+1)
			}
		}
		id := nextReference
		nextReference++
		return id
	}
	typeID := func(kind TypeKind, depth, index int, declaration DeclID) TypeID {
		if depth == 0 {
			if kind == TypeArray {
				return TypeID(0x50000000) | TypeID(declaration)<<8 | TypeID(index+1)
			}
			if kind == TypeTuple {
				return TypeID(0x60000000) | TypeID(declaration)<<8 | TypeID(index+1)
			}
		}
		if depth == 1 {
			if kind == TypeTuple {
				return TypeID(0x52000000) | TypeID(declaration)<<8 | TypeID(index+1)
			}
			if kind == TypeUnion {
				return TypeID(0x54000000) | TypeID(declaration)<<8 | TypeID(index+1)
			}
			if kind == TypeArray {
				return TypeID(0x60000000) | TypeID(declaration)<<8 | TypeID(index+1)
			}
		}
		id := nextType
		nextType++
		return id
	}
	var collect func(interface{}, TypeKind, int, int, DeclID, Origin) (CollectedStructuredElement, bool)
	collect = func(node interface{}, parent TypeKind, depth, index int, declaration DeclID, origin Origin) (CollectedStructuredElement, bool) {
		if reference, ok := node.(*ast.TypeReference); ok && reference != nil {
			id := refID(parent, depth, index, declaration)
			addReference(id, reference)
			return CollectedStructuredElement{Reference: id}, true
		}
		if tuple, ok := node.(*ast.TupleType); ok && tuple != nil {
			id := typeID(parent, depth, index, declaration)
			elements := make([]CollectedStructuredElement, 0, len(tuple.Members))
			for i, member := range tuple.Members {
				var memberType ast.FunctionTypeParameterType
				switch member := member.(type) {
				case *ast.TupleTypeMember:
					memberType = member.Type
				case *ast.LabeledTupleTypeMember:
					memberType = member.Type
				}
				element, ok := collect(memberType, TypeTuple, depth+1, i, declaration, origin)
				if !ok {
					return CollectedStructuredElement{}, false
				}
				elements = append(elements, element)
			}
			references, types := structuredFields(elements)
			result.Structured = append(result.Structured, CollectedStructuredType{ID: id, Declaration: declaration, Kind: TypeTuple, ElementReferences: references, ElementTypes: types, ElementOrder: elements, Origin: origin})
			return CollectedStructuredElement{Type: id}, true
		}
		if union, ok := node.(*ast.UnionType); ok && union != nil {
			if len(union.Members) == 0 {
				result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_INVALID", Message: "union declarations require at least one member", Origin: origin})
				return CollectedStructuredElement{}, false
			}
			id := typeID(parent, depth, index, declaration)
			elements := make([]CollectedStructuredElement, 0, len(union.Members))
			for i, member := range union.Members {
				element, ok := collect(member, TypeUnion, depth+1, i, declaration, origin)
				if !ok {
					return CollectedStructuredElement{}, false
				}
				elements = append(elements, element)
			}
			references, types := structuredFields(elements)
			result.Structured = append(result.Structured, CollectedStructuredType{ID: id, Declaration: declaration, Kind: TypeUnion, ElementReferences: references, ElementTypes: types, ElementOrder: elements, Origin: origin})
			return CollectedStructuredElement{Type: id}, true
		}
		var arrayElement ast.ArrayElementType
		fixed := false
		length := uint64(0)
		if array, ok := node.(*ast.DynamicArrayType); ok && array != nil {
			arrayElement = array.ElementType
		} else if array, ok := node.(*ast.FixedSizeArrayType); ok && array != nil {
			arrayElement = array.ElementType
			size, ok := array.Size.(*ast.IntegerLiteral)
			if !ok || size == nil || size.IntegerValue < 0 {
				result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_INVALID", Message: "fixed array requires a non-negative integer size", Origin: origin})
				return CollectedStructuredElement{}, false
			}
			fixed = true
			length = uint64(size.IntegerValue)
		} else {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_UNSUPPORTED", Message: "structured type members must be type references, tuples, unions, or arrays", Origin: origin})
			return CollectedStructuredElement{}, false
		}
		id := typeID(parent, depth, index, declaration)
		element, ok := collect(arrayElement, TypeArray, depth+1, 0, declaration, origin)
		if !ok {
			return CollectedStructuredElement{}, false
		}
		references, types := structuredFields([]CollectedStructuredElement{element})
		result.Structured = append(result.Structured, CollectedStructuredType{ID: id, Declaration: declaration, Kind: TypeArray, ElementReferences: references, ElementTypes: types, ElementOrder: []CollectedStructuredElement{element}, FixedLength: fixed, Length: length, Origin: origin})
		return CollectedStructuredElement{Type: id}, true
	}
	for _, declaration := range declarations {
		if declaration == nil || declaration.LHS == nil || declaration.LHS.Name == nil {
			continue
		}
		next++
		name := declaration.LHS.Name
		origin := Origin{Primary: Span{File: file, Start: uint32(name.StartOffset), End: uint32(name.StartOffset + name.Length)}}
		irDeclaration := Declaration{ID: DeclID(next), Scope: scope, Name: name.String(), Kind: DeclType, Origin: origin}
		result.Declarations = append(result.Declarations, irDeclaration)
		result.Entries = append(result.Entries, ScopeEntry{Name: irDeclaration.Name, Declaration: irDeclaration.ID})
		if reference, ok := declaration.RHS.(*ast.TypeReference); ok && reference != nil {
			addReference(TypeReferenceID(irDeclaration.ID), reference)
			continue
		}
		var root interface{}
		var kind TypeKind
		switch rhs := declaration.RHS.(type) {
		case *ast.DynamicArrayType:
			root, kind = rhs, TypeArray
		case *ast.FixedSizeArrayType:
			root, kind = rhs, TypeArray
		case *ast.UnionType:
			root, kind = rhs, TypeUnion
		case *ast.TypeTuple:
			if rhs != nil && rhs.TupleType != nil {
				root, kind = rhs.TupleType, TypeTuple
			}
		}
		if root == nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_UNSUPPORTED", Message: "structured type declaration bodies are not yet supported", Origin: origin})
			continue
		}
		_, ok := collect(root, kind, -1, 0, irDeclaration.ID, origin)
		if !ok {
			continue
		}
		rootType := &result.Structured[len(result.Structured)-1]
		if rootType.Declaration == irDeclaration.ID && rootType.Kind == kind {
			rootType.ID = TypeID(irDeclaration.ID)
		}
	}
	return result
}
