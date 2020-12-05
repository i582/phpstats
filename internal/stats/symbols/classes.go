package symbols

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func (c *Classes) GetInterfaceByPartOfName(name string) (*Class, error) {
	ifaceNames, err := c.GetFullClassName(name)
	if err != nil {
		return nil, fmt.Errorf("interface %s not found", name)
	}

	var hasOneClass bool

	for _, ifaceName := range ifaceNames {
		iface, found := c.Get(ifaceName)
		if found && iface.IsInterface {
			return iface, nil
		}
		hasOneClass = true
	}

	if hasOneClass {
		return nil, fmt.Errorf("interface %s was not found, but a class with the same name was found, \n"+
			"use the 'info class %s' command to get information about class", name, name)
	}

	return nil, fmt.Errorf("interface %s not found", name)
}

func (c *Classes) GetClassByPartOfName(name string) (*Class, error) {
	classNames, err := c.GetFullClassName(name)
	if err != nil {
		return nil, fmt.Errorf("class %s not found", name)
	}

	var hasOneInterface bool

	for _, className := range classNames {
		class, found := c.Get(className)
		if found && !class.IsInterface {
			return class, nil
		}
		hasOneInterface = true
	}

	if hasOneInterface {
		return nil, fmt.Errorf("class %s was not found, but a interface with the same name was found, \n"+
			"use the 'info iface %s' command to get information about interface", name, name)
	}

	return nil, fmt.Errorf("class %s not found", name)
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

func (c *Class) ClassName() string {
	parts := strings.Split(c.Name, `\`)
	if len(parts) == 1 {
		return c.Name
	}

	return parts[len(parts)-1]
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
