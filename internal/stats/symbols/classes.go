package symbols

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Classes struct {
	m sync.Mutex

	Classes map[string]*Class
}

func NewClasses() *Classes {
	return &Classes{
		Classes: map[string]*Class{},
	}
}

func (c *Classes) GetAllInterfaces(count int64, offset int64, sorted bool) []*Class {
	return c.GetAll(true, count, offset, sorted)
}

func (c *Classes) GetAllClasses(count int64, offset int64, sorted bool) []*Class {
	return c.GetAll(false, count, offset, sorted)
}

func (c *Classes) GetAll(onlyInterface bool, count int64, offset int64, sorted bool) []*Class {
	var res []*Class
	var index int64

	if offset < 0 {
		offset = 0
	}

	for _, class := range c.Classes {
		if !sorted {
			if index+offset > count && count != -1 {
				break
			}
		}

		if onlyInterface {
			if !class.IsInterface {
				continue
			}
		}

		res = append(res, class)
		index++
	}

	if sorted {
		sort.Slice(res, func(i, j int) bool {
			if res[i].Deps.Len() == res[j].Deps.Len() {
				return res[i].Name < res[j].Name
			}

			return res[i].Deps.Len() > res[j].Deps.Len()
		})

		if count != -1 {
			if count+offset < int64(len(res)) {
				res = res[:count+offset]
			}

			if offset < int64(len(res)) {
				res = res[offset:]
			}
		}
	}

	return res
}

func (c *Classes) Len() int {
	return len(c.Classes)
}

func (c *Classes) GetFullClassName(name string) ([]string, error) {
	var res []string

	for _, class := range c.Classes {
		if class.Name == name {
			return []string{class.Name}, nil
		}

		if strings.Contains(class.Name, name) {
			res = append(res, class.Name)
		}
	}

	if len(res) == 0 {
		return res, fmt.Errorf("class %s not found", name)
	}

	return res, nil
}

func (c *Classes) Add(class *Class) {
	c.m.Lock()
	c.Classes[class.Name] = class
	c.m.Unlock()
}

func (c *Classes) Get(name string) (*Class, bool) {
	c.m.Lock()
	class, ok := c.Classes[name]
	c.m.Unlock()
	return class, ok
}

type Class struct {
	Name string
	File *File

	Implements *Classes
	Extends    *Classes

	IsAbstract  bool
	IsInterface bool

	Fields    *Fields
	Methods   *Functions
	Constants *Constants

	// Зависим от
	Deps *Classes

	// Зависят от нас
	DepsBy *Classes

	Vendor bool

	// metrics
	LcomResolved bool
	Lcom         float64

	Lcom4Resolved bool
	Lcom4         int64
}

func NewClass(name string, file *File) *Class {
	return &Class{
		Name:       name,
		File:       file,
		Methods:    NewFunctions(),
		Fields:     NewFields(),
		Constants:  NewConstants(),
		Implements: NewClasses(),
		Extends:    NewClasses(),
		Deps:       NewClasses(),
		DepsBy:     NewClasses(),
	}
}

func NewInterface(name string, file *File) *Class {
	class := NewClass(name, file)
	class.IsInterface = true

	return class
}

func NewAbstractClass(name string, file *File) *Class {
	class := NewClass(name, file)
	class.IsAbstract = true

	return class
}

func (c *Class) AddMethod(fn *Function) {
	c.Methods.Add(fn)
}

func (c *Class) AddImplements(class *Class) {
	c.Implements.Add(class)
}

func (c *Class) AddExtends(class *Class) {
	c.Extends.Add(class)
}

func (c *Class) AddDeps(class *Class) {
	if c == class {
		return
	}

	c.Deps.Add(class)
}

func (c *Class) AddDepsBy(class *Class) {
	if c == class {
		return
	}

	c.DepsBy.Add(class)
}

// GobEncode is a custom gob marshaller
func (c *Class) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(c.Name)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(c.File)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(c.IsAbstract)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(c.IsInterface)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(c.Vendor)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// GobDecode is custom gob unmarshaller
func (c *Class) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&c.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.File)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.IsAbstract)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.IsInterface)
	if err != nil {
		return err
	}
	err = decoder.Decode(&c.Vendor)
	if err != nil {
		return err
	}
	return nil
}
