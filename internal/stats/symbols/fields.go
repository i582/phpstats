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
	Class *Class

	Used *Functions

	Id int64
}

func NewField(name string, class *Class) *Field {
	atomic.AddInt64(&FieldsCount, 1)
	return &Field{
		Name:  name,
		Class: class,
		Used:  NewFunctions(),
		Id:    FieldsCount,
	}
}

func (f *Field) ID() int64 {
	return f.Id
}

func (f *Field) String() string {
	return f.Class.Name + "::" + f.Name
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
	c.Fields[NewFieldKey(field.Name, field.Class.Name)] = field
	c.m.Unlock()
}

func (c *Fields) Get(key FieldKey) (*Field, bool) {
	c.m.Lock()
	field, ok := c.Fields[key]
	c.m.Unlock()
	return field, ok
}

func (c *Fields) AddMethodAccess(key FieldKey, class *Class, method *Function) {
	field, found := c.Get(key)
	if !found {
		c.Add(NewField(key.Name, class))
		field, _ = c.Get(key)
	}

	if method.Class != nil {
		method.Class.AddDeps(field.Class)
	}
	method.UsedFields.Add(field)

	field.Used.Add(method)
}
