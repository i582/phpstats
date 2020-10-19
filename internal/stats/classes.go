package stats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/VKCOM/noverify/src/meta"

	"github.com/i582/phpstats/internal/utils"
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
		if strings.Contains(class.Name, name) {
			res = append(res, class.Name)
		}
	}

	if len(res) == 0 {
		return res, fmt.Errorf("class %s not found", name)
	}

	return res, nil
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

func (c *Class) GraphvizRecursive(level int64, maxLevel int64, visited map[string]struct{}) string {
	var res string

	if level == 0 {
		res += "digraph test{\n"
	}

	if level > maxLevel {
		return ""
	}

	for _, class := range c.Deps.Classes {
		str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", c.Name, class.Name)

		if _, ok := visited[str]; ok {
			continue
		}
		visited[str] = struct{}{}

		res += str

		res += class.GraphvizRecursive(level+1, maxLevel, visited)
	}

	if level == 0 {
		res += "}\n"
	}

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

	Fields    *Fields
	Methods   *Functions
	Constants *Constants

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

func (c *Class) Lcom4Graph() string {
	var res string

	res += "digraph \"Lcom4" + utils.NameNormalize(c.Name) + "\" {\n"

	showed := map[string]struct{}{}

	for _, method := range c.Methods.Funcs {
		res += fmt.Sprintf("  \"%s\"\n", method.Name.Name)
	}

	for _, method := range c.Methods.Funcs {
		for _, called := range method.Called.Funcs {
			if _, ok := c.Methods.Get(called.Name); ok && method != called {
				str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", method.Name.Name, called.Name.Name)

				if _, ok := showed[str]; ok {
					continue
				}
				showed[str] = struct{}{}

				res += str
			}
		}
	}

	for _, field := range c.Fields.Fields {
		functions := make([]*Function, 0, len(field.Used))

		for used := range field.Used {
			functions = append(functions, used)
		}

		for i := 0; i < len(functions)-1; i++ {
			for j := i + 1; j < len(functions); j++ {
				str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", functions[i].Name.Name, functions[j].Name.Name)

				if _, ok := showed[str]; ok {
					continue
				}
				showed[str] = struct{}{}

				res += str
			}
		}
	}

	res += "}\n"

	return res
}

func (c *Class) AffEffString(full bool) string {
	var res string

	aff, eff, stab := AfferentEfferentStabilityOfClass(c)

	res += fmt.Sprintf(" Afferent:  %.2f\n", aff)
	if full {
		for _, class := range c.DepsBy.Classes {
			res += fmt.Sprintf("%s", class.ExtraShortString(2))
		}
	}

	res += fmt.Sprintf(" Efferent:  %.2f\n", eff)
	if full {
		for _, class := range c.Deps.Classes {
			res += fmt.Sprintf("%s", class.ExtraShortString(2))
		}
	}

	res += fmt.Sprintf(" Stability: %.2f\n", stab)

	lcom, ok := LackOfCohesionInMethodsOfCLass(c)
	if !ok {
		res += fmt.Sprintf(" LCOM: undefined (the number of methods or fields is zero)\n")
	} else {
		res += fmt.Sprintf(" LCOM: %.6f\n", lcom)
	}

	lcom4 := Lcom4(c)
	res += fmt.Sprintf(" LCOM4: %d\n", lcom4)

	return res
}

func (c *Class) FullString(level int, withAff bool) string {
	var res string

	res += c.ShortString(level)

	if withAff {
		res += c.AffEffString(false)
	}

	return res
}

func (c *Class) OnlyMetricsString() string {
	var res string

	res += c.ShortString(0)

	aff, eff, stab := AfferentEfferentStabilityOfClass(c)

	res += fmt.Sprintf(" Afferent:  %.2f\n", aff)
	res += fmt.Sprintf(" Efferent:  %.2f\n", eff)
	res += fmt.Sprintf(" Stability: %.2f\n", stab)

	lcom, ok := LackOfCohesionInMethodsOfCLass(c)
	if !ok {
		res += fmt.Sprintf(" LCOM: undefined (the number of methods or fields is zero)\n")
	} else {
		res += fmt.Sprintf(" LCOM: %.6f\n", lcom)
	}

	lcom4 := Lcom4(c)
	res += fmt.Sprintf(" LCOM4: %d\n", lcom4)

	return res
}

func (c *Class) ExtraFullString(level int) string {
	var res string

	res += c.FullString(level, false)

	res += c.AffEffString(true)

	if c.Implements.Len() != 0 {
		res += fmt.Sprintf(" Implements:\n")
	}
	for _, class := range c.Implements.Classes {
		res += fmt.Sprintf("%s", class.ShortStringWithPrefix(level+1, " ↳ "))
	}

	if c.Extends.Len() != 0 {
		res += fmt.Sprintf(" Extends:\n")
	}
	for _, class := range c.Extends.Classes {
		res += fmt.Sprintf("%s", class.ShortStringWithPrefix(level+1, " ↳ "))
	}

	if c.Methods.Len() != 0 {
		res += fmt.Sprintf(" Methods (%d):\n", c.Methods.Len())
	}
	for _, method := range c.Methods.Funcs {
		res += fmt.Sprintf("   %s\n", method.Name.Name)
	}

	if c.Fields.Len() != 0 {
		res += fmt.Sprintf(" Fields (%d):\n", c.Fields.Len())
	}
	for _, field := range c.Fields.Fields {
		res += fmt.Sprintf("   %s\n", field.Name)
	}

	if c.Constants.Len() != 0 {
		res += fmt.Sprintf(" Constants (%d):\n", c.Constants.Len())
	}
	for _, constant := range c.Constants.Constants {
		res += fmt.Sprintf("   %s\n", constant.Name)
	}

	return res
}

func (c *Class) ExtraShortString(level int) string {
	var res string

	if c.IsInterface {
		res += fmt.Sprintf("%sInterface %s\n", utils.GenIndent(level-1), c.Name)
	} else {
		if c.IsAbstract {
			res += fmt.Sprintf("%sAbstract class %s\n", utils.GenIndent(level-1), c.Name)
		} else {
			res += fmt.Sprintf("%sClass %s\n", utils.GenIndent(level-1), c.Name)
		}
	}

	return res
}

func (c *Class) ShortString(level int) string {
	return c.ShortStringWithPrefix(level, "")
}

func (c *Class) ShortStringWithPrefix(level int, prefix string) string {
	var res string

	if c.IsInterface {
		res += fmt.Sprintf("%sInterface %s\n", utils.GenIndent(level-1), c.Name)
	} else {
		if c.IsAbstract {
			res += fmt.Sprintf("%sAbstract class %s\n", utils.GenIndent(level-1), c.Name)
		} else {
			res += fmt.Sprintf("%sClass %s\n", utils.GenIndent(level-1), c.Name)
		}
	}
	res += fmt.Sprintf("%s Name: %s\n", utils.GenIndent(level), c.Name)
	res += fmt.Sprintf("%s File: %s:0\n", utils.GenIndent(level), c.File.Path)

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
	if c == class {
		return
	}

	_, ok := c.Deps.Get(class.Name)
	if !ok {
		c.Deps.Add(class)
	}
}

func (c *Class) AddDepsBy(class *Class) {
	if c == class {
		return
	}

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
