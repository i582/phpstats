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

	return fn.CalledBy.GetAll(false, false, true, -1, false), nil
}

func (fi *Functions) GetFuncsThatCalledInFunc(name string) ([]*Function, error) {
	fn, ok := fi.Get(NewFuncKey(name))

	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	return fn.Called.GetAll(false, false, true, -1, false), nil
}

func (fi *Functions) GetAll(onlyMethods, onlyFuncs, all bool, count int, sorted bool) []*Function {
	var res []*Function
	var index int

	for key, fn := range fi.Funcs {
		if !sorted {
			if index > count && count != -1 {
				break
			}
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
		if count < len(res) {
			res = res[:count]
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
	Class *Class

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

func (fi *Function) IsMethod() bool {
	return fi.Name.IsMethod()
}

func (fi Function) Equal(fi2 Function) bool {
	return fi.Name.Equal(fi2.Name)
}

func (fi *Function) ShortString() string {
	var res string

	if fi.Name.IsMethod() {
		res += fmt.Sprintf("Метод %s\n", fi.Name)

		// если функция не встроенная
		if !IsEmbeddedFunc(fi.Pos.Filename) {
			res += fmt.Sprintf(" Класс: ")
			if fi.Class != nil {
				res += fi.Class.Name
			} else {
				res += "undefined"
			}
			res += "\n"
		}

	} else {
		res += fmt.Sprintf("Функция %s\n", fi.Name)
	}

	return res
}

func (fi *Function) String() string {
	var res string

	res += "\n"

	res += fi.ShortString()

	if IsEmbeddedFunc(fi.Pos.Filename) {
		res += fmt.Sprintf(" Встроенная функция\n")
	} else {
		res += fmt.Sprintf(" Определена здесь: %s:%d\n", fi.Pos.Filename, fi.Pos.Line)
	}

	res += fmt.Sprintf(" Количество использований: %d\n", fi.UsesCount)

	if len(fi.Called.Funcs) != 0 {
		res += fmt.Sprintf(" Вызываемые функции (%d):\n", len(fi.Called.Funcs))
	}
	for _, fn := range fi.Called.Funcs {
		res += fmt.Sprintf("   %s\n", fn.Name.Name)
	}

	return res
}

func (fi *Function) AddCalled(fn *Function) {
	fi.Called.Add(fn)
}

func (fi *Function) AddCalledBy(fn *Function) {
	fi.CalledBy.Add(fn)
}

func (fi *Function) AddUse() {
	atomic.AddInt64(&fi.UsesCount, 1)
}

func (b *blockChecker) CurFunc() string {
	curClass := b.ctx.ClassParseState().CurrentClass
	curFunction := b.ctx.ClassParseState().CurrentFunction
	curNamespace := b.ctx.ClassParseState().Namespace

	if curFunction == "" {
		return ""
	}

	if curNamespace != "" {
		curNamespace = curNamespace + `\`
	}

	if curClass != "" {
		curClass = curClass + `::`
	}

	if curClass == "" && curNamespace == "" {
		curFunction = `\` + curFunction
	}

	if curClass == "" && curNamespace != "" {
		curFunction = curNamespace + curFunction
	}

	curFqn := curClass + curFunction

	return curFqn
}

// GobDecode is custom gob unmarshaller
func (fi *Function) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&fi.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&fi.Pos)
	if err != nil {
		return err
	}
	return nil
}

// GobEncode is a custom gob marshaller
func (fi *Function) GobEncode() ([]byte, error) {
	if fi.Encoded {
		return nil, nil
	}
	fi.Encoded = true

	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(fi.Name)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(fi.Pos)
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
