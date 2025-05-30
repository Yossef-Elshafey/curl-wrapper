package terminal

import (
	"fmt"
	"wcurl/app/command"
	// provides a terminal wizard like, where i could controll the input while being written
	"golang.org/x/term"
	"os"
)

const (
	CTRL_C    = 3
	ENTER     = 13
	BACKSAPCE = 127
)

type TerminalHandler struct {
	currentCommand string
	histroy        []string
	currentBuffer  []byte
	exit           int
}

func (th *TerminalHandler) ctrlCHandler() {
	if th.currentBuffer[0] == CTRL_C {
		th.exit = 1
	}
}

func (th *TerminalHandler) asciiRangeHandler() {
	buf := th.currentBuffer[0]
	if buf >= 32 && buf <= 126 {
		fmt.Print(string(buf))
		th.currentCommand += string(buf)
	}
}

func (th *TerminalHandler) enterHandler(co command.CommandHandler) {
	if th.currentBuffer[0] == ENTER {
		switch th.currentCommand {
		case "":
			fmt.Printf("\n\r>> ")
		case "exit":
			fmt.Printf("\n\rExiting...\n\r")
			th.exit = 1
		case "clear":
			fmt.Print("\033[2J\033[H")
			fmt.Print(">> ")
		default:
			co.CommandFactory(th.currentCommand)
		}
		th.currentCommand = ""
	}
}

func (th *TerminalHandler) backspaceBehavior() {
	if th.currentBuffer[0] == BACKSAPCE {
		cmdLen := len(th.currentCommand)
		if cmdLen > 0 {
			fmt.Printf("\b")
			fmt.Printf(" ")
			fmt.Printf("\b")
			th.currentCommand = th.currentCommand[0 : cmdLen-1]
		}
	}
}

func (th *TerminalHandler) Start(co command.CommandHandler) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print(">> ")
	for {
		buf := make([]byte, 3)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}

		th.currentBuffer = buf
		th.ctrlCHandler()
		th.asciiRangeHandler()
		th.enterHandler(co)
		th.backspaceBehavior()

		if th.exit == 1 {
			return
		}
	}
}
