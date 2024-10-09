package main

type Callable struct {
	arity int
	call  func(*Interpreter, ...any) any
}

func newCallable(f *FunStmt) *Callable {
	return &Callable{
		arity: len(f.args),
		call: func(i *Interpreter, args ...any) any {
			env := newEnvironment(i.globals)

			for idx, arg := range f.args {
				env.set(arg.lexeme, args[idx])
			}

			i.executeBlock(f.body, env)

			return nil
		},
	}
}
