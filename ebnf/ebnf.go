package ebnf

import (
	"fmt"
	"io/ioutil"

	"strings"
	"unicode"

	"github.com/antoniomralmeida/k2/initializers"
)

func (e *EBNF) GetBase() *Token {
	return e.base
}

func (e *EBNF) FindSymbols(str string, both bool) int {
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
		if pt.GetTokentype() == Control && pt.token == "." && len(*stack) > 0 {
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
				n := e.rules[x.rule_jump].tokens[0]
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
	Token := Token{id: len(rule.tokens) + 1, token: strings.Trim(str, " "), rule_id: rule.id, tokentype: Tokentype}
	for _, jump := range nexts {
		Token.next = append(Token.next, jump)
	}
	rule.tokens = append(rule.tokens, &Token)
}

func (e *EBNF) newJump(node *Token, nexts ...*Token) {
	for _, jump := range nexts {
		node.next = append(node.next, jump)
	}
}

func (e *EBNF) ReadToken(Tokenfile string) int {

	file, err := ioutil.ReadFile(Tokenfile)
	if err != nil {
		initializers.Log("Could not read the file due to this %s error \n"+err.Error(), initializers.Fatal)
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
				case e.FindSymbols(string(c), true) != -1 || c == ' ' || c == '|':
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
		for _, t := range r.tokens {
			if t.GetTokentype() == Reference {
				t.rule_jump = e.findRule(t.token)
				if t.rule_jump == -1 {
					for z := 1; z < t.GetTokentype().Size(); z++ {
						if t.token == Tokentype(z).String() {
							t.tokentype = Tokentype(z)
							break
						}
					}
					if t.GetTokentype() == Reference {
						initializers.Log("Reference not found! "+t.token, initializers.Fatal)
					}
				}

			}
		}
	}
	e.base = e.rules[0].tokens[0]
	return 1
}

func (e *EBNF) findClose(rule *Statement, symb int, Token string, i int, level int) int {
	for j := i + 1; j < len(rule.tokens); j++ {
		if rule.tokens[j].GetTokentype() == Control {
			s := symbols[symb]
			if rule.tokens[j].token == s.end && level == 0 {
				return j
			} else if rule.tokens[j].token == s.begin {
				return e.findClose(rule, symb, Token, j, level+1)
			} else if rule.tokens[j].token == s.end {
				return e.findClose(rule, symb, Token, j, level-1)
			}
		}
	}
	return -1
}

func (e *EBNF) parsingStatement(rule *Statement) {
	var pairs []PAIR
	initializers.Log("Parsing ebnf rule "+rule.name, initializers.Info)
	for i := 0; i < len(rule.tokens); i++ {
		if rule.tokens[i].GetTokentype() == Control {
			s := e.FindSymbols(rule.tokens[i].token, false)
			if s != -1 {
				c := e.findClose(rule, s, rule.tokens[i].token, i, 0)
				if c == -1 {
					msg := fmt.Sprint("Parssing erro in Token ", rule.tokens[i].token, " #", rule.tokens[i].id, s, i)
					initializers.Log(msg, initializers.Fatal)
					return
				}
				pairs = append(pairs, PAIR{i, c})
				if rule.tokens[i].token == "\"" || rule.tokens[i].token == "'" && c != -1 {
					i = c + 1
				}
			}
		}
	}
	for i := 0; i < len(rule.tokens)-1; i++ {
		var p = findPair(pairs, i)
		if rule.tokens[i].GetTokentype() != Jump {
			e.newJump(rule.tokens[i], rule.tokens[i+1])
		} else {
			e.newJump(rule.tokens[i], rule.tokens[pairs[p].end])
			e.newJump(rule.tokens[pairs[p].begin], rule.tokens[i+1])
		}
		if rule.tokens[i].token == "{" {
			e.newJump(rule.tokens[pairs[p].begin], rule.tokens[pairs[p].end])
			e.newJump(rule.tokens[pairs[p].end], rule.tokens[pairs[p].begin])
		}
		if rule.tokens[i].token == "[" {
			e.newJump(rule.tokens[pairs[p].begin], rule.tokens[pairs[p].end])
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
		for _, t := range r.tokens {
			fmt.Println(t.String())
			for _, t2 := range t.next {
				fmt.Println("...jump to ", t2.String())
			}
		}
	}
}
