package mathval

import (
	"strings"
	"testing"
	"unicode"

	. "gopkg.in/check.v1"
)

var (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
	scanner  *Scanner
)

type ScanResult struct {
	token   Token
	literal string
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ScannerSuite struct{}

var _ = Suite(&ScannerSuite{})

func (s *ScannerSuite) SetUpTest(c *C) {
	scanner = NewScanner(strings.NewReader(alphabet))
}

func (t *ScannerSuite) TestNewScanner(c *C) {
	c.Assert(scanner, NotNil)
}

func (t *ScannerSuite) TestReadUnread(c *C) {
	c.Assert(scanner.read(), Equals, rune(alphabet[0]))
	c.Assert(scanner.read(), Equals, rune(alphabet[1]))
	scanner.unread()
	c.Assert(scanner.read(), Equals, rune(alphabet[1]))
}

func (s *ScannerSuite) TestScanKeyword(c *C) {
	token, keyword := scanner.scanKeyword()
	c.Assert(token, Equals, UNKNOWN_KEYWORD)
	c.Assert(keyword, Equals, alphabet)
}

func (s *ScannerSuite) TestScanDigits(c *C) {
	digitStr := "123 456*789"
	scanner = NewScanner(strings.NewReader(digitStr))

	token, digits := scanner.scanDigits()
	c.Assert(token, Equals, DIGITS)
	c.Assert(digits, Equals, "123")
	scanner.read()

	token, digits = scanner.scanDigits()
	c.Assert(token, Equals, DIGITS)
	c.Assert(digits, Equals, "456")
	scanner.read()

	token, digits = scanner.scanDigits()
	c.Assert(token, Equals, DIGITS)
	c.Assert(digits, Equals, "789")
}

func (s *ScannerSuite) TestScanContiguous(c *C) {
	str := "aaa b  \t\ncc"
	scanner = NewScanner(strings.NewReader(str))
	c.Assert(scanner.scanContiguous(unicode.Letter), Equals, "aaa")
	c.Assert(scanner.scanContiguous(unicode.White_Space), Equals, " ")
	c.Assert(scanner.scanContiguous(unicode.Letter), Equals, "b")
	c.Assert(scanner.scanContiguous(unicode.White_Space), Equals, "  \t\n")
	c.Assert(scanner.scanContiguous(unicode.Letter), Equals, "cc")
}

func (s *ScannerSuite) TestScanEOF(c *C) {
	scanner = NewScanner(strings.NewReader(""))
	token, literal := scanner.Scan()
	c.Assert(token, Equals, EOF)
	c.Assert(literal, Equals, "")
}

func (s *ScannerSuite) TestScan(c *C) {
	scan_str := "(49 + 77)*((14-2)/11)\\2"
	scanner = NewScanner(strings.NewReader(scan_str))

	expected := []ScanResult{
		{token: LPAREN, literal: "("},
		{token: DIGITS, literal: "49"},
		{token: WS, literal: " "},
		{token: PLUS, literal: "+"},
		{token: WS, literal: " "},
		{token: DIGITS, literal: "77"},
		{token: RPAREN, literal: ")"},
		{token: MULTIPLY, literal: "*"},
		{token: LPAREN, literal: "("},
		{token: LPAREN, literal: "("},
		{token: DIGITS, literal: "14"},
		{token: MINUS, literal: "-"},
		{token: DIGITS, literal: "2"},
		{token: RPAREN, literal: ")"},
		{token: DIVIDE, literal: "/"},
		{token: DIGITS, literal: "11"},
		{token: RPAREN, literal: ")"},
		{token: INT_DIVIDE, literal: "\\"},
		{token: DIGITS, literal: "2"},
		{token: EOF, literal: ""},
	}

	for _, res := range expected {
		token, literal := scanner.Scan()
		c.Assert(token, Equals, res.token)
		c.Assert(literal, Equals, res.literal)
	}
}
