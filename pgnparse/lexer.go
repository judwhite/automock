package pgnparse

import (
	"bytes"
	"fmt"
	"strings"
)

const eof = 255

var (
	nagChars         = []byte{'!', '?'}
	digits           = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	terminationChars = []byte{'1', '/', '2', '0', '-'}
	dots             = []byte{'.'}

	resultWhiteWins     = []byte("1-0")
	resultBlackWins     = []byte("0-1")
	resultDraw          = []byte("1/2-1/2")
	resultGameContinues = []byte("*")
)

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

type lexer struct {
	input      []byte // the string being scanned
	start      int    // start position of this item
	pos        int    // current position in the input
	items      []item // channel of scanned items
	parenDepth int    // nesting depth of ( ) exprs
}

func lex(input []byte) (*lexer, error) {
	l := &lexer{
		input: input,
		items: make([]item, 0, 4096),
	}
	err := l.run()
	return l, err
}

func (l *lexer) run() error {
	for state := lexStatement; state != nil; {
		state = state(l)
	}

	if lastItem := l.items[len(l.items)-1]; lastItem.typ == itemError {
		max := lastItem.pos
		for ; max < len(l.input); max++ {
			r := l.input[max]
			if isSpace(r) {
				break
			}
		}

		return fmt.Errorf("%s at pos %d '%s'", lastItem.val, lastItem.pos, l.input[:max])
	}

	return nil
}

func lexStatement(l *lexer) stateFn {
loop:
	for {
		switch r := l.next(); {
		case r == eof:
			break loop
		case isSpace(r):
			return lexSpace
		case isNumeric(r):
			return lexMoveNumber
		case isMovePrefixCharacter(r):
			return lexMoveSAN
		case r == '$':
			return lexNAG
		case r == '"':
			return lexQuotedString
		case r == '{':
			return lexComment
		case r == '(':
			l.emit(itemLeftParen)
			l.parenDepth++
		case r == ')':
			l.emit(itemRightParen)
			l.parenDepth--
			if l.parenDepth < 0 {
				return l.errorf("unexpected right paren %#U", r)
			}
		case r == '.':
			return lexDots
		case r == '*':
			return lexGameTerminationMarker
		case r == '[':
			return lexTag
		default:
			return l.errorf("unrecognized character: %#U", r)
		}
	}
	l.emit(itemEOF) // correctly reached EOF.
	return nil      // stop the run loop.
}

func lexMoveSAN(l *lexer) stateFn {
loop:
	for {
		r := l.next()
		if !runeInMoveCharSet(r) {
			if r != eof {
				l.backup()
				if !l.atMoveSANTerminator() {
					return l.errorf("lexMoveSAN: bad character %#U '%s'", r, l.items[len(l.items)-1].val)
				}
			}
			l.emit(itemMoveSAN)
			break loop
		}
	}
	return lexSuffixAnnotation
}

func lexSuffixAnnotation(l *lexer) stateFn {
	l.chomp()

	r := l.peek()
	if r != '!' && r != '?' {
		return lexStatement
	}

	l.acceptRun(nagChars)

	r = l.peek()
	if !isSpace(r) && r != eof && r != '{' && r != '(' && r != ')' {
		return l.errorf("bad suffix annotation syntax: %q", l.input[l.start:l.pos+1])
	}

	// 1    good move (traditional "!")
	// 2    poor move (traditional "?")
	// 3    very good move (traditional "!!")
	// 4    very poor move (traditional "??")
	// 5    speculative move (traditional "!?")
	// 6    questionable move (traditional "?!")

	nag := string(l.input[l.start:l.pos])
	switch nag {
	case "!":
		nag = "$1"
	case "?":
		nag = "$2"
	case "!!":
		nag = "$3"
	case "??":
		nag = "$4"
	case "!?":
		nag = "$5"
	case "?!":
		nag = "$6"
	default:
		return l.errorf("unknown suffix annotation value: %q", l.input[l.start:l.pos+1])
	}

	l.emitMoveNAG(nag)

	return lexStatement
}

func lexNAG(l *lexer) stateFn {
	l.acceptRun(digits)

	r := l.peek()
	if !isSpace(r) && r != eof && r != '{' && r != '(' && r != ')' {
		return l.errorf("bad NAG syntax: %q", l.input[l.start:l.pos+1])
	}

	l.emit(itemMoveNAG)

	return lexStatement
}

func lexTag(l *lexer) stateFn {
	l.ignore()

	for {
		r := l.next()
		if r == eof {
			return l.errorf("unexpected eof")
		}
		if r == ' ' {
			l.backup()
			l.emit(itemTagName)
			break
		}
	}

	l.chomp()

	r := l.next()
	if r != '"' {
		return l.errorf("unexpected character '%c'; expected '\"'", r)
	}
	l.ignore()

	return lexQuotedString
}

func lexGameTerminationMarker(l *lexer) stateFn {
	l.acceptRun(terminationChars)

	val := l.input[l.start:l.pos]
	if !bytes.Equal(val, resultDraw) &&
		!bytes.Equal(val, resultWhiteWins) &&
		!bytes.Equal(val, resultBlackWins) &&
		!bytes.Equal(val, resultGameContinues) {
		return l.errorf("invalid game termination marker '%s'", val)
	}

	l.emit(itemGameTermination)

	return lexStatement
}

func lexMoveNumber(l *lexer) stateFn {
	l.acceptRun(digits)

	// must end with '.'
	if r := l.peek(); r != '.' {
		if r == '/' || r == '-' {
			return lexGameTerminationMarker
		}
		return l.errorf("bad move number syntax: %q", safeSubstring(l.input, l.start, l.pos+2))
	}

	l.emit(itemMoveNumber)

	return lexStatement
}

func safeSubstring(s []byte, start, end int) string {
	if start < 0 {
		start = 0
	}

	if end > len(s) {
		end = len(s)
	}

	return string(s[start:end])
}

func lexDots(l *lexer) stateFn {
	l.acceptRun(dots)
	l.ignore()
	return lexStatement
}

func lexQuotedString(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			return l.errorf("unexpected eof")
		}
		if r == '"' {
			l.emitString()

			r := l.next()
			// next character must be space or closing bracket
			if !isSpace(r) && r != ']' {
				return l.errorf("invalid character %#U %s", r, l.items[len(l.items)-1].val)
			}
			l.ignore()

			break
		}
	}
	return lexStatement
}

func lexComment(l *lexer) stateFn {
	l.ignore()

	var comment string
	for {
		r := l.next()
		if r == eof {
			return l.errorf("unexpected eof")
		}
		if r == '\n' {
			l.backup()
			val := l.input[l.start:l.pos]
			comment += strings.TrimSpace(string(val)) + " "
			l.chomp()
			continue
		}

		if r == '}' {
			l.emitComment(comment)
			break
		}
	}
	return lexStatement
}

func lexSpace(l *lexer) stateFn {
	l.chomp()
	return lexStatement
}

func (l *lexer) chomp() {
	for {
		r := l.next()
		if r == eof {
			break
		}
		if !isSpace(r) {
			l.backup()
			break
		}
	}
	l.ignore()
}

func (l *lexer) emitMoveNAG(val string) {
	l.items = append(l.items, item{
		typ: itemMoveNAG,
		val: []byte(val),
		pos: l.start,
	})
	l.start = l.pos
}

func (l *lexer) emit(t itemType) {
	l.items = append(l.items, item{
		typ: t,
		val: l.input[l.start:l.pos],
		pos: l.start,
	})
	l.start = l.pos
}

func (l *lexer) emitComment(prefix string) {
	val := string(l.input[l.start : l.pos-1])
	val = strings.TrimSpace(prefix + strings.TrimSpace(val))

	l.items = append(l.items, item{
		typ: itemComment,
		val: []byte(val),
		pos: l.start,
	})
	l.start = l.pos
}

func (l *lexer) emitString() {
	val := l.input[l.start:l.pos]

	// equivalent to: val = bytes.Trim(val, `"`)
	for len(val) > 0 && val[0] == '"' {
		val = val[1:]
	}
	for len(val) > 0 && val[len(val)-1] == '"' {
		val = val[:len(val)-1]
	}

	l.items = append(l.items, item{
		typ: itemString,
		val: val,
		pos: l.start,
	})
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) acceptRun(valid []byte) {
	c := l.next()
	for bytes.IndexByte(valid, c) != -1 {
		c = l.next()
	}
	if c != eof {
		l.backup()
	}
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= 1
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() byte {
	// equivalent to l.next(), l.backup()
	if l.pos >= len(l.input) {
		return eof
	}
	return l.input[l.pos]
}

func (l *lexer) next() byte {
	if l.pos >= len(l.input) {
		return eof
	}
	r := l.input[l.pos]
	l.pos += 1
	return r
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items = append(l.items, item{
		typ: itemError,
		val: []byte(fmt.Sprintf(format, args...)),
		pos: l.pos - 1,
	})
	return nil
}

func (l *lexer) atMoveSANTerminator() bool {
	r := l.peek()
	if isSpace(r) {
		return true
	}
	switch r {
	case eof, ')', '{', '?', '!':
		return true
	}
	return false
}

func isSpace(r byte) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isMovePrefixCharacter(r byte) bool {
	if r >= 'a' && r <= 'h' {
		return true
	}
	switch r {
	case 'K', 'Q', 'R', 'N', 'B', 'O', '@', 'P':
		return true
	}
	return false
}

func isNumeric(r byte) bool {
	return r >= '0' && r <= '9'
}

func runeInMoveCharSet(r byte) bool {
	if r >= 'a' && r <= 'h' {
		return true
	}
	if r >= '1' && r <= '8' {
		return true
	}
	switch r {
	case 'K', 'Q', 'R', 'N', 'B', 'O', '-', '+', '#', '=', 'x', '@':
		return true
	}
	return false
}
