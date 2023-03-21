package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

type Mp3Player struct {
	DecoderContext *oto.Context
	ConfigBot      *ConfigBot
	ActiveMp3File  string
}

func NewMp3Player(sampleMp3_file string, config *ConfigBot) *Mp3Player {

	file, err := os.Open(sampleMp3_file)
	if err != nil {
		LogError(ErrorNoMp3File, err)
	}
	d, _ := mp3.NewDecoder(file)
	pContext, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		LogError(ErrorNoMp3File, err)
	}
	return &Mp3Player{DecoderContext: pContext, ConfigBot: config}
}

func (mp *Mp3Player) PlayMP3File(filename string) error {
	if !mp.ConfigBot.UseSound {
		return nil
	}
	// Open the MP3 file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the MP3 file
	d, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}
	// Initialize the player
	p := mp.DecoderContext.NewPlayer()
	if err != nil {
		return err
	}
	defer p.Close()

	// Play the MP3 file
	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	go mp.RemoveFile(mp.ActiveMp3File + "-text-to-speech.mp3")
	return nil
}

func (mp *Mp3Player) CreateUniqueName() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", hex.EncodeToString(b))
}

func (mp *Mp3Player) RemoveFile(file_name string) {
	if mp.ConfigBot.SaveAudioFiles {
		return
	}
	time.Sleep(5000)
	err := os.Remove(file_name)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (mp *Mp3Player) TextToSpeech(text string) error {
	if !mp.ConfigBot.UseSound {
		return nil
	}
	mp.ConfigBot.CheckAWSCredentials()
	mp.ActiveMp3File = mp.CreateUniqueName()
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(mp.ConfigBot.AWSRegion),
		Credentials: credentials.NewSharedCredentials(mp.ConfigBot.AWSKey, mp.ConfigBot.AWSSecret),
	})
	if err != nil {
		// Handle Session creation error
		return err
	}
	// Create Polly client
	svc := polly.New(sess)

	// Convert text to speech
	result, err := svc.SynthesizeSpeech(&polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		Text:         aws.String(text),
		VoiceId:      aws.String("Vicki"),
	})
	if err != nil {
		return err
	}

	// Save to file
	f, err := os.Create(mp.ActiveMp3File + "-text-to-speech.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, result.AudioStream)
	if err != nil {

		return err
	}

	go mp.PlayMP3File(mp.ActiveMp3File + "-text-to-speech.mp3")
	return nil
}
