package main

import (
	"github.com/antoniomralmeida/k2/db"
	"github.com/antoniomralmeida/k2/kb"
)

func main() {
	db.ConnectDB("mongodb://localhost:27017", "K2")
	//ebnf := kb.EBNF{}
	//ebnf.ReadToken("k2.ebnf")
	//ebnf.Parsing("for any MotorElétrico M if the Status of M is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230")
	//ebnf.PrintEBNF()

	//tts := kb.TTS{}
	//tts.SetText("for any MotorElétrico M if the Status of M is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230", voices.English)
	//tts.SetText("Olá sou Manoel", voices.Portuguese)
	//tts.Speech()

	kb1 := kb.KnowledgeBase{}
	kb1.Init()

	kb1.ReadEBNF("k2.ebnf")
	kb1.ReadBK()
	//kb.PrintEBNF()
	//kb.Run()
	//kb.NewRule(90, "for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230")
	//kb.NewRule(100, "initially unconditionally then set the Status of the M01 to PowerOn")

	//web.Run()

	/*
	   c1 := kb.FindClassByName("MotorElétrico")

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

	   		//kb.SaveValue(ao1, "120", kb.User)

	   		fmt.Println(o1.Name, ao1.KbAttribute.Name, ao1.Value())

	   		//	kb.NewRule("for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230")
	*/
	/*
			c1 := kb.KBClass{Name: "MotorElétrico",
				Attributes: []kb.KBAttribute{
					kb.KBAttribute{Name: "Data", AType: kb.KBDate, Sources: []kb.KBSource{kb.User}},
					kb.KBAttribute{Name: "Potência", AType: kb.KBNumber, Sources: []kb.KBSource{kb.User, kb.PLC}},
					kb.KBAttribute{Name: "Name", AType: kb.KBString, Sources: []kb.KBSource{kb.User}},
					kb.KBAttribute{Name: "CurrenPower", AType: kb.KBNumber, Sources: []kb.KBSource{kb.PLC}},
					kb.KBAttribute{Name: "Status", AType: kb.KBList, Options: []string{"PowerOff", "PowerOn"}, Sources: []kb.KBSource{kb.PLC}}}}
			kb1.NewClass(&c1)

		w1 := kb1.NewWorkspace("Chão de fábrica", "chao.jpg")
	*/

	//c1 := kb1.FindClassByName("MotorElétrico", true)

	//o1 := kb1.NewObject(c1, "M01")
	//kb1.NewObject(c1, "M02")

	o1 := kb1.FindObjectByName("M01")

	ao1 := kb1.FindAttributeObject(o1, "Potência")
	kb1.SaveValue(ao1, "121", kb.User)

	kb1.NewRule("for any MotorElétrico M if the Status is PowerOff then inform to the operator that 'O Motor' the Name of M 'parou!' and set the CurrentPower of M = 0.3230", 100)

}
