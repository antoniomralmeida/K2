package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
)

type TokenType byte

const (
	Trace TokenType = iota
	Rule
	Literal
	Control
)

type EBNF struct {
	id        int
	token     string
	tokentype TokenType
	next      []EBNF
}

var g_ebnf []EBNF
var _lastid = 0

func saveWord(str string, tokentype TokenType, nexts ...EBNF) EBNF {
	var ebnf EBNF
	_lastid = _lastid + 1
	ebnf.id = _lastid
	ebnf.token = str
	ebnf.tokentype = tokentype
	for _, jump := range nexts {
		ebnf.next = append(ebnf.next, jump)
	}
	g_ebnf = append(g_ebnf, ebnf)
	return ebnf
}

func newJump(node *EBNF, nexts ...EBNF) {

	for _, jump := range nexts {
		node.next = append(node.next, jump)
	}
}

func readEBNF(ebnffile string) int {

	file, err := ioutil.ReadFile(ebnffile)

	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	ebnf := string(file)
	ebnf = strings.Replace(ebnf, "\r\n", "", -1)
	ebnf = strings.Replace(ebnf, "\\n", "", -1)
	ebnf = strings.Replace(ebnf, "\t", " ", -1)
	for strings.Contains(ebnf, "  ") {
		ebnf = strings.Replace(ebnf, "  ", " ", -1)
	}

	rules := strings.Split(ebnf, ".")
	for _, rule := range rules {
		tokens := strings.Split(rule, "=")

		if len(tokens) > 1 {
			var base = saveWord(tokens[0], Rule)
			var last = base
			tokens[1] = tokens[1] + "."
			var isWord = false
			var isRule = false
			var start = 0

			for i, c := range tokens[1] {
				switch {
				case c == '[' || c == ']' || c == '{' || c == '}' || c == '.' || c == ' ':
					{
						if isWord {
							var tokentype = Literal
							if isRule {
								tokentype = Rule
							}
							last = saveWord(tokens[1][start:i], tokentype)
							isWord = false
							isRule = false
						}

						if c != ' ' {
							saveWord(string(c), Control, last)
						}
					}
				case unicode.IsLower(c) && !isWord:
					{
						start = i
						isWord = true
					}
				case unicode.IsUpper(c) && !isWord:
					{
						start = i
						isWord = true
						isRule = true
					}
				default:
				}
				if c == '.' {
					break
				}
			}
			/*
				for i := len(g_ebnf) - 1; i >= 0; i-- {
					if g_ebnf[i].id == base.id {
						break
					}
					newJump(&g_ebnf[i-1], g_ebnf[i])
					fmt.Println(g_ebnf[i-1])
				}
			*/
		}

	}
	return 1
}

func main() {

	readEBNF("teste.ebnf")
	for _, no := range g_ebnf {
		fmt.Println(no)

	}
}
