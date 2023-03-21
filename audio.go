package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func PlayMP3File(filename string, pContext *oto.Context) error {
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
	p := pContext.NewPlayer()
	if err != nil {
		return err
	}
	defer p.Close()

	// Play the MP3 file
	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func testPlayer(pContext *oto.Context) {
	err := PlayMP3File("text-to-speech.mp3", pContext)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
