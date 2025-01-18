package main

import (
	"log"

	"github.com/mjlefevre/sanoja/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
