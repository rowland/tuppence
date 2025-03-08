package tok

// reservedWords maps identifier strings to reserved token types.
var reservedWords = map[string]TokenType{
	"array":        TokKwArray,
	"break":        TokKwBreak,
	"continue":     TokKwContinue,
	"contract":     TokKwContract,
	"else":         TokKwElse,
	"enum":         TokKwEnum,
	"error":        TokKwError,
	"false":        TokBoolLit,
	"fn":           TokKwFn,
	"fx":           TokKwFx,
	"for":          TokKwFor,
	"if":           TokKwIf,
	"in":           TokKwIn,
	"it":           TokKwIt,
	"import":       TokKwImport,
	"mut":          TokKwMut,
	"return":       TokKwReturn,
	"switch":       TokKwSwitch,
	"true":         TokBoolLit,
	"try":          TokKwTry,
	"try_break":    TokKwTryBreak,
	"try_continue": TokKwTryContinue,
	"type":         TokKwType,
	"typeof":       TokKwTypeof,
	"union":        TokKwUnion,
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
