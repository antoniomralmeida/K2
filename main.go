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
	*/
	w1 := classes.KBWorkspace{Workspace: "Chão de fábrica"}
	kb.NewWorkspace(&w1)

	i := kb.FindWorkspaceByName("Chão de fábrica")
	c2 := kb.GetClass("MotorElétrico")
	o1 := classes.KBObject{Name: "M01", Class: c2.Id}

	kb.NewObject(i, &o1, c2)

	r1 := classes.KBRule{Rule: "for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230"}
	kb.NewRule(&r1)

}
