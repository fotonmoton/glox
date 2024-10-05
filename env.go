package main

type Environment struct {
	values map[string]any
	parent *Environment
}

func newEnvironment(parent *Environment) *Environment {
	return &Environment{values: map[string]any{}, parent: parent}
}

func (env *Environment) get(key string) any {
	if found, ok := env.values[key]; ok {
		return found
	}

	if env.parent != nil {
		return env.parent.get(key)
	}

	return nil
}

func (env *Environment) exists(key string) bool {
	_, ok := env.values[key]

	if !ok && env.parent != nil {
		return env.parent.exists(key)
	}

	return ok

}

func (env *Environment) set(key string, val any) {
	env.values[key] = val
}
