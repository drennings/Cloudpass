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
	Length    int
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
	Length    int
}

// Worker represents a single EC2 instance used for a Job.
type Worker struct {
	Id              string
	Share           int
	PublicIpAddress string
	*Job            // Embedded struct
}

// Work represents a unit of work sent to the worker
type Work struct {
	Id       string `json:"workerId"`
	MasterIp string `json:"masterAddr"`
	Hash     string `json:"hashStr"`
	HashType string `json:"hashType"`
	Share    int    `json:"share"`
	Capacity int    `json:"cap"`
	Length   int    `json:"length"`
}
