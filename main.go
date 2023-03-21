package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

type NullWriter int

const (
	AskColor        = "\033[1;34m%s: \033[0m"
	InfoColor       = "\033[1;34m%s\033[0m\n"
	NoticeColor     = "\033[1;36m%s\033[0m\n"
	NoticeHelpColor = "\033[1;31m%s:\t\033[0;36m%s\033[0m\n"
	WarningColor    = "\033[1;33m%s\033[0m\n"
	ErrorColor      = "\033[1;31m%s\033[0m\n"
	DebugColor      = "\033[0;36m%s\033[0m\n"
	NoneColor       = "\033[0m"
	AssistenColor   = "\n\033[1;35m%s:\n\033[0m%s\n"
	TitleTextToken  = "\033[1;34m%s: \033[1;33m%s\n\033[0m\n"
	TitleIntToken   = "\033[1;34m%s: \033[1;33m%d\n\033[0m\n"
)

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func LogError(message string, err error) {
	if err != nil {
		fmt.Printf(TitleTextToken, message, err.Error())
		return
	}
	fmt.Printf(TitleTextToken, message, "")
}

var gS *GService
var contentGrapper *ContentGrapper
var twitterBot *TwitterBot

func main() {
	log.SetOutput(new(NullWriter))

	// Init Chatbot
	chatBot := NewChatBot(context.Background(), "config.yml")
	// Init Trello Connect
	trelloConnect := NewTrelloConnect(chatBot.Config)
	twitterBot := NewTwitterBot(chatBot.Config)
	defaultEngine := "text-davinci-003"

	if chatBot.Config.UseSound {
		// Init MP3 Player
		p := NewMp3Player(chatBot.Config.PollySampleFile, chatBot.Config)
		chatBot.Mp3Player = p
	}

	go NewAuthConfigWordpress(chatBot.Config)
	go NewAuthConfigGMail(chatBot.Config)

	contentGrapper = NewContentGrapper("", chatBot)

	for !chatBot.Quit {
		chatBot.Config.MaxTokens = 2048
		fmt.Printf(AskColor, Promt)
		if !scanner.Scan() {
			break
		}
		cCli := NewCommandCli(scanner.Text())
		if len(cCli.UserInput) < 1 {
			continue
		}
		chatBot.Command = &cCli

		switch cCli.Operation {
		case Quit:
			chatBot.Quit = true
			chatBot.Die()
			continue
		case Help:
			fmt.Printf(NoticeHelpColor, Quit.Name, Quit.Desc)
			fmt.Printf(NoticeHelpColor, Help.Name, Help.Desc)
			fmt.Printf(NoticeHelpColor, Clear.Name, Clear.Desc)
			fmt.Printf(NoticeHelpColor, Save.Name, Save.Desc)
			fmt.Printf(NoticeHelpColor, Load.Name, Load.Desc)
			fmt.Printf(NoticeHelpColor, Mail.Name, Mail.Desc)
			fmt.Printf(NoticeHelpColor, History.Name, History.Desc)
			fmt.Printf(NoticeHelpColor, Token.Name, Token.Desc)
			fmt.Printf(NoticeHelpColor, Assisten.Name, Assisten.Desc)
			fmt.Printf(NoticeHelpColor, Trello.Name, Trello.Desc)
			fmt.Printf(NoticeHelpColor, Config.Name, Config.Desc)
			continue
		case Clear:
			chatBot.Clear()
			fmt.Printf(AssistenColor, chatBot.ActiveAsi, "Historie gelöscht.")
			continue
		case Save:
			chatBot.SaveToFile(chatBot.LastResponse, chatBot.Command.Params)
			continue
		case Load:
			chatBot.Config.MaxTokens = 3*1024
			fmt.Printf(AssistenColor, chatBot.ActiveAsi, chatBot.LoadFromFile(defaultEngine))
			continue
		case Mail:
			gS = NewGService(chatBot.Config)
			chatBot.ActiveAsi = "mail"
			split := strings.Split(chatBot.Command.UserInput, " ")
			var count int64
			label := "INBOX"

			if len(split) >= 2 {
				count, err = strconv.ParseInt(split[1], 10, 64)
			}
			if err != nil {
				count = 10
			}
			if len(split) < 2 {
				count = 10
			}
			if len(split) > 2 {
				label = split[2]
			}

			mails := gS.GetMails(count, label)
			chatBot.ContentHistory = GetJSONString(mails)
			chatBot.Command.UserInput = "Erstelle eine Zusammenfassung der Emails und zeig die wichtigsten an:"
			fmt.Printf(AssistenColor, chatBot.ActiveAsi, chatBot.GetResponse(chatBot.ActiveAsi, defaultEngine))
			continue
		case History:
			fmt.Printf(TitleTextToken, "LastResponse", chatBot.Response)
			fmt.Printf(TitleTextToken, "History", chatBot.ContentHistory)
			continue
		case Token:
			fmt.Printf(TitleIntToken, "Used Token", len(chatBot.CallPromt))
			continue
		case Assisten:
			chatBot.SetAsi()
			continue
		case Trello:
			if chatBot.Config.UseTrello {
				trelloConnect.AddTrelloCard(chatBot.LastQuestion, chatBot.LastResponse)
			}
			continue
		case Twitter:
			if chatBot.Config.UseTwitter {
				fmt.Printf(NoticeColor, "Twitter Active")
				twitterBot.Test()
			}
			continue
		case Config:
			fmt.Println(chatBot.Config)
			continue
		default:
			switch chatBot.ActiveAsi {
			case "dev":
				chatBot.Config.FrequencyPenalty = 0
				chatBot.Config.PresencePenalty = 0
				chatBot.Config.MaxTokens = 3800
				chatBot.Config.Temperature = 0.7
				fmt.Printf(AssistenColor, chatBot.ActiveAsi, chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine))
			case "mail":
				if gS == nil {
					gS = NewGService(chatBot.Config)
				}
				if gS.Mails != nil && len(gS.Mails) == 0 {
					gS.LoadJSON()
				}
				chatBot.ContentHistory = "E-Mails als JSON Format:" + GetJSONString(gS.Mails)
				fmt.Printf(AssistenColor, chatBot.ActiveAsi, chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine))
				continue
			case "voice":
				voicefunc := func() {
					chatBot.Config.FrequencyPenalty = 1
					chatBot.Config.PresencePenalty = 1
					chatBot.Config.MaxTokens = 512
					chatBot.Config.Temperature = 1
					chatBot.Mp3Player.TextToSpeech(chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine))
				}
				if chatBot.Config.UseAWS {
					go voicefunc()
				}
			case "trello-ausbilder":
				if chatBot.Config.UseTrello {
					trelloConnect.TrelloBoardId = trelloConnect.Config.TrelloAzubiBoard
					chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine)
					trelloConnect.AddTrelloCard(chatBot.LastQuestion, chatBot.LastResponse)
				}
			case "trello-doku":
				if chatBot.Config.UseTrello {
					chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine)
					trelloConnect.TrelloBoardId = trelloConnect.Config.TrelloDokuBoard
					trelloConnect.AddTrelloCard(chatBot.LastQuestion, chatBot.LastResponse)
				}
			case "tw":
				if chatBot.Config.UseTrello {
					trelloConnect.TrelloBoardId = trelloConnect.Config.TrelloRedakBoard
					todo := trelloConnect.GetTodoFromList()
					if len(todo) < 1 {
						fmt.Printf(ErrorColor, "Keine Karte gefunden")
						continue
					}
					chatBot.Command.UserInput = contentGrapper.CleantText("Erstelle, ohne Kommentare und ohne Erklärungen, ein Blogpost zur Zielgruppe Unternehmer mit 800 Wörter für das folgende Thema: " + todo)
					chatBot.Config.MaxTokens = 3 * 1024
					response := chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine)
					if len(response) > 1 {
						chatBot.ContentHistory = ""
						wordpress := NewWordpress(chatBot.Config)
						wordPressData := NewWordpressData(todo, response)
						chatBot.Command.UserInput = contentGrapper.CleantText("Erstelle ohne Kommentare oder Erklärungen 10 Keywörter als Liste komma separiert " + todo)
						responsTag := chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine)
						wordPressData.Tags = strings.Split(responsTag, ",")
						chatBot.Command.UserInput = contentGrapper.CleantText("Erstelle eine Twitter Nachricht zum Beitrag: " + todo)
						responsPreview := chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine)
						wordPressData.PublicizeMessage = responsPreview
						wordPressData.Slug = strings.Join(strings.Split(todo, ","), "-")
						resp := wordpress.CreateBlog(wordPressData)
						fmt.Printf(TitleTextToken, "Blog Post Response: ", resp)
						trelloConnect.MoveCardToNextList(response, resp)
						chatBot.Config.MaxTokens = 2048
					}
				}
			default:

				twitterBot.GetTweets()
				url := contentGrapper.GetUrlFromText(chatBot.Command.UserInput)
				if url != "" {
					content := contentGrapper.GetContent(url)
					chatBot.Config.MaxTokens = 2056
					chatBot.ContentHistory = ""
					chatBot.Command.UserInput = " Webseite: " + content + " : Ohne Kommentare oder Erklärungen: " + chatBot.Command.UserInput
				}
				fmt.Printf(AssistenColor, chatBot.ActiveAsi, chatBot.GetResponse(chatBot.ActiveAsi,defaultEngine))
			}
		}
	}
}

func GetJSONString(data interface{}) string {
	jsonString, _ := json.Marshal(data)
	return string(jsonString)
}
