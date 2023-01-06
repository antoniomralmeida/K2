package util

import (
	"io/ioutil"

	"github.com/antoniomralmeida/k2/initializers"
)

// ReadFile returns the bytes of a file searched in the path and beyond it
func ReadFile(path string) (bytes []byte) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		initializers.Log(err, initializers.Fatal)
	}
	return bytes
}
