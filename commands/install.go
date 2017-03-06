package commands

import (
	"os"
	"flag"
	"go-git-webhook/models"
	"fmt"
	"strings"
)

func Install()  {

	if len(os.Args) > 2 && os.Args[1] == "install" {

		account := flag.String("account","admin","Administrator account.")
		pwd := flag.String("password","","Administrator password.")
		email := flag.String("email","","Administrator email.")

		flag.CommandLine.Parse(os.Args[2:])

		password := strings.TrimSpace(*pwd)

		if(password == ""){
			fmt.Println("Administrator password  is required.")
			os.Exit(2)
		}
		if(*email == ""){
			fmt.Println("Administrator email is required")
			os.Exit(2)
		}
		fmt.Println("|" + password + "|")

		member := models.NewMember()
		member.Account = *account
		member.Password = password
		member.Email = *email

		if err := member.Add();err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		fmt.Println("ok")
		os.Exit(2)
	}
}
