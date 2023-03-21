package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WordpressData struct {
	Title            string   `json:"title"`   //"My New Blog",                                       // Blog Title
	Content          string   `json:"content"` //"My first blog created with the Wordpress REST API", // Blog Description
	Slug             string   `json:"slug"`
	Author           string   `json:"author"`
	PublicizeMessage string   `json:"publicize_message"`
	Status           string   `json:"status"`
	Categories       []string `json:"categories"`
	Tags             []string `json:"tags"`
}

func NewWordpressData(title string, content string) *WordpressData {
	wD := WordpressData{
		Title:            title,
		Content:          content,
		Author:           "muellermh",
		PublicizeMessage: "test",
		Status:           "draft",
		Slug:             "",
		Categories:       []string{"Allgemein"},
		Tags:             []string{},
	}
	return &wD
}

type Wordpress struct {
	Endpoint string `json:"endpoint"`
	Site     string
	Token    string
}

func NewWordpress(config *ConfigBot) *Wordpress {
	if !config.ActiveWordpress {
		return nil
	}
	wp := Wordpress{Endpoint: "https://public-api.wordpress.com/rest/v1.1", Site: "digital-business.blog", Token: config.OauthWordpressToken.AccessToken}
	return &wp
}

func (wp *Wordpress) CreateBlog(data *WordpressData) string {
	oauth := AuthConfig{}
	token, _ := oauth.LoadJSON("wordpress.oauth.json")
	if token != nil && token.Expiry.Before(time.Now()) {
		go Openbrowser("http://localhost:30000/wordpress")
	}
	client := &http.Client{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	}
	req, err := http.NewRequest("POST", wp.Endpoint+"/sites/"+wp.Site+"/posts/new/", bytes.NewBuffer(jsonData))
	if err != nil {
		return err.Error()
	}
	fmt.Println(data)
	reqME, err := http.NewRequest("GET", wp.Endpoint+"/rest/v1/me/", bytes.NewBuffer(jsonData))
	if err != nil {
		return err.Error()
	}
	reqME.Header.Set("Content-Type", "application/json")
	reqME.Header.Add("Authorization", "Bearer "+wp.Token)

	respMe, err := client.Do(reqME)
	if err != nil {
		return err.Error()
	}

	bME, err := io.ReadAll(respMe.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		LogError("konnte antwort nicht lesen", err)
		return err.Error()
	}
	fmt.Printf("Status: %d, Body: %s", respMe.StatusCode, string(bME))

	fmt.Println(wp.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+wp.Token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to create blog")
		return err.Error()
	}
	fmt.Println("Successfully created blog")
	b, err := io.ReadAll(resp.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		LogError("konnte antwort nicht lesen", err)
		return err.Error()
	}
	return fmt.Sprintf("Status: %d, Body: %s", resp.StatusCode, string(b))
}
