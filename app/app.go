package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"wcurl/app/command"
	"wcurl/app/models/wcurl"
)

func setShellCommands(co command.CommandHandler) {
	co.Add("clear", "Clear shell", func() {})
	co.Add("exit", "Exit program", func() {})
}

func Run() {
	ww := wcurl.WcurlWrapper{}
	co := command.CommandHandler{}
	ww.Listener = co

	ww.SetCommand()
	setShellCommands(co)
	co.Init()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">> ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		switch input {
		case "":
			continue
		case "exit":
			fmt.Println("Exiting...")
			return
		case "clear":
			fmt.Print("\033[2J\033[H")
		default:
			co.CommandFactory(input)
		}
	}
}
