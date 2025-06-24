package terminal

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"wcurl/app/command"
	"wcurl/app/utils"
)

const (
	CTRL_C    = 3
	ENTER     = 13
	BACKSAPCE = 127
	UP        = 65
	DOWN      = 66
)

type TerminalHandler struct {
	ch             command.CommandHandler
	currentCommand string
	history        *utils.History
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

func (th *TerminalHandler) enterBehavior() {
	if th.currentBuffer[0] == ENTER {
		switch th.currentCommand {
		case "":
			fmt.Printf("\n\r>> ")
		case "exit":
			th.exit = 1
		case "clear":
			th.history.Save(th.currentCommand)
			fmt.Print("\033[2J\033[H")
			fmt.Print(">> ")
		default:
			th.ch.CommandFactory(th.currentCommand)
			th.history.Save(th.currentCommand)
			fmt.Printf("\n\r>> ")
		}
		th.currentCommand = ""
	}
}

func (th *TerminalHandler) deleteChar() {
	cmdLen := len(th.currentCommand)
	if cmdLen > 0 {
		fmt.Printf("\b")
		fmt.Printf(" ")
		fmt.Printf("\b")
		th.currentCommand = th.currentCommand[0 : cmdLen-1]
	}
}

func (th *TerminalHandler) backspaceBehavior() {
	if th.currentBuffer[0] == BACKSAPCE {
		th.deleteChar()
	}
}

func (th *TerminalHandler) clearPrompt() {
	for range len(th.currentCommand) {
		th.deleteChar()
	}
}

func (th *TerminalHandler) arrowUpBehavior() {
	if th.currentBuffer[2] == UP {
		th.clearPrompt()
		th.currentCommand = th.history.Prev()
		fmt.Printf("%s", th.currentCommand)
	}
}

func (th *TerminalHandler) arrowDownBehavior() {
	if th.currentBuffer[2] == DOWN {
		th.clearPrompt()
		th.currentCommand = th.history.Next()
		fmt.Printf("%s", th.currentCommand)
	}
}

func (th *TerminalHandler) behaviorsWrapper() {
	th.ctrlCBehavior()
	th.asciiRangeBehavior()
	th.enterBehavior()
	th.backspaceBehavior()
	th.arrowUpBehavior()
	th.arrowDownBehavior()
}

func (th *TerminalHandler) Start(co command.CommandHandler) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	if th.history == nil {
		th.history = utils.NewHistory()
	}

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
