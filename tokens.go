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
	PLUS       // +
	MINUS      // -
	MULTIPLY   // *
	DIVIDE     // /
	INT_DIVIDE // \
	POW        // ^
	MODULO     // %
	operators_end

	// Known types/keywords
	keywords_begin
	DIGITS // Contiguous block of digits
	keywords_end

	// Misc characters
	misc_begin
	LPAREN // (
	RPAREN // )
	misc_end
)
