package main

import (
	"fmt"
	"strings"
)

type AstStringer struct {
	str   strings.Builder
	stmts []Stmt
}

func (as *AstStringer) visitGet(g *Get) any {
	as.str.WriteString(fmt.Sprintf("(get %s)", g.name.lexeme))
	return nil
}

func (as *AstStringer) visitSet(s *Set) any {
	as.str.WriteString(fmt.Sprintf("(set %s ", s.name.lexeme))
	s.obj.accept(as)
	as.str.WriteString(")")
	return nil
}

func (as AstStringer) String() string {

	for _, stmt := range as.stmts {
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

func (as *AstStringer) visitLogical(l *Logical) any {
	as.str.WriteString(fmt.Sprintf("(%s ", l.operator.lexeme))
	l.left.accept(as)
	as.str.WriteString(" ")
	l.right.accept(as)
	as.str.WriteString(")")
	return nil
}

func (as *AstStringer) visitCall(c *Call) any {
	as.str.WriteString("(call ")
	c.callee.accept(as)
	if len(c.args) != 0 {
		as.str.WriteString(" ")
	}
	for i, arg := range c.args {
		arg.accept(as)
		if i < len(c.args)-1 {
			as.str.WriteString(" ")
		}
	}
	as.str.WriteString(")")

	return nil
}

func (as *AstStringer) visitLambda(l *Lambda) any {
	as.str.WriteString("(lambda ")
	if len(l.args) != 0 {
		as.str.WriteString("(")
		for i, arg := range l.args {
			as.str.WriteString(arg.lexeme)
			if i < len(l.args)-1 {
				as.str.WriteString(" ")
			}
		}
		as.str.WriteString(")")
	}
	for _, stmt := range l.body {
		stmt.accept(as)
	}
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

func (as *AstStringer) visitBlockStmt(b *BlockStmt) {
	as.str.WriteString("(block ")

	for _, stmt := range b.stmts {
		stmt.accept(as)
	}

	as.str.WriteString(")")

}

func (as *AstStringer) visitIfStmt(i *IfStmt) {
	as.str.WriteString("(if ")
	i.cond.accept(as)
	as.str.WriteString(" ")
	i.then.accept(as)
	if i.or != nil {
		as.str.WriteString(" ")
		i.or.accept(as)
	}
	as.str.WriteString(")")
}

func (as *AstStringer) visitEnvStmt(e *EnvStmt) {
	as.str.WriteString("(env)")
}

func (as *AstStringer) visitWhileStmt(w *WhileStmt) {
	as.str.WriteString("(while ")
	w.cond.accept(as)
	as.str.WriteString(" ")
	w.body.accept(as)
	as.str.WriteString(")")
}

func (as *AstStringer) visitBreakStmt(b *BreakStmt) {
	as.str.WriteString("(break)")
}

func (as *AstStringer) visitFunStmt(f *FunStmt) {
	as.str.WriteString(fmt.Sprintf("(fun %s ", f.name.lexeme))
	if len(f.args) != 0 {
		as.str.WriteString("(")
		for i, arg := range f.args {
			as.str.WriteString(arg.lexeme)
			if i < len(f.args)-1 {
				as.str.WriteString(" ")
			}
		}
		as.str.WriteString(")")
	}
	for _, stmt := range f.body {
		stmt.accept(as)
	}
	as.str.WriteString(")")
}

func (as *AstStringer) visitReturnStmt(r *ReturnStmt) {
	as.str.WriteString("(return ")
	r.value.accept(as)
	as.str.WriteString(")")
}

func (as *AstStringer) visitClassStmt(c *ClassStmt) {
	fmt.Printf("(class %s)", c.name.lexeme)
}
