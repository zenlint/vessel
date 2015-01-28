package main

import (
	"log"
	"os"

	"github.com/dockercn/vessel/modules/core"
)

func main() {
	log.Println("Creating solution...")
	sln, err := core.NewSolutionFromFile(os.Args[1])
	if err != nil {
		log.Fatalf("Error creating solution: %v", err)
	}

	log.Println("Building solution...")
	imageId, err := sln.Build()
	if err != nil {
		log.Fatalf("Error building solution: %v", err)
	}

	log.Println("Starting container...")
	if err = sln.Start(imageId); err != nil {
		log.Fatalf("Error starting solution: %v", err)
	}
}
