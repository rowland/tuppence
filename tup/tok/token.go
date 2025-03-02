package tok

// reservedWords maps identifier strings to reserved token types.
var reservedWords = map[string]TokenType{
	"array":        TokenKeywordArray,
	"break":        TokenKeywordBreak,
	"continue":     TokenKeywordContinue,
	"contract":     TokenKeywordContract,
	"else":         TokenKeywordElse,
	"enum":         TokenKeywordEnum,
	"error":        TokenKeywordError,
	"false":        TokenBooleanLiteral,
	"fn":           TokenKeywordFn,
	"fx":           TokenKeywordFx,
	"for":          TokenKeywordFor,
	"if":           TokenKeywordIf,
	"in":           TokenKeywordIn,
	"it":           TokenKeywordIt,
	"import":       TokenKeywordImport,
	"mut":          TokenKeywordMut,
	"return":       TokenKeywordReturn,
	"switch":       TokenKeywordSwitch,
	"true":         TokenBooleanLiteral,
	"try":          TokenKeywordTry,
	"try_break":    TokenKeywordTryBreak,
	"try_continue": TokenKeywordTryContinue,
	"type":         TokenKeywordType,
	"typeof":       TokenKeywordTypeof,
	"union":        TokenKeywordUnion,
}

// GetReserved returns the reserved token type for a given identifier,
// if it exists.
func GetReserved(word string) (TokenType, bool) {
	t, ok := reservedWords[word]
	return t, ok
}

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Invalid bool
	Value   string
	File    *Source
	Offset  int
}

func (t *Token) Line() int {
	return t.File.Line(t.Offset)
}

func (t *Token) Column() int {
	return t.File.Column(t.Offset)
}

func (t *Token) Position() (int, int) {
	return t.File.Position(t.Offset)
}
