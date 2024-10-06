package main

type ExprVisitor interface {
	visitUnary(u *Unary) any
	visitBinary(b *Binary) any
	visitLiteral(l *Literal) any
	visitGrouping(g *Grouping) any
	visitVariable(v *Variable) any
	visitLogicalOr(l *LogicalOr) any
	visitAssignment(a *Assign) any
	visitLogicalAnd(l *LogicalAnd) any
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

type LogicalOr struct {
	left  Expr
	or    Token
	right Expr
}

type LogicalAnd struct {
	left  Expr
	and   Token
	right Expr
}

func (u *Unary) expr()      {}
func (a *Assign) expr()     {}
func (b *Binary) expr()     {}
func (l *Literal) expr()    {}
func (g *Grouping) expr()   {}
func (v *Variable) expr()   {}
func (l *LogicalOr) expr()  {}
func (l *LogicalAnd) expr() {}

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

func (l *LogicalOr) accept(v ExprVisitor) any {
	return v.visitLogicalOr(l)
}

func (l *LogicalAnd) accept(v ExprVisitor) any {
	return v.visitLogicalAnd(l)
}
