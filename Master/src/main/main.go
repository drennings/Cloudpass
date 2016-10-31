package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

const pollInterval = 10 * time.Second
const apiPort = ":8080"

func mainTest() {
	region := aws.String(os.Getenv("AWS_REGION"))
	man := NewManager(region)
	job := JobFromRecord(newRecordAdded())
	man.StartJob(job)

	// Stop the first worker for this job
	for _, v := range man.Jobs[job.Id].Workers {
		man.runCommand(v, "whoami")
		man.runCommands(v, []string{"touch a", "ls", "pwd", "cd /", "ls"})
		man.stopWorker(v)
		break
	}
	//man.StopJob(job)
}

func main() {
	fmt.Println("Starting...")

	// Create a Manager for the specified region.
	region := aws.String(os.Getenv("AWS_REGION"))
	man := NewManager(region)
	fmt.Printf("Created manager in region %s\n", *region)

	// Create and start the API
	api := NewAPI(apiPort)
	man.AttachAPI(api)
	fmt.Printf("Started webserver on port %s\n", apiPort)

	// Check if new records are added every 'pollInterval'
	ticker := time.NewTicker(pollInterval)
	for time := range ticker.C {
		if record := newRecordAdded(); record != nil {
			// Process new records
			fmt.Printf("Processing record %s at %v.\n", record.Id, time.String())

			// Create and start the job
			job := JobFromRecord(record)
			err := man.StartJob(job)
			if err != nil {
				fmt.Printf("An error occurred when starting job for record %s: %v", record.Id, err.Error())
			}
		}
	}
}

// newRecordAdded is a callback for new rows added to RDS
func newRecordAdded() *Record {
	return &Record{
		Id:       "1",
		Hash:     "233fedii3if90398j3",
		Capacity: 2,
	}
}
