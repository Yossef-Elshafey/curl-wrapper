package app

import (
	"wcurl/app/command"
	"wcurl/app/models/terminal"
	"wcurl/app/models/wcurl"
)

func Run() {
	t := terminal.TerminalHandler{}
	ww := wcurl.WcurlWrapper{}
	co := command.CommandHandler{}
	t.Start(co)
	ww.Init()
	co.Init()
}
