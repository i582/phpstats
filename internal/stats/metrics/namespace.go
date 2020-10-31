package metrics

import (
	"github.com/i582/phpstats/internal/stats/symbols"
)

func AfferentEfferentStabilityOfNamespace(n *symbols.Namespace) (aff, eff, stab float64) {
	if n.MetricsResolved {
		return n.Aff, n.Eff, n.Instab
	}

	var efferent float64
	var afferent float64

	var instability float64

	for _, class := range n.Classes.Classes {
		clAff, clEff, _ := AfferentEfferentStabilityOfClass(class)
		afferent += clAff
		efferent += clEff
	}

	if efferent+afferent == 0 {
		instability = 0
	} else {
		instability = efferent / (efferent + afferent)
	}

	n.MetricsResolved = true
	n.Aff, n.Eff, n.Instab = afferent, efferent, instability

	return afferent, efferent, instability
}
