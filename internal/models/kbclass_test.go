package models

import (
	"testing"

	"github.com/antoniomralmeida/k2/internal/inits"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load("../../configs/.env")
	inits.ConnectDB()
}
func TestClassFactory(t *testing.T) {
	result := ClassFactory("xxx", "", "") //Name <5

	if result != nil {
		t.Errorf("ClassFactory(\"xxx\", \"\", \"\") FAILED. Expcted nil, got %v", result)
	} else {
		t.Logf("ClassFactory(\"xxx\", \"\", \"\") PASSED. Expcted nil, got %v", result)
	}

	result = ClassFactory("Motor Elétrico", "", "") //Normal

	if result == nil {
		t.Errorf("models.ClassFactory(\"Motor Elétrico\", \"\", \"\") FAILED. Expcted *KBClass, got %v", result)
	} else {
		t.Logf("ClassFactory(\"Motor Elétrico\", \"\", \"\") PASSED. Expcted *KBClass, got %v", result)
	}

	result = ClassFactory("Motor Elétrico", "", "") //Duplicado

	if result != nil {
		t.Errorf("models.ClassFactory(\"Motor Elétrico\", \"\", \"\") FAILED. Expcted nil, got %v", result)
	} else {
		t.Logf("ClassFactory(\"Motor Elétrico\", \"\", \"\") PASSED. Expcted nil, got %v", result)
	}
}
