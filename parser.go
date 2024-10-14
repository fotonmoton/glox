package main

import (
	"errors"
	"fmt"
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
func (p *Parser) parse() ([]Stmt, error) {
	defer p.recover()

	stmts := []Stmt{}

	for !p.isAtEnd() {

		if stmt := p.declaration(); stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	return stmts, errors.Join(p.errors...)
}

// declaration -> varDecl | funDecl | statement
func (p *Parser) declaration() Stmt {
	defer p.synchronize()
	if p.match(VAR) {
		return p.varDecl()
	}

	if p.match(FUN) {
		return p.function("function")
	}

	return p.statement()
}

// varDecl -> "var" IDENTIFIER ("=" expression)? ";"
func (p *Parser) varDecl() Stmt {
	name := p.consume(IDENTIFIER, "Expect identifier for variable")

	var initializer Expr = nil
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after expression.")

	return &VarStmt{name, initializer}
}

// funDecl -> "fun" function
// function -> IDENTIFIER "("  parameters? ")" blockStmt
// parameters -> IDENTIFIER ( "," IDENTIFIER )*
func (p *Parser) function(kind string) Stmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))

	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))

	args := []Token{}
	for !p.check(RIGHT_PAREN) {
		args = append(
			args,
			p.consume(
				IDENTIFIER,
				fmt.Sprintf("Expect %s argument.", kind),
			),
		)

		if p.check(COMMA) {
			p.advance()
		}
	}

	p.consume(RIGHT_PAREN, fmt.Sprintf("Expect ')' after %s name.", kind))
	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' after %s arguments.", kind))

	body := p.block()

	return &FunStmt{name, args, body}
}

// statement ->  exprStmt
//
//	| whileStmt
//	| forStmt
//	| printStmt
//	| blockStmt
//	| breakStmt
//	| ifStmt
//	| returnStmt
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

	if p.match(FOR) {
		return p.forStmt()
	}

	if p.match(BREAK) {
		return p.breakStmt()
	}

	if p.match(RETURN) {
		return p.returnStmt()
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

func (p *Parser) block() []Stmt {

	stmts := []Stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Unclosed block: Expected '}'.")

	return stmts
}

// blockStmt -> "{" statement* "}"
func (p *Parser) blockStmt() *BlockStmt {
	return &BlockStmt{p.block()}
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

// for -> "for" ( "(" ( varDecl | exprStmt | ";" ) expression? ";" expression  ")" )? statement
func (p *Parser) forStmt() Stmt {

	if p.check(LEFT_BRACE) {
		return &WhileStmt{&Literal{true}, p.statement()}
	}

	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var init Stmt

	if p.match(SEMICOLON) {
		init = nil
	} else if p.match(VAR) {
		init = p.varDecl()
	} else {
		init = p.exprStmt()
	}

	var cond Expr

	if !p.check(SEMICOLON) {
		cond = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after for loop condition;")

	var incr Expr

	if !p.check(RIGHT_PAREN) {
		incr = p.expression()
	}

	p.consume(RIGHT_PAREN, "Expect ')' after for clauses;")

	var body = p.statement()

	if incr != nil {
		body = &BlockStmt{[]Stmt{body, &ExprStmt{incr}}}
	}

	if cond == nil {
		cond = &Literal{true}
	}

	body = &WhileStmt{cond, body}

	if init != nil {
		body = &BlockStmt{[]Stmt{init, body}}
	}

	return body
}

// env -> "env" ";"
func (p *Parser) envStmt() Stmt {
	p.consume(SEMICOLON, "Expect ';' after 'env'.")
	return &EnvStmt{}
}

// return -> "return" expression ";"
func (p *Parser) returnStmt() Stmt {
	ret := p.expression()
	p.consume(SEMICOLON, "Expect ';' after return;")
	return &ReturnStmt{ret}
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

	return p.call()
}

// call ->  primary ( "(" arguments? ")"  )*
func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.arguments(expr)
		} else {
			break
		}
	}

	return expr
}

// arguments ->  expression ( "," expression )*
func (p *Parser) arguments(callee Expr) Expr {
	arguments := []Expr{}

	if !p.check(RIGHT_PAREN) {
		for {
			arguments = append(arguments, p.expression())

			if !p.match(COMMA) {
				break
			}
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")

	return &Call{callee, paren, arguments}
}

// primary -> IDENTIFIER
//
//	| NUMBER
//	| STRING
//	| "true"
//	| "false"
//	| "nil"
//	| "(" expression ")"
//	| lambda
func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return &Literal{false}
	case p.match(TRUE):
		return &Literal{true}
	case p.match(NIL):
		return &Literal{nil}
	}

	if p.match(FUN) {
		return p.lambda()
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

func (p *Parser) lambda() Expr {
	name := p.previous()

	p.consume(LEFT_PAREN, "Expect '(' before lambda arguments.")

	args := []Token{}
	for !p.check(RIGHT_PAREN) {
		args = append(
			args,
			p.consume(
				IDENTIFIER,
				"Expect lambda argument.",
			),
		)

		if p.check(COMMA) {
			p.advance()
		}
	}

	p.consume(RIGHT_PAREN, "Expect ')' after lambda arguments.")
	p.consume(LEFT_BRACE, "Expect '{' before lambda body.")

	body := p.block()

	return &Lambda{name, args, body}
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
