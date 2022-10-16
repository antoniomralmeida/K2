package main

import (
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
	/*
			a1 := classes.KBAttribute{Name: "Nome", AType: classes.KBString, Sources: []classes.KBSource{classes.User}}
			a2 := classes.KBAttribute{Name: "Data", AType: classes.KBDate, Sources: []classes.KBSource{classes.User}}
			a3 := classes.KBAttribute{Name: "Potência", AType: classes.KBNumber, Sources: []classes.KBSource{classes.User, classes.PLC}}
			c1 := classes.KBClass{Name: "MotorElétrico", Attributes: []classes.KBAttribute{a1, a2, a3}}
			kb.NewClass(&c1)

		c1 := kb.GetClass("MotorElétrico")
		attrs := kb.FindAttributes(c1)
		a1 := classes.KBAttributeObject{Attribute: attrs[0].Id}
		a2 := classes.KBAttributeObject{Attribute: attrs[1].Id}

		o1 := classes.KBObject{Name: "M01", Class: c1.Id, Attributes: []classes.KBAttributeObject{a1, a2}}
		w1 := classes.KBWorkspace{Workspace: "Chão de fábrica", Objects: []classes.KBObject{o1}}
		kb.NewWorkspace(&w1)
	*/
	r1 := classes.KBRule{Rule: "for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230"}
	kb.NewRule(&r1)

}
