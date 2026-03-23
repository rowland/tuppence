package ast

// NodeType represents the type of AST node
type NodeType byte

const (
	// Annotation-related node types
	NodeAnnotations NodeType = iota
	NodeAnnotationValue

	// Contract and union node types
	NodeContractDeclaration
	NodeContractField
	NodeContractFunction
	NodeContractImplementsAnnotation
	NodeContractMember
	NodeContractMembers
	NodeUnionDeclaration
	NodeUnionDeclarationWithError
	NodeUnionMemberDeclaration
	NodeUnionMembers

	// Control flow node types
	NodeBlock
	NodeBlockBody
	NodeBlockParameters
	NodeSwitchCase
	NodeElseBlock
	NodeForBlock
	NodeForExpression
	NodeForHeader
	NodeForInHeader
	NodeIfExpression
	NodeInlineForExpression
	NodeIterableHeader
	NodeSwitchExpression

	// Declaration node types
	NodeAnnotation
	NodeEnumDeclaration
	NodeEnumMember
	NodeEnumMembers
	NodeErrorDeclaration
	NodeFunctionDeclaration
	NodeFunctionDeclarationLHS
	NodeGenericTypeParam
	NodeTypeDeclaration
	NodeTypeDeclarationLHS
	NodeTypeParameter
	NodeTypeParameters

	// Export related node types
	NodeExportAssignment
	NodeExportDeclaration
	NodeExportFunctionDeclaration
	NodeExportFunctionTypeDeclaration
	NodeExportTypeDeclaration
	NodeExportTypeQualifiedDeclaration
	NodeExportTypeQualifiedFunctionDeclaration
	NodeFunctionDeclarationType
	NodeFunctionParameterTypes
	NodeFunctionTypeDeclaration
	NodeTypeQualifiedDeclaration
	NodeTypeQualifiedFunctionDeclaration

	// Expression node types
	NodeArgument
	NodeArguments
	NodeArrayFunctionCall
	NodeBinaryExpression
	NodeBreakExpression
	NodeBuiltinFunctionCall
	NodeChainedExpression
	NodeConstant
	NodeContinueExpression
	NodeFunctionCall
	NodeImportExpression
	NodeIndexedAccess
	NodeMemberAccess
	NodeMetaExpression
	NodeLogicalOrExpression
	NodeLogicalAndExpression
	NodeComparisonExpression
	NodeAddSubExpression
	NodeMulDivExpression
	NodePowExpression
	NodeRelationalComparison
	NodeReturnExpression
	NodeSafeIndexedAccess
	NodeTryExpression
	NodeTupleUpdateExpression
	NodeTypeComparison
	NodeTypeConstructorCall
	NodeTypeofExpression
	NodeUFCSFunctionCall
	NodePrefixedUnaryExpression
	NodeUnaryExpression

	// Identifier node types
	NodeFunctionIdentifier
	NodeIdentifier
	NodeItExpression
	NodeScopedFunctionIdentifier
	NodeScopedIdentifier
	NodeRenameIdentifier
	NodeRenameType
	NodeTypeIdentifier

	// Literal node types
	NodeArrayLiteral
	NodeBooleanLiteral
	NodeFloatLiteral
	NodeIntegerLiteral
	NodeInterpolatedStringLiteral
	NodeInterpolation
	NodeMultiLineStringLiteral
	NodeRawStringLiteral
	NodeRuneLiteral
	NodeStringLiteral
	NodeSymbolLiteral
	NodeTupleLiteral
	NodeTupleMember

	// Miscellaneous node types
	NodeComment
	NodeErrorNode
	NodeModule
	NodeSyntaxTree

	// Operator node types
	NodeAddSubOp
	NodeCheckedArithmeticOp
	NodeCompoundAssignmentOp
	NodeIsOp
	NodeLogicalOrOp
	NodeLogicalAndOp
	NodeMulDivOp
	NodeRelOp
	NodeShortCircuitOp

	// Additional pattern matching node types
	NodeOrdinalAssignmentLHS
	NodeLabeledAssignmentLHS
	NodeLabeledPattern
	NodeLabeledPatternMember
	NodeListMatch

	// Pattern and matching node types
	NodeArrayPattern
	NodeAssignment
	NodeCompoundAssignment
	NodeTuplePattern
	NodeTypedPattern
	NodeWildcardPattern

	// PrimaryExpression node types
	NodeFunctionArguments
	NodeFunctionBlock
	NodeFunctionCallContext
	NodeInitializer
	NodeIterable
	NodeLabeledArgument
	NodeLabeledArguments
	NodePartialApplication
	NodePrimaryExpression
	NodeSpreadArgument
	NodeStatement
	NodeStepExpression
	NodeTopLevelItem

	// Range-related node types
	NodeRange
	NodeRangeBound
	NodeRestOperator

	// Tuple-related node types
	NodeLabeledTuple
	NodeLabeledTupleMembers

	// Type node types
	NodeArrayType
	NodeFunctionType
	NodeTypeArgument
	NodeTypeArgumentList
	NodeGenericType
	NodeInlineUnion
	NodeNamedTuple
	NodeNilableType
	NodeFallibleType
	NodeInferredErrorType
	NodeParameter
	NodeLabeledParameter
	NodeLabeledRestParameter
	NodeErrorTuple
	NodeRestParameter
	NodeReturnType
	NodeTupleType
	NodeTypeTuple
	NodeTupleTypeMember
	NodeLabeledTupleTypeMember
	NodeTypeReference
	NodeUnionType
	NodeUnionWithError
)

func (t NodeType) String() string {
	return NodeTypes[t]
}

// NodeTypes maps node types to their string representations
var NodeTypes = map[NodeType]string{
	// Annotation-related node types
	NodeAnnotations:     "Annotations",
	NodeAnnotationValue: "AnnotationValue",

	// Contract and union node types
	NodeContractDeclaration:          "ContractDeclaration",
	NodeContractField:                "ContractField",
	NodeContractFunction:             "ContractFunction",
	NodeContractImplementsAnnotation: "ContractImplementsAnnotation",
	NodeContractMember:               "ContractMember",
	NodeContractMembers:              "ContractMembers",
	NodeEnumMembers:                  "EnumMembers",
	NodeUnionDeclaration:             "UnionDeclaration",
	NodeUnionDeclarationWithError:    "UnionDeclarationWithError",
	NodeUnionMemberDeclaration:       "UnionMemberDeclaration",
	NodeUnionMembers:                 "UnionMembers",

	// Control flow node types
	NodeBlock:               "Block",
	NodeBlockBody:           "BlockBody",
	NodeBlockParameters:     "BlockParameters",
	NodeSwitchCase:          "SwitchCase",
	NodeElseBlock:           "ElseBlock",
	NodeForBlock:            "ForBlock",
	NodeForExpression:       "ForExpression",
	NodeForHeader:           "ForHeader",
	NodeForInHeader:         "ForInHeader",
	NodeIfExpression:        "IfExpression",
	NodeInlineForExpression: "InlineForExpression",
	NodeIterableHeader:      "IterableHeader",
	NodeSwitchExpression:    "SwitchExpression",

	// Declaration node types
	NodeAnnotation:             "Annotation",
	NodeEnumDeclaration:        "EnumDeclaration",
	NodeEnumMember:             "EnumMember",
	NodeErrorDeclaration:       "ErrorDeclaration",
	NodeFunctionDeclaration:    "FunctionDeclaration",
	NodeFunctionDeclarationLHS: "FunctionDeclarationLHS",
	NodeGenericTypeParam:       "GenericTypeParam",
	NodeTypeDeclaration:        "TypeDeclaration",
	NodeTypeDeclarationLHS:     "TypeDeclarationLHS",
	NodeTypeParameter:          "TypeParameter",
	NodeTypeParameters:         "TypeParameters",

	// Export related node types
	NodeExportAssignment:                       "ExportAssignment",
	NodeExportDeclaration:                      "ExportDeclaration",
	NodeExportFunctionDeclaration:              "ExportFunctionDeclaration",
	NodeExportFunctionTypeDeclaration:          "ExportFunctionTypeDeclaration",
	NodeExportTypeDeclaration:                  "ExportTypeDeclaration",
	NodeExportTypeQualifiedDeclaration:         "ExportTypeQualifiedDeclaration",
	NodeExportTypeQualifiedFunctionDeclaration: "ExportTypeQualifiedFunctionDeclaration",
	NodeFunctionDeclarationType:                "FunctionDeclarationType",
	NodeFunctionParameterTypes:                 "FunctionParameterTypes",
	NodeFunctionTypeDeclaration:                "FunctionTypeDeclaration",
	NodeTypeQualifiedDeclaration:               "TypeQualifiedDeclaration",
	NodeTypeQualifiedFunctionDeclaration:       "TypeQualifiedFunctionDeclaration",

	// Expression node types
	NodeArgument:                "Argument",
	NodeArguments:               "Arguments",
	NodeArrayFunctionCall:       "ArrayFunctionCall",
	NodeBinaryExpression:        "BinaryExpression",
	NodeBreakExpression:         "BreakExpression",
	NodeBuiltinFunctionCall:     "BuiltinFunctionCall",
	NodeChainedExpression:       "ChainedExpression",
	NodeConstant:                "NodeConstant",
	NodeContinueExpression:      "ContinueExpression",
	NodeFunctionCall:            "FunctionCall",
	NodeImportExpression:        "ImportExpression",
	NodeIndexedAccess:           "IndexedAccess",
	NodeMemberAccess:            "MemberAccess",
	NodeMetaExpression:          "MetaExpression",
	NodeLogicalOrExpression:     "LogicalOrExpression",
	NodeLogicalAndExpression:    "LogicalAndExpression",
	NodeComparisonExpression:    "ComparisonExpression",
	NodeAddSubExpression:        "AddSubExpression",
	NodeMulDivExpression:        "MulDivExpression",
	NodePowExpression:           "PowExpression",
	NodeRelationalComparison:    "RelationalComparison",
	NodeReturnExpression:        "ReturnExpression",
	NodeSafeIndexedAccess:       "SafeIndexedAccess",
	NodeTryExpression:           "TryExpression",
	NodeTupleUpdateExpression:   "TupleUpdateExpression",
	NodeTypeComparison:          "TypeComparison",
	NodeTypeConstructorCall:     "TypeConstructorCall",
	NodeTypeofExpression:        "TypeofExpression",
	NodeUFCSFunctionCall:        "UFCSFunctionCall",
	NodePrefixedUnaryExpression: "PrefixedUnaryExpression",
	NodeUnaryExpression:         "UnaryExpression",

	// Identifier node types
	NodeFunctionIdentifier:       "FunctionIdentifier",
	NodeIdentifier:               "Identifier",
	NodeItExpression:             "ItExpression",
	NodeScopedFunctionIdentifier: "ScopedFunctionIdentifier",
	NodeScopedIdentifier:         "ScopedIdentifier",
	NodeRenameIdentifier:         "RenameIdentifier",
	NodeRenameType:               "RenameType",
	NodeTypeIdentifier:           "TypeIdentifier",

	// Literal node types
	NodeArrayLiteral:              "ArrayLiteral",
	NodeBooleanLiteral:            "BooleanLiteral",
	NodeFloatLiteral:              "FloatLiteral",
	NodeIntegerLiteral:            "IntegerLiteral",
	NodeInterpolatedStringLiteral: "InterpolatedStringLiteral",
	NodeInterpolation:             "Interpolation",
	NodeMultiLineStringLiteral:    "MultiLineStringLiteral",
	NodeRawStringLiteral:          "RawStringLiteral",
	NodeRuneLiteral:               "RuneLiteral",
	NodeStringLiteral:             "StringLiteral",
	NodeSymbolLiteral:             "SymbolLiteral",
	NodeTupleLiteral:              "TupleLiteral",
	NodeTupleMember:               "TupleMember",

	// Miscellaneous node types
	NodeComment:    "Comment",
	NodeModule:     "Module",
	NodeSyntaxTree: "SyntaxTree",

	// Operator node types
	NodeAddSubOp:             "AddSubOp",
	NodeCheckedArithmeticOp:  "CheckedArithmeticOp",
	NodeCompoundAssignmentOp: "CompoundAssignmentOp",
	NodeIsOp:                 "IsOp",
	NodeLogicalOrOp:          "LogicalOrOp",
	NodeLogicalAndOp:         "LogicalAndOp",
	NodeMulDivOp:             "MulDivOp",
	NodeRelOp:                "RelOp",
	NodeShortCircuitOp:       "ShortCircuitOp",

	// Additional pattern matching node types
	NodeOrdinalAssignmentLHS: "OrdinalAssignmentLHS",
	NodeLabeledAssignmentLHS: "LabeledAssignmentLhs",
	NodeLabeledPattern:       "LabeledPattern",
	NodeLabeledPatternMember: "LabeledPatternMember",
	NodeListMatch:            "ListMatch",

	// Pattern and matching node types
	NodeArrayPattern:       "ArrayPattern",
	NodeAssignment:         "Assignment",
	NodeCompoundAssignment: "CompoundAssignment",
	NodeTuplePattern:       "TuplePattern",
	NodeTypedPattern:       "TypedPattern",
	NodeWildcardPattern:    "WildcardPattern",

	// PrimaryExpression node types
	NodeFunctionArguments:   "FunctionArguments",
	NodeFunctionBlock:       "FunctionBlock",
	NodeFunctionCallContext: "FunctionCallContext",
	NodeInitializer:         "Initializer",
	NodeIterable:            "Iterable",
	NodeLabeledArgument:     "LabeledArgument",
	NodeLabeledArguments:    "LabeledArguments",
	NodePartialApplication:  "PartialApplication",
	NodePrimaryExpression:   "PrimaryExpression",
	NodeSpreadArgument:      "SpreadArgument",
	NodeStatement:           "Statement",
	NodeStepExpression:      "StepExpression",
	NodeTopLevelItem:        "TopLevelItem",

	// Range-related node types
	NodeRange:        "Range",
	NodeRangeBound:   "RangeBound",
	NodeRestOperator: "RestOperator",

	// Tuple-related node types
	NodeLabeledTuple:        "LabeledTuple",
	NodeLabeledTupleMembers: "LabeledTupleMembers",

	// Type node types
	NodeArrayType:              "ArrayType",
	NodeFunctionType:           "FunctionType",
	NodeTypeArgument:           "TypeArgument",
	NodeTypeArgumentList:       "TypeArgumentList",
	NodeGenericType:            "GenericType",
	NodeInlineUnion:            "InlineUnion",
	NodeParameter:              "Parameter",
	NodeLabeledParameter:       "LabeledParameter",
	NodeLabeledRestParameter:   "LabeledRestParameter",
	NodeNamedTuple:             "NamedTuple",
	NodeNilableType:            "NilableType",
	NodeFallibleType:           "FallibleType",
	NodeInferredErrorType:      "InferredErrorType",
	NodeErrorTuple:             "ErrorTuple",
	NodeRestParameter:          "RestParameter",
	NodeReturnType:             "ReturnType",
	NodeTupleType:              "TupleType",
	NodeTypeTuple:              "TypeTuple",
	NodeTupleTypeMember:        "TupleTypeMember",
	NodeLabeledTupleTypeMember: "LabeledTupleTypeMember",
	NodeTypeReference:          "TypeReference",
	NodeUnionType:              "UnionType",
	NodeUnionWithError:         "UnionWithError",
}
