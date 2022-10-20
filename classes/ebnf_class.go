package classes

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"main/lib"
	"strconv"
	"strings"
	"unicode"
)

type TokenType byte

const (
	Null TokenType = iota
	Reference
	Literal
	Text
	Control
	Jump
	Object
	DynamicReference
	Attribute
	Constant
	Class
	ListType
)

var TokenTypeStr = []string{"", "Reference", "Literal", "Text", "Control", "Jump", "Object", "DynamicReference", "Attribute", "Constant", "Class", "ListType"}

func (me TokenType) String() string {
	return TokenTypeStr[me]
}

func (me TokenType) Size() int {
	return len(TokenTypeStr)
}

type TokenBin byte

const (
	b_null TokenBin = iota
	b_open_par
	b_close_par
	b_iqual
	b_activate
	b_and
	b_any
	b_change
	b_conclude
	b_create
	b_deactivate
	b_delete
	b_different
	b_equal
	b_focus
	b_for
	b_greater
	b_halt
	b_hide
	b_if
	b_inform
	b_initially
	b_insert
	b_invoke
	b_is
	b_less
	b_move
	b_of
	b_operator
	b_or
	b_remove
	b_rotate
	b_set
	b_show
	b_start
	b_than
	b_that
	b_the
	b_then
	b_to
	b_transfer
	b_unconditionally
	b_when
	b_whenever
)

var TokenBinStr = []string{
	"",
	"(",
	")",
	"=",
	"activate",
	"and",
	"any",
	"change",
	"conclude",
	"create",
	"deactivate",
	"delete",
	"different",
	"equal",
	"focus",
	"for",
	"greater",
	"halt",
	"hide",
	"if",
	"inform",
	"initially",
	"insert",
	"invoke",
	"is",
	"less",
	"move",
	"of",
	"operator",
	"or",
	"remove",
	"rotate",
	"set",
	"show",
	"start",
	"than",
	"that",
	"the",
	"then",
	"to",
	"transfer",
	"unconditionally",
	"when",
	"whenever"}

func (me TokenBin) String() string {
	return TokenBinStr[me]
}

func (me TokenBin) Size() int {
	return len(TokenBinStr)
}

type BIN struct {
	tokentype       TokenType
	typebin         TokenBin
	token           string
	class           *KBClass
	object          *KBObject
	attribute       *KBAttribute
	attributeObject *KBAttributeObject
}

type Token struct {
	id        int
	tokentype TokenType
	rule_id   int
	rule_jump int
	token     string
	bin       TokenBin
	next      []*Token
}

func (t *Token) String() string {
	return " #" + strconv.Itoa(t.id) + ",token:" + t.token + ",type:" + t.tokentype.String() + ",bin:" + t.bin.String()
}

func (b *BIN) findTokenBin(i byte, j byte) TokenBin {
	if j >= i {
		avg := (i + j) / 2
		tb := TokenBin(avg)
		if b.token == tb.String() {
			return tb
		} else if b.token >= tb.String() {
			return b.findTokenBin(avg+1, j)
		} else {
			return b.findTokenBin(i, avg-1)
		}
	}
	return TokenBin(0)
}

func (b *BIN) setTokenBin() {
	if b.tokentype == Literal {
		b.typebin = b.findTokenBin(0, byte(b.typebin.Size()-1))
		if b.typebin == b_null {
			log.Fatal("Literal unknown!", b.token)
		}
	}
}

type Statement struct {
	id     int
	name   string
	tokens []*Token
}

type EBNF struct {
	rules []*Statement
	base  *Token
}

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

func (e *EBNF) Parsing(kb *KnowledgeBase, cmd string) ([]*Token, []*BIN, error) {
	cmd = strings.Replace(cmd, "\r\n", "", -1)
	cmd = strings.Replace(cmd, "\\n", "", -1)
	cmd = strings.Replace(cmd, "\t", " ", -1)
	for strings.Contains(cmd, "  ") {
		cmd = strings.Replace(cmd, "  ", " ", -1)
	}
	log.Println("Parsing Prodution Rule: ", cmd)
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
	var bin []*BIN
	for _, x := range tokens {
		var ok = false
		opts = e.findOptions(pt, &stack, 0)
		for _, y := range opts {
			//fmt.Println(x, y)
			if (y.token == x) ||
				(y.tokentype == DynamicReference && len(x) == 1) ||
				((y.tokentype == Object || y.tokentype == Class || y.tokentype == Attribute || y.tokentype == Constant || y.tokentype == Reference) && unicode.IsUpper(rune(x[0]))) ||
				(y.tokentype == Text && (rune(x[0]) == '\'' || rune(x[0]) == '"') ||
					(y.tokentype == Constant && lib.IsNumber(x))) {
				if y.tokentype == Class {
					if kb.FindClassByName(x, false) != nil {
						ok = true
					}
				} else if y.tokentype == Object {
					if kb.FindObjectByName(x) != nil {
						ok = true
					}
				} else {
					ok = true
				}
				if ok {
					pt = y
					break

				}
			}
		}
		if !ok || len(opts) == 0 {
			str := "Compiller error in " + x + " when the expected was: "
			for _, y := range opts {
				str = str + "... " + y.token
			}
			return opts, nil, errors.New(str)
		}
		code := BIN{tokentype: pt.tokentype, token: x}
		code.setTokenBin()
		bin = append(bin, &code)
	}
	for _, y := range pt.next {
		if y.token == "." && y.tokentype == Control {
			log.Println(", compilation successfully!")
			return nil, bin, nil
		}
	}
	opts = e.findOptions(pt, &stack, 0)
	str := "Incomplete sentence when the expected was: "
	for _, y := range opts {
		str = str + "... " + y.token
	}
	return opts, nil, errors.New(str)
}

func (e *EBNF) newStatement(str string) *Statement {
	var rule Statement
	rule.id = len(e.rules) + 1
	rule.name = strings.Trim(str, " ")
	e.rules = append(e.rules, &rule)
	return &rule
}

func (e *EBNF) newToken(rule *Statement, str string, tokentype TokenType, nexts ...*Token) {
	token := Token{id: len(rule.tokens) + 1, token: strings.Trim(str, " "), rule_id: rule.id, tokentype: tokentype}
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
			var nrule = e.newStatement(left)
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
			e.parsingStatement(nrule)
		}
	}

	for _, r := range e.rules {
		for _, t := range r.tokens {
			if t.tokentype == Reference {
				t.rule_jump = e.findRule(t.token)
				if t.rule_jump == -1 {
					for z := 1; z < t.tokentype.Size(); z++ {
						if t.token == TokenType(z).String() {
							t.tokentype = TokenType(z)
							break
						}
					}
					if t.tokentype == Reference {
						log.Fatal("Reference not found! ", t.token)
					}
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

func (e *EBNF) findClose(rule *Statement, symb int, token string, i int, level int) int {
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

func (e *EBNF) parsingStatement(rule *Statement) {
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
