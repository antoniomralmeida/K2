package tests

import (
	"fmt"
	"math/rand"

	"github.com/antoniomralmeida/k2/kb"
)

func Test1(kbase *kb.KnowledgeBased) {
	a := kbase.FindAttributeObjectByName("M01.Potência")
	a.LinearRegression()
}

func Test2(kbase *kb.KnowledgeBased) {
	a := kbase.FindAttributeObjectByName("M01.Potência")
	fmt.Println(a)
	for i := 0; i < 100; i++ {
		a.SetValue(rand.Float64(), kb.IOT, 100.0)
	}
}

func Test3(kbase *kb.KnowledgeBased) {

	json := `{
        "name": "Equipamento",  
        "icon": "eqp.jpg",
        "attributes":    [
			{
			  "name": "Nome",
			  "atype": "String",
			  "sources": [
				"User"
			  ],
			  "keephistory": 0,
			  "validityinterval": 0,
			  "deadline": 0
			}]
    }`
	c := kbase.NewClass(json)
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
			"keephistory": 0,
			"validityinterval": 0,
			"deadline": 5
		  },
		  {
			"name": "Potência",
			"atype": "Number",
			"sources": [
			  "User",
			  "IOT",
			  "Simulation"
			],
			"keephistory": 0,
			"validityinterval": 0,
			"deadline": 0
		  },
		  {
			"name": "CurrenPower",
			"atype": "Number",
			"sources": [
			  "IOT"
			],
			"keephistory": 0,
			"validityinterval": 0,
			"deadline": 0
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
			"validityinterval": 0,
			"deadline": 0
		  }
		]
	  }`
	c = kbase.NewClass(json)
	fmt.Println(c)

}

func Test5(kbase *kb.KnowledgeBased) {
	o := kbase.NewObject("MotorElétrico", "M01")
	fmt.Println(o)
}

func Test6(kbase *kb.KnowledgeBased) {
	kbase.NewRule("for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230", 100, 0)
}

func Test7(kbase *kb.KnowledgeBased) {
	o := kbase.FindObjectByName("M01")
	fmt.Println(o)
}

func Test8(kbase *kb.KnowledgeBased) {
	a := kbase.FindAttributeObjectByName("M01.Status")
	fmt.Println(a)
}
