package main

import (
	"github.com/antoniomralmeida/k2/knowledgebase"
	"github.com/antoniomralmeida/k2/web"
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

	kb := knowledgebase.KnowledgeBase{}

	kb.ConnectDB("mongodb://localhost:27017", "K2")
	kb.ReadEBNF("k2.ebnf")
	kb.ReadBK()
	//kb.PrintEBNF()
	//kb.Run()
	//kb.NewRule(90, "for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230")
	//kb.NewRule(100, "initially unconditionally then set the Status of the M01 to PowerOn")

	web.Run()

	/*
	   c1 := kb.FindClassByName("MotorElétrico")

	   a1 := classes.KBAttribute{Name: "Name", AType: classes.KBString, Sources: []classes.KBSource{classes.User}}
	   a2 := classes.KBAttribute{Name: "CurrenPower", AType: classes.KBNumber, Sources: []classes.KBSource{classes.PLC}}

	   kb.AddAttribute(c1, &a1, &a2)
	   /*w1 := kb.NewWorkspace("Chão de fábrica", "chao.jpg")

	   	c1 := kb.FindClassByName("MotorElétrico")

	   	o1 := kb.NewObject(c1, "M01")
	   	o2 := kb.NewObject(c1, "M02")

	   	kb.LinkObjects(w1, o1, o2)

	   		w1 := kb.FindWorkspaceByName("Chão de fábrica")
	   		//w1 := kb.NewWorkspace("Chão de fábrica", "chao.jpg")

	   		o1 := kb.FindObjectByName(w1, "M02")
	   		ao1 := kb.FindAttributeObject(o1, "Potência")

	   		//c1 := kb.FindClassByName("MotorElétrico")

	   		//o1 := kb.NewObject(w1, c1, "M01")
	   		//kb.NewObject(w1, c1, "M02")
	   		//ao1 := kb.FindAttributeObject(o1, "Potência")

	   		//kb.SaveValue(ao1, "120", classes.User)

	   		fmt.Println(o1.Name, ao1.KbAttribute.Name, ao1.Value())

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
