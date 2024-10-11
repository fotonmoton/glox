package main

type Callable interface {
	arity() int
	call(i *Interpreter, args ...any) (ret any)
}

type Function struct {
	definition *FunStmt
	closure    *Environment
}

func (f *Function) call(i *Interpreter, args ...any) (ret any) {

	defer func() {
		if err := recover(); err != nil {
			re, ok := err.(Return)

			if !ok {
				panic(err)
			}

			ret = re.val
		}
	}()

	env := newEnvironment(f.closure)

	for idx, arg := range f.definition.args {
		env.define(arg.lexeme, args[idx])
	}

	i.executeBlock(f.definition.body, env)

	return nil
}

func (f *Function) arity() int {
	return len(f.definition.args)
}

func newFunction(fun *FunStmt, env *Environment) Callable {
	return &Function{fun, env}
}
