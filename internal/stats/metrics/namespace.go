package metrics

import (
	"github.com/i582/phpstats/internal/stats/symbols"
)

// AfferentEfferentStabilityOfNamespace calculates afferent, efferent and instability
// metrics for the passed namespace.
func AfferentEfferentStabilityOfNamespace(n *symbols.Namespace) (aff, eff, stab float64) {
	if n.MetricsResolved {
		return n.Aff, n.Eff, n.Instab
	}

	var efferent float64
	var afferent float64

	var instability float64

	for _, class := range n.Classes.Classes {
		clAff, clEff, _ := AfferentEfferentInstabilityOfClass(class)
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

// AbstractnessOfNamespace calculates abstractness metrics for the passed namespace.
func AbstractnessOfNamespace(n *symbols.Namespace) float64 {
	abstractClasses, allClasses := abstractnessOfNamespace(n)

	for _, namespace := range n.Childs.Namespaces {
		abstractClassesForChildNamespace, allClassesForChildNamespace := abstractnessOfNamespace(namespace)
		abstractClasses += abstractClassesForChildNamespace
		allClasses += allClassesForChildNamespace
	}

	if allClasses == 0 {
		allClasses = 1
	}

	return abstractClasses / allClasses
}

func abstractnessOfNamespace(n *symbols.Namespace) (abstract float64, all float64) {
	abstractClasses := n.Classes.CountAbstractClasses()
	return float64(abstractClasses), float64(n.Classes.Len())
}
