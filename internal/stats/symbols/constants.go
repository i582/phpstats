package symbols

import (
	"sync"

	"github.com/i582/phpstats/internal/utils"
)

type Constant struct {
	Name  string
	Class *Class

	Used *Functions
}

func NewConstant(name string, class *Class) *Constant {
	return &Constant{
		Name:  name,
		Class: class,
		Used:  NewFunctions(),
	}
}

func NewConstantKey(name string, class *Class) Constant {
	return Constant{
		Name:  name,
		Class: class,
	}
}

func (c *Constant) IsSuperGlobal() bool {
	return utils.IsSuperGlobal(c.Name)
}

func (c *Constant) IsEmbedded() bool {
	return utils.IsEmbeddedConstant(c.Name)
}

func (c *Constant) String() string {
	if c.Class == nil {
		return c.Name
	}

	return c.Class.Name + "::" + c.Name
}

type Constants struct {
	m sync.Mutex

	Constants map[string]*Constant
}

func NewConstants() *Constants {
	return &Constants{
		Constants: map[string]*Constant{},
	}
}

func (c *Constants) Len() int {
	return len(c.Constants)
}

func (c *Constants) Add(constant *Constant) *Constant {
	c.m.Lock()
	c.Constants[constant.String()] = constant
	c.m.Unlock()
	return constant
}

func (c *Constants) Get(constantKey Constant) (*Constant, bool) {
	c.m.Lock()
	constant, ok := c.Constants[constantKey.String()]
	c.m.Unlock()
	return constant, ok
}

func (c *Constants) AddMethodAccess(constantKey Constant, method *Function) {
	constant, found := c.Get(constantKey)
	if !found {
		c.Add(NewConstant(constantKey.Name, constantKey.Class))
		constant, _ = c.Get(constantKey)
	}

	if method.Class != nil {
		method.Class.AddDeps(constant.Class)
	}
	method.UsedConstants.Add(constant)
	constant.Used.Add(method)

	if constant.Class != nil {
		constant.Class.Constants.Add(constant)
	}
}
