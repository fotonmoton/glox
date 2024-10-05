package main

import (
	"fmt"
	"log"
	"reflect"
)

type Interpreter struct {
	env    *Environment
	errors []error
}

type RuntimeError struct {
	token Token
	msg   string
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError [%d][%s] Error: %s", re.token.line, re.token.typ, re.msg)
}

func newInterpreter() *Interpreter {
	return &Interpreter{env: newEnvironment(nil)}
}

func (i *Interpreter) interpret(stmts []Stmt) []error {
	defer i.recover()

	i.errors = []error{}

	for _, stmt := range stmts {
		stmt.accept(i)
	}

	return i.errors
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

	val := i.env.get(v.name.lexeme)

	if val == nil {
		i.panic(&RuntimeError{v.name, fmt.Sprintf("Can't evaluate: Undefined variable '%s'.", v.name.lexeme)})
		return nil
	}

	return val
}

func (i *Interpreter) visitAssignment(a *Assign) any {

	if !i.env.exists(a.variable.lexeme) {
		i.panic(&RuntimeError{a.variable, fmt.Sprintf("Can't assign: undefined variable '%s'.", a.variable.lexeme)})

		return nil
	}

	val := i.evaluate(a.value)

	i.env.set(a.variable.lexeme, val)

	return val
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

	i.env.set(v.name.lexeme, val)
}

func (i *Interpreter) visitBlockStmt(b *BlockStmt) {

	parentEnv := i.env
	i.env = newEnvironment(parentEnv)

	for _, stmt := range b.stmts {
		stmt.accept(i)
	}

	i.env = parentEnv
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
	ltype := reflect.TypeOf(a)
	rtype := reflect.TypeOf(b)

	return ltype.Kind() == rtype.Kind() && ltype.Kind() == reflect.Float64
}

func isStrings(a any, b any) bool {
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
