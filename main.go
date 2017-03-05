package main

import (
	"go-git-webhook/commands"
)


func main()  {

	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RunCommand()

	commands.Run()
}

