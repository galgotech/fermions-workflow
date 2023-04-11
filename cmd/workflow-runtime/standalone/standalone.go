package main

import (
	"log"

	"github.com/galgotech/fermions-workflow/internal/cmd"
)

func main() {
	err := cmd.Standalone()
	if err != nil {
		log.Fatal(err)
	}
}
