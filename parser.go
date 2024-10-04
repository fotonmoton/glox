package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

type ParseError struct {
	token   Token
	message string
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("%s: %s", pe.token.lexeme, pe.message)
}

func newParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) parse() Expr {
	defer p.recover()
	return p.expression()
}

func (p *Parser) recover() {
	if err := recover(); err != nil {
		pe := err.(ParseError)
		printError(pe.token, pe.message)
	}
}

// expression -> equality
func (p *Parser) expression() Expr {
	return p.equality()
}

// equality -> comparison ( ( "==" | "!=" ) comparison )*
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		op := p.previous()
		right := p.comparison()
		expr = &Binary{expr, op, right}
	}

	return expr
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		expr = &Binary{expr, op, right}
	}

	return expr
}

// term -> factor ( ( "-" | "+"  ) factor )*
func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		op := p.previous()
		right := p.factor()
		expr = &Binary{expr, op, right}
	}

	return expr
}

// factor -> unary ( ( "/" | "*"  ) unary )*
func (p *Parser) factor() Expr {
	exp := p.unary()

	for p.match(SLASH, STAR) {
		op := p.previous()
		right := p.unary()
		exp = &Binary{exp, op, right}
	}

	return exp
}

// unary -> ( "!" | "-"  ) unary | primary
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		op := p.previous()
		right := p.unary()
		return &Unary{op, right}
	}

	return p.primary()
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return &Literal{false}
	case p.match(TRUE):
		return &Literal{true}
	case p.match(NIL):
		return &Literal{nil}
	}

	if p.match(NUMBER, STRING) {
		return &Literal{p.previous().literal}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression")
		return &Grouping{expr}
	}

	panic(ParseError{p.peek(), "Expect expression"})

	return nil
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().typ == EOF
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().typ == typ
}

func (p *Parser) match(types ...TokenType) bool {

	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(typ TokenType, mes string) Token {
	if p.check(typ) {
		return p.advance()
	}

	panic(ParseError{p.peek(), mes})

	return Token{}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().typ == SEMICOLON {
			return
		}

		switch p.peek().typ {
		case CLASS, FOR, FUN, IF, PRINT, RETURN, VAR, WHILE:
			return
		}

		p.advance()
	}
}
