package lib

import (
	"bytes"
	"runtime"
	"strconv"
)

const (
	YYYYMMDD = "2006-01-02"
	DDMMYYYY = "02/01/2006"
	MMDDYYYY = "01/02/2006"
)

func IsNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 32)
	return err == nil
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
