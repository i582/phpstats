package getter

import (
	"sort"
	"strings"

	"github.com/i582/phpstats/internal/stats/symbols"
)

type FunctionsGetOptions struct {
	OnlyMethods       bool
	OnlyFuncs         bool
	Count             int64
	Offset            int64
	WithEmbeddedFuncs bool
	SortColumn        int64
	ReverseSort       bool
}

func GetFunctionsByOptions(f *symbols.Functions, opt FunctionsGetOptions) []*symbols.Function {
	funcs := make([]*symbols.Function, 0, f.Len())

	if opt.Offset < 0 {
		opt.Offset = 0
	}

	var all bool
	if !opt.OnlyFuncs && !opt.OnlyMethods {
		all = true
	}

	for key, fn := range f.Funcs {
		if fn.IsVendorFunction() {
			continue
		}

		if !opt.WithEmbeddedFuncs && fn.IsEmbeddedFunc() {
			continue
		}

		if !all {
			if !key.IsMethod() && opt.OnlyMethods {
				continue
			}

			if key.IsMethod() && opt.OnlyFuncs {
				continue
			}
		}

		funcs = append(funcs, fn)
	}

	sort.Slice(funcs, func(i, j int) bool {
		var fun1 int64
		var fun2 int64
		switch opt.SortColumn {
		case 0, 1: // Name
			fun1 := strings.ToLower(funcs[i].Name.Name)
			fun2 := strings.ToLower(funcs[j].Name.Name)
			if opt.ReverseSort {
				fun1, fun2 = fun2, fun1
			}
			return fun1 < fun2

		case 2: // UsesCount
			fun1 = funcs[i].UsesCount
			fun2 = funcs[j].UsesCount
		case 3: // CountDeps
			fun1 = funcs[i].CountDeps()
			fun2 = funcs[j].CountDeps()
		case 4: // CountDepsBy
			fun1 = funcs[i].CountDepsBy()
			fun2 = funcs[j].CountDepsBy()
		case 5: // Called
			fun1 = int64(funcs[i].Called.Len())
			fun2 = int64(funcs[j].Called.Len())
		case 6: // CalledBy
			fun1 = int64(funcs[i].CalledBy.Len())
			fun2 = int64(funcs[j].CalledBy.Len())
		case 7: // CyclomaticComplexity
			fun1 = funcs[i].CyclomaticComplexity
			fun2 = funcs[j].CyclomaticComplexity
		case 8: // CountMagicNumbers
			fun1 = funcs[i].CountMagicNumbers
			fun2 = funcs[j].CountMagicNumbers
		default:
			return i < j
		}

		if opt.ReverseSort {
			fun1, fun2 = fun2, fun1
		}

		return fun1 > fun2
	})

	if opt.Count+opt.Offset < int64(len(funcs)) {
		funcs = funcs[:opt.Count+opt.Offset]
	}

	if opt.Offset < int64(len(funcs)) {
		funcs = funcs[opt.Offset:]
	}

	return funcs
}
