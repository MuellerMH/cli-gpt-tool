
## CommandCLI
Die `CommandCLI`-Struktur ist eine Struktur, welche Informationen über einen eingegebenen Befehl enthält. Sie wird aus einer Zeichenkette erstellt, die ein Benutzer in einem Chatbot eingibt. Sie enthält die Wörter, die Befehl, die Parameter und die `ChatBotOperationen`, die aufgerufen werden können.

Die Struktur wird durch die Funktion `NewCommandCli` erstellt, die den Text liest und die entsprechenden Felder in der Struktur setzt. Dazu wird der Text in Wörter aufgeteilt, der Befehl in Kleinbuchstaben umgewandelt und die Parameter zu einem String zusammengefasst.