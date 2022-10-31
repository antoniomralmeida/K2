package lib

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leemcloughlin/logfile"
)

func IsNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 32)
	return err == nil
}

func LogInit() {

	logFileName := "./log/k2.log"
	logFile, err := logfile.New(
		&logfile.LogFile{
			FileName: logFileName,
			MaxSize:  500 * 1024, // 500K duh!
			Flags:    logfile.FileOnly | logfile.OverWriteOnStart})
	if err != nil {
		log.Fatalf("Failed to create logFile %s: %s\n", logFileName, err)
	}

	log.SetOutput(logFile)
}

func LogFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Log(msg string) {
	log.Println(msg)
}

func IsMainThread() bool {
	return !fiber.IsChild()
}

func CompressFile(name_of_file string) string {

	f, _ := os.Open(name_of_file)

	read := bufio.NewReader(f)

	data, _ := ioutil.ReadAll(read)

	name_of_zip := "./tmp/" + filepath.Base(name_of_file) + ".zip"
	f, _ = os.Create(name_of_zip)
	fmt.Println(name_of_zip)
	w := gzip.NewWriter(f)

	w.Write(data)
	w.Close()
	return name_of_zip
}

func CompressBuffer(input *bufio.ReadWriter) *gzip.Writer {
	f, _ := os.CreateTemp("./tmp/", "tmpfile-")
	w := gzip.NewWriter(f)
	data, _ := ioutil.ReadAll(input)
	w.Write(data)
	w.Close()
	return w
}
