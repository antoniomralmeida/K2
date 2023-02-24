package models

import (
	"encoding/json"
	"testing"
)

func TestBin(t *testing.T) {

	code := BIN{TokenType: Literal, Token: "any"}
	err := code.CheckLiteralBin()
	if err != nil {
		t.Errorf("CheckLiteralBin(%v) => %v", code, err)
	}

	code = BIN{TokenType: Literal, Token: "kkkkk"}
	err = code.CheckLiteralBin()
	if err == nil {
		t.Errorf("CheckLiteralBin(%v) => %v", code, err)
	}
	code = BIN{TokenType: Text, Token: "kkkkk"}
	err = code.CheckLiteralBin()
	if err != nil {
		t.Errorf("CheckLiteralBin(%v) => %v", code, err)
	}

	j := code.String()
	code2 := BIN{}
	err = json.Unmarshal([]byte(j), &code2)
	if err != nil {
		t.Errorf("String(%v) => %v, %v", code, j, err)
	}
	if code.TokenType != code2.TokenType {
		t.Errorf("String(%v) => %v", code, code2)
	}
}
