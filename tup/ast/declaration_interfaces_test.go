package ast

var _ TopLevelItem = &Assignment{}
var _ TopLevelItem = &FunctionDeclaration{}
var _ TopLevelItem = &FunctionTypeDeclaration{}
var _ TopLevelItem = &TypeDeclaration{}
var _ TopLevelItem = &TypeQualifiedDeclaration{}
var _ TopLevelItem = &TypeQualifiedFunctionDeclaration{}
var _ TopLevelItem = &ExportTypeQualifiedFunctionDeclaration{}
var _ TopLevelItem = &ExportTypeQualifiedDeclaration{}
var _ TopLevelItem = &ExportTypeDeclaration{}
var _ TopLevelItem = &ExportFunctionDeclaration{}
var _ TopLevelItem = &ExportAssignment{}

var _ Statement = &Assignment{}
var _ Statement = &FunctionDeclaration{}
var _ Statement = &TypeDeclaration{}
var _ Statement = &TypeQualifiedDeclaration{}
var _ Statement = &TypeQualifiedFunctionDeclaration{}

var _ ExportDeclaration = &ExportTypeQualifiedFunctionDeclaration{}
var _ ExportDeclaration = &ExportTypeQualifiedDeclaration{}
var _ ExportDeclaration = &ExportTypeDeclaration{}
var _ ExportDeclaration = &ExportFunctionDeclaration{}
var _ ExportDeclaration = &ExportAssignment{}

var _ ContractMemberNode = &ContractFunction{}
var _ ContractMemberNode = &ContractField{}

var _ UnionMemberType = &NamedTuple{}
var _ UnionMemberType = &GenericType{}
var _ UnionMemberType = &DynamicArrayType{}
var _ UnionMemberType = &FixedSizeArrayType{}
var _ UnionMemberType = &TypeReference{}
var _ UnionMemberType = &Identifier{}
var _ UnionMemberType = &ContractDeclaration{}

var _ UnionDeclarationMemberType = &NamedTuple{}
var _ UnionDeclarationMemberType = &GenericType{}
var _ UnionDeclarationMemberType = &DynamicArrayType{}
var _ UnionDeclarationMemberType = &FixedSizeArrayType{}
var _ UnionDeclarationMemberType = &TypeReference{}
