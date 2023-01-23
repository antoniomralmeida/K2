package util

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/antoniomralmeida/golibretranslate"
	"github.com/antoniomralmeida/k2/initializers"
	"github.com/antoniomralmeida/k2/k2olivia/locales"
	"github.com/antoniomralmeida/k2/lib"
)

// Message contains the message's tag and its contained matched sentences
type Message struct {
	Tag      string   `json:"tag"`
	Messages []string `json:"messages"`
}

var messages = map[string][]Message{}

func messagesFile(locale string) string {
	return initializers.GetHomeDir() + "/k2olivia/res/locales/" + locale + "/messages.json"
}

func translateMessages(_messages *[]Message, locale string) (err error) {
	for i := range *_messages {
		var trans string

		for j := range (*_messages)[i].Messages {
			trans, err = golibretranslate.Translate((*_messages)[i].Messages[j], locales.Locale_default, locale)
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
		tmpFile := messagesFile(locales.Locale_default)
		err := json.Unmarshal(ReadFile(tmpFile), &messages_tmp)
		initializers.Log(err, initializers.Fatal)
		err = translateMessages(&messages_tmp, locale)
		initializers.Log(err, initializers.Fatal)
		js, err := json.Marshal(messages_tmp)
		initializers.Log(err, initializers.Error)
		f, err := os.Create(msgFile)
		initializers.Log(err, initializers.Error)
		f.WriteString(string(js))
		f.Close()
	}

	err := json.Unmarshal(ReadFile(msgFile), &currentMessages)
	if err != nil {
		initializers.Log(err, initializers.Error)
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
