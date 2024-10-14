package main

type Scope map[string]bool

type Resolver struct {
	interpreter *Interpreter
	scopes      Stack[Scope]
}

type ResolveError struct {
	msg string
}

func (r *ResolveError) Error() string {
	return r.msg
}

func newResolver(i *Interpreter) *Resolver {
	return &Resolver{i, NewStack[Scope]()}
}

func (r *Resolver) resolveStmts(stmts ...Stmt) error {
	for _, stmt := range stmts {
		stmt.accept(r)
	}

	return nil
}

func (r *Resolver) resolveExprs(exprs ...Expr) error {
	for _, expr := range exprs {
		expr.accept(r)
	}

	return nil
}

func (r *Resolver) beginScope() {
	r.scopes.Push(map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(token Token) {
	if !r.scopes.Empty() {
		r.scopes.Peek()[token.lexeme] = false
	}
}

func (r *Resolver) define(token Token) {
	if !r.scopes.Empty() {
		r.scopes.Peek()[token.lexeme] = true
	}
}

func (r *Resolver) visitBlockStmt(b *BlockStmt) {
	r.beginScope()
	r.resolveStmts(b.stmts...)
	r.endScope()
}

func (r *Resolver) visitVarStmt(v *VarStmt) {
	r.declare(v.name)
	if v.initializer != nil {
		r.resolveExprs(v.initializer)
	}
	r.define(v.name)
}

func (r *Resolver) visitVariable(v *Variable) any {
	if !r.scopes.Empty() {
		defined, declared := r.scopes.Peek()[v.name.lexeme]

		if declared && !defined {
			panic(&ResolveError{"Can't read local variable in its own initializer."})
		}
	}

	r.resolveLocal(v, v.name)
	return nil
}

func (r *Resolver) visitAssignment(a *Assign) any {
	r.resolveExprs(a.value)
	r.resolveLocal(a, a.variable)
	return nil
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		if _, exists := r.scopes.At(i)[name.lexeme]; exists {
			r.interpreter.resolve(expr, r.scopes.Size()-1-i)
			return
		}
	}
}

func (r *Resolver) visitFunStmt(fun *FunStmt) {
	r.declare(fun.name)
	r.define(fun.name)
	r.resolveFun(fun)
}

func (r *Resolver) resolveFun(fun *FunStmt) {
	r.beginScope()
	for _, arg := range fun.args {
		r.declare(arg)
		r.define(arg)
	}
	r.resolveStmts(fun.body...)
	r.endScope()
}

func (r *Resolver) visitExprStmt(es *ExprStmt) {
	r.resolveExprs(es.expr)
}

func (r *Resolver) visitBreakStmt(b *BreakStmt) {}
func (r *Resolver) visitEnvStmt(b *EnvStmt)     {}
func (r *Resolver) visitIfStmt(ifs *IfStmt) {
	r.resolveExprs(ifs.cond)
	r.resolveStmts(ifs.then)
	if ifs.or != nil {
		r.resolveStmts(ifs.or)
	}
}

func (r *Resolver) visitPrintStmt(p *PrintStmt) {
	r.resolveExprs(p.val)
}

func (r *Resolver) visitReturnStmt(ret *ReturnStmt) {
	if ret.value != nil {
		r.resolveExprs(ret.value)
	}
}

func (r *Resolver) visitWhileStmt(w *WhileStmt) {
	r.resolveExprs(w.cond)
	r.resolveStmts(w.body)
}

func (r *Resolver) visitBinary(b *Binary) any {
	r.resolveExprs(b.left)
	r.resolveExprs(b.right)
	return nil
}

func (r *Resolver) visitCall(c *Call) any {
	r.resolveExprs(c.callee)
	for _, arg := range c.args {
		r.resolveExprs(arg)
	}
	return nil
}

func (r *Resolver) visitGrouping(g *Grouping) any {
	r.resolveExprs(g.expression)
	return nil
}

func (r *Resolver) visitLambda(l *Lambda) any {
	r.beginScope()
	for _, arg := range l.args {
		r.declare(arg)
		r.define(arg)
	}
	r.resolveStmts(l.body...)
	r.endScope()
	return nil
}

func (r *Resolver) visitLiteral(l *Literal) any {
	return nil
}

func (r *Resolver) visitLogical(l *Logical) any {
	r.resolveExprs(l.left)
	r.resolveExprs(l.right)
	return nil
}

func (r *Resolver) visitUnary(u *Unary) any {
	r.resolveExprs(u.right)
	return nil
}
