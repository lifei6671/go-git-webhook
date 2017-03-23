package main

import (
	"go-git-webhook/commands"
	"go-git-webhook/modules/goclient"
	"fmt"
	"net/url"
)

func main()  {

	commands.RegisterLogger()
	commands.RegisterDataBase()
	commands.RegisterModel()
	commands.RunCommand()

	commands.RegisterTaskQueue()

	//commands.Run()

	token,err := goclient.GetToken("http://localhost:8081/token","root","123456")

	if err != nil {
		fmt.Println(err)
	}

	u := url.URL{Scheme: "ws", Host: "localhost:8081" , Path: "/socket"}

	client,err := goclient.Connection(u.String(),token)

	if err != nil {
		fmt.Println("Client:",err)
	}
	client.Send([]byte("aaaaaa"))

	res,_ := client.Read()

	fmt.Println(string(res))
}