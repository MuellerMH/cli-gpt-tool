package main

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2/clientcredentials"
)

type TwitterBot struct {
	Token  string
	Config *ConfigBot
	Client *twitter.Client
	Tweets []twitter.Tweet
}

func NewTwitterBot(configBot *ConfigBot) *TwitterBot {
	tb := &TwitterBot{}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     configBot.TwitterKey,
		ClientSecret: configBot.TwitterSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(context.Background())

	// Twitter client
	tb.Client = twitter.NewClient(httpClient)

	tb.Tweets = []twitter.Tweet{}
	tb.Config = configBot
	return tb
}

func (tb *TwitterBot) Test() {
	// user timeline
	userTimelineParams := &twitter.HomeTimelineParams{}
	tweets, resp, err := tb.Client.Timelines.HomeTimeline(userTimelineParams)
	tb.Tweets = tweets
	if err != nil {
		LogError("Twitter Timeline", err)
	}
	if resp != nil {
		fmt.Printf("Twitter Response %s\n", GetJSONString(resp))
	}
	fmt.Printf("USER TIMELINE:\n%+v\n", tb.Tweets)
	for _, t := range tweets {
		// status show
		statusShowParams := &twitter.StatusShowParams{}
		tweet, _, _ := tb.Client.Statuses.Show(t.ID, statusShowParams)
		fmt.Printf("STATUSES SHOW:\n%+v\n", tweet)

		// statuses lookup
		statusLookupParams := &twitter.StatusLookupParams{ID: []int64{20}, TweetMode: "extended"}
		tweets, _, _ := tb.Client.Statuses.Lookup([]int64{t.ID}, statusLookupParams)
		fmt.Printf("STATUSES LOOKUP:\n%+v\n", tweets)

		// oEmbed status
		statusOembedParams := &twitter.StatusOEmbedParams{ID: t.ID, MaxWidth: 500}
		oembed, _, _ := tb.Client.Statuses.OEmbed(statusOembedParams)
		fmt.Printf("OEMBED TWEET:\n%+v\n", oembed)

	}
}

func (tb *TwitterBot) GetTweets() {
	tb.Tweets, _, _ = tb.Client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{})
	for _, tweet := range tb.Tweets {
		fmt.Println(tweet.Text)
	}
}
