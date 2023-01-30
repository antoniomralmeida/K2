package util

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/antoniomralmeida/golibretranslate"
	"github.com/antoniomralmeida/k2/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
)

// Message contains the message's tag and its contained matched sentences
type Message struct {
	Tag      string   `json:"tag"`
	Messages []string `json:"messages"`
}

var messages = map[string][]Message{}

func messagesFile(locale string) string {
	return inits.GetHomeDir() + "/k2olivia/res/locales/" + locale + "/messages.json"
}

func translateMessages(_messages *[]Message, locale string) (err error) {
	for i := range *_messages {
		var trans string
		trans, err = golibretranslate.Translate((*_messages)[i].Tag, inits.DefaultLocale, locale)
		if err == nil {
			(*_messages)[i].Tag = trans
		} else {
			return
		}
		for j := range (*_messages)[i].Messages {
			trans, err = golibretranslate.Translate((*_messages)[i].Messages[j], inits.DefaultLocale, locale)
			if err == nil {
				(*_messages)[i].Messages[j] = trans
			} else {
				return
			}
			time.Sleep(time.Millisecond * 10) //429 Too Many Requests Error
		}
	}
	return
}

// SerializeMessages serializes the content of `res/datasets/messages.json` in JSON
func SerializeMessages(locale string) []Message {
	var currentMessages []Message

	msgFile := messagesFile(locale)

	if ok, _ := lib.Exists(msgFile); !ok {
		messages_tmp := []Message{}
		tmpFile := messagesFile(inits.DefaultLocale)
		err := json.Unmarshal(ReadFile(tmpFile), &messages_tmp)
		inits.Log(err, inits.Fatal)
		err = translateMessages(&messages_tmp, locale)
		inits.Log(err, inits.Fatal)
		js, err := json.Marshal(messages_tmp)
		inits.Log(err, inits.Error)
		f, err := os.Create(msgFile)
		inits.Log(err, inits.Error)
		f.WriteString(string(js))
		f.Close()
	}

	err := json.Unmarshal(ReadFile(msgFile), &currentMessages)
	if err != nil {
		inits.Log(err, inits.Error)
	}
	messages[locale] = currentMessages
	return currentMessages
}

// GetMessages returns the cached messages for the given locale
func GetMessages(locale string) []Message {
	return messages[locale]
}

// GetMessageByTag returns a message found by the given tag and locale
func GetMessageByTag(tag, locale string) Message {
	for _, message := range messages[locale] {
		if tag != message.Tag {
			continue
		}
		return message
	}
	return Message{}
}

// GetMessage retrieves a message tag and returns a random message chose from res/datasets/messages.json
func GetMessage(locale, tag string) string {
	for _, message := range messages[locale] {
		// Find the message with the right tag
		if message.Tag != tag {
			continue
		}

		// Returns the only element if there aren't more
		if len(message.Messages) == 1 {
			return message.Messages[0]
		}

		// Returns a random sentence
		rand.Seed(time.Now().UnixNano())
		return message.Messages[rand.Intn(len(message.Messages))]
	}

	return ""
}
