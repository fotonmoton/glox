package main

import "fmt"

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func newEnvironment(enclosing *Environment) *Environment {
	return &Environment{map[string]any{}, enclosing}
}

func (env *Environment) get(key string) any {
	if found, ok := env.values[key]; ok {
		return found
	}

	if env.enclosing != nil {
		return env.enclosing.get(key)
	}

	return nil
}

func (env *Environment) exists(key string) bool {
	_, ok := env.values[key]
	return ok
}

func (env *Environment) define(key string, val any) {
	env.values[key] = val
}

func (env *Environment) assign(key Token, val any) *RuntimeError {
	if env.exists(key.lexeme) {
		env.values[key.lexeme] = val
		return nil
	}

	if env.enclosing == nil {
		return &RuntimeError{key, fmt.Sprintf("Can't assign: undefined variable '%s'.", key.lexeme)}
	}

	return env.enclosing.assign(key, val)
}

func (env *Environment) getAt(distance int, key string) any {
	return env.ancestor(distance).get(key)
}

func (env *Environment) assignAt(distance int, key Token, val any) {
	env.ancestor(distance).values[key.lexeme] = val
}

func (env *Environment) ancestor(distance int) *Environment {
	parent := env
	for i := 0; i < distance; i++ {
		parent = parent.enclosing
	}

	return parent
}
