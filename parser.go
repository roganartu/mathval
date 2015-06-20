package mathval

type Token int

const (
	ILLEGAL Token = iota

	// Special tokens
	special_tokens_beg
	EOF
	WS // Whitespace
	special_tokens_end

	// String literals
	literals_beg
	IDENTIFIER
	literals_end

	// Arithmetic operators
	operators_beg
	PLUS       // +
	MINUS      // -
	MULTIPLY   // * or x
	DIVIDE     // /
	INT_DIVIDE // \
	POW        // ^
	MODULO     // % or mod depending on context
	PERCENT    // % depending on context
	operators_end

	// Misc characters
	misc_beg
	LPAREN // (
	RPAREN // )
	misc_end
)
