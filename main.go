package main

import (
	"fmt"
	"main/classes"
)

func main() {
	//ebnf := classes.EBNF{}
	//ebnf.ReadToken("k2.ebnf")
	//ebnf.Parsing("for any MotorElétrico M if the Status of M is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230")
	//ebnf.PrintEBNF()

	//tts := classes.TTS{}
	//tts.SetText("for any MotorElétrico M if the Status of M is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230", voices.English)
	//tts.SetText("Olá sou Manoel", voices.Portuguese)
	//tts.Speech()

	kb := classes.KnowledgeBase{}

	kb.ConnectDB("mongodb://localhost:27017", "K2")
	kb.ReadBK()

	w1 := kb.FindWorkspaceByName("Chão de fábrica")

	o1 := kb.FindObjectByName(w1, "M01")
	ao1 := kb.FindAttributeObject(o1, "Potência")

	fmt.Println(ao1.KbAttribute.Name, ao1.KbHistory.Value)

	//	kb.NewRule("for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230")

	/*
					a1 := classes.KBAttribute{Name: "Nome", AType: classes.KBString, Sources: []classes.KBSource{classes.User}}
					a2 := classes.KBAttribute{Name: "Data", AType: classes.KBDate, Sources: []classes.KBSource{classes.User}}
					a3 := classes.KBAttribute{Name: "Potência", AType: classes.KBNumber, Sources: []classes.KBSource{classes.User, classes.PLC}}
					c1 := classes.KBClass{Name: "MotorElétrico", Attributes: []classes.KBAttribute{a1, a2, a3}}
					kb.NewClass(&c1)
				kb.NewRule("for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230")


			w1 := kb.NewWorkspace("Chão de fábrica", "chao.jpg")

		w1 := kb.FindWorkspaceByName("Chão de fábrica")

		o1 := kb.FindObjectByName(w1, "M01")
		ao1 := kb.FindAttributeObject(o1, "Potência")
		kb.SaveValue(ao1, "120", classes.User)
	*/
}
