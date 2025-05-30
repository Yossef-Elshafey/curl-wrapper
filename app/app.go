package app

import (
	"wcurl/app/command"
	"wcurl/app/models/wcurl"
	"wcurl/app/terminal"
)

func Run() {
	t := terminal.TerminalHandler{}
	ww := wcurl.WcurlWrapper{}
	co := command.CommandHandler{}
	t.Start(co)
	ww.Init()
	co.Init()
}
