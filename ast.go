package main

import (
	"fmt"
	"strings"
)

type Expr interface {
	expr()
	accept(v Visitor)
}

type Unary struct {
	op    Token
	right Expr
}

func (u *Unary) expr() {}
func (u *Unary) accept(v Visitor) {
	v.visitUnary(u)
}

type Grouping struct {
	expression Expr
}

func (g *Grouping) expr() {}
func (g *Grouping) accept(v Visitor) {
	v.visitGrouping(g)
}

type Literal struct {
	value any
}

func (l *Literal) expr() {}
func (l *Literal) accept(v Visitor) {
	v.visitLiteral(l)
}

type Binary struct {
	left  Expr
	op    Token
	right Expr
}

func (b *Binary) expr() {}
func (b *Binary) accept(v Visitor) {
	v.visitBinary(b)
}

type Visitor interface {
	visitBinary(b *Binary)
	visitLiteral(l *Literal)
	visitGrouping(g *Grouping)
	visitUnary(u *Unary)
}

type AstStringer struct {
	str strings.Builder
}

func (as AstStringer) String(expr Expr) string {
	expr.accept(&as)
	return as.str.String()
}

func (as *AstStringer) visitBinary(b *Binary) {
	as.str.WriteString("(")
	as.str.WriteString(b.op.lexeme)
	as.str.WriteString(" ")
	b.left.accept(as)
	as.str.WriteString(" ")
	b.right.accept(as)
	as.str.WriteString(")")

}

func (as *AstStringer) visitLiteral(l *Literal) {
	as.str.WriteString(fmt.Sprintf("%v", l.value))
}

func (as *AstStringer) visitGrouping(g *Grouping) {
	as.str.WriteString("(group ")
	g.expression.accept(as)
	as.str.WriteString(")")
}

func (as *AstStringer) visitUnary(u *Unary) {
	as.str.WriteString(fmt.Sprintf("(%s ", u.op.lexeme))
	u.right.accept(as)
	as.str.WriteString(")")
}
