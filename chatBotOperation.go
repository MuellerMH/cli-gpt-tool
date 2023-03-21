package main

type ChatBotOperationen struct {
	Name string
	Desc string
}

var (
	Quit     = ChatBotOperationen{Name: "quit", Desc: "Bot beendent, Short: q"}
	Help     = ChatBotOperationen{Name: "help", Desc: "Hilfe Aufrufen, Short: h"}
	Clear    = ChatBotOperationen{Name: "clear", Desc: "Löscht den aktuellen Verlauf Short: c"}
	Save     = ChatBotOperationen{Name: "save", Desc: "Speichert die letzte antwort: use [ save FileName.txt ]"}
	Load     = ChatBotOperationen{Name: "load", Desc: "Läd eine Datei und führ den Inhalt als Request aus"}
	Mail     = ChatBotOperationen{Name: "mail", Desc: ""}
	History  = ChatBotOperationen{Name: "history", Desc: "Zeigt den aktuellen Chatverlauf an"}
	Token    = ChatBotOperationen{Name: "token", Desc: "Zeigt die aktuellen Token des Request an"}
	Assisten = ChatBotOperationen{Name: "asi", Desc: "Wechsel den Assistenen [use asi NAME ]"}
	Trello   = ChatBotOperationen{Name: "trello", Desc: "Speichert die Letzte Antwort als Trello Karte und nimmt die Frage als Titel"}
	Config   = ChatBotOperationen{Name: "config", Desc: "Zeigt die Config"}
	Default  = ChatBotOperationen{Name: "default", Desc: "Default"}
	Twitter  = ChatBotOperationen{Name: "twitter", Desc: "Interagiert mit Twitter"}
)

func OperationFromString(s string) ChatBotOperationen {
	switch s {
	case "q":
		return Quit
	case "quit":
		return Quit
	case "h":
		return Help
	case "help":
		return Help
	case "c":
		return Clear
	case "clear":
		return Clear
	case "save":
		return Save
	case "load":
		return Load
	case "mail":
		return Mail
	case "history":
		return History
	case "token":
		return Token
	case "asi":
		return Assisten
	case "trello":
		return Trello
	case "twitter":
		return Twitter
	case "config":
		return Config
	default:
		return Default
	}
}
