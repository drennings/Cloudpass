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
	fmt.Println("Start time: %v", time.Now())
	// Create a Manager for the specified region.
	region := aws.String(os.Getenv("AWS_REGION"))
	man := NewManager(region)
	fmt.Printf("Created manager in region %s\n", *region)

	go func() {
		if record := newRecordAdded(); record != nil {
			// Process new records
			fmt.Printf("Processing record %s at %v.\n", record.Id, time.Now())

			// Create and start the job
			job := JobFromRecord(record)
			err := man.StartJob(job)
			if err != nil {
				fmt.Printf("An error occurred when starting job for record %s: %v", record.Id, err.Error())
			}
		}
		//}
	}()

	// Create and start the API
	api := NewAPI(apiPort, man)
	err := api.Serve()
	if err != nil {
		fmt.Printf("An error occurred: %v", err)
	}
	fmt.Printf("Started webserver on port %s\n", apiPort)

}

// newRecordAdded is a callback for new rows added to RDS
func newRecordAdded() *Record {
	return &Record{
		Id:       "1",
		Hash:     "c8c605999f3d8352d7bb792cf3fdb25b",
		HashType: "md5",
		Name:     "Bob",
		Capacity: 8,
		Length:   9,
	}
}
