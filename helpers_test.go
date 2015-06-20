package mathval

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

// TODO these three tests can be simplified to a single one
func (t *TestSuite) TestIsWhitespace(c *C) {
	whitespace := []int32{'\t', '\n', '\v', '\f', '\r', ' ', 0x0085, 0x00A0}
	for _, ch := range whitespace {
		c.Assert(isWhitespace(rune(ch)), Equals, true)
	}
}

// TODO these three tests can be simplified to a single one
func (t *TestSuite) TestIsLetter(c *C) {
	alpha := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, ch := range alpha {
		c.Assert(isLetter(rune(ch)), Equals, true)
	}
}

// TODO these three tests can be simplified to a single one
func (t *TestSuite) TestIsDigit(c *C) {
	number := "0123456789"
	for _, ch := range number {
		c.Assert(isDigit(rune(ch)), Equals, true)
	}
}
