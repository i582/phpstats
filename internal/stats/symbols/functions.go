package symbols

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

func (f FuncKey) IsMethod() bool {
	return f.ClassName != ""
}

func (f FuncKey) String() string {
	var res string

	if f.ClassName != "" {
		res += f.ClassName + "::" + f.Name
	} else {
		res += f.Name
	}

	return res
}

func (f FuncKey) Equal(fk2 FuncKey) bool {
	return f.Name == fk2.Name && f.ClassName == fk2.ClassName
}

type Functions struct {
	m sync.Mutex

	Funcs map[FuncKey]*Function
}

func (f *Functions) Len() int {
	return len(f.Funcs)
}

func (f *Functions) CyclomaticComplexity() int64 {
	var res int64
	for _, fn := range f.Funcs {
		res += fn.CyclomaticComplexity
	}
	return res
}

func (f *Functions) CountMagicNumbers() int64 {
	var count int64
	for _, fn := range f.Funcs {
		count += fn.CountMagicNumbers
	}
	return count
}

func (f *Functions) MaxMinAvgMethodCyclomaticComplexity() (max, min, avg float64) {
	return f.maxMinAvgCyclomaticComplexity(true, false)
}

func (f *Functions) MaxMinAvgFunctionsCyclomaticComplexity() (max, min, avg float64) {
	return f.maxMinAvgCyclomaticComplexity(false, true)
}

func (f *Functions) maxMinAvgCyclomaticComplexity(onlyMethods, onlyFunctions bool) (max, min, avg float64) {
	const maxValue = 100000.0
	var count float64

	max = 0
	min = maxValue

	for _, fn := range f.Funcs {
		if onlyMethods && !fn.IsMethod() {
			continue
		}

		if onlyFunctions && (fn.IsMethod() || fn.IsEmbeddedFunc()) {
			continue
		}

		count += float64(fn.CyclomaticComplexity)

		if float64(fn.CyclomaticComplexity) < min {
			min = float64(fn.CyclomaticComplexity)
		}

		if float64(fn.CyclomaticComplexity) > max {
			max = float64(fn.CyclomaticComplexity)
		}
	}

	if min == maxValue {
		min = 0
	}

	if onlyMethods && f.CountMethods() != 0 {
		avg = count / float64(f.CountMethods())
	} else if onlyFunctions && f.CountFunctions() != 0 {
		avg = count / float64(f.CountMethods())
	} else if f.Len() != 0 {
		avg = count / float64(f.Len())
	}

	return max, min, avg
}

func (f *Functions) MaxMinAvgMethodCountMagicNumbers() (max, min, avg int64) {
	return f.maxMinAvgCountMagicNumbers(true, false)
}

func (f *Functions) MaxMinAvgFunctionsCountMagicNumbers() (max, min, avg int64) {
	return f.maxMinAvgCountMagicNumbers(false, true)
}

func (f *Functions) maxMinAvgCountMagicNumbers(onlyMethods, onlyFunctions bool) (max, min, avg int64) {
	const maxValue = 100000
	var count int64

	max = 0
	min = maxValue

	for _, fn := range f.Funcs {
		if onlyMethods && !fn.IsMethod() {
			continue
		}

		if onlyFunctions && (fn.IsMethod() || fn.IsEmbeddedFunc()) {
			continue
		}

		count += fn.CountMagicNumbers

		if fn.CountMagicNumbers < min {
			min = fn.CountMagicNumbers
		}

		if fn.CountMagicNumbers > max {
			max = fn.CountMagicNumbers
		}
	}

	if min == maxValue {
		min = 0
	}

	if onlyMethods && f.CountMethods() != 0 {
		avg = count / f.CountMethods()
	} else if onlyFunctions && f.CountFunctions() != 0 {
		avg = count / f.CountFunctions()
	} else if f.Len() != 0 {
		avg = count / int64(f.Len())
	}

	return max, min, avg
}

func (f *Functions) CountFunctions() int64 {
	var count int64
	for _, fn := range f.Funcs {
		if !fn.IsMethod() && !fn.IsEmbeddedFunc() {
			count++
		}
	}
	return count
}

func (f *Functions) CountMethods() int64 {
	var count int64
	for _, fn := range f.Funcs {
		if fn.IsMethod() {
			count++
		}
	}
	return count
}

func (f *Functions) GetFunctionByPartOfName(name string) (*Function, error) {
	funcs, err := f.GetFullFuncName(name)
	if err != nil {
		return nil, err
	}

	fun, found := f.Get(funcs[0])
	if !found {
		return nil, fmt.Errorf("function %s not found", name)
	}
	return fun, nil
}

func (f *Functions) GetFullFuncName(name string) ([]FuncKey, error) {
	var res []FuncKey

	for _, fn := range f.Funcs {
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

func (f *Functions) GetAll(onlyMethods, onlyFuncs, all bool, count int64, offset int64, sorted bool, withEmbeddedFuncs bool) []*Function {
	var res = make([]*Function, 0, len(f.Funcs))
	var index int64

	if offset < 0 {
		offset = 0
	}

	for key, fn := range f.Funcs {
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

func (f *Functions) Add(fn *Function) {
	f.m.Lock()
	f.Funcs[fn.Name] = fn
	f.m.Unlock()
}

func (f *Functions) Get(fn FuncKey) (*Function, bool) {
	f.m.Lock()
	fun, ok := f.Funcs[fn]
	f.m.Unlock()
	return fun, ok
}

var FunctionCount int64

type Function struct {
	Name FuncKey
	Pos  meta.ElementPosition

	Namespace *Namespace

	Called   *Functions
	CalledBy *Functions

	UsedFields    *Fields
	UsedConstants *Constants

	UsesCount int64

	depsResolved bool
	deps         *Classes

	depsByResolved bool
	depsBy         *Classes

	CyclomaticComplexity int64
	CountMagicNumbers    int64

	// Method part
	Class *Class

	Id int64
}

func (f *Function) ID() int64 {
	return f.Id
}

func NewFunction(name FuncKey, pos meta.ElementPosition) *Function {
	atomic.AddInt64(&FunctionCount, 1)
	return &Function{
		Name:          name,
		Called:        NewFunctions(),
		CalledBy:      NewFunctions(),
		UsedFields:    NewFields(),
		UsedConstants: NewConstants(),
		deps:          NewClasses(),
		depsBy:        NewClasses(),
		Pos:           pos,
		Id:            FunctionCount,
	}
}

func NewMethod(name FuncKey, pos meta.ElementPosition, class *Class) *Function {
	method := NewFunction(name, pos)
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
	err = decoder.Decode(&f.CountMagicNumbers)
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
	err = encoder.Encode(f.CountMagicNumbers)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
