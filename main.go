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

	kb := classes.KB{}

	kb.ReadBK("mongodb://localhost:27017", "K2")

}
