package main

type Class struct {
	name string
}

func (c *Class) String() string {
	return c.name
}
