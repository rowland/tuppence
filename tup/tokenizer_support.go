package main

// isHexDigit returns true if c is a valid hexadecimal digit
func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'f') ||
		(c >= 'A' && c <= 'F')
}

// isSimpleEscape returns true if c is a valid single-character escape sequence
func isSimpleEscape(c byte) bool {
	return c == 'n' || c == 't' || c == '"' || c == '\'' || c == '\\' ||
		c == 'r' || c == 'b' || c == 'f' || c == 'v' || c == '0'
}

// isOctalDigit returns true if c is a valid octal digit
func isOctalDigit(c byte) bool {
	return c >= '0' && c <= '7'
}

// isDecimalDigit returns true if c is a valid decimal digit
func isDecimalDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isLetter returns true if c is a letter (A-Z or a-z)
func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// isIdentifierStart returns true if c is a letter or underscore
func isIdentifierStart(c byte) bool {
	return isLetter(c) || c == '_'
}

// isInvalidNumberLetter returns true if c is a letter that would make a number invalid
// This excludes 'b', 'o', and 'x' which are handled separately as valid number prefixes
func isInvalidNumberLetter(c byte) bool {
	return isLetter(c) && c != 'b' && c != 'o' && c != 'x'
}

// isInvalidIntegerLetter returns true if c is a letter that would make an integer invalid
// This excludes 'e' which is handled separately as an exponent marker
func isInvalidIntegerLetter(c byte) bool {
	return isLetter(c) && c != 'e'
}

// isInvalidExponentSignChar returns true if c is an invalid character after an exponent sign
// This includes letters, underscore, and signs (+ or -) since we've already handled the sign
func isInvalidExponentSignChar(c byte) bool {
	return isLetter(c) || c == '_' || c == '+' || c == '-'
}

// isInvalidExponentIntChar returns true if c is an invalid character in an exponent's integer part
// This includes letters and underscore since only digits are allowed
func isInvalidExponentIntChar(c byte) bool {
	return isLetter(c) || c == '_'
}

// isInvalidBinaryFirstChar returns true if c is an invalid character for the first position of a binary number
// This includes any character that is not 0 or 1
func isInvalidBinaryFirstChar(c byte) bool {
	return (c >= '2' && c <= '9') || isLetter(c) || c == '_' || c == '.'
}

// isInvalidBinaryChar returns true if c is an invalid character for a binary number
// This includes any character that is not 0, 1, or underscore
func isInvalidBinaryChar(c byte) bool {
	return (c >= '2' && c <= '9') || isLetter(c)
}

// isInvalidOctalChar returns true if c is an invalid character for an octal number
// This includes any character that is not an octal digit (0-7) or underscore
func isInvalidOctalChar(c byte) bool {
	return (c >= '8' && c <= '9') || isLetter(c)
}

// isInvalidOctalFirstChar returns true if c is an invalid character for the first position of an octal number
// This includes any character that is not an octal digit (0-7)
func isInvalidOctalFirstChar(c byte) bool {
	return (c >= '8' && c <= '9') || isLetter(c) || c == '_' || c == '.'
}

// isInvalidHexadecimalChar returns true if c is a letter that would make a hexadecimal number invalid
// This includes letters G-Z and g-z, since A-F and a-f are valid hex digits
func isInvalidHexadecimalChar(c byte) bool {
	return (c >= 'G' && c <= 'Z') || (c >= 'g' && c <= 'z')
}
