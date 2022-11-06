package fuzzy

import (
	"fmt"
	"strings"
)

/*
Fuzzy logic operators - Zadeh Basic Operators
Fuzzy logic works with membership values in a way that mimics Boolean logic. To this end, replacements for basic operators AND, OR, NOT must be available.

Boolean		Fuzzy		GOLANG expression
AND(x,y)	MIN(x,y)	(x<y?x:y)
OR(x,y)		MAX(x,y)	(x>y?x:y)
NOT(x)		1 – x		** not implemented
*/
func SliceToString(slice []string) string {
	ret := ""
	for x := range slice {
		if slice[x] == "and" || slice[x] == "or" {
			ret = ret + " " + slice[x] + " "
		} else {
			ret = ret + slice[x]
		}
	}
	return ret
}

func FuzzyLogicalOperation(exp string) (string, bool) {
	s := strings.Split(exp, " ")
	for x := range s {
		switch s[x] {
		case "or":
			s[x] = s[x-1]
			s[x-1] = "(" + s[x-1] + ">" + s[x+1] + "?"
			s[x+1] = ":" + s[x+1] + ")"
			return SliceToString(s), true
		}
	}
	for x := range s {
		switch s[x] {
		case "and":
			s[x] = s[x-1]
			s[x-1] = "(" + s[x-1] + "<" + s[x+1] + "?"
			s[x+1] = ":" + s[x+1] + ")"
			return SliceToString(s), true
		}
	}
	return exp, false
}

func FuzzyLogicalInference(exp string) string {
	stack := make(map[string]string)
	i_stack := 100
	key := "base"
	stack[key] = exp
	//Find subexpressions
oulter:
	for {
		i := -1
		j := -1
		np := 0
	inter:
		for x := range stack[key] {
			switch stack[key][x] {
			case '(':
				np++
				if np == 1 {
					i = x
				}
			case ')':
				np--
				if np == 0 {
					j = x
					break inter
				}
			}
		}
		if i != -1 {
			i_stack++
			k := fmt.Sprintf("zz%v", i_stack)
			stack[k] = stack[key][i+1 : j]
			stack[key] = strings.Replace(stack[key], "("+stack[k]+")", k, 1)
			if !strings.Contains(stack[key], "(") {
				todo := false
			todo:
				for k := range stack {
					if strings.Contains(stack[k], "(") {
						key = k
						todo = true
						break todo
					}
				}
				if !todo {
					break oulter
				}
			}

		} else {
			break
		}
	}

	//keys of subexpressions
	keys := []string{}
	for key := range stack {
		keys = append(keys, key)
	}
	//Application of fuzzy logic operations
	for {
		var found bool
		for i := len(keys) - 1; i >= 0; i-- {
			var f bool
			stack[keys[i]], f = FuzzyLogicalOperation(stack[keys[i]])
			found = found || f
		}
		if !found {
			break
		}
	}
	//substitution of subexpressions
	exp = stack["base"]
	for {
		found := false
		for i := len(keys) - 1; i > 0; i-- {
			if strings.Contains(exp, keys[i]) {
				exp = strings.Replace(exp, keys[i], stack[keys[i]], -1)
				found = true
			}
		}
		if !found {
			break
		}
	}

	/* Recomended use // https://github.com/PaesslerAG/gval to evaluate fuzzy logical expression
	use: i, err := gval.Evaluate(exp, nil)
	*/
	return exp
}
