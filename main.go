package main

import (
	"fmt"

	"github.com/bluecolor/tractor/commands"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	err := commands.TractorCmd.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
	}
}
