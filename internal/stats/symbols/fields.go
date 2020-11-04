package symbols

import (
	"sync"
	"sync/atomic"
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

var FieldsCount int64

type Field struct {
	Name  string
	Class string

	Used map[*Function]struct{}

	Id int64
}

func NewField(name, class string) *Field {
	atomic.AddInt64(&FieldsCount, 1)
	return &Field{
		Name:  name,
		Class: class,
		Used:  map[*Function]struct{}{},
		Id:    FieldsCount,
	}
}

func (f *Field) ID() int64 {
	return f.Id
}

func (f *Field) String() string {
	return f.Class + "::" + f.Name
}

type Fields struct {
	m sync.Mutex

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
	c.m.Lock()
	c.Fields[NewFieldKey(field.Name, field.Class)] = field
	c.m.Unlock()
}

func (c *Fields) Get(key FieldKey) (*Field, bool) {
	c.m.Lock()
	field, ok := c.Fields[key]
	c.m.Unlock()
	return field, ok
}

func (c *Fields) AddMethodAccess(key FieldKey, method *Function) {
	c.m.Lock()
	field, ok := c.Fields[key]
	c.m.Unlock()
	if ok {
		field.Used[method] = struct{}{}
	}
}
