package main

type Callable struct {
	arity int
	call  func(*Interpreter, ...any) any
}

func newCallable(f *FunStmt) *Callable {
	return &Callable{
		arity: len(f.args),
		call: func(i *Interpreter, args ...any) (ret any) {

			defer func() {
				if err := recover(); err != nil {
					re, ok := err.(Return)

					if !ok {
						panic(err)
					}

					ret = re.val
				}
			}()

			env := newEnvironment(i.globals)

			for idx, arg := range f.args {
				env.define(arg.lexeme, args[idx])
			}

			i.executeBlock(f.body, env)

			return nil
		},
	}
}
