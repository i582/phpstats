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
		if class.IsVendor {
			continue
		}

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

func (c *Classes) CountClasses() int64 {
	var count int64
	for _, class := range c.Classes {
		if !class.IsInterface && !class.IsAbstract {
			count++
		}
	}
	return count
}

func (c *Classes) CountAbstractClasses() int64 {
	var count int64
	for _, class := range c.Classes {
		if class.IsAbstract {
			count++
		}
	}
	return count
}

func (c *Classes) CountIfaces() int64 {
	var count int64
	for _, class := range c.Classes {
		if class.IsInterface {
			count++
		}
	}
	return count
}

func (c *Classes) MaxMinAvgCyclomaticComplexity() (max, min, avg float64) {
	const maxValue = 10000000
	var count float64

	max = 0
	min = maxValue

	for _, class := range c.Classes {
		countMN := class.Methods.CyclomaticComplexity()
		count += float64(countMN)

		if count < min {
			min = count
		}

		if count > max {
			max = count
		}
	}

	if min == maxValue {
		min = 0
	}

	if c.Len() != 0 {
		avg = count / float64(c.Len())
	}

	return max, min, avg
}

func (c *Classes) MaxMinAvgCountMagicNumbers() (max, min, avg int64) {
	const maxValue = 10000000
	var count int64

	max = 0
	min = maxValue

	for _, class := range c.Classes {
		countMN := class.Methods.CountMagicNumbers()
		count += countMN

		if count < min {
			min = count
		}

		if count > max {
			max = count
		}
	}

	if min == maxValue {
		min = 0
	}

	if c.Len() != 0 {
		avg = count / int64(c.Len())
	}

	return max, min, avg
}

func (c *Classes) GetClassByPartOfName(name string) (*Class, error) {
	classes, err := c.GetFullClassName(name)
	if err != nil {
		return nil, err
	}

	class, found := c.Get(classes[0])
	if !found {
		return nil, fmt.Errorf("class %s not found", name)
	}
	return class, nil
}

func (c *Classes) GetFullClassName(name string) ([]string, error) {
	var res []string

	if !strings.HasPrefix(name, `\`) {
		name = `\` + name
	}

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
	if class == nil {
		return
	}

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

	Namespace *Namespace

	Implements *Classes
	Extends    *Classes

	ImplementsBy *Classes
	ExtendsBy    *Classes

	IsAbstract  bool
	IsInterface bool
	IsTrait     bool

	Fields    *Fields
	Methods   *Functions
	Constants *Constants

	UsedConstants *Constants

	// Зависим от
	Deps *Classes

	// Зависят от нас
	DepsBy *Classes

	IsVendor bool

	// metrics
	LcomResolved bool
	Lcom         float64

	Lcom4Resolved bool
	Lcom4         int64
}

func NewClass(name string, file *File) *Class {
	return &Class{
		Name:          name,
		File:          file,
		Methods:       NewFunctions(),
		Fields:        NewFields(),
		Constants:     NewConstants(),
		UsedConstants: NewConstants(),
		Implements:    NewClasses(),
		Extends:       NewClasses(),
		ImplementsBy:  NewClasses(),
		ExtendsBy:     NewClasses(),
		Deps:          NewClasses(),
		DepsBy:        NewClasses(),
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

func NewTrait(name string, file *File) *Class {
	class := NewClass(name, file)
	class.IsTrait = true

	return class
}

func (c *Class) AddMethod(fn *Function) {
	c.Methods.Add(fn)
}

func (c *Class) AddImplements(class *Class) {
	c.Implements.Add(class)
	class.ImplementsBy.Add(c)
}

func (c *Class) AddExtends(class *Class) {
	c.Extends.Add(class)
	class.ExtendsBy.Add(c)
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

func (c *Class) Type() string {
	if c.IsInterface {
		return "interface"
	}
	if c.IsAbstract {
		return "abstract class"
	}
	if c.IsTrait {
		return "trait"
	}

	return "class"
}

func (c *Class) NamespaceName() string {
	parts := strings.Split(c.Name, `\`)
	if len(parts) == 1 {
		return c.Name
	}

	return strings.Join(parts[0:len(parts)-1], `\`)
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
	err = encoder.Encode(c.IsVendor)
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
	err = decoder.Decode(&c.IsVendor)
	if err != nil {
		return err
	}
	return nil
}
