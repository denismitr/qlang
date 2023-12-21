package tokenizer

type Token int

const (
	// Stop tokens

	Illegal Token = iota
	EOF

	// Single-character tokens

	Assign
	Colon
	Comma
	Divide
	Dot
	Gt
	LBrace
	LBracket
	LParen
	Lt
	Minus
	Modulo
	Plus
	RBrace
	RBracket
	RParen
	Times

	// Two-character tokens

	Equal
	Gte
	Lte
	NotEqual

	// Three-character tokens

	Ellipsis

	// Keywords

	And
	Else
	False
	For
	If
	In
	Nil
	Not
	Or
	Return
	True
	While
	Func

	// Literals and identifiers

	Int
	Name
	String
)

var tokenNames = map[Token]string{
	Illegal: "Illegal",
	EOF:     "EOF",

	Assign:   "=",
	Colon:    ":",
	Comma:    ",",
	Divide:   "/",
	Dot:      ".",
	Gt:       ">",
	LBrace:   "{",
	LBracket: "[",
	LParen:   "(",
	Lt:       "<",
	Minus:    "-",
	Modulo:   "%",
	Plus:     "+",
	RBrace:   "}",
	RBracket: "]",
	RParen:   ")",
	Times:    "*",

	Equal:    "==",
	Gte:      ">=",
	Lte:      "<=",
	NotEqual: "!=",

	Ellipsis: "...",

	And:    "and",
	Else:   "else",
	False:  "false",
	For:    "for",
	Func:   "func",
	If:     "if",
	In:     "in",
	Nil:    "nil",
	Not:    "not",
	Or:     "or",
	Return: "return",
	True:   "true",
	While:  "while",

	Int:    "int",
	Name:   "name",
	String: "str",
}

var keywordTokens = map[string]Token{
	"and":    And,
	"else":   Else,
	"false":  False,
	"for":    For,
	"func":   Func,
	"if":     If,
	"in":     In,
	"nil":    Nil,
	"not":    Not,
	"or":     Or,
	"return": Return,
	"true":   True,
	"while":  While,
}

func (t Token) String() string {
	return tokenNames[t]
}
