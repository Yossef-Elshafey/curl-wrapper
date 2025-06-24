package app

import (
	"wcurl/app/command"
	"wcurl/app/models"
	"wcurl/app/terminal"
)

func Run() {
	t := terminal.TerminalHandler{}
	ww := models.WcurlWrapper{}
	co := command.CommandHandler{}
	ww.Init()
	co.Init()
	t.Start(co)
}
