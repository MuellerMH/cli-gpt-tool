package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GService struct {
	Config    *ConfigBot
	Token     *oauth2.Token
	Ctx       context.Context
	Service   *gmail.Service
	User      string
	Mails     []MailContent
	MailLimit int
}
type TokenSource struct {
	ConfigBot *ConfigBot
}

func NewGService(configBot *ConfigBot) *GService {
	// Create a Gmail service using the token
	if !configBot.ActiveGoogle {
		return nil
	}
	tokenSource := oauth2.StaticTokenSource(configBot.OauthGoogleToken)
	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		LogError("no gmail Client", err)
	}
	gS := GService{
		Token:     configBot.OauthGoogleToken,
		Ctx:       context.Background(),
		Service:   gmailService,
		User:      "me",
		MailLimit: 10,
		Config:    configBot,
		Mails:     []MailContent{},
	}
	return &gS
}

type MailHeader struct {
	To      string
	From    string
	Date    string
	Subject string
}

func (m *MailHeader) SetValue(name string, value string) {
	switch name {
	case "To":
		m.To = value
	case "Subject":
		m.Subject = value
	case "From":
		m.From = value
	case "Date":
		m.Date = value
	}
}

type MailContent struct {
	Header MailHeader
	Body   string
}

func (gs *GService) GetMails(count int64, labels string) []MailContent {
	if gs.Config.OauthGoogleToken == nil {
		go NewAuthConfigGMail(gs.Config)
		LogError("Keine aktiver Token", nil)
		return nil
	}
	if !gs.Config.ActiveGoogle {
		return nil
	}
	if gs.Service == nil {
		LogError("No GService avaible: %v", err)
		return nil
	}
	mails, err := gs.Service.Users.Messages.List(gs.User).Q("label:" + labels).MaxResults(count).Do()
	if err != nil {
		LogError("no gmail Client", err)
		return nil
	}
	for _, mailMessage := range mails.Messages {
		mail, err := gs.Service.Users.Messages.Get(gs.User, mailMessage.Id).Do()
		if err != nil {
			LogError("no gmail Client", err)
			return nil
		}
		body, _ := base64.StdEncoding.DecodeString(mail.Raw)
		mailHeader := MailHeader{}
		for _, head := range mail.Payload.Headers {
			mailHeader.SetValue(head.Name, head.Value)
		}
		gs.Mails = append(gs.Mails, MailContent{Header: mailHeader, Body: string(body)})
	}
	gs.SaveJSON(gs.Mails)
	return gs.Mails
}

// SaveJSON saves the token information in a JSON file.
func (gs *GService) SaveJSON(mails []MailContent) error {
	data, err := json.MarshalIndent(mails, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(time.Now().Format("2006-01-02T15-04-05")+"-mails.json", data, 0644); err != nil {
		return err
	}

	return nil
}

// LoadJSON loads the token information from a JSON file.
func (gs *GService) LoadJSON() ([]MailContent, error) {
	data, err := ioutil.ReadFile("mails.json")
	mails := []MailContent{}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &mails); err != nil {
		return nil, err
	}
	gs.Mails = mails
	return mails, nil
}
func (gs *GService) GetLabels() {
	if !gs.Config.ActiveGoogle {
		return
	}
	if gs.Service == nil {
		LogError("No GService avaible: %v", err)
		return
	}
	if gs.Service.Users == nil {
		fmt.Println("No *gmail.UsersService found.")
		return
	}
	if gs.Service.Users == nil {
		fmt.Println("No *gmail.UsersLabelsService found.")
		return
	}
	if gs.User != "me" {
		fmt.Println("No *gmail.UsersLabelsService found.")
		return
	}
	r, err := gs.Service.Users.Labels.List(gs.User).Do()
	if err != nil {
		LogError("Unable to retrieve labels: %v", err)
		return
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}
