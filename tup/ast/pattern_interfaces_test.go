package ast

var _ MatchCondition = &Constant{}
var _ MatchCondition = &Range{}
var _ MatchCondition = &InferredErrorType{}
var _ MatchCondition = &TypeReference{}
var _ MatchCondition = &ListMatch{}
var _ MatchCondition = &WildcardPattern{}
var _ MatchCondition = &TypedPattern{}
var _ MatchCondition = &LabeledPattern{}
var _ MatchCondition = &TuplePattern{}
var _ MatchCondition = &ArrayPattern{}

var _ Pattern = &Constant{}
var _ Pattern = &Range{}
var _ Pattern = &InferredErrorType{}
var _ Pattern = &TypeReference{}
var _ Pattern = &WildcardPattern{}
var _ Pattern = &TypedPattern{}
var _ Pattern = &LabeledPattern{}
var _ Pattern = &TuplePattern{}
var _ Pattern = &ArrayPattern{}

var _ MatchElement = &Constant{}
var _ MatchElement = &Range{}
var _ MatchElement = &InferredErrorType{}
var _ MatchElement = &TypeReference{}
