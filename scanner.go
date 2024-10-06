package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=TokenType
type TokenType int

const (
	// one char
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// one or two chars
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	CLASS
	ENV
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
	"env":    ENV,
}

type Token struct {
	typ     TokenType
	lexeme  string
	literal any
	line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s - %s - %v", t.typ, t.lexeme, t.literal)
}

type ScanError struct {
	line    int
	message string
}

func (se *ScanError) Error() string {
	return fmt.Sprintf("ScanError [%d] Error: %s", se.line, se.message)
}

type Scanner struct {
	source  []byte
	tokens  []Token
	start   int
	current int
	line    int
	err     error
}

func newScanner(source []byte) *Scanner {
	return &Scanner{source: source, start: 0, current: 0, line: 1, err: nil}
}

func (s *Scanner) scan() ([]Token, error) {

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "EOF", struct{}{}, s.line})

	return s.tokens, s.err
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, struct{}{})
	case ')':
		s.addToken(RIGHT_PAREN, struct{}{})
	case '{':
		s.addToken(LEFT_BRACE, struct{}{})
	case '}':
		s.addToken(RIGHT_BRACE, struct{}{})
	case ',':
		s.addToken(COMMA, struct{}{})
	case '.':
		s.addToken(DOT, struct{}{})
	case '-':
		s.addToken(MINUS, struct{}{})
	case '+':
		s.addToken(PLUS, struct{}{})
	case ';':
		s.addToken(SEMICOLON, struct{}{})
	case '*':
		s.addToken(STAR, struct{}{})

	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, struct{}{})
		} else {
			s.addToken(BANG, struct{}{})
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, struct{}{})
		} else {
			s.addToken(EQUAL, struct{}{})
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, struct{}{})
		} else {
			s.addToken(LESS, struct{}{})
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, struct{}{})
		} else {
			s.addToken(GREATER, struct{}{})
		}

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, struct{}{})
		}
	case '"':
		s.string()
	case ' ':
	case '\t':
	case '\r':
		break
	case '\n':
		s.line++
	default:
		if unicode.IsDigit(c) {
			s.number()
			break
		}

		if s.isAlpha(c) {
			s.identifier()
			break
		}

		s.error(&ScanError{s.line, fmt.Sprintf("Unexpected character %s", string(c))})
	}
}

func (s *Scanner) addToken(typ TokenType, literal any) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, Token{typ: typ, lexeme: text, literal: literal, line: s.line})
}

func (s *Scanner) identifier() {
	for unicode.IsDigit(s.peek()) || s.isAlpha(s.peek()) {
		s.advance()
	}

	str := string(s.source[s.start:s.current])

	if id, found := keywords[string(str)]; found {
		s.addToken(id, str)
	} else {
		s.addToken(IDENTIFIER, str)
	}

}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.error(&ScanError{s.line, "Unterminated string"})
		return
	}

	s.advance()
	str := string(s.source[s.start+1 : s.current-1])
	s.addToken(STRING, str)
}

func (s *Scanner) number() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()
	}

	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	num, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)

	if err != nil {
		s.error(&ScanError{s.line, err.Error()})
	}

	s.addToken(NUMBER, num)
}

func (s *Scanner) isAlpha(ch rune) bool {
	return regexp.MustCompile(`^[A-Za-z_]+$`).MatchString(string(ch))
}

func (s *Scanner) advance() rune {
	char, size := utf8.DecodeRune(s.source[s.current:])
	s.current += size
	return char
}

func (s *Scanner) peek() rune {
	char, _ := utf8.DecodeRune(s.source[s.current:])
	return char
}

func (s *Scanner) peekNext() rune {
	_, size := utf8.DecodeRune(s.source[s.current+1:])
	if s.current+size >= len(s.source) {
		return '\000'
	}

	next, _ := utf8.DecodeRune(s.source[s.current+size:])
	return next
}

func (s *Scanner) match(ch rune) bool {
	if s.isAtEnd() {
		return false
	}

	decoded, size := utf8.DecodeRune(s.source[s.current:])

	if decoded != ch {
		return false
	}

	s.current += size
	return true
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) error(err *ScanError) {
	log.Print(err)
	s.err = err
}
