package main

import (
	"cli_todo_application/commands"
	"os"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
