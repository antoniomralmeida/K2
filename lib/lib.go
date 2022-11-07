package lib

import (
	"log"
	"strconv"
)

func IsNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 32)
	return err == nil
}

func LogFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Log(msg string) {
	log.Println(msg)
}
