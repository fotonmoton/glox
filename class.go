package main

import "fmt"

type Class struct {
	name string
}

type ClassInstance struct {
	klass *Class
	props map[string]any
}

func (c *ClassInstance) String() string {
	return fmt.Sprintf("instance of %s", c.klass.name)
}

func (c *ClassInstance) get(name string) (any, bool) {
	val, ok := c.props[name]
	return val, ok
}

func (c *ClassInstance) set(name string, val any) {
	c.props[name] = val
}

func (c *Class) arity() int {
	return 0
}

func (c *Class) call(i *Interpreter, args ...any) (ret any) {
	return &ClassInstance{c, map[string]any{}}
}

func (c *Class) String() string {
	return c.name
}
