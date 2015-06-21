package mathval

import (
	"unicode"
)

var eof = rune(0)

// isWhitespace returns true if the rune is a unicode whitespace character. Just a wrapper around unicode.IsSpace
func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

// isLetter returns true if the rune is a character. Just a wrapper around unicode.IsLetter
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isNumber returns true if the rune is a decimal digit. Just a wrapper around unicode.IsDigit
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
