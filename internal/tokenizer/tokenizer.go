package tokenizer

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

var (
	ErrEOF = errors.New("EOF")
)

type Position struct {
	Line   int
	Column int
}

func Tokenize(input string) ([]Item, error) {
	t := NewTokenizer([]byte(input))
	var items []Item
	for {
		item, err := t.Next()
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}
			return nil, fmt.Errorf("at pos %v, %s token: %w", item.Position, item.Token, err)
		}

		items = append(items, item)
	}

	return items, nil
}

type Tokenizer struct {
	input   []byte
	offset  int
	current rune
	err     error
	pos     Position
	nextPos Position
}

func NewTokenizer(input []byte) *Tokenizer {
	t := new(Tokenizer)
	t.input = input
	t.nextPos.Line = 1
	t.nextPos.Column = 1
	t.next()
	return t
}

func (t *Tokenizer) next() {
	t.pos = t.nextPos
	ch, size := utf8.DecodeRune(t.input[t.offset:])
	if size == 0 {
		t.current = -1
		return
	}

	if ch == utf8.RuneError {
		t.current = -1
		t.err = fmt.Errorf("invalid UTF-8 byte 0x%02x", t.input[t.offset])
		return
	}

	if ch == '\n' {
		t.nextPos.Line++
		t.nextPos.Column = 1
	} else {
		t.nextPos.Column++
	}
	t.current = ch
	t.offset += size
}

func (t *Tokenizer) skip() {
	for {
		for t.current == ' ' || t.current == '\t' || t.current == '\r' || t.current == '\n' {
			t.next()
		}

		if !(t.current == '/' && t.offset < len(t.input) && t.input[t.offset] == '/') {
			break
		}

		// Skip //-prefixed comment (to end of line or end of input)
		t.next()
		t.next()
		for t.current != '\n' && t.current >= 0 {
			t.next()
		}
		t.next()
	}
}

func isAlphaUnderscore(ch rune) bool {
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNumeric(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

type Item struct {
	Token    Token
	Position Position
	Value    string
}

// Next returns token type, the position and token value of the next token
func (t *Tokenizer) Next() (Item, error) {
	t.skip()
	if t.current <= 0 {
		if t.err != nil {
			return Item{Illegal, t.pos, ""}, t.err
		}
		return Item{EOF, t.pos, ""}, ErrEOF
	}

	pos := t.pos
	token := Illegal
	value := ""
	current := t.current

	t.next()
	if isAlphaUnderscore(current) {
		runes := []rune{current}
		for isAlphaUnderscore(t.current) || isNumeric(t.current) {
			runes = append(runes, t.current)
			t.next()
		}
		name := string(runes)
		token, isKeyword := keywordTokens[name]
		if !isKeyword {
			token = Name
			value = name
		}
		return Item{token, pos, value}, nil
	}

	switch current {
	case ':':
		token = Colon
	case ',':
		token = Comma
	case '/':
		token = Divide
	case '{':
		token = LBrace
	case '[':
		token = LBracket
	case '(':
		token = LParen
	case '-':
		token = Minus
	case '%':
		token = Modulo
	case '+':
		token = Plus
	case '}':
		token = RBrace
	case ']':
		token = RBracket
	case ')':
		token = RParen
	case '*':
		token = Times
	case '=':
		if t.current == '=' {
			t.next()
			token = Equal
		} else {
			token = Assign
		}
	case '!':
		if t.current == '=' {
			t.next()
			token = NotEqual
		} else {
			return Item{Illegal, pos, ""}, fmt.Errorf("expected != instead of !%c", t.current)
		}
	case '<':
		if t.current == '=' {
			t.next()
			token = Lte
		} else {
			token = Lt
		}
	case '>':
		if t.current == '=' {
			t.next()
			token = Gte
		} else {
			token = Gt
		}
	case '.':
		if t.current == '.' {
			t.next()
			if t.current != '.' {
				return Item{Illegal, pos, ".."}, fmt.Errorf("unexpected %s", "..")
			}
			t.next()
			token = Ellipsis
		} else {
			token = Dot
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		runes := []rune{current}
		for isNumeric(t.current) {
			runes = append(runes, t.current)
			t.next()
		}
		token = Int
		value = string(runes)
	case '"':
		var runes []rune
		for t.current != '"' {
			if t.current < 0 {
				return Item{Illegal, pos, ""}, errors.New("didn't find end quote in string")
			}
			if t.current == '\n' || t.current == '\r' {
				return Item{Illegal, pos, ""}, errors.New("can't have newline in string")
			}
			runes = append(runes, t.current)
			t.next()
		}
		t.next()
		token = String
		value = string(runes)
	default:
		return Item{Illegal, pos, string(current)}, fmt.Errorf("unexpected %c", current)
	}

	return Item{token, pos, value}, nil
}
