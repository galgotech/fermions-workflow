package main

import (
	"log"

	"github.com/galgotech/fermions-workflow/internal/cmd"
)

func main() {
	err := cmd.Server()
	if err != nil {
		log.Fatal(err)
	}
}
