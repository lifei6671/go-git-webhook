package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lifei6671/go-git-webhook/routers"
	_ "github.com/lifei6671/go-git-webhook/modules/filters"
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