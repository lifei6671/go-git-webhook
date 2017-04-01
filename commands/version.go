package commands

import (
	"os"
	"fmt"
)

// Version 获取当前版本.
func Version()  {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("v0.1")
		os.Exit(2)
	}
}
