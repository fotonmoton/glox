package main

type ExprVisitor interface {
	visitCall(c *Call) any
	visitUnary(u *Unary) any
	visitLambda(l *Lambda) any
	visitBinary(b *Binary) any
	visitLiteral(l *Literal) any
	visitGrouping(g *Grouping) any
	visitVariable(v *Variable) any
	visitLogical(l *Logical) any
	visitAssignment(a *Assign) any
}

type Expr interface {
	expr()
	accept(v ExprVisitor) any
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

type Variable struct {
	name Token
}

type Assign struct {
	variable Token
	value    Expr
}

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

type Call struct {
	callee Expr
	paren  Token
	args   []Expr
}

type Lambda struct {
	name Token
	args []Token
	body []Stmt
}

func (c *Call) expr()     {}
func (u *Unary) expr()    {}
func (a *Assign) expr()   {}
func (b *Binary) expr()   {}
func (l *Lambda) expr()   {}
func (l *Literal) expr()  {}
func (g *Grouping) expr() {}
func (v *Variable) expr() {}
func (l *Logical) expr()  {}

func (u *Unary) accept(v ExprVisitor) any {
	return v.visitUnary(u)
}

func (b *Binary) accept(v ExprVisitor) any {
	return v.visitBinary(b)
}

func (l *Literal) accept(v ExprVisitor) any {
	return v.visitLiteral(l)
}

func (g *Grouping) accept(v ExprVisitor) any {
	return v.visitGrouping(g)
}

func (va *Variable) accept(v ExprVisitor) any {
	return v.visitVariable(va)
}

func (a *Assign) accept(v ExprVisitor) any {
	return v.visitAssignment(a)
}

func (l *Logical) accept(v ExprVisitor) any {
	return v.visitLogical(l)
}

func (c *Call) accept(v ExprVisitor) any {
	return v.visitCall(c)
}

func (l *Lambda) accept(v ExprVisitor) any {
	return v.visitLambda(l)
}
