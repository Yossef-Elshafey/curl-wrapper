package models

import (
	"fmt"
	"regexp"
	"strings"
	"wcurl/app/command"
	"wcurl/app/utils"
)

type Endpoint struct {
	Ep             map[string][]string `json:"endpoints"`
	storage        utils.Storage
	commandHandler command.CommandHandler
}

func (e Endpoint) AddEndpoint() Endpoint {
	ex := e.ExtractEndpoint()
	if e.Ep == nil {
		e.Ep = make(map[string][]string)
	}

	for _, v := range e.Ep[ex] {
		if strings.Compare(v, e.commandHandler.GetUserInput()) == 0 {
			return e
		}
	}

	e.Ep[ex] = append(e.Ep[ex], e.commandHandler.GetUserInput())
	return e
}

func (e Endpoint) ListEndPoints() {
	sep := strings.Repeat("=", 10)
	i := 0
	keys := make([]string, 0)
	for endpoint, commands := range e.Ep {
		keys = append(keys, endpoint)
		fmt.Printf("\n\r%s %s %s\n\r", sep, endpoint, sep)
		i += 1
		for j, command := range commands {
			fmt.Printf("\n\r%d) -> %s\n\r", j, command)
		}
	}
}

func (e *Endpoint) ExtractEndpoint() string {
	// Returns: endpoint from curl user input, localhost:3000/admin -> /admin

	input := e.commandHandler.GetUserInput()
	re := regexp.MustCompile(`(?:^|\s)(?:https?://)?([a-zA-Z0-9.-]+(?::[0-9]+)?)(/.+?)(?:\s|$)`)
	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 {
		return "/"
	}

	return matches[0][2]
}
