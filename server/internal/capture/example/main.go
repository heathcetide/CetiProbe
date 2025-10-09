package main

import (
	"log"
	"probe/internal/capture"
)

func main() {
	capturer, err := capture.NewCapturer("en1")
	if err != nil {
		log.Fatalf("Error when creating capturer: %v", err)
	}

	err = capturer.Start()
	if err != nil {
		log.Println("Error when starting capturer:", err)
	}
}
