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
	NodeEnumMembers
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
	NodeConstDeclaration
	NodeEnumDeclaration
	NodeEnumMember
	NodeErrorDeclaration
	NodeFunctionDeclaration
	NodeGenericTypeParam
	NodeGlobalDeclaration
	NodeImportDeclaration
	NodeNamespaceDeclaration
	NodeTypeAlias
	NodeTypeDeclaration
	NodeVariableDeclaration

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
	NodeBlockComment
	NodeComment
	NodeDirective
	NodeDocComment
	NodeEmptyStatement
	NodeErrorNode
	NodeLineComment
	NodeMetadata
	NodeModule
	NodeSyntaxTree

	// String interpolation related node types
	NodeStringInterpolationEscape

	// Operator node types
	NodeAddSubOp
	NodeCheckedArithmeticOp
	NodeCompoundAssignmentOp
	NodeIsOp
	NodeLogicalOp
	NodeMulDivOp
	NodeRelOp
	NodeShortCircuitOp

	// Additional pattern matching node types
	NodeAssignmentLhs
	NodeLabeledAssignmentLhs
	NodeLabeledPattern
	NodeListMatch
	NodeMatchCondition
	NodePatternMatch
	NodeStructuredMatch

	// Pattern and matching node types
	NodeArrayPattern
	NodeAssignment
	NodeDestructuringAssignment
	NodeDestructuringPattern
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
	NodeErrorTuple
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
	NodeAnnotation:           "Annotation",
	NodeConstDeclaration:     "ConstDeclaration",
	NodeEnumDeclaration:      "EnumDeclaration",
	NodeEnumMember:           "EnumMember",
	NodeErrorDeclaration:     "ErrorDeclaration",
	NodeFunctionDeclaration:  "FunctionDeclaration",
	NodeGenericTypeParam:     "GenericTypeParam",
	NodeGlobalDeclaration:    "GlobalDeclaration",
	NodeImportDeclaration:    "ImportDeclaration",
	NodeNamespaceDeclaration: "NamespaceDeclaration",
	NodeTypeAlias:            "TypeAlias",
	NodeTypeDeclaration:      "TypeDeclaration",
	NodeVariableDeclaration:  "VariableDeclaration",

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
	NodeBlockComment:   "BlockComment",
	NodeComment:        "Comment",
	NodeDirective:      "Directive",
	NodeDocComment:     "DocComment",
	NodeEmptyStatement: "EmptyStatement",
	NodeErrorNode:      "ErrorNode",
	NodeLineComment:    "LineComment",
	NodeMetadata:       "Metadata",
	NodeModule:         "Module",
	NodeSyntaxTree:     "SyntaxTree",

	// String interpolation related node types
	NodeStringInterpolationEscape: "StringInterpolationEscape",

	// Operator node types
	NodeAddSubOp:             "AddSubOp",
	NodeCheckedArithmeticOp:  "CheckedArithmeticOp",
	NodeCompoundAssignmentOp: "CompoundAssignmentOp",
	NodeIsOp:                 "IsOp",
	NodeLogicalOp:            "LogicalOp",
	NodeMulDivOp:             "MulDivOp",
	NodeRelOp:                "RelOp",
	NodeShortCircuitOp:       "ShortCircuitOp",

	// Additional pattern matching node types
	NodeAssignmentLhs:        "AssignmentLhs",
	NodeLabeledAssignmentLhs: "LabeledAssignmentLhs",
	NodeLabeledPattern:       "LabeledPattern",
	NodeListMatch:            "ListMatch",
	NodeMatchCondition:       "MatchCondition",
	NodePatternMatch:         "PatternMatch",
	NodeStructuredMatch:      "StructuredMatch",

	// Pattern and matching node types
	NodeArrayPattern:            "ArrayPattern",
	NodeAssignment:              "Assignment",
	NodeDestructuringAssignment: "DestructuringAssignment",
	NodeDestructuringPattern:    "DestructuringPattern",
	NodeLiteralPattern:          "LiteralPattern",
	NodeMatchCase:               "MatchCase",
	NodeMatchExpression:         "MatchExpression",
	NodePatternIdentifier:       "PatternIdentifier",
	NodeTuplePattern:            "TuplePattern",
	NodeTypePattern:             "TypePattern",
	NodeWildcardPattern:         "WildcardPattern",

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
	NodeErrorTuple:    "ErrorTuple",
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
