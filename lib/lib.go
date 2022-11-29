package lib

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

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

type item struct {
	idx int
	max int
}

type Cartesian struct {
	itens   map[string]*item
	keys    []string
	current int
}

func (c *Cartesian) AddItem(key string, max int) {
	if c.itens == nil {
		c.itens = make(map[string]*item)
	}
	c.itens[key] = &item{max: max}
	c.keys = []string{}
	c.current = 0
	for key := range c.itens {
		c.keys = append(c.keys, key)
	}
}

func (c *Cartesian) next(i int) bool {
	if i >= len(c.keys) {
		return false
	}
	if c.itens[c.keys[i]].idx < c.itens[c.keys[i]].max {
		c.itens[c.keys[i]].idx = c.itens[c.keys[i]].idx + 1
		return true
	} else if i < len(c.keys) {
		for k := i; k <= i; k++ {
			c.itens[c.keys[k]].idx = 0
		}
		return c.next(i + 1)
	}
	return false
}

func (c *Cartesian) GetCombination() (end bool, idxs map[string]int) {
	idxs = make(map[string]int)
	for key := range c.itens {
		idxs[key] = c.itens[key].idx
	}
	end = c.next(c.current)
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
	dst = "./k2web/pub/img/" + uuid.New().String() + filepath.Ext(src)
	_, err = copy(src, dst)
	dst = "./img/" + filepath.Base(dst)
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
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	_, err = net.Dial("tcp", u.Host)
	return err
}
