package main

import (
	"io/ioutil"

	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type ConfigBot struct {
	UseAWS              bool   `yaml:"use.aws"`
	UseSound            bool   `yaml:"polly.active"`
	UseTrello           bool   `yaml:"trello.active"`
	UseTwitter          bool   `yaml:"twitter.active"`
	OpenAiApiKey        string `yaml:"openai.key"`
	TrelloKey           string `yaml:"trello.key"`
	TrelloSecret        string `yaml:"trello.secret"`
	TrelloAzubiBoard    string `yaml:"trello.azubi-board"`
	TrelloDokuBoard     string `yaml:"trello.doku-board"`
	TrelloRedakBoard    string `yaml:"trello.redaktions-board"`
	AWSRegion           string `yaml:"aws.region"`
	AWSKey              string `yaml:"aws.key"`
	AWSSecret           string `yaml:"aws.secret"`
	PollySampleFile     string `yaml:"polly.sample"`
	SaveAudioFiles      bool   `yaml:"polly.safeAudio"`
	OauthWordpressToken *oauth2.Token
	OauthGoogleToken    *oauth2.Token
	OauthTwitterToken   *oauth2.Token
	MaxTokens           int      `yaml:"bot.MaxTokens"`
	Temperature         float32  `yaml:"bot.Temperature"`
	PresencePenalty     float32  `yaml:"bot.PresencePenalty"`
	FrequencyPenalty    float32  `yaml:"bot.FrequencyPenalty"`
	N                   int      `yaml:"bot.N"`
	Stop                []string `yaml:"bot.Stop"`
	ActiveWordpress     bool     `yaml:"use.wordpress"`
	ActiveGoogle        bool     `yaml:"use.google"`
	ActiveTwitter       bool     `yaml:"use.twitter"`
	TwitterToken        string   `yaml:"twitter.bearer"`
	TwitterKey          string   `yaml:"twitter.key"`
	TwitterSecret       string   `yaml:"twitter.secret"`
}

func GetConfigFromYAMLFile(file string) (*ConfigBot, error) {
	var c *ConfigBot

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ConfigBot) SetToken(authToken *oauth2.Token, loginUrl string) {
	switch loginUrl {
	case "/google":
		c.OauthGoogleToken = authToken
	case "/wordpress":
		c.OauthWordpressToken = authToken
	default:
		LogError("no used OauthToken", nil)
	}

}

func (c *ConfigBot) CheckApiKey() {
	if c.OpenAiApiKey == "" {
		panic("Missing API KEY")
	}
}

func (c *ConfigBot) CheckAWSCredentials() {
	if c.AWSKey == "" && c.AWSSecret == "" {
		panic("Missing AWS Credentials")
	}
}
