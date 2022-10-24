package classes

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
)

type TTS struct {
	text     string
	language string
}

func (e *TTS) SetText(txt string, language string) {
	e.text = txt
	e.language = language
}

func (e *TTS) Speech() {
	speech := htgotts.Speech{Folder: "audio", Language: e.language, Handler: &handlers.Native{}}
	f, err := speech.CreateSpeechFile(e.text, "tmpfile-"+uuid.New().String())
	if err != nil {
		fmt.Println("CreateSpeechFile fail %v", err)
	}
	fmt.Println(f)
	speech.PlaySpeechFile(f)
	err = os.Remove(f)
	if err != nil {
		fmt.Println(err)
		return
	}

}
