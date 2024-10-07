package main

type StmtVisitor interface {
	visitIfStmt(i *IfStmt)
	visitVarStmt(v *VarStmt)
	visitExprStmt(es *ExprStmt)
	visitPrintStmt(p *PrintStmt)
	visitBlockStmt(b *BlockStmt)
	visitEnvStmt(e *EnvStmt)
	visitWhileStmt(w *WhileStmt)
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

type EnvStmt struct{}

type IfStmt struct {
	name Token
	expr Expr
	then Stmt
	or   Stmt
}

type WhileStmt struct {
	cond Expr
	body Stmt
}

func (i *IfStmt) stmt()    {}
func (e *EnvStmt) stmt()   {}
func (vs *VarStmt) stmt()  {}
func (es *ExprStmt) stmt() {}
func (p *PrintStmt) stmt() {}
func (b *BlockStmt) stmt() {}
func (w *WhileStmt) stmt() {}

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

func (i *IfStmt) accept(v StmtVisitor) {
	v.visitIfStmt(i)
}

func (e *EnvStmt) accept(v StmtVisitor) {
	v.visitEnvStmt(e)
}

func (w *WhileStmt) accept(v StmtVisitor) {
	v.visitWhileStmt(w)
}
