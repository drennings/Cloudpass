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
	Id              string
	Share           int
	PublicIpAddress string
	*Job            // Embedded struct
}

// Work represents a unit of work sent to the worker
type Work struct {
	Id       string `json:"worker_id"`
	MasterIp string `json:"master_addr"`
	Hash     string `json:"hash_str"`
	HashType string `json:"hash_type"`
	Share    int    `json:"share"`
	Capacity int    `json:"cap"`
}
