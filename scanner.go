package mathval

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

// Scanner is a lexical scanner
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new Scanner instance
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the Reader
// Returns rune(0) if an error (or io.EOF) occurs
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// Scan returns the next token and its literal value
func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return WS, s.scanContiguous(unicode.White_Space)
	} else if isLetter(ch) {
		s.unread()
		return s.scanKeyword()
	} else if isDigit(ch) {
		s.unread()
		return s.scanDigits()
	}

	// Single-character token
	switch ch {
	case eof:
		return EOF, ""
	// Operators
	case '+':
		return PLUS, string(ch)
	case '-':
		return MINUS, string(ch)
	case '*':
		return MULTIPLY, string(ch)
	case '/':
		return DIVIDE, string(ch)
	case '\\':
		return INT_DIVIDE, string(ch)
	case '^':
		return POW, string(ch)
	case '%':
		return MODULO, string(ch)
	// Misc characters
	case '(':
		return LPAREN, string(ch)
	case ')':
		return RPAREN, string(ch)
	case '.':
		return DOT, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanKeyword consumes all contiguous character runes and checks whether they are a known keyword
func (s *Scanner) scanKeyword() (Token, string) {
	keyword := s.scanContiguous(unicode.Letter)

	// TODO iterate over a defined list of functions
	/*switch keyword {
	case "func":
		return FUNCTION, keyword
	}*/

	// Otherwise return as a regular identifier.
	return UNKNOWN_KEYWORD, keyword
}

// scanDigits consumes all contiguous decimal digit runes
func (s *Scanner) scanDigits() (Token, string) {
	keyword := s.scanContiguous(unicode.Digit)
	return DIGITS, keyword
}

// scanContiguous consumes all contiguous runes from the current rune to the first that isn't in
// the given unicode.RangeTable
func (s *Scanner) scanContiguous(table *unicode.RangeTable) string {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.Is(table, ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}
