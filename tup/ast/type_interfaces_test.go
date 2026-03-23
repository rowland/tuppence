package ast

var _ LocalTypeReference = &TypeReference{}
var _ LocalTypeReference = &Identifier{}

var _ FunctionParameterType = &TypeReference{}
var _ FunctionParameterType = &Identifier{}
var _ FunctionParameterType = &NilableType{}
var _ FunctionParameterType = &FallibleType{}
var _ FunctionParameterType = &DynamicArrayType{}
var _ FunctionParameterType = &FixedSizeArrayType{}

var _ TypeDeclarationRHS = &NilableType{}
var _ TypeDeclarationRHS = &TypeTuple{}
var _ TypeDeclarationRHS = &ErrorTuple{}
var _ TypeDeclarationRHS = &DynamicArrayType{}
var _ TypeDeclarationRHS = &FixedSizeArrayType{}
var _ TypeDeclarationRHS = &UnionType{}
var _ TypeDeclarationRHS = &UnionDeclaration{}
var _ TypeDeclarationRHS = &EnumDeclaration{}
var _ TypeDeclarationRHS = &ContractDeclaration{}
var _ TypeDeclarationRHS = &TypeReference{}

var _ ArrayElementType = &TypeReference{}
var _ ArrayElementType = &DynamicArrayType{}
var _ ArrayElementType = &FixedSizeArrayType{}

var _ FunctionTypeParameter = &Parameter{}
var _ FunctionTypeParameter = &LabeledParameter{}
var _ FunctionTypeParameter = &RestParameter{}
var _ FunctionTypeParameter = &LabeledRestParameter{}

var _ FunctionTypeParameterType = &NilableType{}
var _ FunctionTypeParameterType = &DynamicArrayType{}
var _ FunctionTypeParameterType = &FixedSizeArrayType{}
var _ FunctionTypeParameterType = &FunctionType{}
var _ FunctionTypeParameterType = &ErrorTuple{}
var _ FunctionTypeParameterType = &TupleType{}
var _ FunctionTypeParameterType = &GenericType{}
var _ FunctionTypeParameterType = &TypeReference{}
var _ FunctionTypeParameterType = &Identifier{}
var _ FunctionTypeParameterType = &InlineUnion{}
var _ FunctionTypeParameterType = &UnionType{}
var _ FunctionTypeParameterType = &UnionDeclaration{}
var _ FunctionTypeParameterType = &FloatLiteral{}
var _ FunctionTypeParameterType = &IntegerLiteral{}
var _ FunctionTypeParameterType = &BooleanLiteral{}
var _ FunctionTypeParameterType = &StringLiteral{}
var _ FunctionTypeParameterType = &InterpolatedStringLiteral{}
var _ FunctionTypeParameterType = &RawStringLiteral{}
var _ FunctionTypeParameterType = &MultiLineStringLiteral{}
var _ FunctionTypeParameterType = &TupleLiteral{}
var _ FunctionTypeParameterType = &ArrayLiteral{}
var _ FunctionTypeParameterType = &SymbolLiteral{}
var _ FunctionTypeParameterType = &RuneLiteral{}

var _ TupleTypeMemberNode = &TupleTypeMember{}
var _ TupleTypeMemberNode = &LabeledTupleTypeMember{}

var _ TypeNode = &TypeReference{}
var _ TypeNode = &Identifier{}
var _ TypeNode = &DynamicArrayType{}
var _ TypeNode = &FixedSizeArrayType{}
var _ TypeNode = &FunctionType{}
var _ TypeNode = &ErrorTuple{}
var _ TypeNode = &TupleType{}
var _ TypeNode = &GenericType{}
var _ TypeNode = &InlineUnion{}

var _ TypePredicate = &TypeReference{}
var _ TypePredicate = &InlineUnion{}

var _ ContractFieldType = &NilableType{}
var _ ContractFieldType = &TypeReference{}
var _ ContractFieldType = &Identifier{}
var _ ContractFieldType = &DynamicArrayType{}
var _ ContractFieldType = &FixedSizeArrayType{}
var _ ContractFieldType = &FunctionType{}
var _ ContractFieldType = &ErrorTuple{}
var _ ContractFieldType = &TupleType{}
var _ ContractFieldType = &GenericType{}
var _ ContractFieldType = &InlineUnion{}

var _ ReturnTypeValue = &UnionWithError{}
var _ ReturnTypeValue = &UnionDeclarationWithError{}
var _ ReturnTypeValue = &NilableType{}
var _ ReturnTypeValue = &InferredErrorType{}
var _ ReturnTypeValue = &TypeReference{}
var _ ReturnTypeValue = &Identifier{}
var _ ReturnTypeValue = &DynamicArrayType{}
var _ ReturnTypeValue = &FixedSizeArrayType{}
var _ ReturnTypeValue = &FunctionType{}
var _ ReturnTypeValue = &ErrorTuple{}
var _ ReturnTypeValue = &TupleType{}
var _ ReturnTypeValue = &GenericType{}
var _ ReturnTypeValue = &InlineUnion{}

var _ InterpolatedStringPart = &StringLiteral{}
var _ InterpolatedStringPart = &Interpolation{}

var _ ConstantValue = &FloatLiteral{}
var _ ConstantValue = &IntegerLiteral{}
var _ ConstantValue = &BooleanLiteral{}
var _ ConstantValue = &StringLiteral{}
var _ ConstantValue = &InterpolatedStringLiteral{}
var _ ConstantValue = &RawStringLiteral{}
var _ ConstantValue = &MultiLineStringLiteral{}
var _ ConstantValue = &TupleLiteral{}
var _ ConstantValue = &ArrayLiteral{}
var _ ConstantValue = &SymbolLiteral{}
var _ ConstantValue = &RuneLiteral{}
var _ ConstantValue = &ScopedIdentifier{}
