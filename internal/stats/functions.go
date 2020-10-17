package stats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/VKCOM/noverify/src/meta"
)

type FuncKey struct {
	Name      string
	ClassName string
}

func NewFuncKey(name string) FuncKey {
	return FuncKey{
		Name: name,
	}
}

func NewMethodKey(name, className string) FuncKey {
	return FuncKey{
		Name:      name,
		ClassName: className,
	}
}

func (fk FuncKey) IsMethod() bool {
	return fk.ClassName != ""
}

func (fk FuncKey) String() string {
	var res string

	if fk.ClassName != "" {
		res += fk.ClassName + "::" + fk.Name
	} else {
		res += fk.Name
	}

	return res
}

func (fk FuncKey) Equal(fk2 FuncKey) bool {
	return fk.Name == fk2.Name && fk.ClassName == fk2.ClassName
}

type Functions struct {
	sync.Mutex

	Funcs map[FuncKey]*Function
}

func (fi *Functions) Len() int {
	return len(fi.Funcs)
}

func (fi *Functions) GetFullFuncName(name string) ([]FuncKey, error) {
	var res []FuncKey

	for _, fn := range fi.Funcs {
		if strings.Contains(fn.Name.String(), name) {
			res = append(res, fn.Name)
		}
	}

	if len(res) == 0 {
		return res, fmt.Errorf("function %s not found", name)
	}

	return res, nil
}

func NewFunctionsInfo() *Functions {
	return &Functions{
		Funcs: map[FuncKey]*Function{},
	}
}

func (fi *Functions) GetFuncsThatCalledFunc(name string) ([]*Function, error) {
	fn, ok := fi.Get(NewFuncKey(name))

	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	return fn.CalledBy.GetAll(false, false, true, -1, 0, false, true), nil
}

func (fi *Functions) GetFuncsThatCalledInFunc(name string) ([]*Function, error) {
	fn, ok := fi.Get(NewFuncKey(name))

	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	return fn.Called.GetAll(false, false, true, -1, 0, false, true), nil
}

func (fi *Functions) GetAll(onlyMethods, onlyFuncs, all bool, count int64, offset int64, sorted bool, withEmbeddedFuncs bool) []*Function {
	var res []*Function
	var index int64

	if offset < 0 {
		offset = 0
	}

	for key, fn := range fi.Funcs {
		if !sorted {
			if index > count+offset && count != -1 {
				break
			}
		}

		if !withEmbeddedFuncs && fn.IsEmbeddedFunc() {
			continue
		}

		if !all {
			if !key.IsMethod() && onlyMethods {
				continue
			}

			if key.IsMethod() && onlyFuncs {
				continue
			}
		}

		res = append(res, fn)
		index++
	}

	if sorted {
		sort.Slice(res, func(i, j int) bool {
			return res[i].UsesCount > res[j].UsesCount
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

func (fi *Functions) Add(fn *Function) {
	fi.Lock()
	defer fi.Unlock()
	fi.Funcs[fn.Name] = fn
}

func (fi *Functions) Get(fn FuncKey) (*Function, bool) {
	fi.Lock()
	defer fi.Unlock()
	f, ok := fi.Funcs[fn]
	return f, ok
}

func (fi *Functions) GetOrCreateFunction(fn FuncKey, pos meta.ElementPosition) *Function {
	fi.Lock()
	f, ok := fi.Funcs[fn]
	fi.Unlock()
	if !ok {
		f = NewFunctionInfo(fn, pos)
		GlobalCtx.Funcs.Add(f)
	}

	return f
}

func (fi *Functions) GetOrCreateMethod(fn FuncKey, pos meta.ElementPosition, class *Class) *Function {
	fi.Lock()
	f, ok := fi.Funcs[fn]
	fi.Unlock()
	if !ok {
		f = NewMethodInfo(fn, pos, class)
		GlobalCtx.Funcs.Add(f)
	}

	return f
}

type Function struct {
	Name FuncKey
	Pos  meta.ElementPosition

	Called   *Functions
	CalledBy *Functions

	UsesCount int64

	// Method part
	Class   *Class
	Encoded bool
}

func NewFunctionInfo(name FuncKey, pos meta.ElementPosition) *Function {
	return &Function{
		Name:     name,
		Called:   NewFunctionsInfo(),
		CalledBy: NewFunctionsInfo(),
		Pos:      pos,
	}
}

func NewMethodInfo(name FuncKey, pos meta.ElementPosition, class *Class) *Function {
	return &Function{
		Name:     name,
		Called:   NewFunctionsInfo(),
		CalledBy: NewFunctionsInfo(),
		Pos:      pos,
		Class:    class,
	}
}

func IsEmbeddedFunc(name string) bool {
	return strings.Contains(name, "phpstorm-stubs")
}

func (f *Function) IsEmbeddedFunc() bool {
	return IsEmbeddedFunc(f.Pos.Filename)
}

func (f *Function) IsMethod() bool {
	return f.Name.IsMethod()
}

func (f Function) Equal(fi2 Function) bool {
	return f.Name.Equal(fi2.Name)
}

func (f *Function) GraphvizRecursive(level int64, maxLevel int64, visited map[string]struct{}) string {
	var res string

	if level > maxLevel {
		return ""
	}

	classes := NewClasses()

	for _, called := range f.Called.Funcs {
		if called.Class == nil {
			continue
		}

		// res += called.Class.Deps.Graphviz()

		res += called.GraphvizRecursive(level+1, maxLevel, visited)

		classes.Add(called.Class)
	}

	// graph func -o test2.gv \VK\API\Library\DeprecatedWrappers::wrapComments -show -r 2

	// res += fmt.Sprintf("   \"%s\" [shape=\"rectangle\"]\n", f.Name)

	for _, class := range classes.Classes {
		str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", f.Name, class.Name)

		if _, ok := visited[str]; ok {
			continue
		}
		visited[str] = struct{}{}

		res += str
	}

	return res
}

func (f *Function) ShortString() string {
	var res string

	if f.Name.IsMethod() {
		res += fmt.Sprintf("Метод %s\n", f.Name)

		// если функция не встроенная
		if !IsEmbeddedFunc(f.Pos.Filename) {
			res += fmt.Sprintf(" Класс: ")
			if f.Class != nil {
				res += f.Class.Name
			} else {
				res += "undefined"
			}
			res += "\n"
		}

	} else {
		res += fmt.Sprintf("Функция %s\n", f.Name)
	}

	return res
}

func (f *Function) FullString() string {
	var res string

	res += f.ShortString()

	if IsEmbeddedFunc(f.Pos.Filename) {
		res += fmt.Sprintf(" Встроенная функция\n")
	} else {
		res += fmt.Sprintf(" Определена здесь: %s:%d\n", f.Pos.Filename, f.Pos.Line)
	}

	res += fmt.Sprintf(" Количество использований: %d\n", f.UsesCount)

	if len(f.Called.Funcs) != 0 {
		res += fmt.Sprintf(" Вызываемые функции (%d):\n", len(f.Called.Funcs))
	}
	for _, fn := range f.Called.Funcs {
		res += fmt.Sprintf("   %s\n", fn.Name.Name)
	}

	return res
}

func (f *Function) AddCalled(fn *Function) {
	f.Called.Add(fn)
}

func (f *Function) AddCalledBy(fn *Function) {
	f.CalledBy.Add(fn)
}

func (f *Function) AddUse() {
	atomic.AddInt64(&f.UsesCount, 1)
}

// GobDecode is custom gob unmarshaller
func (f *Function) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&f.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&f.Pos)
	if err != nil {
		return err
	}
	return nil
}

// GobEncode is a custom gob marshaller
func (f *Function) GobEncode() ([]byte, error) {
	if f.Encoded {
		return nil, nil
	}
	f.Encoded = true

	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(f.Name)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(f.Pos)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// GobDecode is custom gob unmarshaller
func (fi *Functions) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&fi.Funcs)
	if err != nil {
		return err
	}
	return nil
}

// GobEncode is a custom gob marshaller
func (fi *Functions) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(fi.Funcs)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
