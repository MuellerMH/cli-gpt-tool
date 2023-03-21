package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
)

type OauthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

type AuthConfig struct {
	Config        *ConfigBot
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	Endpoint      oauth2.Endpoint
	LocalhostPort string
	LocalURl      string
	LoginURL      string
	Scopes        []string
}

func NewAuthConfigWordpress(configBot *ConfigBot) {
	if !configBot.ActiveWordpress {
		return
	}
	conf := AuthConfig{
		ClientID:     "83777",
		ClientSecret: "aDOP7mGMn8QwuI3u66VQcGpwvT9kj6adXwwJmFuFzdgU2d5kLlSieX7KjxYMk9Zr",
		RedirectURL:  "/wordpress/callback",
		LoginURL:     "/wordpress",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://public-api.wordpress.com/oauth2/authorize",
			TokenURL: "https://public-api.wordpress.com/oauth2/token",
		},
		LocalURl:      "http://localhost",
		LocalhostPort: "30000",
		Config:        configBot,
	}
	go Auth(&conf)
}

func NewAuthConfigGMail(configBot *ConfigBot) {
	if !configBot.ActiveGoogle {
		return
	}
	conf := AuthConfig{
		ClientID:     "721337092891-68fcjaaatr44u0o4nbofd6iu27h3cfbd.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "/google/callback",
		LoginURL:     "/google",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		LocalURl:      "http://localhost",
		LocalhostPort: "30001",
		Config:        configBot,
		Scopes: []string{
			"https://www.googleapis.com/auth/plus.login",
			"https://www.googleapis.com/auth/calendar",
			"https://www.googleapis.com/auth/calendar.readonly",
			"https://www.googleapis.com/auth/contacts",
			"https://www.googleapis.com/auth/contacts.readonly",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/user.addresses.read",
			"https://www.googleapis.com/auth/user.birthday.read",
			"https://www.googleapis.com/auth/user.emails.read",
			"https://www.googleapis.com/auth/user.phonenumbers.read",
			"https://www.googleapis.com/auth/gmail.labels",
			"https://www.googleapis.com/auth/gmail.send",
			"https://www.googleapis.com/auth/gmail.readonly",
		},
	}
	go Auth(&conf)
}

func NewAuthConfigTwitter(configBot *ConfigBot) {
	if !configBot.ActiveTwitter {
		return
	}
	conf := AuthConfig{
		ClientID:     "jpHVL1Z0fnHFMTLeEqdrsK6Nl",
		ClientSecret: "Ar9JPnA934HTmtA5NU4AnCYOAgSBbDShVtdgk8SkTH8U7UeFcf",
		RedirectURL:  "/twitter/callback",
		LoginURL:     "/twitter",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.twitter.com/oauth2/token",
			TokenURL: "https://api.twitter.com/oauth2/token",
		},
		LocalURl:      "http://localhost",
		LocalhostPort: "30002",
		Config:        configBot,
	}
	go Auth(&conf)
}

func Auth(authConfig *AuthConfig) {

	token, err := authConfig.LoadJSON("." + authConfig.LoginURL + ".oauth.json")
	if err != nil {
		LogError("token load error", err)
	}
	if token != nil && token.Expiry.After(time.Now()) {
		authConfig.Config.SetToken(token, authConfig.LoginURL)
		return
	}

	context := context.Background()
	conf := &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  authConfig.LocalURl + ":" + authConfig.LocalhostPort + authConfig.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authConfig.Endpoint.AuthURL,
			TokenURL: authConfig.Endpoint.TokenURL,
		},
	}
	if authConfig.LoginURL == "/google" {
		conf.Scopes = authConfig.Scopes
	}

	http.HandleFunc(authConfig.LoginURL, func(w http.ResponseWriter, r *http.Request) {
		url := conf.AuthCodeURL("state")
		http.Redirect(w, r, url, http.StatusFound)
	})

	http.HandleFunc(authConfig.RedirectURL, func(w http.ResponseWriter, r *http.Request) {
		// Handle the exchange code to initiate a transport.
		token, oauthErr := conf.Exchange(context, r.FormValue("code"))
		if oauthErr != nil {
			LogError("Callback", oauthErr)
		}
		authConfig.Config.SetToken(token, authConfig.LoginURL)
		authConfig.SaveJSON("."+authConfig.LoginURL+".oauth.json", token)
		w.Write([]byte("<html><body><h1>Token erfolgreich erstellt</h1><script type=\"javascript\">(function() {window.close();})();</script></body></html>"))
	})
	go Openbrowser(authConfig.LocalURl + ":" + authConfig.LocalhostPort + authConfig.LoginURL)
	http.ListenAndServe(":"+authConfig.LocalhostPort, nil)
}

// SaveJSON saves the token information in a JSON file.
func (a *AuthConfig) SaveJSON(filename string, token *oauth2.Token) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}

// LoadJSON loads the token information from a JSON file.
func (a *AuthConfig) LoadJSON(filename string) (*oauth2.Token, error) {
	data, err := ioutil.ReadFile(filename)
	token := oauth2.Token{}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}
