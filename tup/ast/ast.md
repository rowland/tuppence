# AST Nodes Punch List

## Expressions
- [x] ArrayFunctionCall
- [x] BinaryExpression
- [x] BreakExpression
- [x] BuiltinFunctionCall
- [x] ChainedExpression
- [x] Constant
- [x] ContinueExpression
- [x] FunctionCall
- [x] IndexedAccess
- [x] Interpolation
- [x] MemberAccess
- [x] MetaExpression
- [x] PrimaryExpression
- [x] RelationalComparison
- [x] ReturnExpression
- [x] SafeIndexedAccess
- [x] TryExpression
- [x] TypeComparison
- [x] TypeConstructorCall
- [x] TypeofExpression
- [x] UFCSFunctionCall
- [x] UnaryExpression

## Control Flow
- [x] Block
- [x] BlockBody
- [x] BlockParameters
- [x] ForExpression
- [x] IfExpression
- [x] InlineForExpression
- [x] SwitchStatement

## Literals
- [x] ArrayLiteral
- [x] BooleanLiteral
- [x] FixedSizeArrayLiteral
- [x] FloatLiteral
- [x] IntegerLiteral (Decimal, Binary, Hexadecimal, Octal)
- [x] InterpolatedStringLiteral
- [x] MultiLineStringLiteral
- [x] RawStringLiteral
- [x] RuneLiteral
- [x] StringLiteral
- [x] SymbolLiteral
- [x] TupleLiteral

## Declarations
- [x] ExportDeclaration
    - [x] ExportAssignment
    - [x] ExportFunctionDeclaration
    - [x] ExportTypeDeclaration
    - [x] ExportTypeQualifiedDeclaration
    - [x] ExportTypeQualifiedFunctionDeclaration
- [x] FunctionBlock
- [x] FunctionCallContext
- [x] FunctionDeclaration
- [x] FunctionDeclarationType
- [x] FunctionParameterTypes
- [x] FunctionTypeDeclaration
- [x] TypeDeclaration
- [x] TypeQualifiedDeclaration
- [x] TypeQualifiedFunctionDeclaration

## Types
- [x] ArrayType (Fixed Size, Dynamic)
- [x] ContractDeclaration
- [x] ContractField
- [x] ContractFunction
- [x] ContractMember
- [x] ContractMembers
- [x] EnumDeclaration
- [x] EnumMember
- [x] EnumMembers
- [x] ErrorTuple
- [x] FunctionType
- [x] GenericType
- [x] InlineUnion
- [x] NamedTuple
- [x] NilableType
- [x] TupleType
- [x] TypeArgumentList
- [x] TypeParameter
- [x] TypeParameters
- [x] TypeReference
- [x] UnionDeclaration
- [x] UnionMemberDeclaration
- [x] UnionMembers
- [x] UnionType
- [x] UnionWithError

## Patterns and Matching
- [x] ArrayPattern
- [x] LabeledPattern
- [x] ListMatch
- [x] MatchCondition
- [x] MatchExpression
- [x] MatchCase
- [x] PatternMatch
- [x] StructuredMatch
- [x] TuplePattern

## Assignments and Parameters
- [x] Assignment
- [x] AssignmentLhs
- [x] CompoundAssignment
- [x] ForHeader
- [x] ForInHeader
- [x] FunctionArguments
- [x] Initializer
- [x] Iterable
- [x] IterableHeader
- [x] LabeledArgument
- [x] LabeledAssignmentLhs
- [x] LabeledParameter
- [x] LabeledRestParameter
- [x] Parameter
- [x] PartialApplication
- [x] RestParameter
- [x] SpreadArgument
- [x] StepExpression

## Operators
- [x] AddSubOp
- [x] CheckedArithmeticOp
- [x] CompoundAssignmentOp
- [x] IsOp
- [x] LogicalOp
- [x] MulDivOp
- [x] RelOp
- [x] ShortCircuitOp

## Misc
- [x] Annotation
- [x] AnnotationValue
- [x] Annotations
- [x] CaseBlock
- [x] Comment
- [x] ContractImplementsAnnotation
- [x] ElseBlock
- [x] FunctionIdentifier
- [x] Identifier
- [x] LabeledTuple
- [x] LabeledTupleMembers
- [x] LabeledTupleTypeMember
- [x] Module
- [x] Range
- [x] RangeBound
- [x] RenameIdentifier
- [x] RenameType
- [x] RestOperator
- [x] Statement
- [x] TopLevelItem
- [x] TupleUpdateExpression
- [x] TypeIdentifier

## Additional Implemented Nodes
- [x] DestructuringAssignment
- [x] DestructuringPattern
- [x] Directive
- [x] ErrorDeclaration
- [x] ErrorNode
- [x] LiteralPattern
- [x] PatternIdentifier
- [x] ReturnType
- [x] SyntaxTree
- [x] TupleMember
- [x] TupleTypeMember
- [x] TypePattern
- [x] WildcardPattern
