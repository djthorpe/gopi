package mdns

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

////////////////////////////////////////////////////////////////////////////////
// TOKENIZER

const (
	st_start = iota
	st_backslash
	st_digit
)

type tokenizer struct {
	State int
	in    []rune
	out   []byte
	store string
	pos   int
}

const (
	EOF rune = -(iota + 1)
)

var (
	ErrParseError = errors.New("Parse Error")
)

func NewTokenizer(src string) *tokenizer {
	this := new(tokenizer)
	this.in = []rune(src)
	this.out = make([]byte, 0)
	this.pos = 0
	this.State = st_start
	return this
}

func (this *tokenizer) Next() rune {
	if this.pos >= len(this.in) {
		return EOF
	} else {
		r := this.in[this.pos]
		this.pos++
		return r
	}
}

func (this *tokenizer) Conv(tok rune) rune {
	switch {
	case tok == 'n':
		return '\n'
	case tok == 't':
		return '\t'
	case tok == 'r':
		return '\r'
	case tok == 'f':
		return '\f'
	case unicode.IsDigit(tok):
		return rune(0)
	default:
		return tok
	}
}

func (this *tokenizer) NewState(state int) error {
	if state != this.State && (state == st_start || state == st_backslash) && len(this.store) > 0 {
		// Change xFF into 0xFF
		if strings.HasPrefix(this.store, "x") {
			this.store = "0" + this.store
		}
		// Convert to a byte
		if value, err := strconv.ParseUint(this.store, 10, 32); err != nil {
			return ErrParseError
		} else if value > 0xFF {
			return ErrParseError
		} else {
			this.out = append(this.out, byte(value))
		}
		// Reset the store
		this.store = ""
	}
	this.State = state
	return nil
}

func (this *tokenizer) Append(tok rune, state int) error {
	if err := this.NewState(state); err != nil {
		return err
	} else {
		this.out = append(this.out, []byte(string(tok))...)
		return nil
	}
}

func (this *tokenizer) Store(tok rune, state int) error {
	if err := this.NewState(state); err != nil {
		return err
	} else {
		this.store += string(tok)
		return nil
	}
}

func (this *tokenizer) String() string {
	// Only return the parsed version if valid
	if utf8.Valid(this.out) {
		return string(this.out)
	} else {
		return string(this.in)
	}
}

////////////////////////////////////////////////////////////////////////////////
// QUOTE

func Quote(src string) string {
	// Tokenizer
	t := NewTokenizer(src)
	for tok := t.Next(); tok != EOF; tok = t.Next() {
		switch {
		case tok == ' ', tok == '@', tok == '.', tok == '\\':
			t.Append('\\', t.State)
			t.Append(tok, t.State)
		case tok == '\n':
			t.Append('\\', t.State)
			t.Append('n', t.State)
		case tok == '\r':
			t.Append('\\', t.State)
			t.Append('r', t.State)
		case tok == '\t':
			t.Append('\\', t.State)
			t.Append('t', t.State)
		case tok == '\f':
			t.Append('\\', t.State)
			t.Append('f', t.State)
		case unicode.IsDigit(tok), unicode.IsLetter(tok):
			fallthrough
		case tok == '_', tok == '-', tok == ':':
			t.Append(tok, t.State)
		default:
			for _, b := range []byte(string(tok)) {
				t.Append('\\', t.State)
				for _, d := range []rune(fmt.Sprintf("%03d", b)) {
					t.Append(d, t.State)
				}
			}
		}
	}
	return t.String()
}

////////////////////////////////////////////////////////////////////////////////
// UNQUOTE

// Unquote returns a bare string without quoted characters. The following
// Coversions happen:
// \\ \n \f \t \r Happen as normal
// \xFF returns a byte from hex
// \123 returns a byte from decimal
// \0123 returns a bype from octal
func Unquote(src string) (string, error) {
	// The simple case without any backslashes
	if strings.ContainsAny(src, "\\") == false {
		return src, nil
	}

	// Tokenizer
	t := NewTokenizer(src)
	for tok := t.Next(); tok != EOF; tok = t.Next() {
		switch t.State {
		case st_start:
			if tok == '\\' {
				if err := t.NewState(st_backslash); err != nil {
					return src, err
				}
			} else if err := t.Append(tok, st_start); err != nil {
				return src, err
			}
		case st_backslash:
			if ch := t.Conv(tok); ch != rune(0) {
				if err := t.Append(ch, st_start); err != nil {
					return src, err
				}
			} else if unicode.IsDigit(tok) || tok == 'x' {
				if err := t.Store(tok, st_digit); err != nil {
					return src, err
				}
			}
		case st_digit:
			if unicode.IsDigit(tok) {
				if err := t.Store(tok, st_digit); err != nil {
					return src, err
				}
			} else if tok == '\\' {
				if err := t.NewState(st_backslash); err != nil {
					return src, err
				}
			} else if err := t.Append(tok, st_start); err != nil {
				return src, err
			}
		}
	}

	// We should always end up in st_start state
	if t.State != st_start {
		return src, ErrParseError
	}

	// Eject last value
	if err := t.NewState(st_start); err != nil {
		return src, err
	}

	// Success
	return t.String(), nil
}
