package main

type Visitor interface {
	visitUnary(u *Unary) any
	visitBinary(b *Binary) any
	visitLiteral(l *Literal) any
	visitGrouping(g *Grouping) any
}

type Expr interface {
	expr()
	accept(v Visitor) any
}

type Unary struct {
	op    Token
	right Expr
}

type Binary struct {
	left  Expr
	op    Token
	right Expr
}

type Literal struct {
	value any
}

type Grouping struct {
	expression Expr
}

func (u *Unary) expr()    {}
func (b *Binary) expr()   {}
func (l *Literal) expr()  {}
func (g *Grouping) expr() {}

func (u *Unary) accept(v Visitor) any {
	return v.visitUnary(u)
}

func (b *Binary) accept(v Visitor) any {
	return v.visitBinary(b)
}

func (l *Literal) accept(v Visitor) any {
	return v.visitLiteral(l)
}

func (g *Grouping) accept(v Visitor) any {
	return v.visitGrouping(g)
}
