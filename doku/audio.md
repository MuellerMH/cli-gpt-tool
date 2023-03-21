
func main() {
	pContext, err := oto.NewContext(44100, 1, 2, 8192)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer pContext.Close()
	testPlayer(pContext)
}

Dieses Programm implementiert eine Funktion zum Abspielen einer MP3-Datei. Es öffnet die Datei, decodiert sie und initialisiert den Player, um sie abzuspielen. Der Player wird dann mit der decodierten Datei gefüllt und schließlich abgespielt.