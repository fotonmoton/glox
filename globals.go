package main

import "time"

func defineGlobals(env *Environment) {

	env.set("clock", &Callable{
		arity: 0,
		call: func(i *Interpreter, arg ...any) any {
			return time.Now().Unix()
		},
	})
}