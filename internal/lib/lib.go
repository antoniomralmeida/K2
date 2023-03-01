package lib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/antoniomralmeida/k2/pkg/dsn"
	"github.com/google/uuid"
	"github.com/mattn/go-tty"
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

func KeyPress() byte {
	tty, err := tty.Open()
	if err != nil {
		return 0
	}
	defer tty.Close()

	r, err := tty.ReadRune()

	if err != nil {
		return 0
	} else {
		return byte(r)
	}
}

func LoadImage(src string) (dst string, err error) {
	dst = "./web/upload/" + uuid.New().String() + filepath.Ext(src)
	_, err = copy(src, dst)
	dst = "./upload/" + filepath.Base(dst)
	return
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func Ping(uri string) error {
	if runtime.GOOS == "windows" {
		parts := dsn.Decode(uri)
		_, err := net.Dial("tcp", parts.Socket())
		return err
	} else {
		return nil
	}
}

func isDebugging() bool {
	filename := filepath.Base(os.Args[0])
	return strings.Contains(filename, "debug")
}

func GetWorkDir() string {
	wd, _ := os.Getwd()
	return wd
}

func Identify(name string) string {
	var specialCharSet = "'\"!@#$%&*+-/ "
	for _, c := range specialCharSet {
		name = strings.ReplaceAll(name, string(c), "")
	}
	return name
}

func TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

func DownloadFile(fullURLFile string, dirPath string) (error, string) {

	// Build fileName from fullPath
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return err, ""
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := dirPath + segments[len(segments)-1]
	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		return err, ""
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return err, ""
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status), ""
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	defer file.Close()
	return err, fileName
}

// exists returns whether the given file or directory exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
