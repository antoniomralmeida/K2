package main

import (
	"fmt"
	"main/classes"
)

func main() {
	ebnf := classes.EBNF{}
	ebnf.ReadToken("k2.ebnf")
	ebnf.PrintEBNF()
	fmt.Println(ebnf.FindToken("then"))
}
