package main

import (
	"fmt"
	"os"

	cmd "github.com/vimek-go/server-faker/cmd"
)

func main() {
	if err := cmd.PrepareCommand(); err != nil {
		fmt.Println("error preparing command", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
