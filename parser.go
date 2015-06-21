package mathval

import (
	"errors"
	"io"
	"math/big"
)

// Parser is a parser including a Scanner and a buffer
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size. Currently max=1 as no lookahead
	}
}

// NewParser returns a new instance of Parser with the defined lookahead length
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// peek returns the next token in the scanner. Whitespace is ignored
func (p *Parser) peek() (Token, string) {
	if p.buf.n == 0 {
		p.buf.tok, p.buf.lit = p.s.Scan()
		p.buf.n = 1
	}

	// Ignore whitespace
	if p.buf.tok == WS {
		p.buf.tok, p.buf.lit = p.s.Scan()
	}

	return p.buf.tok, p.buf.lit
}

// scanIgnoreWhitespace scans the next non-whitespace token
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// Parse parses the output from the Scanner
func (p *Parser) Parse() (*Expression, error) {
	return p.parseExpression()
}

// parseExpression recursively parses an Expression starting at the next Token
func (p *Parser) parseExpression() (exp *Expression, err error) {
	if tok, _ := p.peek(); tok == EOF {
		return nil, errors.New("Unexpected EOF")
	}

	exp = &Expression{}
	exp.factor, err = p.parseFactor()
	if err != nil {
		return
	}

	// Check for an additive operator
	if tok, _ := p.peek(); tok >= additive_begin && tok <= additive_end {
		exp.op, err = p.parseAddOp()
		if err != nil {
			return
		}
		exp.expression, err = p.parseExpression()
	}
	return
}

// parseFactor recursively parses a Factor starting at the next Token
func (p *Parser) parseFactor() (fac *Factor, err error) {
	if tok, _ := p.peek(); tok == EOF {
		return nil, errors.New("Unexpected EOF")
	}

	fac = &Factor{}
	fac.power, err = p.parsePower()
	if err != nil {
		return
	}

	// Check for a multiplicative operator
	if tok, _ := p.peek(); tok >= multiplicative_begin && tok <= multiplicative_end {
		fac.op, err = p.parseMultiplyOp()
		if err != nil {
			return
		}
		fac.factor, err = p.parseFactor()
	}
	return
}

// parsePower recursively parses a Power starting at the next Token
func (p *Parser) parsePower() (pow *Power, err error) {
	if tok, _ := p.peek(); tok == EOF {
		return nil, errors.New("Unexpected EOF")
	}

	pow = &Power{}
	pow.term, err = p.parseTerm()
	if err != nil {
		return
	}

	// Check for an exponentiation operator
	if tok, _ := p.peek(); tok >= exponentiation_begin && tok <= exponentiation_end {
		pow.op, err = p.parseExponentOp()
		if err != nil {
			return
		}
		pow.power, err = p.parsePower()
	}
	return
}

// parseTerm recursively parses a Term starting at the next Token
func (p *Parser) parseTerm() (term *Term, err error) {
	if tok, _ := p.peek(); tok == EOF {
		return nil, errors.New("Unexpected EOF")
	}

	term = &Term{}

	// '(' EXPRESSION ')'
	if tok, _ := p.peek(); tok == LPAREN {
		p.scanIgnoreWhitespace()
		term.exp, err = p.parseExpression()
		if tok, _ = p.peek(); tok != RPAREN {
			return term, errors.New("Expected RPAREN")
		}
		p.scanIgnoreWhitespace()
	} else if tok == DIGITS {
		term.number, err = p.parseNumber()
	} else {
		// TODO add helper in tokens.go to convert token values to names for errors
		return nil, errors.New("Unexpected TOKEN")
	}
	return
}

// parseNumber parses the number represented by the next Token
func (p *Parser) parseNumber() (num *Number, err error) {
	if tok, _ := p.peek(); tok == EOF {
		return nil, errors.New("Unexpected EOF")
	}

	num = &Number{}
	if tok, _ := p.peek(); tok != DIGITS {
		return nil, errors.New("Expected decimal or floating digits")
	}

	_, integral := p.scanIgnoreWhitespace()
	fractional := ""
	if tok, _ := p.peek(); tok == DOT {
		p.scanIgnoreWhitespace()
		if tok, _ = p.peek(); tok != DIGITS {
			return nil, errors.New("Expected fractional digits")
		}
		_, fractional = p.scanIgnoreWhitespace()
	}

	num.str = integral
	if fractional != "" {
		num.str += "." + fractional
	}

	num.val = new(big.Rat)
	if _, ok := num.val.SetString(num.str); !ok {
		return num, errors.New("Error parsing value")
	}
	return
}

// parseAddOp recursively parses an AddOp starting at the next Token
func (p *Parser) parseAddOp() (add *AddOp, err error) {
	if tok, _ := p.peek(); tok == EOF {
		return nil, errors.New("Unexpected EOF")
	}

	add = &AddOp{}
	if tok, _ := p.peek(); tok < additive_begin || tok > additive_end {
		return nil, errors.New("Expected additive operator")
	}
	add.op, _ = p.scanIgnoreWhitespace()
	return
}

// parseMultiplyOp recursively parses a MultiplyOp starting at the next Token
func (p *Parser) parseMultiplyOp() (mul *MultiplyOp, err error) {
	if tok, _ := p.peek(); tok == EOF {
		return nil, errors.New("Unexpected EOF")
	}

	mul = &MultiplyOp{}
	if tok, _ := p.peek(); tok < multiplicative_begin || tok > multiplicative_end {
		return nil, errors.New("Expected multiplicative operator")
	}
	mul.op, _ = p.scanIgnoreWhitespace()
	return
}

// parseExponentOp recursively parses an ExponentOp starting at the next Token
func (p *Parser) parseExponentOp() (exp *ExponentOp, err error) {
	if tok, _ := p.peek(); tok == EOF {
		return nil, errors.New("Unexpected EOF")
	}

	exp = &ExponentOp{}
	if tok, _ := p.peek(); tok < exponentiation_begin || tok > exponentiation_end {
		return nil, errors.New("Expected exponentiation operator")
	}
	exp.op, _ = p.scanIgnoreWhitespace()
	return
}
