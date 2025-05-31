package terminal

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"wcurl/app/command"
)

const (
	CTRL_C    = 3
	ENTER     = 13
	BACKSAPCE = 127
	UP        = 65
)

type TerminalHandler struct {
	ch             command.CommandHandler
	currentCommand string
	history        []string
	currentBuffer  []byte
	exit           int
}

func (th *TerminalHandler) ctrlCBehavior() {
	if th.currentBuffer[0] == CTRL_C {
		th.exit = 1
	}
}

func (th *TerminalHandler) asciiRangeBehavior() {
	buf := th.currentBuffer[0]
	if buf >= 32 && buf <= 126 {
		fmt.Print(string(buf))
		th.currentCommand += string(buf)
	}
}

func (th *TerminalHandler) saveToHistory() {
	th.history = append(th.history, th.currentCommand)
}

func (th *TerminalHandler) enterBehavior() {
	if th.currentBuffer[0] == ENTER {
		switch th.currentCommand {
		case "":
			fmt.Printf("\n\r>> ")
		case "exit":
			th.exit = 1
		case "clear":
			th.saveToHistory()
			fmt.Print("\033[2J\033[H")
			fmt.Print(">> ")
		default:
			th.ch.CommandFactory(th.currentCommand)
			th.saveToHistory()
			fmt.Printf("\n\r>> ")
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

func (th *TerminalHandler) arrowUpBehavior() {
	fmt.Printf("\r\n%v", th.currentBuffer)
	if th.currentBuffer[0] == UP {
		fmt.Printf("\n\rArrow up detected")
	}
}

func (th *TerminalHandler) behaviorsWrapper() {
	th.ctrlCBehavior()
	th.asciiRangeBehavior()
	th.enterBehavior()
	th.backspaceBehavior()
	th.arrowUpBehavior()
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
		th.behaviorsWrapper()

		if th.exit == 1 {
			fmt.Printf("\n\rExiting...\n\r")
			return
		}
	}
}
