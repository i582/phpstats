package stats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/VKCOM/noverify/src/meta"
)

type Classes struct {
	sync.Mutex

	Classes map[string]*Class
}

func NewClasses() *Classes {
	return &Classes{
		Classes: map[string]*Class{},
	}
}

var AlreadyShown = map[string]struct{}{}

func (c *Classes) Len() int {
	return len(c.Classes)
}

func (c *Classes) GetUsedClassesInClass(name string) (*Classes, bool) {
	class, ok := c.Get(name)
	if !ok {
		return nil, false
	}

	res := NewClasses()

	for _, method := range class.Methods.Funcs {
		for _, called := range method.Called.Funcs {
			if called.Class == nil {
				continue
			}
			if called.Class.File == method.Class.File {
				continue
			}

			res.Add(called.Class)
		}
	}

	return res, true
}

func (c *Classes) CalculateClassDeps() {
	for _, class := range c.Classes {
		for _, method := range class.Methods.Funcs {
			for _, called := range method.Called.Funcs {
				if called.Class != nil && method.Class != nil && called.Class.ShortString(0) != method.Class.ShortString(0) {
					class.Deps.Add(called.Class)
					called.Class.DepsBy.Add(class)
				}
			}
		}
	}
}

func (c *Classes) Graphviz() string {
	var res string
	res += "digraph test{\n"

	var count int

Outer:
	for _, class := range c.Classes {
		// if !strings.Contains(class.File.Path, `VK\API`) {
		// 	continue
		// }
		for _, method := range class.Methods.Funcs {
			for _, called := range method.Called.Funcs {

				if called.Class != nil && method.Class != nil && called.Class.File != method.Class.File {

					str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", method.Class.Name, called.Class.Name)
					if _, ok := AlreadyShown[str]; ok {
						continue
					}

					res += str
					AlreadyShown[str] = struct{}{}
					count++

					if count > 1000000 {
						break Outer
					}
				}

			}
		}
	}
	res += "}"
	return res
}

func (c *Classes) Add(class *Class) {
	c.Lock()
	defer c.Unlock()
	c.Classes[class.Name] = class
}

func (c *Classes) Get(name string) (*Class, bool) {
	c.Lock()
	defer c.Unlock()
	class, ok := c.Classes[name]
	return class, ok
}

type Class struct {
	Name string
	File *File

	Implements *Classes
	Extends    *Classes

	IsAbstract  bool
	IsInterface bool

	Methods *Functions

	// Зависим от
	Deps *Classes

	// Зависят от нас
	DepsBy *Classes
}

func NewClass(name string, file *File) *Class {
	return &Class{
		Name:       name,
		File:       file,
		Methods:    NewFunctionsInfo(),
		Implements: NewClasses(),
		Extends:    NewClasses(),
		Deps:       NewClasses(),
		DepsBy:     NewClasses(),
	}
}

func NewInterface(name string, file *File) *Class {
	return &Class{
		Name:        name,
		File:        file,
		Methods:     NewFunctionsInfo(),
		Implements:  NewClasses(),
		Extends:     NewClasses(),
		Deps:        NewClasses(),
		DepsBy:      NewClasses(),
		IsInterface: true,
	}
}

func NewAbstractClass(name string, file *File) *Class {
	return &Class{
		Name:       name,
		File:       file,
		Methods:    NewFunctionsInfo(),
		Implements: NewClasses(),
		Extends:    NewClasses(),
		Deps:       NewClasses(),
		DepsBy:     NewClasses(),
		IsAbstract: true,
	}
}

func (c *Class) FullString(level int) string {
	var res string

	res += c.ShortString(level)

	efferent := float64(len(c.Deps.Classes))
	afferent := float64(len(c.DepsBy.Classes))

	var stability float64
	if efferent+afferent == 0 {
		stability = 0
	} else {
		stability = efferent / (efferent + afferent)
	}

	res += fmt.Sprintf(" Афферентность: %.2f\n", afferent)
	res += fmt.Sprintf(" Эфферентность: %.2f\n", efferent)
	res += fmt.Sprintf(" Стабильность:  %.2f\n", stability)

	if c.Implements.Len() != 0 {
		res += fmt.Sprintf(" Реализует:\n")
	}
	for _, class := range c.Implements.Classes {
		res += fmt.Sprintf("%s", class.ShortStringWithPrefix(level+1, " ↳ "))
	}
	// info class -f AppPost
	// info class -f PeriodToOneOf
	if c.Extends.Len() != 0 {
		res += fmt.Sprintf(" Расширяет:\n")
	}
	for _, class := range c.Extends.Classes {
		res += fmt.Sprintf("%s", class.ShortStringWithPrefix(level+1, " ↳ "))
	}

	if c.Methods.Len() != 0 {
		res += fmt.Sprintf(" Методы (%d):\n", c.Methods.Len())
	}
	for _, method := range c.Methods.Funcs {
		res += fmt.Sprintf("   %s\n", method.Name.Name)
	}

	return res
}

func (c *Class) ShortString(level int) string {
	return c.ShortStringWithPrefix(level, "")
}

func (c *Class) ShortStringWithPrefix(level int, prefix string) string {
	var res string

	if c.IsInterface {
		res += fmt.Sprintf("%s%sИнтерфейс %s\n", prefix, genIndent(level-1), c.Name)
	} else {
		if c.IsAbstract {
			res += fmt.Sprintf("%s%sАбстрактный класс %s\n", prefix, genIndent(level-1), c.Name)
		} else {
			res += fmt.Sprintf("%s%sКласс %s\n", prefix, genIndent(level-1), c.Name)
		}
	}
	res += fmt.Sprintf("%s Имя:  %s\n", genIndent(level), c.Name)
	res += fmt.Sprintf("%s Файл: %s:0\n", genIndent(level), c.File.Path)

	return res
}

func (c *Class) AddMethod(fn *Function) {
	_, ok := c.Methods.Get(fn.Name)
	if !ok {
		c.Methods.Add(fn)
	}
}

func (c *Class) AddImplements(class *Class) {
	_, ok := c.Implements.Get(class.Name)
	if !ok {
		c.Implements.Add(class)
	}
}

func (c *Class) AddExtends(class *Class) {
	_, ok := c.Extends.Get(class.Name)
	if !ok {
		c.Extends.Add(class)
	}
}

func (c *Class) AddDeps(class *Class) {
	_, ok := c.Deps.Get(class.Name)
	if !ok {
		c.Deps.Add(class)
	}
}

func (c *Class) AddDepsBy(class *Class) {
	_, ok := c.DepsBy.Get(class.Name)
	if !ok {
		c.DepsBy.Add(class)
	}
}

func (c *Class) GetOrCreateMethod(fn FuncKey, pos meta.ElementPosition) *Function {
	method := c.Methods.GetOrCreateMethod(fn, pos, c)

	c.Methods.Add(method)
	return method
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
	return nil
}

// GobDecode is custom gob unmarshaller
func (c *Classes) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&c.Classes)
	if err != nil {
		return err
	}
	return nil
}

// GobEncode is a custom gob marshaller
func (c *Classes) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(c.Classes)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
