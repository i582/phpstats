package grapher

import (
	"github.com/i582/phpstats/internal/grapher/unl"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func (g *Grapher) Namespace(ns *symbols.Namespace) string {
	var res string

	res += graphHeader
	res += g.namespace(ns)

	return g.graphWrapper(res, utils.NameToIdentifier(ns.FullName))
}

func (g *Grapher) Namespaces(ns *symbols.Namespaces) string {
	var res string

	res += graphHeader

	for _, n := range ns.Namespaces {
		res += g.namespace(n)
	}

	return g.graphWrapper(res, "allNamespacesGraph")
}

func (g *Grapher) namespace(ns *symbols.Namespace) string {
	var res string

	for _, n := range ns.Childs.Namespaces {
		res += g.namespace(n)
	}
	if ns.Childs.Len() == 0 {
		res += uml.GetUmlForNamespace(ns)
	}

	return g.subGraphWrapper(res, utils.NormalizeSlashes(ns.FullName))
}
