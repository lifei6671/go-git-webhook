package main

import (
	"go-git-webhook/commands"
)

func main()  {

	commands.RegisterLogger()
	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RunCommand()

	commands.RegisterTaskQueue()

	commands.Run()
}