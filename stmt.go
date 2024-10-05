package main

type StmtVisitor interface {
	visitVarStmt(v *VarStmt)
	visitExprStmt(es *ExprStmt)
	visitPrintStmt(p *PrintStmt)
	visitBlockStmt(b *BlockStmt)
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

type BlockStmt struct {
	stmts []Stmt
}

func (vs *VarStmt) stmt()  {}
func (es *ExprStmt) stmt() {}
func (p *PrintStmt) stmt() {}
func (b *BlockStmt) stmt() {}

func (p *PrintStmt) accept(v StmtVisitor) {
	v.visitPrintStmt(p)
}

func (se *ExprStmt) accept(v StmtVisitor) {
	v.visitExprStmt(se)
}

func (vs *VarStmt) accept(v StmtVisitor) {
	v.visitVarStmt(vs)
}

func (b *BlockStmt) accept(v StmtVisitor) {
	v.visitBlockStmt(b)
}
