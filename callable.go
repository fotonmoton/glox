package main

type Callable struct {
	arity int
	call  func(*Interpreter, ...any) any
}
