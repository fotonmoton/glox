package main

type StmtVisitor interface {
	visitPrintStmt(p *PrintStmt)
	visitExprStmt(es *ExprStmt)
	visitVarStmt(v *VarStmt)
}

type Stmt interface {
	stmt()
	accept(v StmtVisitor)
}

type PrintStmt struct {
	val Expr
}

type ExprStmt struct {
	expr Expr
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func (p *PrintStmt) stmt() {}
func (es *ExprStmt) stmt() {}
func (vs *VarStmt) stmt()  {}

func (p *PrintStmt) accept(v StmtVisitor) {
	v.visitPrintStmt(p)
}

func (se *ExprStmt) accept(v StmtVisitor) {
	v.visitExprStmt(se)
}

func (vs *VarStmt) accept(v StmtVisitor) {
	v.visitVarStmt(vs)
}
