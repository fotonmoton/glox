package main

type StmtVisitor interface {
	visitIfStmt(i *IfStmt)
	visitVarStmt(v *VarStmt)
	visitEnvStmt(e *EnvStmt)
	visitFunStmt(f *FunStmt)
	visitExprStmt(es *ExprStmt)
	visitPrintStmt(p *PrintStmt)
	visitBlockStmt(b *BlockStmt)
	visitWhileStmt(w *WhileStmt)
	visitBreakStmt(b *BreakStmt)
	visitReturnStmt(r *ReturnStmt)
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
	cond Expr
	then Stmt
	or   Stmt
}

type WhileStmt struct {
	cond Expr
	body Stmt
}

type BreakStmt struct{}

type FunStmt struct {
	name Token
	args []Token
	body []Stmt
}

type ReturnStmt struct {
	value Expr
}

func (i *IfStmt) stmt()     {}
func (f *FunStmt) stmt()    {}
func (e *EnvStmt) stmt()    {}
func (vs *VarStmt) stmt()   {}
func (es *ExprStmt) stmt()  {}
func (p *PrintStmt) stmt()  {}
func (b *BlockStmt) stmt()  {}
func (w *WhileStmt) stmt()  {}
func (b *BreakStmt) stmt()  {}
func (r *ReturnStmt) stmt() {}

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

func (b *BreakStmt) accept(v StmtVisitor) {
	v.visitBreakStmt(b)
}

func (f *FunStmt) accept(v StmtVisitor) {
	v.visitFunStmt(f)
}

func (r *ReturnStmt) accept(v StmtVisitor) {
	v.visitReturnStmt(r)
}
