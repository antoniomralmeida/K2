package kb

func (ao *KBAttributeObject) Value() string {
	if ao.KbHistory != nil {
		return ao.KbHistory.Value
	} else {
		return ""
	}
}
