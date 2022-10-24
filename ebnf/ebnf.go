package ebnf

import (
	"fmt"
	"io/ioutil"
	"log"

	"strings"
	"unicode"
)

func (e *EBNF) GetBase() *Token {
	return e.base
}

func (e *EBNF) findSymbols(str string, both bool) int {
	for i, x := range symbols {
		if x.begin == str {
			return i
		}
		if x.end == str && both {
			return i
		}
	}
	return -1
}

func (e *EBNF) FindOptions(pt *Token, stack *[]*Token, level int) []*Token {
	var ret []*Token
	if level < 10 {
		if pt.GetTokentype() == Control && pt.Token == "." && len(*stack) > 0 {
			pt = (*stack)[len(*stack)-1]
			x := (*stack)[:len(*stack)-1]
			stack = &x
		}

		for _, x := range pt.next {
			if x.GetTokentype() == Control || x.GetTokentype() == Jump {
				for _, k := range e.FindOptions(x, stack, level+1) {
					if !isElementExist(ret, k) {
						ret = append(ret, k)
					}
				}
			} else if x.GetTokentype() == Reference {
				*stack = append(*stack, x)
				n := e.rules[x.rule_jump].Tokens[0]
				for _, k := range e.FindOptions(n, stack, level+1) {
					if !isElementExist(ret, k) {
						ret = append(ret, k)
					}
				}
			} else {
				if !isElementExist(ret, x) {
					ret = append(ret, x)
				}
			}
		}
	}
	return ret
}

func (e *EBNF) newStatement(str string) *Statement {
	var rule Statement
	rule.id = len(e.rules) + 1
	rule.name = strings.Trim(str, " ")
	e.rules = append(e.rules, &rule)
	return &rule
}

func (e *EBNF) newToken(rule *Statement, str string, Tokentype Tokentype, nexts ...*Token) {
	Token := Token{id: len(rule.Tokens) + 1, Token: strings.Trim(str, " "), rule_id: rule.id, Tokentype: Tokentype}
	for _, jump := range nexts {
		Token.next = append(Token.next, jump)
	}
	rule.Tokens = append(rule.Tokens, &Token)
}

func (e *EBNF) newJump(node *Token, nexts ...*Token) {
	for _, jump := range nexts {
		node.next = append(node.next, jump)
	}
}

func (e *EBNF) ReadToken(Tokenfile string) int {

	file, err := ioutil.ReadFile(Tokenfile)
	if err != nil {
		log.Println("Could not read the file due to this %s error \n", err)
	}
	ebnf_txt := string(file)
	ebnf_txt = strings.Replace(ebnf_txt, "\r\n", "", -1)
	ebnf_txt = strings.Replace(ebnf_txt, "\\n", "", -1)
	ebnf_txt = strings.Replace(ebnf_txt, "\t", " ", -1)
	for strings.Contains(ebnf_txt, "  ") {
		ebnf_txt = strings.Replace(ebnf_txt, "  ", " ", -1)
	}

	rules := strings.Split(ebnf_txt, ".")
	for _, rule := range rules {

		var left string
		var right string
		for i := 0; i < len(rule); i++ {
			if rule[i] == '=' {
				left = rule[0:i]
				right = rule[i:] + "."
				break
			}
		}
		if len(left) > 0 {
			var nrule = e.newStatement(left)
			var inWord = false
			var inString = false
			var inRule = false
			var start = 0
			for i, c := range right {
				switch {
				case e.findSymbols(string(c), true) != -1 || c == ' ' || c == '|':
					if inString {
						if c == '"' {
							e.newToken(nrule, right[start:i], Literal)
							inString = false
						}
					} else if inWord {
						var Tokentype = Literal
						if inRule {
							Tokentype = Reference
						}
						e.newToken(nrule, right[start:i], Tokentype)
						inWord = false
						inRule = false
					} else {
						if c == '"' {
							start = i + 1
							inString = true
						}
					}
					if c != ' ' && c != '"' && c != '\'' && !inString {
						if c == '|' {
							e.newToken(nrule, string(c), Jump)
						} else {
							e.newToken(nrule, string(c), Control)
						}
					}
				case unicode.IsLower(c) && !inWord && !inString:
					start = i
					inWord = true
				case unicode.IsUpper(c) && !inWord && !inString:
					start = i
					inWord = true
					inRule = true
				default:
				}
			}
			e.parsingStatement(nrule)
		}
	}

	for _, r := range e.rules {
		for _, t := range r.Tokens {
			if t.GetTokentype() == Reference {
				t.rule_jump = e.findRule(t.Token)
				if t.rule_jump == -1 {
					for z := 1; z < t.GetTokentype().Size(); z++ {
						if t.Token == Tokentype(z).String() {
							t.GetTokentype() = Tokentype(z)
							break
						}
					}
					if t.GetTokentype() == Reference {
						log.Fatal("Reference not found! ", t.Token)
					}
				}

			}
		}
	}
	e.base = e.rules[0].Tokens[0]
	return 1
}

func (e *EBNF) findClose(rule *Statement, symb int, Token string, i int, level int) int {
	for j := i + 1; j < len(rule.Tokens); j++ {
		if rule.Tokens[j].GetTokentype() == Control {
			s := symbols[symb]
			if rule.Tokens[j].Token == s.end && level == 0 {
				return j
			} else if rule.Tokens[j].Token == s.begin {
				return e.findClose(rule, symb, Token, j, level+1)
			} else if rule.Tokens[j].Token == s.end {
				return e.findClose(rule, symb, Token, j, level-1)
			}
		}
	}
	return -1
}

func (e *EBNF) parsingStatement(rule *Statement) {
	var pairs []PAIR
	log.Println("Parsing rule ", rule.name)
	for i := 0; i < len(rule.Tokens); i++ {
		if rule.Tokens[i].GetTokentype() == Control {
			s := e.findSymbols(rule.Tokens[i].Token, false)
			if s != -1 {
				c := e.findClose(rule, s, rule.Tokens[i].Token, i, 0)
				if c == -1 {
					log.Println("Parssing erro in Token ", rule.Tokens[i].Token, " #", rule.Tokens[i].id, s, i)
					return
				}
				pairs = append(pairs, PAIR{i, c})
				if rule.Tokens[i].Token == "\"" || rule.Tokens[i].Token == "'" && c != -1 {
					i = c + 1
				}
			}
		}
	}
	for i := 0; i < len(rule.Tokens)-1; i++ {
		var p = findPair(pairs, i)
		if rule.Tokens[i].GetTokentype() != Jump {
			e.newJump(rule.Tokens[i], rule.Tokens[i+1])
		} else {
			e.newJump(rule.Tokens[i], rule.Tokens[pairs[p].end])
			e.newJump(rule.Tokens[pairs[p].begin], rule.Tokens[i+1])
		}
		if rule.Tokens[i].Token == "{" {
			e.newJump(rule.Tokens[pairs[p].begin], rule.Tokens[pairs[p].end])
			e.newJump(rule.Tokens[pairs[p].end], rule.Tokens[pairs[p].begin])
		}
		if rule.Tokens[i].Token == "[" {
			e.newJump(rule.Tokens[pairs[p].begin], rule.Tokens[pairs[p].end])
		}
	}
}

func (e *EBNF) findRule(key string) int {
	for i, r := range e.rules {
		if r.name == key {
			return i
		}
	}
	return -1
}

func (e *EBNF) PrintEBNF() {
	fmt.Println("----------EBNF tree--------------")
	for _, r := range e.rules {
		fmt.Println("====> Rule: ", r.name)
		for _, t := range r.Tokens {
			fmt.Println(t.String())
			for _, t2 := range t.next {
				fmt.Println("...jump to ", t2.String())
			}
		}
	}
}
