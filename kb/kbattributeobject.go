package kb

func (ao *KBAttributeObject) Value() any {
	if ao.KbHistory != nil {
		return ao.KbHistory.Value
	} else {
		return ""
	}
}
