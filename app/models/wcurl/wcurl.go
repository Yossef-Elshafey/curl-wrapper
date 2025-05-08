package wcurl

import (
	"fmt"
	"wcurl/app/command"
	"wcurl/app/models/endpoiot"
)

type WcurlWrapper struct {
	Data     map[string]endpoint.Endpoint
	Listener command.CommandHandler
}

func (w *WcurlWrapper) SetCommand() {
	w.Listener.Add("init", "initialize a new project", w.NewProject)
	w.Listener.Add("ep", "does something", w.NewProject)
}

func (w *WcurlWrapper) NewProject() {
	fmt.Println("hello idiot")
}
