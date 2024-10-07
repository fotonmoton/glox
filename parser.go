package main

import (
	"fmt"
	"log"
)

type Parser struct {
	tokens  []Token
	current int
	errors  []error
}

type ParseError struct {
	token   Token
	message string
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("ParseError [%d][%s]: %s", pe.token.line, pe.token.typ, pe.message)
}

func newParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// program -> declaration* EOF
func (p *Parser) parse() ([]Stmt, []error) {
	defer p.recover()

	stmts := []Stmt{}

	for !p.isAtEnd() {

		if stmt := p.declaration(); stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	return stmts, p.errors
}

// declaration -> varDecl | statement
func (p *Parser) declaration() Stmt {
	defer p.synchronize()
	if p.match(VAR) {
		return p.varDecl()
	}
	return p.statement()
}

// varDecl -> "var" IDENTIFIER ("=" expression)? ";"
func (p *Parser) varDecl() Stmt {
	name := p.consume(IDENTIFIER, "expect identifier for variable")

	var initializer Expr = nil
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after expression.")

	return &VarStmt{name, initializer}
}

// statement ->  exprStmt
//
//	| whileStmt
//	| printStmt
//	| blockStmt
//	| breakStmt
//	| ifStmt
//	| env
func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStmt()
	}

	if p.match(LEFT_BRACE) {
		return p.blockStmt()
	}

	if p.match(IF) {
		return p.ifStmt()
	}

	if p.match(ENV) {
		return p.envStmt()
	}

	if p.match(WHILE) {
		return p.whileStmt()
	}

	if p.match(BREAK) {
		return p.breakStmt()
	}

	return p.exprStmt()
}

// exprStmt -> expression ";"
func (p *Parser) exprStmt() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")

	if expr == nil {
		return nil
	}

	return &ExprStmt{expr}
}

// printStmt -> "print" expression ";"
func (p *Parser) printStmt() Stmt {
	expr := p.expression()

	if expr == nil {
		p.panic(&ParseError{p.previous(), "Expect expression after 'print'"})
	}

	p.consume(SEMICOLON, "Expect ';' after expression.")
	return &PrintStmt{expr}
}

// blockStmt -> "{" statement* "}"
func (p *Parser) blockStmt() Stmt {

	stmts := []Stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Unclosed block: Expected '}'.")

	return &BlockStmt{stmts}
}

// breakStmt -> break ";"
func (p *Parser) breakStmt() Stmt {
	p.consume(SEMICOLON, "Expect ';' after break.")
	return &BreakStmt{}
}

// if -> "if" "(" expression ")" statement ("else" statement)?
func (p *Parser) ifStmt() Stmt {
	name := p.previous()
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	expr := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after 'if' condition.")
	then := p.statement()

	var or Stmt = nil
	if p.match(ELSE) {
		or = p.statement()
	}

	return &IfStmt{name, expr, then, or}
}

// while -> "while" "(" expression ")" statement
func (p *Parser) whileStmt() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	cond := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after 'while' expression.")
	body := p.statement()

	return &WhileStmt{cond, body}
}

// env -> "env" ";"
func (p *Parser) envStmt() Stmt {
	p.consume(SEMICOLON, "Expect ';' after 'env'.")
	return &EnvStmt{}
}

// expression -> assignment
func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment | or
func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		eq := p.previous()
		val := p.assignment()

		if variable, ok := expr.(*Variable); ok {
			return &Assign{variable.name, val}
		}

		p.panic(&ParseError{eq, "Invalid assignment target."})
	}

	return expr
}

// or -> and ( "or" and )*
func (p *Parser) or() Expr {
	left := p.and()

	for p.match(OR) {
		or := p.previous()
		right := p.and()
		left = &Logical{left, or, right}
	}

	return left
}

// and -> equality ( "and" equality )*
func (p *Parser) and() Expr {
	left := p.equality()

	for p.match(AND) {
		or := p.previous()
		right := p.equality()

		left = &Logical{left, or, right}
	}

	return left
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

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER
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

	if p.match(IDENTIFIER) {
		return &Variable{p.previous()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression")
		return &Grouping{expr}
	}

	p.panic(&ParseError{p.peek(), "Expect expression"})

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

	p.panic(&ParseError{p.peek(), mes})

	return Token{}
}

func (p *Parser) synchronize() {
	err := recover()

	pe := p.isParseError(err)

	if pe == nil {
		return
	}

	p.advance()

	for !p.isAtEnd() {
		if p.previous().typ == SEMICOLON {
			return
		}

		switch p.peek().typ {
		case CLASS, FOR, FUN, IF, PRINT, RETURN, VAR, WHILE, ENV:
			return
		}

		p.advance()
	}

}

func (p *Parser) recover() {
	p.isParseError(recover())
}

func (p *Parser) panic(pe *ParseError) {
	p.errors = append(p.errors, pe)
	log.Println(pe)
	panic(pe)
}

func (p *Parser) isParseError(err any) *ParseError {
	if err == nil {
		return nil
	}

	pe, ok := err.(*ParseError)

	if !ok {
		panic(err)
	}

	return pe
}
