package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

const pollInterval = 10 * time.Second
const apiPort = ":80"

func main() {
	fmt.Println("Starting...")

	// Create a Manager for the specified region.
	region := aws.String(os.Getenv("AWS_REGION"))
	man := NewManager(region)
	fmt.Printf("Created manager in region %s\n", *region)

	// Create and start the API
	api := NewAPI(apiPort, man)
	err := api.Serve()

	if err != nil {
		fmt.Printf("An error occurred: %v", err)
	}
	fmt.Printf("Started webserver on port %s\n", apiPort)

	// Check if new records are added every 'pollInterval'
	//ticker := time.NewTicker(pollInterval)
	//for time := range ticker.C {
	if record := newRecordAdded(); record != nil {
		// Process new records
		//fmt.Printf("Processing record %s at %v.\n", record.Id, time.String())

		// Create and start the job
		job := JobFromRecord(record)
		err := man.StartJob(job)
		if err != nil {
			fmt.Printf("An error occurred when starting job for record %s: %v", record.Id, err.Error())
		}
	}
	//}
}

// newRecordAdded is a callback for new rows added to RDS
func newRecordAdded() *Record {
	return &Record{
		Id:       "1",
		Hash:     "ef775988943825d2871e1cfa75473ec0",
		HashType: "md5",
		Name:     "Bob",
		Capacity: 1,
	}
}
