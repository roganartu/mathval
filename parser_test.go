package mathval

import (
	"fmt"
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

type expressionChecker struct {
	info *CheckerInfo
}
type opChecker expressionChecker
type factorChecker expressionChecker
type powerChecker expressionChecker
type termChecker expressionChecker
type numberChecker expressionChecker

var ExpressionEquals Checker = &expressionChecker{
	&CheckerInfo{Name: "ExpressionEquals", Params: []string{"obtained", "expected"}},
}

var OperatorEquals Checker = &opChecker{
	&CheckerInfo{Name: "Operator", Params: []string{"obtained", "expected"}},
}

func (p *ParserSuite) TestParse(c *C) {
	parser = NewParser(strings.NewReader("10^2*4+1"))
	expected := &Expression{
		factor: &Factor{
			power: &Power{
				term:  &Term{number: &Number{str: "10", val: big.NewRat(10, 1)}},
				op:    &ExponentOp{op: POW},
				power: &Power{term: &Term{number: &Number{str: "2", val: big.NewRat(2, 1)}}},
			},
			op:     &MultiplyOp{op: MULTIPLY},
			factor: &Factor{power: &Power{term: &Term{number: &Number{str: "4", val: big.NewRat(4, 1)}}}},
		},
		op:         &AddOp{op: PLUS},
		expression: &Expression{factor: &Factor{power: &Power{term: &Term{number: &Number{str: "1", val: big.NewRat(1, 1)}}}}},
	}
	exp, err := parser.Parse()
	c.Assert(err, IsNil)
	c.Assert(exp, ExpressionEquals, expected)
}

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

// Checker methods
func (e *expressionChecker) Check(params []interface{}, names []string) (result bool, err string) {
	var obtained, expected *Expression
	var ok bool
	if obtained, ok = params[0].(*Expression); !ok {
		return false, "Not an Expression"
	}
	if expected, ok = params[1].(*Expression); !ok {
		return false, "Not an Expression"
	}

	// Nil checks
	if ((obtained == nil || expected == nil) && obtained != expected) ||
		((obtained.expression == nil || expected.expression == nil) && obtained.expression != expected.expression) ||
		((obtained.factor == nil || expected.factor == nil) && obtained.factor != expected.factor) ||
		((obtained.op == nil || expected.op == nil) && obtained.op != expected.op) {
		return false, "Expressions not equal"
	}

	if obtained.expression != nil {
		if result, err = e.Check([]interface{}{obtained.expression, expected.expression}, names); !result {
			return
		}
	}

	if obtained.factor != nil {
		f := &factorChecker{}
		if result, err = f.Check([]interface{}{obtained.factor, expected.factor}, names); !result {
			return
		}
	}

	if obtained.op != nil {
		o := &opChecker{}
		if result, err = o.Check([]interface{}{obtained.op, expected.op}, names); !result {
			return
		}
	}

	return true, ""
}
func (e *expressionChecker) Info() *CheckerInfo {
	return e.info
}

func (f *factorChecker) Check(params []interface{}, names []string) (result bool, err string) {
	var obtained, expected *Factor
	var ok bool
	if obtained, ok = params[0].(*Factor); !ok {
		return false, "Not a Factor"
	}
	if expected, ok = params[1].(*Factor); !ok {
		return false, "Not a Factor"
	}

	// Nil checks
	if ((obtained == nil || expected == nil) && obtained != expected) ||
		((obtained.power == nil || expected.power == nil) && obtained.power != expected.power) ||
		((obtained.factor == nil || expected.factor == nil) && obtained.factor != expected.factor) ||
		((obtained.op == nil || expected.op == nil) && obtained.op != expected.op) {
		return false, "Factors not equal"
	}

	if obtained.factor != nil {
		if result, err = f.Check([]interface{}{obtained.factor, expected.factor}, names); !result {
			return
		}
	}

	if obtained.power != nil {
		p := &powerChecker{}
		if result, err = p.Check([]interface{}{obtained.power, expected.power}, names); !result {
			return
		}
	}

	if obtained.op != nil {
		o := &opChecker{}
		if result, err = o.Check([]interface{}{obtained.op, expected.op}, names); !result {
			return
		}
	}

	return true, ""
}
func (f *factorChecker) Info() *CheckerInfo {
	return f.info
}

func (p *powerChecker) Check(params []interface{}, names []string) (result bool, err string) {
	var obtained, expected *Power
	var ok bool
	if obtained, ok = params[0].(*Power); !ok {
		return false, "Not a Power"
	}
	if expected, ok = params[1].(*Power); !ok {
		return false, "Not a Power"
	}

	// Nil checks
	if ((obtained == nil || expected == nil) && obtained != expected) ||
		((obtained.term == nil || expected.term == nil) && obtained.term != expected.term) ||
		((obtained.power == nil || expected.power == nil) && obtained.power != expected.power) ||
		((obtained.op == nil || expected.op == nil) && obtained.op != expected.op) {
		return false, "Powers not equal"
	}

	if obtained.power != nil {
		if result, err = p.Check([]interface{}{obtained.power, expected.power}, names); !result {
			return
		}
	}

	if obtained.term != nil {
		t := &termChecker{}
		if result, err = t.Check([]interface{}{obtained.term, expected.term}, names); !result {
			return
		}
	}

	if obtained.op != nil {
		o := &opChecker{}
		if result, err = o.Check([]interface{}{obtained.op, expected.op}, names); !result {
			return
		}
	}

	return true, ""
}
func (p *powerChecker) Info() *CheckerInfo {
	return p.info
}

func (t *termChecker) Check(params []interface{}, names []string) (result bool, err string) {
	var obtained, expected *Term
	var ok bool
	if obtained, ok = params[0].(*Term); !ok {
		return false, "Not a Term"
	}
	if expected, ok = params[1].(*Term); !ok {
		return false, "Not a Term"
	}

	// Nil checks
	if ((obtained == nil || expected == nil) && obtained != expected) ||
		((obtained.exp == nil || expected.exp == nil) && obtained.exp != expected.exp) ||
		((obtained.number == nil || expected.number == nil) && obtained.number != expected.number) {
		return false, "Terms not equal"
	}

	if obtained.exp != nil {
		e := &expressionChecker{}
		if result, err = e.Check([]interface{}{obtained.exp, expected.exp}, names); !result {
			return
		}
	}

	if obtained.number != nil {
		n := &numberChecker{}
		if result, err = n.Check([]interface{}{obtained.number, expected.number}, names); !result {
			return
		}
	}

	return true, ""
}
func (t *termChecker) Info() *CheckerInfo {
	return t.info
}

func (n *numberChecker) Check(params []interface{}, names []string) (result bool, err string) {
	var obtained, expected *Number
	var ok bool
	if obtained, ok = params[0].(*Number); !ok {
		return false, "Not a Number"
	}
	if expected, ok = params[1].(*Number); !ok {
		return false, "Not a Number"
	}

	// Nil checks
	if ((obtained == nil || expected == nil) && obtained != expected) ||
		((obtained.val == nil || expected.val == nil) && obtained.val != expected.val) {
		return false, "Numbers not equal"
	}

	if obtained.str != expected.str || obtained.val.Cmp(expected.val) != 0 {
		return false, "Numbers not equal"
	}

	return true, ""
}
func (n *numberChecker) Info() *CheckerInfo {
	return n.info
}

func (o *opChecker) Check(params []interface{}, names []string) (result bool, err string) {
	var obtained, expected *Operator
	if a, ok := params[0].(*AddOp); ok {
		obtained = (*Operator)(a)
		expected = (*Operator)(params[1].(*AddOp))
	} else if a, ok := params[0].(*MultiplyOp); ok {
		obtained = (*Operator)(a)
		expected = (*Operator)(params[1].(*MultiplyOp))
	} else if a, ok := params[0].(*ExponentOp); ok {
		obtained = (*Operator)(a)
		expected = (*Operator)(params[1].(*ExponentOp))
	} else {
		return false, "Not an Operator"
	}

	// Nil checks
	if ((obtained == nil || expected == nil) && obtained != expected) ||
		(obtained.op != expected.op) {
		return false, fmt.Sprintf("Expected: %d\nObtained: %d", expected.op, obtained.op)
	}

	return true, ""
}
func (o *opChecker) Info() *CheckerInfo {
	return o.info
}
