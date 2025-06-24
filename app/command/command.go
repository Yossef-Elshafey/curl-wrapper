package command

import (
	"fmt"
	"strings"
)

type CommandHandler struct {
	Description string
	Executer    func()
	Command     string
}

var commandMap = make(map[string]*CommandHandler)
var userInput string

func (c *CommandHandler) Init() {
	c.Add("help", "Print help", c.printHelp)
}

func (c *CommandHandler) share(s string) {
	userInput = s
}

func (c *CommandHandler) GetUserInput() string {
	return userInput
}

func (c *CommandHandler) Add(command string, desc string, f func()) {
	if _, ok := commandMap[command]; ok {
		err := fmt.Sprintf("Trying to add an existing command: %s", command)
		panic(err)
	}

	commandMap[command] = &CommandHandler{
		Command:     command,
		Description: desc,
		Executer:    f,
	}

}

func (c *CommandHandler) shift() int {
	shift := 0
	for command := range commandMap {
		if len(command) > shift {
			shift = len(command)
		}
	}
	return shift
}

func (c *CommandHandler) printHelp() {
	for _, v := range commandMap {
		spaces := strings.Repeat(" ", c.shift()-len(v.Command))
		fmt.Printf("\n\r%s %s %s", v.Command, spaces, v.Description)
	}
}

func (c *CommandHandler) CommandFactory(inp string) {
	c.share(inp)
	reqAction := strings.Split(inp, " ")
	if action, ok := commandMap[reqAction[0]]; ok {
		action.Executer()
	} else {
		fmt.Printf("\n\rUnknown command: %s", reqAction[0])
	}
}
