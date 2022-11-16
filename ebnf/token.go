package ebnf

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (t *Token) GetToken() string {
	return t.Token
}

func (t *Token) GetTokentype() Tokentype {
	return t.Tokentype
}

func (t *Token) GetNexts() []*Token {
	return t.Nexts
}

func (t *Token) String() string {
	return "#" + strconv.Itoa(t.Id) + ", token: " + t.Token + ", type:" + t.Tokentype.String()
}

func (t *Token) MarshalJSON() ([]byte, error) {
	var result map[string]string = make(map[string]string)
	result["Id"] = strconv.Itoa(t.Id)
	result["Tokentype"] = t.Tokentype.String()
	result["Rule_id"] = strconv.Itoa(t.Rule_id)
	result["Rule_jump"] = strconv.Itoa(t.Rule_jump)
	result["Token"] = t.Token
	result["Nexts"] = fmt.Sprintf("%v", t.Nexts)
	return json.Marshal(&result)
}
