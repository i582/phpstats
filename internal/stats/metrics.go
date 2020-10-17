package stats

func AfferentEfferentStabilityOfClass(c *Class) (aff, eff, stab float64) {
	efferent := float64(len(c.Deps.Classes))
	afferent := float64(len(c.DepsBy.Classes))

	var stability float64
	if efferent+afferent == 0 {
		stability = 0
	} else {
		stability = efferent / (efferent + afferent)
	}

	return afferent, efferent, stability
}

func LackOfCohesionInMethodsOfCLass(c *Class) (float64, bool) {
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
