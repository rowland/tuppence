package ast

var _ Expression = &FloatLiteral{}
var _ Expression = &IntegerLiteral{}
var _ Expression = &BooleanLiteral{}
var _ Expression = &StringLiteral{}
var _ Expression = &InterpolatedStringLiteral{}
var _ Expression = &RawStringLiteral{}
var _ Expression = &MultiLineStringLiteral{}
var _ Expression = &TupleLiteral{}
var _ Expression = &ArrayLiteral{}
var _ Expression = &SymbolLiteral{}
var _ Expression = &RuneLiteral{}
var _ Expression = &FixedSizeArrayLiteral{}
