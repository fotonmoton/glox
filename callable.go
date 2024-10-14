package main

type Callable interface {
	arity() int
	call(i *Interpreter, args ...any) (ret any)
}

type Function struct {
	name    Token
	args    []Token
	body    []Stmt
	closure *Environment
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

	for idx, arg := range f.args {
		env.define(arg.lexeme, args[idx])
	}

	i.executeBlock(f.body, env)

	return nil
}

func (f *Function) arity() int {
	return len(f.args)
}

func newFunction(name Token, args []Token, body []Stmt, env *Environment) Callable {
	return &Function{name, args, body, env}
}
