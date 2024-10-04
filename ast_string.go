package main

import (
	"fmt"
	"strings"
)

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

func (as *AstStringer) visitBinary(b *Binary) any {
	as.str.WriteString("(")
	as.str.WriteString(b.op.lexeme)
	as.str.WriteString(" ")
	b.left.accept(as)
	as.str.WriteString(" ")
	b.right.accept(as)
	as.str.WriteString(")")
	return nil

}

func (as *AstStringer) visitLiteral(l *Literal) any {
	as.str.WriteString(fmt.Sprintf("%v", l.value))
	return nil
}

func (as *AstStringer) visitGrouping(g *Grouping) any {
	as.str.WriteString("(group ")
	g.expression.accept(as)
	as.str.WriteString(")")
	return nil
}

func (as *AstStringer) visitUnary(u *Unary) any {
	as.str.WriteString(fmt.Sprintf("(%s ", u.op.lexeme))
	u.right.accept(as)
	as.str.WriteString(")")
	return nil
}
