package tok

// reservedWords maps identifier strings to reserved token types.
var reservedWords = map[string]TokenType{
	"array":        TokKeywordArray,
	"break":        TokKeywordBreak,
	"continue":     TokKeywordContinue,
	"contract":     TokKeywordContract,
	"else":         TokKeywordElse,
	"enum":         TokKeywordEnum,
	"error":        TokKeywordError,
	"false":        TokBoolLit,
	"fn":           TokKeywordFn,
	"fx":           TokKeywordFx,
	"for":          TokKeywordFor,
	"if":           TokKeywordIf,
	"in":           TokKeywordIn,
	"it":           TokKeywordIt,
	"import":       TokKeywordImport,
	"mut":          TokKeywordMut,
	"return":       TokKeywordReturn,
	"switch":       TokKeywordSwitch,
	"true":         TokBoolLit,
	"try":          TokKeywordTry,
	"try_break":    TokKeywordTryBreak,
	"try_continue": TokKeywordTryContinue,
	"type":         TokKeywordType,
	"typeof":       TokKeywordTypeof,
	"union":        TokKeywordUnion,
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
