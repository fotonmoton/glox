package main

import (
	"fmt"
	"strings"
)

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

func (as *AstToRPN) visitBinary(b *Binary) any {
	b.left.accept(as)
	as.str.WriteString(" ")
	b.right.accept(as)
	as.str.WriteString(" ")
	as.str.WriteString(b.op.lexeme)
	return nil
}

func (as *AstToRPN) visitLiteral(l *Literal) any {
	as.str.WriteString(fmt.Sprintf("%v", l.value))
	return nil
}

func (as *AstToRPN) visitGrouping(g *Grouping) any {
	g.expression.accept(as)
	as.str.WriteString(" group")
	return nil
}

func (as *AstToRPN) visitUnary(u *Unary) any {
	u.right.accept(as)
	as.str.WriteString(fmt.Sprintf(" %s", u.op.lexeme))
	return nil
}
