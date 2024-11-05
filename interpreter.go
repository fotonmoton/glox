package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
)

type Interpreter struct {
	env     *Environment
	globals *Environment
	locals  map[Expr]int
	errors  []error
	brk     bool
}

type RuntimeError struct {
	token Token
	msg   string
}

type Return struct {
	val any
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError [%d][%s] Error: %s", re.token.line, re.token.typ, re.msg)
}

func newInterpreter() *Interpreter {

	globals := newEnvironment(nil)

	defineGlobals(globals)

	return &Interpreter{
		env:     globals,
		globals: globals,
		locals:  map[Expr]int{},
		errors:  []error{},
		brk:     false,
	}
}

func (i *Interpreter) interpret(stmts []Stmt) error {
	defer i.recover()

	i.errors = []error{}

	for _, stmt := range stmts {
		stmt.accept(i)
	}

	return errors.Join(i.errors...)
}

func (i *Interpreter) recover() {
	if err := recover(); err != nil {
		_, ok := err.(*RuntimeError)

		if !ok {
			panic(err)
		}
	}
}

func (i *Interpreter) evaluate(e Expr) any {
	return e.accept(i)
}

func (i *Interpreter) visitBinary(b *Binary) any {
	left := i.evaluate(b.left)
	right := i.evaluate(b.right)

	switch b.op.typ {
	case MINUS:
		i.checkIfFloats(b.op, left, right)
		return left.(float64) - right.(float64)
	case SLASH:
		i.checkIfFloats(b.op, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		i.checkIfFloats(b.op, left, right)
		return left.(float64) * right.(float64)
	case GREATER:
		i.checkIfFloats(b.op, left, right)
		return left.(float64) > right.(float64)
	case LESS:
		i.checkIfFloats(b.op, left, right)
		return left.(float64) < right.(float64)
	case GREATER_EQUAL:
		i.checkIfFloats(b.op, left, right)
		return left.(float64) >= right.(float64)
	case LESS_EQUAL:
		i.checkIfFloats(b.op, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !reflect.DeepEqual(left, right)
	case EQUAL_EQUAL:
		return reflect.DeepEqual(left, right)
	case PLUS:
		if isFloats(left, right) {
			return left.(float64) + right.(float64)
		}

		if isStrings(left, right) {
			return left.(string) + right.(string)
		}
	}

	i.panic(&RuntimeError{b.op, fmt.Sprintf("Operands must be numbers or strings: %v %s %v", left, b.op.lexeme, right)})

	return nil
}

func (i *Interpreter) visitLiteral(l *Literal) any {
	return l.value
}

func (i *Interpreter) visitGrouping(g *Grouping) any {
	return i.evaluate(g.expression)
}

func (i *Interpreter) visitUnary(u *Unary) any {
	val := i.evaluate(u.right)

	switch u.op.typ {
	case MINUS:
		i.checkIfFloat(u.op, val)
		return -val.(float64)
	case BANG:
		return !isTruthy(val)
	}

	return nil
}

func (i *Interpreter) visitVariable(v *Variable) any {
	return i.lookUpVariable(v.name, v)
}

func (i *Interpreter) visitAssignment(a *Assign) any {
	val := i.evaluate(a.value)
	distance, isLocal := i.locals[a]

	if isLocal {
		i.env.assignAt(distance, a.variable, val)
		return val
	}

	err := i.globals.assign(a.variable, val)
	if err != nil {
		i.panic(err)
	}
	return val
}

func (i *Interpreter) visitLogical(lo *Logical) any {

	left := i.evaluate(lo.left)

	shortOr := lo.operator.typ == OR && isTruthy(left)
	shortAnd := lo.operator.typ == AND && !isTruthy(left)

	if shortOr || shortAnd {
		return left
	}

	return i.evaluate(lo.right)
}

func (i *Interpreter) visitCall(c *Call) any {

	callee := i.evaluate(c.callee)

	args := []any{}

	for _, arg := range c.args {
		args = append(args, i.evaluate(arg))
	}

	callable, ok := callee.(Callable)

	if !ok {
		i.panic(&RuntimeError{c.paren, "Can only call function and classes."})
	}

	if callable.arity() != len(args) {
		i.panic(&RuntimeError{
			c.paren,
			fmt.Sprintf(
				"Expected %d arguments  but got %d",
				callable.arity(),
				len(args),
			),
		})
	}

	return callable.call(i, args...)
}

func (i *Interpreter) visitFunStmt(f *FunStmt) {
	i.env.define(f.name.lexeme, newFunction(f.name, f.args, f.body, i.env))
}

func (i *Interpreter) visitClassStmt(c *ClassStmt) {
	i.env.define(c.name.lexeme, nil)
	class := &Class{c.name.lexeme}
	i.env.assign(c.name, class)

}
func (i *Interpreter) visitLambda(l *Lambda) any {
	return newFunction(l.name, l.args, l.body, i.env)
}

func (i *Interpreter) visitReturnStmt(r *ReturnStmt) {
	var value any

	if r.value != nil {
		value = i.evaluate(r.value)
	}

	panic(Return{value})
}

func (i *Interpreter) visitPrintStmt(p *PrintStmt) {
	fmt.Printf("%v\n", i.evaluate(p.val))
}

func (i *Interpreter) visitExprStmt(se *ExprStmt) {
	i.evaluate(se.expr)
}

func (i *Interpreter) visitVarStmt(v *VarStmt) {

	var val any = nil

	if v.initializer != nil {
		val = i.evaluate(v.initializer)
	}

	i.env.define(v.name.lexeme, val)
}

func (i *Interpreter) visitBlockStmt(b *BlockStmt) {
	i.executeBlock(b.stmts, newEnvironment(i.env))
}

func (i *Interpreter) executeBlock(stmts []Stmt, current *Environment) {

	parentEnv := i.env
	i.env = current

	// need to restore environment after
	// panic(Return) in visitReturnStmt
	defer func() {
		i.env = parentEnv
	}()

	for _, stmt := range stmts {

		if i.brk {
			break
		}

		stmt.accept(i)
	}

}

func (i *Interpreter) visitBreakStmt(b *BreakStmt) {
	i.brk = true
}

func (i *Interpreter) visitIfStmt(iff *IfStmt) {
	if isTruthy(i.evaluate(iff.cond)) {
		iff.then.accept(i)

	} else if iff.or != nil {
		iff.or.accept(i)
	}
}

func (i *Interpreter) visitEnvStmt(e *EnvStmt) {

	walker := i.env

	flatten := []*Environment{}

	for walker != nil {
		flatten = slices.Insert(flatten, 0, walker)
		walker = walker.enclosing
	}

	fmt.Printf("globals: %+v\n", *i.globals)

	for ident, e := range flatten {
		fmt.Printf("%*s", ident, "")
		fmt.Printf("%+v\n", *e)
	}
}

func (i *Interpreter) visitWhileStmt(w *WhileStmt) {
	for isTruthy(i.evaluate(w.cond)) {

		if i.brk {
			i.brk = false
			break
		}

		w.body.accept(i)
	}
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) lookUpVariable(name Token, expr Expr) any {
	distance, isLocal := i.locals[expr]

	if !isLocal {
		return i.globals.get(name.lexeme)
	}

	return i.env.getAt(distance, name.lexeme)
}

func (i *Interpreter) panic(re *RuntimeError) {
	i.errors = append(i.errors, re)
	log.Println(re)
	panic(re)
}

func (i *Interpreter) checkIfFloat(op Token, val any) {
	if _, ok := val.(float64); ok {
		return
	}

	i.panic(&RuntimeError{op, "value must be a number."})
}

func (i *Interpreter) checkIfFloats(op Token, a any, b any) {
	if isFloats(a, b) {
		return
	}

	i.panic(&RuntimeError{op, fmt.Sprintf("Operands must be numbers: %v %s %v", a, op.lexeme, b)})
}

func isFloats(a any, b any) bool {

	if a == nil || b == nil {
		return false
	}

	ltype := reflect.TypeOf(a)
	rtype := reflect.TypeOf(b)

	return ltype.Kind() == rtype.Kind() && ltype.Kind() == reflect.Float64
}

func isStrings(a any, b any) bool {

	if a == nil || b == nil {
		return false
	}

	ltype := reflect.TypeOf(a)
	rtype := reflect.TypeOf(b)

	return ltype.Kind() == rtype.Kind() && ltype.Kind() == reflect.String
}

func isTruthy(val any) bool {
	if val == nil {
		return false
	}

	if boolean, ok := val.(bool); ok {
		return boolean
	}

	return true
}
