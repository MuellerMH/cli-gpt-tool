package main

import (
	"fmt"

	"github.com/adlio/trello"
)

var err error

type TrelloConnect struct {
	TrelloClient  *trello.Client
	TrelloBoard   *trello.Board
	TrelloListe   []*trello.List
	Config        *ConfigBot
	TrelloBoardId string
	LastUsedCard  *trello.Card
}

func NewTrelloConnect(config *ConfigBot) *TrelloConnect {
	tC := TrelloConnect{Config: config}
	tC.TrelloBoardId = config.TrelloAzubiBoard
	return &tC
}

// Funktion zum Hinzufügen einer Trello-Karte
func (t *TrelloConnect) AddTrelloCard(title string, description string) {
	if !t.Config.UseTrello {
		return
	}
	// Eine neue Trello-Karte erstellen
	card := &trello.Card{
		Name: title,
		Desc: description,
	}
	fmt.Printf("Card Inhalt %s", GetJSONString(card))
	// Trello-Karte hinzufügen
	trelloListen, err := t.GetLists()
	if err != nil {
		LogError("Trello Liste", err)
		return
	}
	if trelloListen == nil {
		LogError("Konnte Trello Liste nicht laden", err)
		return
	}
	list := trelloListen[0]
	listC, _ := t.TrelloClient.GetList(list.ID, trello.Defaults())
	err = listC.AddCard(card, trello.Defaults())
	if err != nil {
		LogError("Trello Liste", err)
		return
	}
	// Erfolgsmeldung anzeigen
	fmt.Println("Trello-Karte erfolgreich hinzugefügt!")
}

func (t *TrelloConnect) MoveCardToNextList(description string, response string) string {
	if t.TrelloClient == nil {
		t.LoginToTrelloBoard()
	}
	board, err := t.GetBoard()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if board == nil {
		return ""
	}
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if lists == nil {
		return ""
	}
	t.LastUsedCard = nil
	for i, list := range lists {
		if i == 0 {
			cards, err := list.GetCards(trello.Defaults())
			if err != nil {
				fmt.Println(err)
				return ""
			}
			for _, card := range cards {
				t.LastUsedCard = card
				break
			}
		}
		if i == 1 {
			if t.LastUsedCard == nil {
				break
			}
			t.LastUsedCard.Desc = description
			t.LastUsedCard.AddComment(response, trello.Defaults())
			t.LastUsedCard.MoveToList(list.ID, trello.Defaults())
		}
	}
	return ""
}
func (t *TrelloConnect) GetTodoFromList() string {
	if t.TrelloClient == nil {
		t.LoginToTrelloBoard()
	}
	board, err := t.GetBoard()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if board == nil {
		return ""
	}
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if lists == nil {
		return ""
	}
	for _, list := range lists {
		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			fmt.Println(err)
			return ""
		}
		for _, card := range cards {
			t.LastUsedCard = card
			return card.Name
		}
	}

	return ""
}

// Funktion zum Login in ein Trello-Board
func (t *TrelloConnect) LoginToTrelloBoard() *trello.Client {
	if !t.Config.UseTrello {
		return nil
	}
	if t.TrelloClient != nil {
		return t.TrelloClient
	}
	return trello.NewClient(t.Config.TrelloKey, t.Config.TrelloSecret)
}

func (t *TrelloConnect) GetBoard() (*trello.Board, error) {
	t.TrelloClient = t.LoginToTrelloBoard()
	if t.TrelloClient == nil {
		fmt.Println("Kein Client forhanden")
		return nil, nil
	}
	fmt.Printf("Versuche Board %s zu finden.", t.TrelloBoardId)
	t.TrelloBoard, err = t.TrelloClient.GetBoard(t.TrelloBoardId, trello.Defaults())
	if err != nil {
		fmt.Println("No Board found for " + t.TrelloBoardId)
		fmt.Println(err.Error())
		return nil, err
	}
	return t.TrelloBoard, nil
}

func (t *TrelloConnect) GetLists() ([]*trello.List, error) {
	trelloBoard, err := t.GetBoard()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	trelloListe, err := trelloBoard.GetLists(trello.Defaults())
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return trelloListe, nil
}
