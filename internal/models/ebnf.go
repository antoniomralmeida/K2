package models

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/antoniomralmeida/k2/internal/lib"
	"github.com/antoniomralmeida/k2/pkg/randnonlinear"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tjarratt/babble"
)

type EBNF struct {
	Rules []*Statement `json:"rules"`
	Base  *Token       `json:"-"`
	last  int          `json:"-"`
}

func EBNFFactory(ebnfFile string) *EBNF {
	ebnf := new(EBNF)
	ebnf.grammarLoad(ebnfFile)
	return ebnf
}
func (e *EBNF) GetBase() *Token {
	return e.Base
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
		if pt.GetTokentype() == Control && pt.Token == "." && len(*stack) > 0 {
			pt = (*stack)[len(*stack)-1]
			x := (*stack)[:len(*stack)-1]
			stack = &x
		}

		for _, x := range pt.Nexts {
			if x.GetTokentype() == Control || x.GetTokentype() == Jump {
				for _, k := range e.FindOptions(x, stack, level+1) {
					if !isElementExist(ret, k) {
						ret = append(ret, k)
					}
				}
			} else if x.GetTokentype() == Reference {
				*stack = append(*stack, x)
				n := e.Rules[x.Rule_jump].Tokens[0]
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
	rule.Id = len(e.Rules) + 1
	rule.Name = strings.Trim(str, " ")
	e.Rules = append(e.Rules, &rule)
	return &rule
}

func (e *EBNF) newToken(rule *Statement, str string, tokentype Tokentype, nexts ...*Token) {
	//Token := Token{Id: len(rule.Tokens) + 1, Token: strings.Trim(str, " "), Rule_id: rule.Id, Tokentype: tokentype}
	e.last++
	Token := Token{Id: e.last, Token: strings.Trim(str, " "), Rule_id: rule.Id, Tokentype: tokentype}
	Token.Nexts = append(Token.Nexts, nexts...)
	rule.Tokens = append(rule.Tokens, &Token)
}

func (e *EBNF) newJump(node *Token, before bool, nexts ...*Token) {
	if before {
		node.Nexts = append(nexts, node.Nexts...)
	} else {
		node.Nexts = append(node.Nexts, nexts...)
	}
}

func (e *EBNF) grammarLoad(ebnfFile string) int {

	file, err := ioutil.ReadFile(ebnfFile)
	if err != nil {
		inits.Log("Could not read the file due to this %s error \n"+err.Error(), inits.Fatal)
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
	errorFatal := false
	for _, r := range e.Rules {
		for _, t := range r.Tokens {
			if t.GetTokentype() == Reference {
				t.Rule_jump = e.findRule(t.Token)
				if t.Rule_jump == -1 {
					for z := 1; z < t.GetTokentype().Size(); z++ {
						if t.Token == Tokentype(z).String() {
							t.Tokentype = Tokentype(z)
							break
						}
					}
					if t.GetTokentype() == Reference {
						errorFatal = true
						inits.Log("Reference not found! "+t.Token, inits.Error)
					}
				}

			}
			if t.Tokentype == Literal {
				if _, ok := LiteralBinStr[t.Token]; !ok {
					errorFatal = true
					inits.Log("Literal not found! "+t.Token, inits.Error)
				}
			}
		}
	}
	if errorFatal {
		inits.Log("Fatal error(s) in EBNF parsing!", inits.Fatal)
	}
	e.Base = e.Rules[0].Tokens[0]
	data, err := os.Create(ebnfFile + ".json")
	if err != nil {
		inits.Log(err, inits.Error)
	} else {
		io.Copy(data, strings.NewReader(e.String()))
	}
	data.Close()
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
	inits.Log("Parsing ebnf rule "+rule.Name, inits.Info)
	//Finding Pair of Symbols
	for i := 0; i < len(rule.Tokens); i++ {
		if rule.Tokens[i].GetTokentype() == Control {
			s := e.FindSymbols(rule.Tokens[i].Token, false)
			if s != -1 {
				c := e.findClose(rule, s, rule.Tokens[i].Token, i, 0)
				if c == -1 {
					msg := fmt.Sprint("Parsing error in Token ", rule.Tokens[i].Token, " #", rule.Tokens[i].Id, s, i)
					inits.Log(msg, inits.Fatal)
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
			e.newJump(rule.Tokens[i], false, rule.Tokens[i+1])
		} else {
			e.newJump(rule.Tokens[i], false, rule.Tokens[pairs[p].end])
			e.newJump(rule.Tokens[pairs[p].begin], false, rule.Tokens[i+1])
		}
		if rule.Tokens[i].Token == "{" {
			e.newJump(rule.Tokens[pairs[p].begin], true, rule.Tokens[pairs[p].end])
			e.newJump(rule.Tokens[pairs[p].end], false, rule.Tokens[pairs[p].begin])
		}
		if rule.Tokens[i].Token == "[" {
			e.newJump(rule.Tokens[pairs[p].begin], false, rule.Tokens[pairs[p].end])
		}
	}
}

func (e *EBNF) findRule(key string) int {
	for i, r := range e.Rules {
		if r.Name == key {
			return i
		}
	}
	return -1
}

func (e *EBNF) String() string {
	ret, err := json.MarshalIndent(e.Rules, "", "    ")
	inits.Log(err, inits.Error)
	return string(ret)
}

func (e *EBNF) GrammarSample() string {
	stack := new([]*Token)
	return e.grammarSampleToken(e.Base, stack, 0)
}

func (e *EBNF) grammarSampleToken(t *Token, stack *[]*Token, npar int) string {
	var token string
	switch t.Tokentype {
	case Constant:
		token = strconv.Itoa(rand.Int()) + " "
	case Text:
		babbler := babble.NewBabbler()
		babbler.Separator = " "
		babbler.Count = rand.Intn(5) + 1
		token = "'" + babbler.Babble() + "' "
	case DynamicReference:
		token = lib.GeneratePassword(2, 0, 0, 2) + " "
	case ListType:
		n := rand.Intn(5) + 1
		babbler := babble.NewBabbler()
		babbler.Separator = " "
		babbler.Count = 1
		token = "('" + babbler.Babble() + "'"
		for i := 1; i < n; i++ {
			token += ",'" + babbler.Babble() + "'"
		}
		token += ") "
	case Literal:
		token = t.Token + " "
	case Class, Object, Attribute, Workspace:
		token = t.Token + primitive.NewObjectID().Hex() + " "
	case Control:
		if token == "(" {
			npar++
		} else if token == ")" {
			npar--
		}
	}
	opts := e.FindOptions(t, stack, 0)
	max := len(opts)
	for n := 0; n < max*2; n++ {
		i := randnonlinear.IntnNL(max, func(idx int) int {
			return int(math.Pow(2, float64(max-idx)))
		})
		if opts[i].Token == ")" && npar == 0 {
			continue
		}
		return token + e.grammarSampleToken(opts[i], stack, npar)
	}
	for npar > 0 {
		npar--
		token += token + ")"
	}
	return token
}
