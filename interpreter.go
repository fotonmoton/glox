package main

import (
	"fmt"
	"reflect"
)

type Interpreter struct{}

type RuntimeError struct {
	token Token
	msg   string
}

func (re RuntimeError) Error() string {
	return re.msg
}

func newInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) evaluate(e Expr) any {
	defer i.recover()
	return e.accept(i)
}

func (i *Interpreter) recover() {
	if err := recover(); err != nil {
		pe, ok := err.(RuntimeError)

		if !ok {
			panic(err)
		}

		reportRuntimeError(pe.token, pe.msg)
		hadRuntimeError = true
	}
}

func (i *Interpreter) visitBinary(b *Binary) any {
	left := i.evaluate(b.left)
	right := i.evaluate(b.right)

	switch b.op.typ {
	case MINUS:
		checkIfFloats(b.op, left, right)
		return left.(float64) - right.(float64)
	case SLASH:
		checkIfFloats(b.op, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		checkIfFloats(b.op, left, right)
		return left.(float64) * right.(float64)
	case GREATER:
		checkIfFloats(b.op, left, right)
		return left.(float64) > right.(float64)
	case LESS:
		checkIfFloats(b.op, left, right)
		return left.(float64) < right.(float64)
	case GREATER_EQUAL:
		checkIfFloats(b.op, left, right)
		return left.(float64) >= right.(float64)
	case LESS_EQUAL:
		checkIfFloats(b.op, left, right)
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

	panic(RuntimeError{b.op, fmt.Sprintf("Operands must be numbers or strings: %v %s %v", left, b.op.lexeme, right)})

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
		checkIfFloat(u.op, val)
		return -val.(float64)
	case BANG:
		return !isTruthy(val)
	}

	return nil
}

func checkIfFloat(op Token, val any) {
	if _, ok := val.(float64); ok {
		return
	}

	panic(RuntimeError{op, "value must ne a number."})
}

func checkIfFloats(op Token, a any, b any) {
	if isFloats(a, b) {
		return
	}

	panic(RuntimeError{op, fmt.Sprintf("Operands must be numbers: %v %s %v", a, op.lexeme, b)})
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
