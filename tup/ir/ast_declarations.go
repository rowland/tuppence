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
			result.References = append(result.References, TypeReference{
				ID:     TypeReferenceID(irDeclaration.ID),
				Scope:  scope,
				Name:   reference.String(),
				Origin: Origin{Primary: Span{File: file, Start: uint32(reference.StartOffset), End: uint32(reference.StartOffset + reference.Length)}},
			})
		} else if array, ok := declaration.RHS.(*ast.DynamicArrayType); ok {
			element, ok := array.ElementType.(*ast.TypeReference)
			if !ok || element == nil {
				result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_UNSUPPORTED", Message: "dynamic array elements must be type references", Origin: origin})
				continue
			}
			referenceID := TypeReferenceID(0x80000000) | TypeReferenceID(irDeclaration.ID)
			result.References = append(result.References, TypeReference{
				ID:     referenceID,
				Scope:  scope,
				Name:   element.String(),
				Origin: Origin{Primary: Span{File: file, Start: uint32(element.StartOffset), End: uint32(element.StartOffset + element.Length)}},
			})
			result.Structured = append(result.Structured, CollectedStructuredType{ID: TypeID(irDeclaration.ID), Declaration: irDeclaration.ID, Kind: TypeArray, ElementReferences: []TypeReferenceID{referenceID}, ElementOrder: []CollectedStructuredElement{{Reference: referenceID}}, Origin: origin})
		} else if array, ok := declaration.RHS.(*ast.FixedSizeArrayType); ok {
			element, elementOK := array.ElementType.(*ast.TypeReference)
			size, sizeOK := array.Size.(*ast.IntegerLiteral)
			if !elementOK || element == nil || !sizeOK || size == nil || size.IntegerValue < 0 {
				result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_INVALID", Message: "fixed array requires a non-negative integer size and type-reference element", Origin: origin})
				continue
			}
			referenceID := TypeReferenceID(0x80000000) | TypeReferenceID(irDeclaration.ID)
			result.References = append(result.References, TypeReference{
				ID:     referenceID,
				Scope:  scope,
				Name:   element.String(),
				Origin: Origin{Primary: Span{File: file, Start: uint32(element.StartOffset), End: uint32(element.StartOffset + element.Length)}},
			})
			result.Structured = append(result.Structured, CollectedStructuredType{ID: TypeID(irDeclaration.ID), Declaration: irDeclaration.ID, Kind: TypeArray, ElementReferences: []TypeReferenceID{referenceID}, ElementOrder: []CollectedStructuredElement{{Reference: referenceID}}, FixedLength: true, Length: uint64(size.IntegerValue), Origin: origin})
		} else if union, ok := declaration.RHS.(*ast.UnionType); ok {
			if len(union.Members) == 0 {
				result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_INVALID", Message: "union declarations require at least one member", Origin: origin})
				continue
			}
			referenceIDs := make([]TypeReferenceID, 0, len(union.Members))
			elementTypeIDs := make([]TypeID, 0, len(union.Members))
			elementOrder := make([]CollectedStructuredElement, 0, len(union.Members))
			for i, member := range union.Members {
				reference, ok := member.(*ast.TypeReference)
				if ok && reference != nil {
					referenceID := TypeReferenceID(0x20000000) | (TypeReferenceID(irDeclaration.ID) << 8) | TypeReferenceID(i+1)
					referenceIDs = append(referenceIDs, referenceID)
					elementOrder = append(elementOrder, CollectedStructuredElement{Reference: referenceID})
					result.References = append(result.References, TypeReference{ID: referenceID, Scope: scope, Name: reference.String(), Origin: Origin{Primary: Span{File: file, Start: uint32(reference.StartOffset), End: uint32(reference.StartOffset + reference.Length)}}})
					continue
				}
				array, arrayOK := member.(*ast.DynamicArrayType)
				fixed, fixedOK := member.(*ast.FixedSizeArrayType)
				var arrayElement ast.ArrayElementType
				fixedLength := uint64(0)
				if arrayOK {
					arrayElement = array.ElementType
				} else if fixedOK {
					arrayElement = fixed.ElementType
					size, sizeOK := fixed.Size.(*ast.IntegerLiteral)
					if !sizeOK || size == nil || size.IntegerValue < 0 {
						fixedOK = false
					} else {
						fixedLength = uint64(size.IntegerValue)
					}
				}
				arrayReference, arrayReferenceOK := arrayElement.(*ast.TypeReference)
				if (!arrayOK && !fixedOK) || !arrayReferenceOK || arrayReference == nil {
					result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_UNSUPPORTED", Message: "union members must be type references or arrays of type references", Origin: origin})
					referenceIDs = nil
					elementTypeIDs = nil
					break
				}
				childID := TypeID(0x50000000) | (TypeID(irDeclaration.ID) << 8) | TypeID(i+1)
				childReferenceID := TypeReferenceID(0x51000000) | (TypeReferenceID(irDeclaration.ID) << 8) | TypeReferenceID(i+1)
				result.References = append(result.References, TypeReference{ID: childReferenceID, Scope: scope, Name: arrayReference.String(), Origin: Origin{Primary: Span{File: file, Start: uint32(arrayReference.StartOffset), End: uint32(arrayReference.StartOffset + arrayReference.Length)}}})
				result.Structured = append(result.Structured, CollectedStructuredType{ID: childID, Declaration: irDeclaration.ID, Kind: TypeArray, ElementReferences: []TypeReferenceID{childReferenceID}, ElementOrder: []CollectedStructuredElement{{Reference: childReferenceID}}, FixedLength: fixedOK, Length: fixedLength, Origin: origin})
				elementTypeIDs = append(elementTypeIDs, childID)
				elementOrder = append(elementOrder, CollectedStructuredElement{Type: childID})
			}
			if referenceIDs != nil {
				result.Structured = append(result.Structured, CollectedStructuredType{ID: TypeID(irDeclaration.ID), Declaration: irDeclaration.ID, Kind: TypeUnion, ElementReferences: referenceIDs, ElementTypes: elementTypeIDs, ElementOrder: elementOrder, Origin: origin})
			}
		} else if tuple, ok := declaration.RHS.(*ast.TypeTuple); ok && tuple.TupleType != nil {
			referenceIDs := make([]TypeReferenceID, 0, len(tuple.TupleType.Members))
			elementTypeIDs := make([]TypeID, 0, len(tuple.TupleType.Members))
			elementOrder := make([]CollectedStructuredElement, 0, len(tuple.TupleType.Members))
			for i, member := range tuple.TupleType.Members {
				var element ast.FunctionTypeParameterType
				switch member := member.(type) {
				case *ast.TupleTypeMember:
					element = member.Type
				case *ast.LabeledTupleTypeMember:
					element = member.Type
				}
				reference, ok := element.(*ast.TypeReference)
				if ok && reference != nil {
					referenceID := TypeReferenceID(0x40000000) | (TypeReferenceID(irDeclaration.ID) << 8) | TypeReferenceID(i+1)
					referenceIDs = append(referenceIDs, referenceID)
					elementOrder = append(elementOrder, CollectedStructuredElement{Reference: referenceID})
					result.References = append(result.References, TypeReference{ID: referenceID, Scope: scope, Name: reference.String(), Origin: Origin{Primary: Span{File: file, Start: uint32(reference.StartOffset), End: uint32(reference.StartOffset + reference.Length)}}})
					continue
				}
				if nested, nestedOK := element.(*ast.TupleType); nestedOK {
					childID := TypeID(0x52000000) | (TypeID(irDeclaration.ID) << 8) | TypeID(i+1)
					childRefs := make([]TypeReferenceID, 0, len(nested.Members))
					childOrder := make([]CollectedStructuredElement, 0, len(nested.Members))
					for j, nestedMember := range nested.Members {
						var nestedElement ast.FunctionTypeParameterType
						switch nestedMember := nestedMember.(type) {
						case *ast.TupleTypeMember:
							nestedElement = nestedMember.Type
						case *ast.LabeledTupleTypeMember:
							nestedElement = nestedMember.Type
						}
						nestedReference, nestedReferenceOK := nestedElement.(*ast.TypeReference)
						if !nestedReferenceOK || nestedReference == nil {
							result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_UNSUPPORTED", Message: "nested tuple members must be type references", Origin: origin})
							referenceIDs = nil
							break
						}
						nestedReferenceID := TypeReferenceID(0x53000000) | (TypeReferenceID(irDeclaration.ID) << 8) | TypeReferenceID(i*16+j+1)
						childRefs = append(childRefs, nestedReferenceID)
						childOrder = append(childOrder, CollectedStructuredElement{Reference: nestedReferenceID})
						result.References = append(result.References, TypeReference{ID: nestedReferenceID, Scope: scope, Name: nestedReference.String(), Origin: Origin{Primary: Span{File: file, Start: uint32(nestedReference.StartOffset), End: uint32(nestedReference.StartOffset + nestedReference.Length)}}})
					}
					if referenceIDs == nil {
						break
					}
					result.Structured = append(result.Structured, CollectedStructuredType{ID: childID, Declaration: irDeclaration.ID, Kind: TypeTuple, ElementReferences: childRefs, ElementOrder: childOrder, Origin: origin})
					elementTypeIDs = append(elementTypeIDs, childID)
					elementOrder = append(elementOrder, CollectedStructuredElement{Type: childID})
					continue
				}
				if nested, nestedOK := element.(*ast.UnionType); nestedOK {
					childID := TypeID(0x54000000) | (TypeID(irDeclaration.ID) << 8) | TypeID(i+1)
					childRefs := make([]TypeReferenceID, 0, len(nested.Members))
					childOrder := make([]CollectedStructuredElement, 0, len(nested.Members))
					for j, nestedMember := range nested.Members {
						nestedReference, nestedReferenceOK := nestedMember.(*ast.TypeReference)
						if !nestedReferenceOK || nestedReference == nil {
							result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_UNSUPPORTED", Message: "nested union members must be type references", Origin: origin})
							referenceIDs = nil
							break
						}
						nestedReferenceID := TypeReferenceID(0x55000000) | (TypeReferenceID(irDeclaration.ID) << 8) | TypeReferenceID(i*16+j+1)
						childRefs = append(childRefs, nestedReferenceID)
						childOrder = append(childOrder, CollectedStructuredElement{Reference: nestedReferenceID})
						result.References = append(result.References, TypeReference{ID: nestedReferenceID, Scope: scope, Name: nestedReference.String(), Origin: Origin{Primary: Span{File: file, Start: uint32(nestedReference.StartOffset), End: uint32(nestedReference.StartOffset + nestedReference.Length)}}})
					}
					if referenceIDs == nil {
						break
					}
					result.Structured = append(result.Structured, CollectedStructuredType{ID: childID, Declaration: irDeclaration.ID, Kind: TypeUnion, ElementReferences: childRefs, ElementOrder: childOrder, Origin: origin})
					elementTypeIDs = append(elementTypeIDs, childID)
					elementOrder = append(elementOrder, CollectedStructuredElement{Type: childID})
					continue
				}
				array, arrayOK := element.(*ast.DynamicArrayType)
				fixed, fixedOK := element.(*ast.FixedSizeArrayType)
				var arrayElement ast.ArrayElementType
				fixedLength := uint64(0)
				if arrayOK {
					arrayElement = array.ElementType
				} else if fixedOK {
					arrayElement = fixed.ElementType
					size, sizeOK := fixed.Size.(*ast.IntegerLiteral)
					if !sizeOK || size == nil || size.IntegerValue < 0 {
						fixedOK = false
					} else {
						fixedLength = uint64(size.IntegerValue)
					}
				}
				arrayReference, arrayReferenceOK := arrayElement.(*ast.TypeReference)
				if (!arrayOK && !fixedOK) || !arrayReferenceOK || arrayReference == nil {
					result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_UNSUPPORTED", Message: "tuple members must be type references or arrays of type references", Origin: origin})
					referenceIDs = nil
					elementTypeIDs = nil
					break
				}
				childID := TypeID(0x60000000) | (TypeID(irDeclaration.ID) << 8) | TypeID(i+1)
				childReferenceID := TypeReferenceID(0x70000000) | (TypeReferenceID(irDeclaration.ID) << 8) | TypeReferenceID(i+1)
				result.References = append(result.References, TypeReference{ID: childReferenceID, Scope: scope, Name: arrayReference.String(), Origin: Origin{Primary: Span{File: file, Start: uint32(arrayReference.StartOffset), End: uint32(arrayReference.StartOffset + arrayReference.Length)}}})
				result.Structured = append(result.Structured, CollectedStructuredType{ID: childID, Declaration: irDeclaration.ID, Kind: TypeArray, ElementReferences: []TypeReferenceID{childReferenceID}, FixedLength: fixedOK, Length: fixedLength, Origin: origin})
				elementTypeIDs = append(elementTypeIDs, childID)
				elementOrder = append(elementOrder, CollectedStructuredElement{Type: childID})
			}
			if referenceIDs != nil {
				result.Structured = append(result.Structured, CollectedStructuredType{ID: TypeID(irDeclaration.ID), Declaration: irDeclaration.ID, Kind: TypeTuple, ElementReferences: referenceIDs, ElementTypes: elementTypeIDs, ElementOrder: elementOrder, Origin: origin})
			}
		} else {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: SeverityError, Code: "TYPE_BODY_UNSUPPORTED", Message: "structured type declaration bodies are not yet supported", Origin: origin})
		}
	}
	return result
}
