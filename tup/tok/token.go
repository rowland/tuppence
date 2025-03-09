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
	Value       string
	File        *Source
	Offset      int32
	ErrorOffset int32 // Position where error was detected within the token, 0 if no error
	Type        TokenType
	Invalid     bool
}

func (t *Token) Line() int {
	return t.File.Line(int(t.Offset))
}

func (t *Token) Column() int {
	return t.File.Column(int(t.Offset))
}

func (t *Token) Position() (int, int) {
	return t.File.Position(int(t.Offset))
}

// ErrorPosition returns the line and column where the error occurred within the token.
// If the token is not invalid, it returns the same as Position().
func (t *Token) ErrorPosition() (int, int) {
	if !t.Invalid {
		return t.Position()
	}
	return t.File.Position(int(t.ErrorOffset))
}

// ErrorLine returns the line where the error occurred within the token.
// If the token is not invalid, it returns the same as Line().
func (t *Token) ErrorLine() int {
	if !t.Invalid {
		return t.Line()
	}
	return t.File.Line(int(t.ErrorOffset))
}

// ErrorColumn returns the column where the error occurred within the token.
// If the token is not invalid, it returns the same as Column().
func (t *Token) ErrorColumn() int {
	if !t.Invalid {
		return t.Column()
	}
	return t.File.Column(int(t.ErrorOffset))
}
