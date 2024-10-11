package main

import "time"

type ClockFun struct{}

func (cf *ClockFun) call(i *Interpreter, args ...any) any {
	return time.Now().Unix()
}

func (cf *ClockFun) arity() int {
	return 0
}

func defineGlobals(env *Environment) {
	env.define("clock", &ClockFun{})
}
