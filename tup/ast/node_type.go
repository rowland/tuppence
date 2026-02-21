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
	NodeUnionMemberDeclaration
	NodeUnionMembers

	// Control flow node types
	NodeBlock
	NodeBlockBody
	NodeBlockParameters
	NodeCaseBlock
	NodeElseBlock
	NodeForExpression
	NodeIfExpression
	NodeInlineForExpression
	NodeSwitchStatement

	// Declaration node types
	NodeAnnotation
	NodeEnumDeclaration
	NodeEnumMember
	NodeEnumMembers
	NodeErrorDeclaration
	NodeFunctionDeclaration
	NodeGenericTypeParam
	NodeTypeDeclaration

	// Export related node types
	NodeExportAssignment
	NodeExportDeclaration
	NodeExportFunctionDeclaration
	NodeExportTypeDeclaration
	NodeExportTypeQualifiedDeclaration
	NodeExportTypeQualifiedFunctionDeclaration
	NodeFunctionDeclarationType
	NodeFunctionParameterTypes
	NodeFunctionTypeDeclaration
	NodeTypeQualifiedDeclaration
	NodeTypeQualifiedFunctionDeclaration

	// Expression node types
	NodeArrayFunctionCall
	NodeBinaryExpression
	NodeBreakExpression
	NodeBuiltinFunctionCall
	NodeChainedExpression
	NodeConstant
	NodeContinueExpression
	NodeFunctionCall
	NodeIndexedAccess
	NodeMemberAccess
	NodeMetaExpression
	NodeRelationalComparison
	NodeReturnExpression
	NodeSafeIndexedAccess
	NodeTryExpression
	NodeTupleUpdateExpression
	NodeTypeComparison
	NodeTypeConstructorCall
	NodeTypeofExpression
	NodeUFCSFunctionCall
	NodeUnaryExpression

	// Identifier node types
	NodeFunctionIdentifier
	NodeIdentifier
	NodeRenameIdentifier
	NodeRenameType
	NodeTypeIdentifier

	// Literal node types
	NodeArrayLiteral
	NodeBooleanLiteral
	NodeFixedSizeArrayLiteral
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
	NodeListMatch
	NodeMatchCondition
	NodePatternMatch
	NodeStructuredMatch

	// Pattern and matching node types
	NodeArrayPattern
	NodeAssignment
	NodeLiteralPattern
	NodeMatchCase
	NodeMatchExpression
	NodePatternIdentifier
	NodeTuplePattern
	NodeTypePattern
	NodeWildcardPattern

	// PrimaryExpression node types
	NodeFunctionArguments
	NodeFunctionBlock
	NodeFunctionCallContext
	NodeInitializer
	NodeIterable
	NodeLabeledArgument
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
	NodeGenericType
	NodeInlineUnion
	NodeNamedTuple
	NodeNilableType
	NodeRestParameter
	NodeReturnType
	NodeTupleType
	NodeTypeReference
	NodeUnionType
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
	NodeUnionMemberDeclaration:       "UnionMemberDeclaration",
	NodeUnionMembers:                 "UnionMembers",

	// Control flow node types
	NodeBlock:               "Block",
	NodeBlockBody:           "BlockBody",
	NodeBlockParameters:     "BlockParameters",
	NodeCaseBlock:           "CaseBlock",
	NodeElseBlock:           "ElseBlock",
	NodeForExpression:       "ForExpression",
	NodeIfExpression:        "IfExpression",
	NodeInlineForExpression: "InlineForExpression",
	NodeSwitchStatement:     "SwitchStatement",

	// Declaration node types
	NodeAnnotation:          "Annotation",
	NodeEnumDeclaration:     "EnumDeclaration",
	NodeEnumMember:          "EnumMember",
	NodeErrorDeclaration:    "ErrorDeclaration",
	NodeFunctionDeclaration: "FunctionDeclaration",
	NodeGenericTypeParam:    "GenericTypeParam",
	NodeTypeDeclaration:     "TypeDeclaration",

	// Export related node types
	NodeExportAssignment:                       "ExportAssignment",
	NodeExportDeclaration:                      "ExportDeclaration",
	NodeExportFunctionDeclaration:              "ExportFunctionDeclaration",
	NodeExportTypeDeclaration:                  "ExportTypeDeclaration",
	NodeExportTypeQualifiedDeclaration:         "ExportTypeQualifiedDeclaration",
	NodeExportTypeQualifiedFunctionDeclaration: "ExportTypeQualifiedFunctionDeclaration",
	NodeFunctionDeclarationType:                "FunctionDeclarationType",
	NodeFunctionParameterTypes:                 "FunctionParameterTypes",
	NodeFunctionTypeDeclaration:                "FunctionTypeDeclaration",
	NodeTypeQualifiedDeclaration:               "TypeQualifiedDeclaration",
	NodeTypeQualifiedFunctionDeclaration:       "TypeQualifiedFunctionDeclaration",

	// Expression node types
	NodeArrayFunctionCall:     "ArrayFunctionCall",
	NodeBinaryExpression:      "BinaryExpression",
	NodeBreakExpression:       "BreakExpression",
	NodeBuiltinFunctionCall:   "BuiltinFunctionCall",
	NodeChainedExpression:     "ChainedExpression",
	NodeConstant:              "NodeConstant",
	NodeContinueExpression:    "ContinueExpression",
	NodeFunctionCall:          "FunctionCall",
	NodeIndexedAccess:         "IndexedAccess",
	NodeMemberAccess:          "MemberAccess",
	NodeMetaExpression:        "MetaExpression",
	NodeRelationalComparison:  "RelationalComparison",
	NodeReturnExpression:      "ReturnExpression",
	NodeSafeIndexedAccess:     "SafeIndexedAccess",
	NodeTryExpression:         "TryExpression",
	NodeTupleUpdateExpression: "TupleUpdateExpression",
	NodeTypeComparison:        "TypeComparison",
	NodeTypeConstructorCall:   "TypeConstructorCall",
	NodeTypeofExpression:      "TypeofExpression",
	NodeUFCSFunctionCall:      "UFCSFunctionCall",
	NodeUnaryExpression:       "UnaryExpression",

	// Identifier node types
	NodeFunctionIdentifier: "FunctionIdentifier",
	NodeIdentifier:         "Identifier",
	NodeRenameIdentifier:   "RenameIdentifier",
	NodeRenameType:         "RenameType",
	NodeTypeIdentifier:     "TypeIdentifier",

	// Literal node types
	NodeArrayLiteral:              "ArrayLiteral",
	NodeBooleanLiteral:            "BooleanLiteral",
	NodeFixedSizeArrayLiteral:     "FixedSizeArrayLiteral",
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
	NodeListMatch:            "ListMatch",
	NodeMatchCondition:       "MatchCondition",
	NodePatternMatch:         "PatternMatch",
	NodeStructuredMatch:      "StructuredMatch",

	// Pattern and matching node types
	NodeArrayPattern:      "ArrayPattern",
	NodeAssignment:        "Assignment",
	NodeLiteralPattern:    "LiteralPattern",
	NodeMatchCase:         "MatchCase",
	NodeMatchExpression:   "MatchExpression",
	NodePatternIdentifier: "PatternIdentifier",
	NodeTuplePattern:      "TuplePattern",
	NodeTypePattern:       "TypePattern",
	NodeWildcardPattern:   "WildcardPattern",

	// PrimaryExpression node types
	NodeFunctionArguments:   "FunctionArguments",
	NodeFunctionBlock:       "FunctionBlock",
	NodeFunctionCallContext: "FunctionCallContext",
	NodeInitializer:         "Initializer",
	NodeIterable:            "Iterable",
	NodeLabeledArgument:     "LabeledArgument",
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
	NodeArrayType:     "ArrayType",
	NodeFunctionType:  "FunctionType",
	NodeGenericType:   "GenericType",
	NodeInlineUnion:   "InlineUnion",
	NodeNamedTuple:    "NamedTuple",
	NodeNilableType:   "NilableType",
	NodeRestParameter: "RestParameter",
	NodeReturnType:    "ReturnType",
	NodeTupleType:     "TupleType",
	NodeTypeReference: "TypeReference",
	NodeUnionType:     "UnionType",
}
