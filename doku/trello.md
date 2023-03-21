
Das Package main implementiert Funktionen, die es ermöglichen, mit dem [Trello-API](https://developers.trello.com/reference) zu interagieren. 

Es wird ein `TrelloConnect`-Struct erstellt, in dem das `trello.Client`-Objekt, ein `Board`-Objekt und eine Liste von `trello.List`-Objekten gespeichert werden. Darüber hinaus enthält das Struct ein Konfigurationsobjekt und ein String, der die ID des Trello-Boards enthält. 

Es gibt Funktionen, um eine Trello-Karte hinzuzufügen, eine Karte in die nächste Liste zu verschieben und eine To-Do-Liste auszulesen. Außerdem gibt es eine Funktion, um sich in ein Trello-Board einzuloggen. 

Der Code ermöglicht es, mit dem Trello-API zu interagieren und stellt Funktionen bereit, um Karten hinzuzufügen, zu verschieben und auszulesen.