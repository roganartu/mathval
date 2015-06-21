package mathval

import (
	"math/big"
)

/*  The parseable language described using EBNF. Order of Operations is maintained by expanding expressions
	through the highest precendence first (ie: exponentiation then multiplication then addition)

EXPRESSION  = FACTOR | FACTOR ADD_OP EXPRESSION ;
FACTOR      = POWER | POWER MULTIPLY_OP FACTOR ;
POWER       = TERM | TERM EXPONENT_OP POWER ;
TERM        = '(' EXPRESSION ')' | NUMBER ;
NUMBER      = { DIGIT } | { DIGIT } '.' { DIGIT }
DIGIT       = '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'
ADD_OP      = '+' | '-' ;
MULTIPLY_OP = '*' | '/' | '%' ;
EXPONENT_OP = '^' ;

*/

// Expression represents an EXPRESSION in the EBNF grammar
// EXPRESSION = POWER | POWER ADD_OP EXPRESSION
type Expression struct {
	factor     *Factor
	op         *AddOp
	expression *Expression
}

// Factor repsents a FACTOR in the EBNF grammar
// FACTOR = POWER | POWER MULTIPLY_OP FACTOR
type Factor struct {
	power  *Power
	op     *MultiplyOp
	factor *Factor
}

// Power represents a POWER in the EBNF grammar
// POWER = TERM | TERM EXPONENT_OP POWER
type Power struct {
	term  *Term
	op    *ExponentOp
	power *Power
}

// Term represents a TERM in the EBNF grammar
// TERM = '(' EXPRESSION ')' | NUMBER
type Term struct {
	exp    *Expression
	number *Number
}

type Number struct {
	str string
	val *big.Rat
}

// Operator represents the different groups of operators in the EBNF grammar
type Operator struct {
	op Token
}

// AddOp represents an ADD_OP in the EBNF grammar
// ADD_OP = '+' | '-'
type AddOp Operator

// MultiplyOp represents a MULTIPLY_OP in the EBNF grammar
// MULTIPLY_OP = '*' | '/' | '%'
type MultiplyOp Operator

// ExponentOp represents an ExponentOp in the EBNF grammar
// EXPONENT_OP = '^'
type ExponentOp Operator
