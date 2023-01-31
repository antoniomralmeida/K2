package models

type SYMBOL struct {
	begin string
	end   string
}

var symbols = []SYMBOL{{"=", "."}, {"{", "}"}, {"[", "]"}, {"(", ")"}, {"\"", "\""}, {"'", "'"}}
