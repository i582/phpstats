package symbols

import (
	"sync"
)

type Constant struct {
	Name  string
	Class string
}

func NewConstant(name, class string) *Constant {
	return &Constant{
		Name:  name,
		Class: class,
	}
}

func (c *Constant) String() string {
	return c.Class + "::" + c.Name
}

type Constants struct {
	m sync.Mutex

	Constants map[Constant]*Constant
}

func NewConstants() *Constants {
	return &Constants{
		Constants: map[Constant]*Constant{},
	}
}

func (c *Constants) Len() int {
	return len(c.Constants)
}

func (c *Constants) Add(constant *Constant) {
	c.m.Lock()
	c.Constants[*constant] = constant
	c.m.Unlock()
}

func (c *Constants) Get(constantKey Constant) (*Constant, bool) {
	c.m.Lock()
	constant, ok := c.Constants[constantKey]
	c.m.Unlock()
	return constant, ok
}
