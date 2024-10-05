package main

import (
	"fmt"
	"strings"
)

type AstStringer struct {
	str strings.Builder
}

func (as AstStringer) String(stmts []Stmt) string {

	for _, stmt := range stmts {
		stmt.accept(&as)
	}

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

func (as *AstStringer) visitVariable(va *Variable) any {
	as.str.WriteString(va.name.lexeme)
	return nil
}

func (as *AstStringer) visitAssignment(a *Assign) any {
	as.str.WriteString(fmt.Sprintf("(= %s ", a.variable.lexeme))
	a.value.accept(as)
	as.str.WriteString(")")
	return nil
}

func (as *AstStringer) visitPrintStmt(p *PrintStmt) {
	as.str.WriteString("(print ")
	p.val.accept(as)
	as.str.WriteString(")")
}

func (as *AstStringer) visitExprStmt(se *ExprStmt) {
	se.expr.accept(as)
}

func (as *AstStringer) visitVarStmt(vs *VarStmt) {
	if vs.initializer != nil {
		as.str.WriteString(fmt.Sprintf("(var %v ", vs.name.literal))
		vs.initializer.accept(as)
		as.str.WriteString(")")
	} else {
		as.str.WriteString(fmt.Sprintf("(var %v)", vs.name.literal))
	}
}
