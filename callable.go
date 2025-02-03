package main

type Callable interface {
	arity() int
	call(i *Interpreter, args ...any) (ret any)
}
