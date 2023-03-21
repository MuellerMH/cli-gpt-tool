package main

import "strings"

type CommandCLI struct {
	Words     []string
	Command   string
	Params    string
	Operation ChatBotOperationen
	UserInput string
}

func NewCommandCli(scannerText string) CommandCLI {
	cC := CommandCLI{UserInput: scannerText}
	cC.Words = strings.Fields(scannerText)
	if len(cC.Words) > 0 {
		cC.Command = strings.ToLower(cC.Words[0])
		cC.Operation = OperationFromString(strings.ToLower(cC.Words[0]))
		cC.Params = strings.Join(cC.Words[len(cC.Words)-1:], " ")
	}
	return cC
}
