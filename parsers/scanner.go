package parsers

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	reader *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{reader: bufio.NewReader(reader)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.reader.UnreadRune() }

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if ch == eof {
		return EOF, ""
	}

	s.unread()
	if isWhitespace(ch) {
		return s.scanWhitespace()
	}
	return s.scanWord()
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var i int
	for {
		if i >= 300 {
			// prevent against accidental infinite loop
			log.Println("infinite loop in scanner.scanWord detected")
			break
		}
		i++
		ch := s.read()
		if ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WHITESPACE, buf.String()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

// scanWord consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanWord() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var i int
	for {
		if i >= 300 {
			// prevent against accidental infinite loop
			log.Println("infinite loop in scanner.scanWord detected")
			break
		}
		i++
		ch := s.read()
		if ch == eof {
			break
		} else if isWhitespace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	stringBuf := buf.String()
	switch stringBuf {
	case "OPEN_ALL":
		return OPEN_ALL, stringBuf
	}

	return WORD, stringBuf
}
