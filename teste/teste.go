package main

import (
	"fmt"
	"unicode"
)

func main() {
	str := "Manoel"
	for _, c := range str {
		switch {
		case c == '[':
		case c == '{':
		case c == ']':
		case c == '}':
		case c == '.':
		case unicode.IsLower(c):
			fmt.Println(c)
		case unicode.IsUpper(c):

		default:
		}
	}
}
