package stats

import (
	"bytes"
	"encoding/gob"
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
	sync.Mutex

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
	c.Lock()
	defer c.Unlock()
	c.Constants[*constant] = constant
}

func (c *Constants) Get(constantKey Constant) (*Constant, bool) {
	c.Lock()
	defer c.Unlock()
	constant, ok := c.Constants[constantKey]
	return constant, ok
}

// GobDecode is custom gob unmarshaller
func (c *Constants) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&c.Constants)
	if err != nil {
		return err
	}
	return nil
}

// GobEncode is a custom gob marshaller
func (c *Constants) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(c.Constants)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
