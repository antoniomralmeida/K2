package classes

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"
)

type TokenType byte

const (
	Null TokenType = iota
	Reference
	Literal
	String
	Control
	Jump
	Object
)

func (me TokenType) String() string {
	return [...]string{"", "Reference", "Literal", "String", "Control", "Jump", "Object"}[me]
}

type Token struct {
	id        int
	tokentype TokenType
	rule_id   int
	rule_jump int
	token     string
	next      []*Token
}

func (t *Token) String() string {
	return " #" + strconv.Itoa(t.id) + ",token=" + t.token + ",type=" + t.tokentype.String()
}

type Rule struct {
	id     int
	name   string
	tokens []*Token
}

type EBNF struct {
	rules []*Rule
	base  *Token
}

var _debug int

type SYMBOL struct {
	begin string
	end   string
}

var symbols = []SYMBOL{SYMBOL{"=", "."}, SYMBOL{"{", "}"}, SYMBOL{"[", "]"}, SYMBOL{"(", ")"}, SYMBOL{"\"", "\""}, SYMBOL{"'", "'"}}

func (e *EBNF) findSymbol(str string, both bool) int {
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

func isElementExist(s []*Token, str *Token) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (e *EBNF) findOptions(pt *Token, stack *[]*Token, level int) []*Token {
	var ret []*Token
	if level < 10 {
		if pt.tokentype == Control && pt.token == "." && len(*stack) > 0 {
			pt = (*stack)[len(*stack)-1]
			x := (*stack)[:len(*stack)-1]
			stack = &x
		}

		for _, x := range pt.next {
			if x.tokentype == Control || x.tokentype == Jump {
				for _, k := range e.findOptions(x, stack, level+1) {
					if !isElementExist(ret, k) {
						ret = append(ret, k)
					}
				}
			} else if x.tokentype == Reference {
				*stack = append(*stack, x)
				n := e.rules[x.rule_jump].tokens[0]
				for _, k := range e.findOptions(n, stack, level+1) {
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

func (e *EBNF) Parsing(cmd string) []*Token {
	cmd = strings.Replace(cmd, "\r\n", "", -1)
	cmd = strings.Replace(cmd, "\\n", "", -1)
	cmd = strings.Replace(cmd, "\t", " ", -1)
	for strings.Contains(cmd, "  ") {
		cmd = strings.Replace(cmd, "  ", " ", -1)
	}
	log.Println("Parsing command...")
	var inWord = false
	var inString = false
	var inNumber = false
	var start = 0
	var tokens []string
	const endline = '春'
	cmd = cmd + string(endline)
	for i, c := range cmd {
		switch {
		case c == '春' || c == ' ' || e.findSymbol(string(c), true) != -1:
			if inNumber && c != '.' {
				tokens = append(tokens, cmd[start:i])
				inNumber = false
			} else if inString {
				if c == '"' || c == '\'' {
					tokens = append(tokens, cmd[start:i+1])
					inString = false
				}
			} else if inWord {
				tokens = append(tokens, cmd[start:i])
				inWord = false
			} else {
				if c == '"' || c == '\'' {
					start = i
					inString = true
				} else if c != ' ' && c != '.' && c != endline {
					tokens = append(tokens, string(c))
				}
			}
		case unicode.IsLower(c) && !inWord && !inString && !inNumber:
			start = i
			inWord = true
		case unicode.IsUpper(c) && !inWord && !inString && !inNumber:
			start = i
			inWord = true
		case unicode.IsNumber(c) && !inNumber && !inString && !inWord:
			start = i
			inNumber = true
		default:
		}
	}
	var pt = e.base
	var stack []*Token
	var opts []*Token

	//tokens = tokens[:len(tokens)-1] // remove endline
	for _, x := range tokens {
		fmt.Print(x, " ")
		var ok = false
		opts = e.findOptions(pt, &stack, 0)
		for _, y := range opts {
			if (y.token == x) ||
				(y.tokentype == Object && y.token == "DynamicReference" && len(x) == 1) ||
				(y.tokentype == Object && len(x) > 1) {
				ok = true
				pt = y
				break
			}
		}
		if !ok || len(opts) == 0 {
			log.Println(", compiller error in ", x, " when the expected was: ")
			for _, y := range opts {
				log.Println("... ", y.token)
			}
			return opts
		}
	}
	for _, y := range pt.next {
		if y.token == "." && y.tokentype == Control {
			log.Println(", compilation successfully!")
			return opts
		}
	}
	opts = e.findOptions(pt, &stack, 0)
	log.Println(", incomplete sentence when the expected was: ")
	for _, y := range opts {
		log.Println("... ", y.token)
	}
	return opts
}

func (e *EBNF) newRule(str string) *Rule {
	var rule Rule
	rule.id = len(e.rules) + 1
	rule.name = strings.Trim(str, " ")
	e.rules = append(e.rules, &rule)
	return &rule
}

func (e *EBNF) newToken(rule *Rule, str string, tokentype TokenType, nexts ...*Token) {
	var token Token
	token.id = len(rule.tokens)
	token.token = strings.Trim(str, " ")
	token.rule_id = rule.id
	token.tokentype = tokentype

	for _, jump := range nexts {
		token.next = append(token.next, jump)
	}
	rule.tokens = append(rule.tokens, &token)
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
			var nrule = e.newRule(left)
			var inWord = false
			var inString = false
			var inRule = false
			var start = 0
			for i, c := range right {
				switch {
				case e.findSymbol(string(c), true) != -1 || c == ' ' || c == '|':
					if inString {
						if c == '"' {
							e.newToken(nrule, right[start:i], Literal)
							inString = false
						}
					} else if inWord {
						var tokentype = Literal
						if inRule {
							tokentype = Reference
						}
						e.newToken(nrule, right[start:i], tokentype)
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
			e.parsingRule(nrule)
		}
	}

	for _, r := range e.rules {
		for _, t := range r.tokens {
			if t.tokentype == Reference {
				t.rule_jump = e.findRule(t.token)
				if t.rule_jump == -1 {
					t.tokentype = Object
				}
			}
		}
	}
	e.base = e.rules[0].tokens[0]
	return 1
}

type PAIR struct {
	begin int
	end   int
}

func findPair(p []PAIR, i int) int {
	var ret = 0
	for k, x := range p {
		if x.begin <= i && x.end >= i && (p[k].begin > p[ret].begin || p[k].end < p[ret].end) {
			ret = k
		}
	}
	return ret
}

func (e *EBNF) findClose(rule *Rule, symb int, token string, i int, level int) int {
	for j := i + 1; j < len(rule.tokens); j++ {
		if rule.tokens[j].tokentype == Control {
			s := symbols[symb]
			if rule.tokens[j].token == s.end && level == 0 {
				return j
			} else if rule.tokens[j].token == s.begin {
				return e.findClose(rule, symb, token, j, level+1)
			} else if rule.tokens[j].token == s.end {
				return e.findClose(rule, symb, token, j, level-1)
			}
		}
	}
	return -1
}

func (e *EBNF) parsingRule(rule *Rule) {
	var pairs []PAIR
	log.Println("Parsing rule ", rule.name)
	for i := 0; i < len(rule.tokens); i++ {
		if rule.tokens[i].tokentype == Control {
			s := e.findSymbol(rule.tokens[i].token, false)
			if s != -1 {
				c := e.findClose(rule, s, rule.tokens[i].token, i, 0)
				if c == -1 {
					log.Println("Parssing erro in token ", rule.tokens[i].token, " #", rule.tokens[i].id, s, i)
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
		if rule.tokens[i].tokentype != Jump {
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
