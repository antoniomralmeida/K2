package classes

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
	Reference
	Literal
	Control
)

type Token struct {
	id        int
	token     string
	tokentype TokenType
	next      []*Token
}

type EBNF struct {
	tokens []Token
	rules  []Token
	_last  int
}

type EBNF_STACK struct {
	itens []Token
}

func (s *EBNF_STACK) push(item Token) {
	s.itens = append(s.itens, item)
}

func (s *EBNF_STACK) get() *Token {
	return &s.itens[len(s.itens)-1]
}

func (s *EBNF_STACK) pop() *Token {
	var item = s.itens[len(s.itens)-1]
	s.itens = s.itens[:len(s.itens)-1]
	return &item
}

func (e *EBNF) saveToken(str string, tokentype TokenType, nexts ...Token) Token {
	var token Token
	e._last++
	token.id = e._last
	token.token = str
	token.tokentype = tokentype
	for _, jump := range nexts {
		token.next = append(token.next, &jump)
	}
	e.tokens = append(e.tokens, token)
	if tokentype == Rule {
		if e.findRule(token.token) == -1 {
			e.rules = append(e.rules, token)
		}
	}
	return token
}

func (e *EBNF) newJump(node *Token, nexts ...Token) {
	for _, jump := range nexts {
		node.next = append(node.next, &jump)
	}
}

func (e *EBNF) ReadToken(Tokenfile string) int {

	file, err := ioutil.ReadFile(Tokenfile)
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	Token := string(file)
	Token = strings.Replace(Token, "\r\n", "", -1)
	Token = strings.Replace(Token, "\\n", "", -1)
	Token = strings.Replace(Token, "\t", " ", -1)
	for strings.Contains(Token, "  ") {
		Token = strings.Replace(Token, "  ", " ", -1)
	}

	rules := strings.Split(Token, ".")
	for _, rule := range rules {
		tokens := strings.Split(rule, "=")
		if len(tokens) > 1 {
			var base = e.saveToken(tokens[0], Rule)
			tokens[1] = tokens[1] + "."
			var isWord = false
			var isRule = false
			var start = 0
			for i, c := range tokens[1] {
				switch {
				case c == '[' || c == ']' || c == '{' || c == '}' || c == '.' || c == '|' || c == ' ' || c == '"' || c == '(' || c == ')':
					if isWord {
						var tokentype = Literal
						if isRule {
							tokentype = Reference
						}
						e.saveToken(tokens[1][start:i], tokentype)
						isWord = false
						isRule = false
					} else {
						if c == '"' {
							start = i
							isWord = true
						}
					}
					if c != ' ' {
						e.saveToken(string(c), Control)
					}
				case unicode.IsLower(c) && !isWord:
					start = i
					isWord = true
				case unicode.IsUpper(c) && !isWord:
					start = i
					isWord = true
					isRule = true
				default:
				}
				if c == '.' {
					break
				}
			}
			var control_points = EBNF_STACK{}
			var jumper_points = EBNF_STACK{}

			for i := len(e.tokens) - 1; i > 0; i-- {
				var t = e.tokens[i].token
				if e.tokens[i].id == base.id {
					break
				}
				switch {
				case t == ".":
					control_points.push(e.tokens[i])
				case t == "]":
					control_points.push(e.tokens[i])
					e.newJump(&e.tokens[i-1], e.tokens[i])
				case t == "|":
					e.newJump(&e.tokens[i-1], e.tokens[i])
					e.newJump(&e.tokens[i], *control_points.get())
					jumper_points.push(e.tokens[i+1])
				case t == "[":
					e.newJump(&e.tokens[i], *control_points.pop())
					for len(jumper_points.itens) > 0 {
						e.newJump(&e.tokens[i], *jumper_points.pop())
					}
					e.newJump(&e.tokens[i-1], e.tokens[i])
				case t == "}":
					control_points.push(e.tokens[i])
					e.newJump(&e.tokens[i-1], e.tokens[i])
				case t == "{":
					e.newJump(control_points.pop(), e.tokens[i])
					e.newJump(&e.tokens[i-1], e.tokens[i])
				case t == ")":
					control_points.push(e.tokens[i])
					e.newJump(&e.tokens[i-1], e.tokens[i])
				case t == "(":
					control_points.pop()
					for len(jumper_points.itens) > 0 {
						e.newJump(&e.tokens[i], *jumper_points.pop())
					}
					e.newJump(&e.tokens[i-1], e.tokens[i])
				default:
					if e.tokens[i-1].token != "|" {
						e.newJump(&e.tokens[i-1], e.tokens[i])
					}
				}
			}
		}
	}
	return 1
}

func (e *EBNF) FindToken(key string) int {
	for i, t := range e.tokens {
		if t.token == key {
			return i
		}
	}
	return -1
}

func (e *EBNF) findRule(key string) int {
	for i, t := range e.rules {
		if t.token == key {
			return i
		}
	}
	return -1
}

func (e *EBNF) PrintEBNF() {
	fmt.Println("----------EBNF tree--------------")
	for _, t := range e.tokens {
		fmt.Println("id: ", t.id, " token: ", t.token, " type: ", t.tokentype)
		for _, t2 := range t.next {
			fmt.Println("...jump to #", t2.id)
		}
	}
}
