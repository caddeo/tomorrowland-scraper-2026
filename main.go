package main

import (
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() (err error) {
	return nil
}
