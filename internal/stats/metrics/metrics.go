package metrics

import (
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"github.com/i582/phpstats/internal/stats"
)

func AfferentEfferentStabilityOfClass(c *stats.Class) (aff, eff, stab float64) {
	efferent := float64(c.Deps.Len())
	afferent := float64(c.DepsBy.Len())

	var instability float64
	if efferent+afferent == 0 {
		instability = 0
	} else {
		instability = efferent / (efferent + afferent)
	}

	return afferent, efferent, instability
}

func LackOfCohesionInMethodsOfCLass(c *stats.Class) (float64, bool) {
	var usedSum int
	for _, field := range c.Fields.Fields {
		usedSum += len(field.Used)
	}

	allFieldMethod := c.Fields.Len() * c.Methods.Len()

	if allFieldMethod != 0 {
		return 1 - float64(usedSum)/float64(allFieldMethod), true
	}

	return -1, false
}

func Lcom4(c *stats.Class) int64 {
	g := simple.NewUndirectedGraph()

	for _, method := range c.Methods.Funcs {
		g.AddNode(method)
	}

	for _, method := range c.Methods.Funcs {
		for _, called := range method.Called.Funcs {
			if _, ok := c.Methods.Get(called.Name); ok && method != called {
				g.SetEdge(simple.Edge{
					F: method,
					T: called,
				})
			}
		}
	}

	for _, field := range c.Fields.Fields {
		functions := make([]*stats.Function, 0, len(field.Used))

		for used := range field.Used {
			functions = append(functions, used)
		}

		for i := 0; i < len(functions)-1; i++ {
			for j := i + 1; j < len(functions); j++ {
				g.SetEdge(simple.Edge{
					F: functions[i],
					T: functions[j],
				})
			}
		}
	}

	connectedComponents := topo.ConnectedComponents(g)
	return int64(len(connectedComponents))
}
