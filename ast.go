package main

import (
	"fmt"
	"strings"
)

type Visitor interface {
	visitBinary(b *Binary)
	visitLiteral(l *Literal)
	visitGrouping(g *Grouping)
	visitUnary(u *Unary)
}

type Expr interface {
	expr()
	accept(v Visitor)
}

type Binary struct {
	left  Expr
	op    Token
	right Expr
}

type Unary struct {
	op    Token
	right Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value any
}

func (u *Unary) expr()    {}
func (g *Grouping) expr() {}
func (l *Literal) expr()  {}
func (b *Binary) expr()   {}

func (u *Unary) accept(v Visitor) {
	v.visitUnary(u)
}

func (g *Grouping) accept(v Visitor) {
	v.visitGrouping(g)
}

func (l *Literal) accept(v Visitor) {
	v.visitLiteral(l)
}

func (b *Binary) accept(v Visitor) {
	v.visitBinary(b)
}

type AstStringer struct {
	str strings.Builder
}

func (as AstStringer) String(expr Expr) string {

	if expr == nil {
		return ""
	}

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

type AstToRPN struct {
	str strings.Builder
}

func (as AstToRPN) String(expr Expr) string {

	if expr == nil {
		return ""
	}

	expr.accept(&as)
	return as.str.String()
}

func (as *AstToRPN) visitBinary(b *Binary) {
	b.left.accept(as)
	as.str.WriteString(" ")
	b.right.accept(as)
	as.str.WriteString(" ")
	as.str.WriteString(b.op.lexeme)

}

func (as *AstToRPN) visitLiteral(l *Literal) {
	as.str.WriteString(fmt.Sprintf("%v", l.value))
}

func (as *AstToRPN) visitGrouping(g *Grouping) {
	g.expression.accept(as)
	as.str.WriteString(" group")
}

func (as *AstToRPN) visitUnary(u *Unary) {
	u.right.accept(as)
	as.str.WriteString(fmt.Sprintf(" %s", u.op.lexeme))
}
