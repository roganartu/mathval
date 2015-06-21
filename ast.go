package mathval

/*  The parseable language described using EBNF. Order of Operations is maintained by expanding expressions
	through the highest precendence first (ie: exponentiation then multiplication then addition)

EXPRESSION  = POWER | POWER ADD_OP EXPRESSION ;
POWER       = TERM | TERM EXPONENT_OP POWER ;
TERM        = FACTOR | FACTOR MULTIPLY_OP TERM ;
FACTOR      = '(' EXPRESSION ')' | NUMBER ;
ADD_OP      = '+' | '-' ;
MULTIPLY_OP = '*' | '/' | '%' ;
EXPONENT_OP = '^' ;

*/

// Expression represents an EXPRESSION in the EBNF grammar
// EXPRESSION = POWER | POWER ADD_OP EXPRESSION
type Expression struct {
	power      Power
	op         AddOp
	expression Expression
}

// Power represents a POWER in the EBNF grammar
// POWER = TERM | TERM EXPONENT_OP POWER
type Power struct {
	term  Term
	op    ExponentOp
	power Power
}

// Term repsents a TERM in the EBNF grammar
// TERM = FACTOR | FACTOR MULTIPLY_OP TERM
type Term struct {
	factor Factor
	op     MultiplyOp
	term   Term
}

// Factor represents a FACTOR in the EBNF grammar
// FACTOR = '(' EXPRESSION ')' | NUMBER
type Factor struct {
	exp    Expression
	number Number
}

// AddOp represents an ADD_OP in the EBNF grammar
// ADD_OP = '+' | '-'
type AddOp struct {
	op Token
}

// ExponentOp represents an ExponentOp in the EBNF grammar
// EXPONENT_OP = '^'
type ExponentOp struct {
	op Token
}

// MultiplyOp represents a MULTIPLY_OP in the EBNF grammar
// MULTIPLY_OP = '*' | '/' | '%'
type MultiplyOp struct {
	op Token
}
