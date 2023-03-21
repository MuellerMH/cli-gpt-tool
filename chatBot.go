package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
)

type ChatBot struct {
	Response       string
	LastResponse   string
	ActiveAsi      string
	CallPromt      string
	LastQuestion   string
	ContentHistory string
	Scanner        *bufio.Scanner
	Quit           bool
	Context        context.Context
	Config         *ConfigBot
	Client         gpt3.Client
	Mp3Player      *Mp3Player
	Command        *CommandCLI
	ConfigFile     string
}

func NewChatBot(ctx context.Context, config_file string) *ChatBot {
	cb := ChatBot{Context: ctx}
	cb.LoadConfig(config_file)
	cb.ActiveAsi = "default"
	return &cb
}

func (cb *ChatBot) Clear() {
	cb.ContentHistory = ""
	cb.Response = ""
}

func (cb *ChatBot) SetAsi() {
	if len(cb.Command.Params) >= 1 {
		cb.SaveBrain()
		cb.ActiveAsi = cb.Command.Params
		cb.ContentHistory = cb.LoadBrain()
		fmt.Println("History",cb.ActiveAsi,"geladen:",cb.ContentHistory)
	}
	fmt.Printf(NoticeColor, cb.ActiveAsi)
}

func (cb *ChatBot) LoadConfig(config_file string) bool {
	cb.ConfigFile = config_file
	config, err := GetConfigFromYAMLFile(config_file)
	if err != nil {
		return false
	}
	config.CheckApiKey()
	cb.Config = config
	cb.Client = gpt3.NewClient(cb.Config.OpenAiApiKey)
	return true
}

func (cb *ChatBot) Die() {
	if len(cb.Command.Words) > 256 {
		cb.ContentHistory = cb.GetResponse("clean","ada")
		cb.SaveBrain()
	} else {
		cb.SaveBrain()
	}
}

func (cb *ChatBot) GetResponse(act string, engineTypes string) string {
	if cb.Client == nil {
		cb.LoadConfig(cb.ConfigFile)
		cb.Config.CheckApiKey()
		cb.Client = gpt3.NewClient(cb.Config.OpenAiApiKey)
	}
	for i := 0; i < 3; i++ {
		cb.Response = ""
		actString := cb.LoadFile("actor/" + cb.ActiveAsi)
		switch act {
		case "clean":
			cb.CallPromt = actString + "\nManuel: " + cb.ContentHistory + "\nAssistent:"
		default:
			words := strings.Fields(cb.ContentHistory)
			if len(words) > 256 {
				cb.ContentHistory = words[256]
			}
			cb.CallPromt = actString + "\nAssistent: " + cb.ContentHistory + "\nManuel: " + cb.Command.UserInput + "\nAssistent:"
		}
		cb.LastQuestion = cb.Command.UserInput
		err := cb.Client.CompletionStreamWithEngine(cb.Context, engineTypes, gpt3.CompletionRequest{
			Prompt: []string{
				cb.CallPromt,
			},
			MaxTokens:        gpt3.IntPtr(cb.Config.MaxTokens),
			Temperature:      gpt3.Float32Ptr(cb.Config.Temperature),
			PresencePenalty:  float32(cb.Config.PresencePenalty),
			FrequencyPenalty: float32(cb.Config.FrequencyPenalty),
			N:                gpt3.IntPtr(cb.Config.N),
			Stop:             cb.Config.Stop,
		}, func(resp *gpt3.CompletionResponse) {
			cb.Response += resp.Choices[0].Text
		})
		if err == nil {
			cb.LastResponse = cb.Response
			cb.ContentHistory += cb.Response
			return cb.Response
		} else {
			fmt.Println("System überlastet, versuch es erneut.")
			fmt.Printf("%s\n", err)
		}
	}
	fmt.Println("Error occured, Please try again Late")
	return ""
}

func (cb *ChatBot) TryOpen(filename string) {
	fmt.Println("unimplemented")
}

func (cb *ChatBot) SaveToFile(text string, file_name string) {
	if len(file_name) < 1 {
		return
	}
	if len(text) < 1 {
		return
	}
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Response saved as: ", cb.Command.Params)
}
func (cb *ChatBot) LoadBrain() string {
	home := os.Getenv("HOME")
	file_name := home + "/.chat-brain." + cb.ActiveAsi
	fmt.Println("Lade History von", file_name, home, cb.ActiveAsi)
	return cb.LoadFile(file_name)

}
func (cb *ChatBot) LoadFile(file_name string) string {
	dat, err := ioutil.ReadFile(file_name)
	if err != nil {
		fmt.Println("Datei konnte nicht geladen werden", file_name)
		return ""
	}
	return string(dat)
}

func (cb *ChatBot) LoadFromFile(engineModel string) string {
	cli := NewCommandCli(cb.LoadFile(cb.Command.Params))
	cb.Command = &cli
	return cb.GetResponse(cb.ActiveAsi,engineModel)
}

func (cb *ChatBot) SaveBrain() {
	home := os.Getenv("HOME")
	file_name := home + "/.chat-brain." + cb.ActiveAsi
	fmt.Println("Speichere History von", file_name, home, cb.ActiveAsi)
	cb.SaveToFile(cb.ContentHistory, file_name)
}
