package getter

import (
	"sort"
	"strings"

	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/stats/metrics"
	"github.com/i582/phpstats/internal/stats/symbols"
)

type NamespacesGetOptions struct {
	Level       int64
	Count       int64
	Offset      int64
	SortColumn  int64
	ReverseSort bool
}

func GetNamespacesByOptions(n *symbols.Namespaces, opt NamespacesGetOptions) []*symbols.Namespace {
	namespaces := n.GetNamespacesWithSpecificLevel(opt.Level, 100000, 0)

	if opt.Offset < 0 {
		opt.Offset = 0
	}

	sort.Slice(namespaces, func(i, j int) bool {
		metrics.AfferentEfferentStabilityOfNamespace(namespaces[i])
		metrics.AfferentEfferentStabilityOfNamespace(namespaces[j])

		var namespace1 float64
		var namespace2 float64
		switch opt.SortColumn {
		case 0, 1: // Name
			namespace1 := strings.ToLower(namespaces[i].Name)
			namespace2 := strings.ToLower(namespaces[j].Name)
			if opt.ReverseSort {
				namespace1, namespace2 = namespace2, namespace1
			}
			return namespace1 < namespace2

		case 2: // Files
			namespace1 = float64(namespaces[i].Files.Len())
			namespace2 = float64(namespaces[j].Files.Len())
		case 3: // All classes
			namespace1 = float64(representator.NamespaceToData(namespaces[i]).Classes)
			namespace2 = float64(representator.NamespaceToData(namespaces[j]).Classes)
		case 4: // Own classes
			namespace1 = float64(representator.NamespaceToData(namespaces[i]).OwnClasses)
			namespace2 = float64(representator.NamespaceToData(namespaces[j]).OwnClasses)
		case 5: // Afferent
			namespace1 = namespaces[i].Aff
			namespace2 = namespaces[j].Aff
		case 6: // Efferent
			namespace1 = namespaces[i].Eff
			namespace2 = namespaces[j].Eff
		case 7: // Instability
			namespace1 = namespaces[i].Instab
			namespace2 = namespaces[j].Instab
		case 8: // Abstractness
			namespace1 = metrics.AbstractnessOfNamespace(namespaces[i])
			namespace2 = metrics.AbstractnessOfNamespace(namespaces[j])
		case 9: // Childs
			namespace1 = float64(namespaces[i].Childs.Len())
			namespace2 = float64(namespaces[j].Childs.Len())
		default:
			return i < j
		}

		if opt.ReverseSort {
			namespace1, namespace2 = namespace2, namespace1
		}

		return namespace1 > namespace2
	})

	if opt.Count+opt.Offset < int64(len(namespaces)) {
		namespaces = namespaces[:opt.Count+opt.Offset]
	}

	if opt.Offset < int64(len(namespaces)) {
		namespaces = namespaces[opt.Offset:]
	}

	return namespaces
}
