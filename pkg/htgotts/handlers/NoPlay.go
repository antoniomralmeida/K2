package handlers

type NoPlay struct {
	FileName string
}

func (h *NoPlay) Play(fileName string) error {
	h.FileName = fileName
	return nil
}
