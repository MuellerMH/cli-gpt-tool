package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	var input string
	var totalTasksCompleted int = 0
	createTodoList(input,totalTasksCompleted)
	os.Exit(0)
}

func runPomodoroCycle(todoList []string, input string, totalTasksCompleted int) {
	pomodoroDuration := time.Duration(25)
	pomodoroBreakDuration := time.Duration(5)
	timer := time.NewTimer(pomodoroDuration * time.Minute)
	for _, item := range todoList {
		fmt.Println("Start working on:", item)
		<-timer.C
		totalTasksCompleted+=1
		todoList = todoList[1:]
		fmt.Println("Aufgabe erledigt:", item)
		fmt.Println("Aufgaben erledigt:", totalTasksCompleted)
		if (len(todoList) > 0) {
			fmt.Println("Du hast noch ", len(todoList)," tasks zu erledigen.")
			fmt.Println("Nimm Dir 5 Minuten Pause!")
			timer.Reset(pomodoroBreakDuration * time.Minute)
			<-timer.C
		} else {
			fmt.Println("Füge neue Pomodoros zu deiner ToDo Liste hinzu.")
		}
		timer.Reset(pomodoroDuration * time.Minute)
	}
	createTodoList(input, totalTasksCompleted)
}

func createTodoList(input string, totalTasksCompleted int) {
	todoList := []string{}
	fmt.Println("Füge neue Aufgaben deiner TodoListe hinzu (q to quit): ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input = scanner.Text()
		if(input == "q") {
			break
		}
		todoList = append(todoList, input)
	}

	// Den Timer zurücksetzen und den Benutzer erinnern, eine Pause einzulegen
	// Warten, bis der Timer abläuft
	runPomodoroCycle(todoList, input, totalTasksCompleted)
}

