package metrics

import (
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"github.com/i582/phpstats/internal/stats/symbols"
)

// AfferentEfferentInstabilityOfClass calculates afferent, efferent and instability
// metrics for the passed class.
func AfferentEfferentInstabilityOfClass(c *symbols.Class) (aff, eff, stab float64) {
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

// LackOfCohesionInMethods calculates the Lack Of Cohesion In Methods metric for the passed class.
func LackOfCohesionInMethods(c *symbols.Class) (float64, bool) {
	if c.LcomResolved {
		return c.Lcom, true
	}

	var usedSum int
	for _, field := range c.Fields.Fields {
		usedSum += field.Used.Len()
	}

	allFieldMethod := c.Fields.Len() * c.Methods.Len()

	if allFieldMethod != 0 {
		c.LcomResolved = true
		c.Lcom = 1 - float64(usedSum)/float64(allFieldMethod)

		return c.Lcom, true
	}

	c.LcomResolved = true
	c.Lcom = -1

	return -1, false
}

// LackOfCohesionInMethods4 calculates the Lack Of Cohesion In Methods 4 metric for the passed class.
func LackOfCohesionInMethods4(c *symbols.Class) int64 {
	if c.Lcom4Resolved {
		return c.Lcom4
	}

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
		functions := make([]*symbols.Function, 0, field.Used.Len())

		for _, used := range field.Used.Funcs {
			functions = append(functions, used)
		}

		for i := 0; i < len(functions)-1; i++ {
			for j := i + 1; j < len(functions); j++ {
				if functions[i].Id == functions[j].Id {
					continue
				}

				g.SetEdge(simple.Edge{
					F: functions[i],
					T: functions[j],
				})
			}
		}
	}

	connectedComponents := topo.ConnectedComponents(g)

	c.Lcom4Resolved = true
	c.Lcom4 = int64(len(connectedComponents))

	return c.Lcom4
}
