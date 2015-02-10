package main

import (
	"log"
	"os"

	"github.com/dockercn/vessel/modules/flow"
	"github.com/dockercn/vessel/modules/sln"
)

func main() {
	log.Println("Creating solution...")
	sln, err := sln.NewSolutionFromFile(os.Args[1])
	if err != nil {
		log.Fatalf("Fail to create solution: %v", err)
	}

	stage := flow.NewStage()
	stage.SetJob(sln)
	if err = stage.Run(); err != nil {
		log.Fatalf("Fail to run stage: %v", err)
	}

	// log.Println("Starting container...")
	// if err = sln.Start(imageId); err != nil {
	// 	log.Fatalf("Error starting solution: %v", err)
	// }
}
