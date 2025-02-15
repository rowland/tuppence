package main

// reservedWords maps identifier strings to reserved token types.
var reservedWords = map[string]TokenType{
	"break":  TokenKeywordBreak,
	"else":   TokenKeywordElse,
	"enum":   TokenKeywordEnum,
	"error":  TokenKeywordError,
	"false":  TokenBooleanLiteral,
	"fn":     TokenKeywordFn,
	"fx":     TokenKeywordFx,
	"for":    TokenKeywordFor,
	"if":     TokenKeywordIf,
	"import": TokenKeywordImport,
	"mut":    TokenKeywordMut,
	"return": TokenKeywordReturn,
	"switch": TokenKeywordSwitch,
	"true":   TokenBooleanLiteral,
	"try":    TokenKeywordTry,
	"type":   TokenKeywordType,
	"typeof": TokenKeywordTypeof,
}

// GetReserved returns the reserved token type for a given identifier,
// if it exists.
func GetReserved(word string) (TokenType, bool) {
	t, ok := reservedWords[word]
	return t, ok
}

// Token represents a lexical token.
type Token struct {
	Type     TokenType
	Invalid  bool
	Line     int
	Column   int
	Value    string
	Filename string
}
