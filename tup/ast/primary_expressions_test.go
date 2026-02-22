package ast

var _ PrimaryExpression = &Block{}
var _ PrimaryExpression = &IfExpression{}
var _ PrimaryExpression = &ForExpression{}
var _ PrimaryExpression = &InlineForExpression{}
var _ PrimaryExpression = &ArrayFunctionCall{}
var _ PrimaryExpression = &TypeofExpression{}
var _ PrimaryExpression = &FunctionCall{}
var _ PrimaryExpression = &TypeConstructorCall{}
var _ PrimaryExpression = &MemberAccess{}
var _ PrimaryExpression = &TupleUpdateExpression{}
var _ PrimaryExpression = &SafeIndexedAccess{}
var _ PrimaryExpression = &IndexedAccess{}
var _ PrimaryExpression = &Identifier{}

// var _ PrimaryExpression = &ImportExpression{}
// var _ PrimaryExpression = &ReturnExpression{}
// var _ PrimaryExpression = &BreakExpression{}
// var _ PrimaryExpression = &ContinueExpression{}
// var _ PrimaryExpression = &Range{}
