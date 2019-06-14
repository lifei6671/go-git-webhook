package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lifei6671/go-git-webhook/commands"
	_ "github.com/lifei6671/go-git-webhook/modules/filters"
	_ "github.com/lifei6671/go-git-webhook/routers"
)

func main() {

	commands.RegisterLogger()
	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RunCommand()

	commands.RegisterTaskQueue()

	commands.Run()
}
