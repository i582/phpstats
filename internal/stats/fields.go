package stats

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type FieldKey struct {
	Name  string
	Class string
}

func NewFieldKey(name, class string) FieldKey {
	return FieldKey{
		Name:  name,
		Class: class,
	}
}

type Field struct {
	Name  string
	Class string

	Used map[*Function]struct{}
}

func NewField(name, class string) *Field {
	return &Field{
		Name:  name,
		Class: class,
		Used:  map[*Function]struct{}{},
	}
}

func (c *Field) String() string {
	return c.Class + "::" + c.Name
}

type Fields struct {
	sync.Mutex

	Fields map[FieldKey]*Field
}

func NewFields() *Fields {
	return &Fields{
		Fields: map[FieldKey]*Field{},
	}
}

func (c *Fields) Len() int {
	return len(c.Fields)
}

func (c *Fields) Add(field *Field) {
	c.Lock()
	defer c.Unlock()
	c.Fields[NewFieldKey(field.Name, field.Class)] = field
}

func (c *Fields) Get(key FieldKey) (*Field, bool) {
	c.Lock()
	defer c.Unlock()
	field, ok := c.Fields[key]
	return field, ok
}

func (c *Fields) AddMethodAccess(key FieldKey, method *Function) {
	c.Lock()
	field, ok := c.Fields[key]
	c.Unlock()
	if ok {
		field.Used[method] = struct{}{}
	}
}

// GobDecode is custom gob unmarshaller
func (c *Fields) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&c.Fields)
	if err != nil {
		return err
	}
	return nil
}

// GobEncode is a custom gob marshaller
func (c *Fields) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(c.Fields)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
