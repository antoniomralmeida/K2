package tts

import (
	htgotts "github.com/hegedustibor/htgo-tts"
)

type NoPlay struct {
	FileName string
}

func (h *NoPlay) Play(fileName string) error {
	h.FileName = fileName
	return nil
}

func TTSToFile(text, language, pathDir string) (string, error) {
	handler := new(NoPlay)
	speech := htgotts.Speech{Folder: pathDir, Language: language, Handler: handler}
	if err := speech.Speak(text); err != nil {
		return "", err
	} else {
		return handler.FileName, nil
	}
}
