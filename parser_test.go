package mathval

import (
	"math/big"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

var (
	parser *Parser
)

type ParseResult struct {
	token   Token
	literal string
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ParserSuite struct{}

var _ = Suite(&ParserSuite{})

func (t *ParserSuite) TestNewParser(c *C) {
	parser = NewParser(strings.NewReader("test"))
	c.Assert(parser, NotNil)
	c.Assert(parser.s, NotNil)
}

func (p *ParserSuite) TestScanUnscan(c *C) {
	scan_str := "(49 + 77)*((14-2)/11)\\2"
	parser = NewParser(strings.NewReader(scan_str))

	token, literal := parser.scan()
	c.Assert(token, Equals, LPAREN)
	c.Assert(literal, Equals, "(")

	token, literal = parser.scan()
	c.Assert(token, Equals, DIGITS)
	c.Assert(literal, Equals, "49")

	parser.unscan()
	token, literal = parser.scan()
	c.Assert(token, Equals, DIGITS)
	c.Assert(literal, Equals, "49")
}

func (p *ParserSuite) TestScanIgnoreWhitespace(c *C) {
	scan_str := " 10 + 12 "
	parser = NewParser(strings.NewReader(scan_str))

	expected := []ParseResult{
		{token: DIGITS, literal: "10"},
		{token: PLUS, literal: "+"},
		{token: DIGITS, literal: "12"},
		{token: EOF, literal: ""},
	}

	for _, res := range expected {
		token, literal := parser.scanIgnoreWhitespace()
		c.Assert(token, Equals, res.token)
		c.Assert(literal, Equals, res.literal)
	}
}

func (t *ParserSuite) TestPeek(c *C) {
	parser = NewParser(strings.NewReader("1+2"))

	tok, lit := parser.peek()
	c.Assert(tok, Equals, DIGITS)
	c.Assert(lit, Equals, "1")

	tok, lit = parser.scan()
	c.Assert(tok, Equals, DIGITS)
	c.Assert(lit, Equals, "1")

	tok, lit = parser.peek()
	c.Assert(tok, Equals, PLUS)
	c.Assert(lit, Equals, "+")
}

func (p *ParserSuite) TestParseExponentOp(c *C) {
	// Assert EOF checking
	parser = NewParser(strings.NewReader(""))
	exp, err := parser.parseExponentOp()
	c.Assert(exp, IsNil)
	c.Assert(err, NotNil)

	// Normal inputs
	parser = NewParser(strings.NewReader("^"))
	exp, err = parser.parseExponentOp()
	c.Assert(exp, NotNil)
	c.Assert(exp.op, Equals, POW)
	c.Assert(err, IsNil)

	// Non-ExponentOp
	parser = NewParser(strings.NewReader(" "))
	_, err = parser.parseExponentOp()
	c.Assert(err, NotNil)
}

func (p *ParserSuite) TestParseMultiplyOp(c *C) {
	// Assert EOF checking
	parser = NewParser(strings.NewReader(""))
	mul, err := parser.parseMultiplyOp()
	c.Assert(mul, IsNil)
	c.Assert(err, NotNil)

	// Normal inputs
	expected := []MultiplyOp{
		{op: MULTIPLY},
		{op: DIVIDE},
		{op: MODULO},
		{op: INT_DIVIDE},
	}
	parser = NewParser(strings.NewReader("*/%\\"))
	for _, res := range expected {
		mul, err = parser.parseMultiplyOp()
		c.Assert(mul, NotNil)
		c.Assert(mul.op, Equals, res.op)
		c.Assert(err, IsNil)
	}

	// Non-MultiplyOp
	parser = NewParser(strings.NewReader(" "))
	_, err = parser.parseMultiplyOp()
	c.Assert(err, NotNil)
}

func (p *ParserSuite) TestParseAddOp(c *C) {
	// Assert EOF checking
	parser = NewParser(strings.NewReader(""))
	add, err := parser.parseAddOp()
	c.Assert(add, IsNil)
	c.Assert(err, NotNil)

	// Normal inputs
	expected := []AddOp{
		{op: PLUS},
		{op: MINUS},
	}
	parser = NewParser(strings.NewReader("+-"))
	for _, res := range expected {
		add, err = parser.parseAddOp()
		c.Assert(add, NotNil)
		c.Assert(add.op, Equals, res.op)
		c.Assert(err, IsNil)
	}

	// Non-AddOp
	parser = NewParser(strings.NewReader(" "))
	_, err = parser.parseAddOp()
	c.Assert(err, NotNil)
}

func (p *ParserSuite) TestParseNumber(c *C) {
	// Assert EOF checking
	parser = NewParser(strings.NewReader(""))
	num, err := parser.parseNumber()
	c.Assert(num, IsNil)
	c.Assert(err, NotNil)

	// Normal inputs
	expected := []Number{
		{str: "10", val: big.NewRat(10, 1)},
		{str: "0.5", val: big.NewRat(1, 2)},
		{str: "2.5", val: big.NewRat(5, 2)},
	}
	for _, res := range expected {
		parser = NewParser(strings.NewReader(res.str))
		num, err = parser.parseNumber()
		c.Assert(num, NotNil)
		c.Assert(num.str, Equals, res.str)
		c.Assert(res.val.Cmp(num.val), Equals, 0)
		c.Assert(err, IsNil)
	}

	// Non-Digit
	parser = NewParser(strings.NewReader(" "))
	_, err = parser.parseNumber()
	c.Assert(err, NotNil)
}

func (p *ParserSuite) TestParseTerm(c *C) {
	// Assert EOF checking
	parser = NewParser(strings.NewReader(""))
	term, err := parser.parseTerm()
	c.Assert(term, IsNil)
	c.Assert(err, NotNil)

	// Normal inputs
	// "10"
	parser = NewParser(strings.NewReader("10"))
	term, err = parser.parseTerm()
	c.Assert(term.exp, IsNil)
	c.Assert(term.number.str, Equals, "10")
	// "(10)"
	parser = NewParser(strings.NewReader("(10)"))
	term, err = parser.parseTerm()
	c.Assert(term.exp.factor.power.term.number.str, Equals, "10")
	c.Assert(term.number, IsNil)

	// No RPAREN
	parser = NewParser(strings.NewReader("(10"))
	_, err = parser.parseTerm()
	c.Assert(err, NotNil)
}

func (p *ParserSuite) TestParsePower(c *C) {
	// Assert EOF checking
	parser = NewParser(strings.NewReader(""))
	pow, err := parser.parsePower()
	c.Assert(pow, IsNil)
	c.Assert(err, NotNil)

	// Normal inputs
	// "10"
	parser = NewParser(strings.NewReader("10"))
	pow, err = parser.parsePower()
	c.Assert(pow.op, IsNil)
	c.Assert(pow.power, IsNil)
	c.Assert(pow.term.number.str, Equals, "10")
	// "10^2"
	parser = NewParser(strings.NewReader("10^2"))
	pow, err = parser.parsePower()
	c.Assert(pow.term.number.str, Equals, "10")
	c.Assert(pow.op.op, Equals, POW)
	c.Assert(pow.power.term.number.str, Equals, "2")

	// No POWER after POW
	parser = NewParser(strings.NewReader("10^"))
	_, err = parser.parsePower()
	c.Assert(err, NotNil)
}

func (p *ParserSuite) TestParseFactor(c *C) {
	// Assert EOF checking
	parser = NewParser(strings.NewReader(""))
	fac, err := parser.parseFactor()
	c.Assert(fac, IsNil)
	c.Assert(err, NotNil)

	// Normal inputs
	// "10"
	parser = NewParser(strings.NewReader("10"))
	fac, err = parser.parseFactor()
	c.Assert(fac.op, IsNil)
	c.Assert(fac.factor, IsNil)
	c.Assert(fac.power.term.number.str, Equals, "10")
	// "10*2"
	parser = NewParser(strings.NewReader("10*2"))
	fac, err = parser.parseFactor()
	c.Assert(fac.power.term.number.str, Equals, "10")
	c.Assert(fac.op.op, Equals, MULTIPLY)
	c.Assert(fac.factor.power.term.number.str, Equals, "2")

	// No POWER after POW
	parser = NewParser(strings.NewReader("10/"))
	_, err = parser.parseFactor()
	c.Assert(err, NotNil)
}

func (p *ParserSuite) TestParseExpression(c *C) {
	// Assert EOF checking
	parser = NewParser(strings.NewReader(""))
	exp, err := parser.parseExpression()
	c.Assert(exp, IsNil)
	c.Assert(err, NotNil)

	// Normal inputs
	// "10"
	parser = NewParser(strings.NewReader("10"))
	exp, err = parser.parseExpression()
	c.Assert(exp.op, IsNil)
	c.Assert(exp.expression, IsNil)
	c.Assert(exp.factor.power.term.number.str, Equals, "10")
	// "10+2"
	parser = NewParser(strings.NewReader("10+2"))
	exp, err = parser.parseExpression()
	c.Assert(exp.factor.power.term.number.str, Equals, "10")
	c.Assert(exp.op.op, Equals, PLUS)
	c.Assert(exp.expression.factor.power.term.number.str, Equals, "2")

	// No POWER after POW
	parser = NewParser(strings.NewReader("10+"))
	_, err = parser.parseExpression()
	c.Assert(err, NotNil)
}
