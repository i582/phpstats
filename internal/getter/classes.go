package getter

import (
	"sort"
	"strings"

	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

type ClassesGetOptions struct {
	OnlyInterfaces bool
	OnlyClasses    bool
	OnlyTraits     bool
	Count          int64
	Offset         int64
	SortColumn     int64
	ReverseSort    bool
}

func GetClassesByOption(c *symbols.Classes, opt ClassesGetOptions) []*symbols.Class {
	classes := make([]*symbols.Class, 0, c.Len())

	if opt.Offset < 0 {
		opt.Offset = 0
	}

	var all bool
	if !opt.OnlyClasses && !opt.OnlyInterfaces && !opt.OnlyTraits {
		all = true
	}

	for _, class := range c.Classes {
		if class.IsVendor {
			continue
		}

		if !all {
			if !class.IsInterface && opt.OnlyInterfaces {
				continue
			}

			if !class.IsTrait && opt.OnlyTraits {
				continue
			}

			if class.IsInterface && !opt.OnlyInterfaces {
				continue
			}

			if class.IsTrait && !opt.OnlyTraits {
				continue
			}
		}

		classes = append(classes, class)
	}

	sort.Slice(classes, func(i, j int) bool {
		var class1 float64
		var class2 float64

		var addition func() bool = nil

		switch opt.SortColumn {
		case 0, 1: // Name
			fun1 := strings.ToLower(classes[i].ClassName())
			fun2 := strings.ToLower(classes[j].ClassName())
			if opt.ReverseSort {
				fun1, fun2 = fun2, fun1
			}
			return fun1 < fun2

		case 2: // Afferent
			class1 = representator.ClassToData(classes[i]).Afferent
			class2 = representator.ClassToData(classes[j]).Afferent
		case 3: // Efferent
			class1 = representator.ClassToData(classes[i]).Efferent
			class2 = representator.ClassToData(classes[j]).Efferent
		case 4: // Instability
			class1 = representator.ClassToData(classes[i]).Instability
			class2 = representator.ClassToData(classes[j]).Instability
		case 5: // Lcom
			class1 = representator.ClassToData(classes[i]).Lcom
			class2 = representator.ClassToData(classes[j]).Lcom
		case 6: // Lcom4
			class1 = float64(representator.ClassToData(classes[i]).Lcom4)
			class2 = float64(representator.ClassToData(classes[j]).Lcom4)
		case 7: // CountDeps
			class1 = float64(representator.ClassToData(classes[i]).CountDeps)
			class2 = float64(representator.ClassToData(classes[j]).CountDeps)
		case 8: // CountDepsBy
			class1 = float64(representator.ClassToData(classes[i]).CountDepsBy)
			class2 = float64(representator.ClassToData(classes[j]).CountDepsBy)
		case 9: //  Count fully typed methods
			class1 = utils.Percent(representator.ClassToData(classes[i]).CountFullyTypedMethods, int64(classes[i].Methods.Len()))
			class2 = utils.Percent(representator.ClassToData(classes[j]).CountFullyTypedMethods, int64(classes[j].Methods.Len()))

			addition = func() bool {
				return classes[i].Methods.Len() > classes[j].Methods.Len()
			}
		default:
			return i < j
		}

		if opt.ReverseSort {
			class1, class2 = class2, class1
		}

		if class1 == class2 && addition != nil {
			return addition()
		}

		return class1 > class2
	})

	if opt.Count+opt.Offset < int64(len(classes)) {
		classes = classes[:opt.Count+opt.Offset]
	}

	if opt.Offset < int64(len(classes)) {
		classes = classes[opt.Offset:]
	}

	return classes
}
