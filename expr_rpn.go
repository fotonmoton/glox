package main

import (
	"fmt"
	"strings"
)

type ExprToRPN struct {
	str strings.Builder
}

func (as ExprToRPN) String(expr Expr) string {

	if expr == nil {
		return ""
	}

	expr.accept(&as)
	return as.str.String()
}

func (as *ExprToRPN) visitBinary(b *Binary) any {
	b.left.accept(as)
	as.str.WriteString(" ")
	b.right.accept(as)
	as.str.WriteString(" ")
	as.str.WriteString(b.op.lexeme)
	return nil
}

func (as *ExprToRPN) visitLiteral(l *Literal) any {
	as.str.WriteString(fmt.Sprintf("%v", l.value))
	return nil
}

func (as *ExprToRPN) visitGrouping(g *Grouping) any {
	g.expression.accept(as)
	as.str.WriteString(" group")
	return nil
}

func (as *ExprToRPN) visitUnary(u *Unary) any {
	u.right.accept(as)
	as.str.WriteString(fmt.Sprintf(" %s", u.op.lexeme))
	return nil
}

func (as *ExprToRPN) visitVariable(va *Variable) any {
	as.str.WriteString(va.name.lexeme)
	return nil
}

func (as *ExprToRPN) visitAssignment(a *Assign) any {
	as.str.WriteString(fmt.Sprintf("%v %s =", a.value, a.variable.lexeme))
	return nil
}

func (as *ExprToRPN) visitLogical(lo *Logical) any {
	lo.left.accept(as)
	lo.right.accept(as)
	as.str.WriteString(" or")
	return nil
}

func (as *ExprToRPN) visitCall(c *Call) any {
	for _, arg := range c.arguments {
		arg.accept(as)
	}
	c.callee.accept(as)
	as.str.WriteString(" call")
	return nil
}
