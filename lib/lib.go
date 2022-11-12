package lib

import (
	"bytes"
	"fmt"
	"os/exec"
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

func Openbrowser(url string) (err error) {

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}
