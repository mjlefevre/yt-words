package main

import (
	"log"

	"github.com/mjlefevre/sanoja/cmd/cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}
