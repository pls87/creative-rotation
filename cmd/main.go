package main

import (
	"fmt"
	"os"

	"github.com/pls87/creative-rotation/cmd/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Printf("Couldn't run app: %s", err)
		os.Exit(1)
	}
}
