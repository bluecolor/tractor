package main

import (
	"github.com/bluecolor/tractor/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
