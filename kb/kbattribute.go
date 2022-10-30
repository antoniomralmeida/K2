package kb

func (a *KBAttribute) addAntecedentRules(r *KBRule) {
	found := false
	for i := range a.antecedentRules {
		if a.antecedentRules[i] == r {
			found = true
			break
		}
	}
	if !found {
		a.antecedentRules = append(a.antecedentRules, r)
	}
}

func (a *KBAttribute) addConsequentRules(r *KBRule) {
	found := false
	for i := range a.consequentRules {
		if a.consequentRules[i] == r {
			found = true
			break
		}
	}
	if !found {
		a.consequentRules = append(a.consequentRules, r)
	}
}

func (a *KBAttribute) isSource(s KBSource) bool {
	for i := range a.Sources {
		if a.Sources[i] == s {
			return true
		}
	}
	return false
}
