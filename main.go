package main

import (
	"github.com/lifei6671/go-git-webhook/commands"
)

func main()  {

	commands.RegisterLogger()
	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RunCommand()

	commands.RegisterTaskQueue()

	commands.Run()
}