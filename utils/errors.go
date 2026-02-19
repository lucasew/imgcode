package utils

import (
	"log"
	"os"
)

// Check handles errors by logging them and exiting.
func Check(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
