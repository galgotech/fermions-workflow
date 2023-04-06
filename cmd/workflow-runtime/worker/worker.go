package main

import (
	"log"

	"github.com/galgotech/fermions-workflow/internal/cmd"
)

func main() {
	err := cmd.WorkerStandalone()
	if err != nil {
		log.Fatal(err)
	}
}
