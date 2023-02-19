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
	"time"

	"github.com/antoniomralmeida/k2/pkg/dsn"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mattn/go-tty"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int, alfaNumeric bool) string {
	var password strings.Builder
	var (
		lowerCharSet   = "abcdedfghijklmnopqrst"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialCharSet = "!@#$%&*"
		numberSet      = "0123456789"
		allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	)
	time.Sleep(time.Microsecond)
	//Set special character
	for i := 0; i < minSpecialChar && !alfaNumeric; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

const SignedKey = "uzUWld6Y0Ad6yUF8GU2gJGg8Q4wZaNNv"

func CreateJWTToken(id primitive.ObjectID, name string) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id.Hex()
	claims["name"] = name
	claims["exp"] = exp
	t, err := token.SignedString([]byte(SignedKey))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func DecodeToken(jwtcookie string) map[string]interface{} {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(jwtcookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SignedKey), nil
	})
	if err != nil {
		return nil
	}
	return claims
}

func ValidateToken(jwtcookie string) (err error) {

	if jwtcookie == "" {
		return errors.New("Invalid Cookie")
	}

	token, err := jwt.Parse(jwtcookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return []byte(SignedKey), nil
	})
	if err != nil {
		return
	}

	if token == nil {
		return errors.New("Invalid Token")
	}

	return nil
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
