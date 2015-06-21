package mathval

type Token int

const (
	errors_begin Token = iota
	ILLEGAL
	UNKNOWN_KEYWORD
	errors_end

	// Special tokens
	special_tokens_begin
	EOF
	WS // Whitespace
	special_tokens_end

	// Arithmetic operators
	operators_begin

	additive_begin
	PLUS  // +
	MINUS // -
	additive_end

	multiplicative_begin
	MULTIPLY   // *
	DIVIDE     // /
	INT_DIVIDE // \
	MODULO     // %
	multiplicative_end

	exponentiation_begin
	POW // ^
	exponentiation_end

	operators_end

	// Known types/keywords
	keywords_begin
	DIGITS // Contiguous block of digits
	keywords_end

	// Misc characters
	misc_begin
	LPAREN // (
	RPAREN // )
	DOT    // .
	misc_end
)
