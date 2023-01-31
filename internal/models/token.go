package models

type Token struct {
	Id        int       `json:"id"`
	Tokentype Tokentype `json:"tokentype"`
	Rule_id   int       `json:"rule_id"`
	Rule_jump int       `json:"rule_jump"`
	Token     string    `json:"token"`
	Nexts     []*Token  `json:"-"`
}
