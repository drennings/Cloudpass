package main

import "time"

// InstanceCfg contains the configuration for an EC2 instance.
type InstanceCfg struct {
	ImageId          *string
	InstanceType     *string
	KeyName          *string
	SecurityGroupIds []*string
}

// Record defines a record as it is stored in RDS
type Record struct {
	Id        string
	Name      string
	Email     string
	Capacity  int
	Timelimit time.Time
	Hash      string
	HashType  string
}

// Job represents a single hash to be cracked.
type Job struct {
	Id        string
	Name      string
	Email     string
	Workers   map[string]*Worker
	Capacity  int
	Timelimit time.Time
	Hash      string
	HashType  string
}

// Worker represents a single EC2 instance used for a Job.
type Worker struct {
	Id    string
	Share int
	*Job  // Embedded struct
}
