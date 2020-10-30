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
		if fn.Name.String() == name {
			return []FuncKey{fn.Name}, nil
		}

		if strings.Contains(fn.Name.String(), name) {
			res = append(res, fn.Name)
		}
	}

	if len(res) == 0 {
		return res, fmt.Errorf("function %s not found", name)
	}

	return res, nil
}

func NewFunctions() *Functions {
	return &Functions{
		Funcs: map[FuncKey]*Function{},
	}
}

func (fi *Functions) GetAll(onlyMethods, onlyFuncs, all bool, count int64, offset int64, sorted bool, withEmbeddedFuncs bool) []*Function {
	var res = make([]*Function, 0, len(fi.Funcs))
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
	fi.Funcs[fn.Name] = fn
	fi.Unlock()
}

func (fi *Functions) Get(fn FuncKey) (*Function, bool) {
	fi.Lock()
	f, ok := fi.Funcs[fn]
	fi.Unlock()
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

var FunctionCount int64

type Function struct {
	Name FuncKey
	Pos  meta.ElementPosition

	Called   *Functions
	CalledBy *Functions

	UsesCount int64

	depsResolved bool
	deps         *Classes

	depsByResolved bool
	depsBy         *Classes

	CyclomaticComplexity int64

	// Method part
	Class *Class

	Id int64
}

func (f *Function) ID() int64 {
	return f.Id
}

func NewFunctionInfo(name FuncKey, pos meta.ElementPosition) *Function {
	atomic.AddInt64(&FunctionCount, 1)
	return &Function{
		Name:     name,
		Called:   NewFunctions(),
		CalledBy: NewFunctions(),
		deps:     NewClasses(),
		depsBy:   NewClasses(),
		Pos:      pos,
		Id:       FunctionCount,
	}
}

func NewMethodInfo(name FuncKey, pos meta.ElementPosition, class *Class) *Function {
	method := NewFunctionInfo(name, pos)
	method.Class = class
	return method
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

func (f *Function) Deps() *Classes {
	if f.depsResolved {
		return f.deps
	}

	for _, called := range f.Called.Funcs {
		if called.Class == nil {
			continue
		}
		if called.Class == f.Class {
			continue
		}

		f.deps.Add(called.Class)
	}

	f.depsResolved = true

	return f.deps
}

func (f *Function) CountDeps() int64 {
	return int64(f.Deps().Len())
}

func (f *Function) DepsBy() *Classes {
	if f.depsByResolved {
		return f.depsBy
	}

	for _, called := range f.CalledBy.Funcs {
		if called.Class == nil {
			continue
		}
		if called.Class == f.Class {
			continue
		}

		f.depsBy.Add(called.Class)
	}

	f.depsByResolved = true

	return f.depsBy
}

func (f *Function) CountDepsBy() int64 {
	return int64(f.DepsBy().Len())
}

func (f *Function) AddCalled(fn *Function) {
	f.Called.Add(fn)

	if f.Class == nil || fn.Class == nil {
		return
	}

	f.Class.AddDeps(fn.Class)
}

func (f *Function) AddCalledBy(fn *Function) {
	f.CalledBy.Add(fn)

	if f.Class == nil || fn.Class == nil {
		return
	}

	f.Class.AddDepsBy(fn.Class)
	f.AddUse()
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
	err = decoder.Decode(&f.CyclomaticComplexity)
	if err != nil {
		return err
	}
	return nil
}

// GobEncode is a custom gob marshaller
func (f *Function) GobEncode() ([]byte, error) {
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
	err = encoder.Encode(f.CyclomaticComplexity)
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
