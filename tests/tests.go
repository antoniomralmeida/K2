package tests

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/antoniomralmeida/k2/kb"
)

func TestConvert(t *testing.T) {
	Test2()
	t.Fail()
}

func Test1() {

	json := `{
        "name": "Equipamento",  
        "icon": "eqp.jpg",
        "attributes":    [
			{
			  "name": "Nome",
			  "atype": "String",
			  "keephistory": 0,
			  "validityinterval": 0,
			  "sources": [
				"User"
			  ]
			}]
    }`
	c := kb.GKB.NewClass(json)
	if c != nil {
		fmt.Println(c)
	}

	json = `{
		"name": "MotorElétrico",
		"icon": "motor.jpg",
		"parent" : "Equipamento", 
		"attributes": [
		  {
			"name": "Data",
			"atype": "Date",
			"sources": [
			  "User"
			],
			"keephistory": 5,
			"validityinterval": 0
		  },
		  {
			"name": "Potência",
			"atype": "Number",
			"sources": [
			  "User",
			  "IOT",
			  "Simulation"
			],
			"simulation": "NormalDistribution",
			"keephistory": 50,
			"validityinterval": 500
		  },
		  {
			"name": "CurrenPower",
			"atype": "Number",
			"sources": [
			  "IOT"
			],
			"keephistory": 0,
			"validityinterval": 0
		  },
		  {
			"name": "Status",
			"atype": "List",
			"options": [
			  "PowerOff",
			  "PowerOn"
			],
			"sources": [
			  "User"
			],
			"keephistory": 0,
			"validityinterval": 0
		  }
		]
	  }`
	c = kb.GKB.NewClass(json)
	fmt.Println(c)

}

func Test2() {
	o := kb.GKB.NewObject("MotorElétrico", "M01")
	fmt.Println(o)
}

func Test3() {
	a := kb.GKB.FindAttributeObjectByName("M01.Potência")
	fmt.Println(a)
	for i := 0; i < 100; i++ {
		a.SetValue(rand.Float64(), kb.IOT, 100.0)
	}
}

func Test4() {
	o := kb.GKB.FindObjectByName("M01")
	fmt.Println(o)
}

func Test5() {
	a := kb.GKB.FindAttributeObjectByName("M01.Potência")
	a.LinearRegression()
}

func Test6() {
	kb.GKB.NewRule("for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230", 100, 0)
}

func Test7() {
	w := kb.GKB.NewWorkspace("Painel Elétrico", "C:\\Users\\manoe\\OneDrive\\Imagens\\painel.jpg")
	fmt.Println(w)
}
