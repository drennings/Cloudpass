package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Manager manages the workers and the jobs.
type Manager struct {
	EC2Svc *ec2.EC2
	API    *API
	Jobs   map[string]*Job
	m      sync.Mutex // Guards jobs
}

// NewManager returns a new Manager struct.
func NewManager(region *string) *Manager {
	svc := ec2.New(session.New(), &aws.Config{Region: region})
	return &Manager{
		EC2Svc: svc,
		Jobs:   make(map[string]*Job),
	}
}

// createWorker starts a new EC2 instance.
// It returns when the worker is ready.
func (man *Manager) createWorker(job *Job, share int) (*Worker, error) {
	inst, err := man.createInstance()
	if err != nil {
		return nil, err
	}

	// Wait until the instance is up and running
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("instance-id"),
				Values: []*string{inst.InstanceId},
			},
		},
	}

	fmt.Printf("%s: waiting to be ready...\n", *inst.InstanceId)
	man.EC2Svc.WaitUntilInstanceRunning(params)
	fmt.Printf("%s: instance ready.\n", *inst.InstanceId)

	return &Worker{
		Id:    *inst.InstanceId,
		Share: share,
		Job:   job,
	}, nil
}

// startWorker sets up the Woker and starts working
func (man *Manager) startWorker(w *Worker) error {
	err := man.runCommands(w, []string{
		"sudo apt-get update",
		"sudo apt-get -y install git",
		"git clone https://github.com/drennings/Cloudpass",
		"cd Cloudpass/Worker && chmod +x make.sh && sudo ./make.sh &",
	})
	if err != nil {
		return err
	}
	return nil
}

// stopWorker stops a worker (running EC2 instance).
func (man *Manager) stopWorker(worker *Worker) error {
	fmt.Printf("%s: stopping worker.\n", worker.Id)
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(worker.Id)},
	}
	_, err := man.EC2Svc.TerminateInstances(input)
	return err
}

// StartJob starts a job and starts the necessary workers
func (man *Manager) StartJob(job *Job) error {
	fmt.Printf("%s: starting job.\n", job.Id)

	// Create and start a number of Workers equal to the capacity of the job
	for i := 0; i < job.Capacity; i++ {
		go man.createStartSubmit(job, i)
	}
	return nil
}

func (man *Manager) createStartSubmit(job *Job, i int) {
	worker, err := man.createWorker(job, i)
	fmt.Println("Starting worker...")
	if err != nil {
		fmt.Printf("Error creating worker: %v", err)
		return
	}
	err = man.startWorker(worker)
	if err != nil {
		fmt.Printf("Error starting worker: %v", err)
		return
	}

	fmt.Println("Submitting work...")
	// Submit the work to the worker
	err = man.submitWork(worker, i, job)
	if err != nil {
		fmt.Printf("Error submitting work: %v", err)
		return
	}

	// Add it to the Workers map for tracking
	man.m.Lock()
	defer man.m.Unlock()
	job.Workers[worker.Id] = worker

	// Job was successfully started, add it to the manager
	man.Jobs[job.Id] = job
}

// StopJob stops a job and stops all its associated workers.
func (man *Manager) StopJob(job *Job) error {
	fmt.Printf("%s: stopping job.\n", job.Id)
	var errors []string

	// Stop all workers
	for _, v := range job.Workers {
		err := /*go*/ man.stopWorker(v)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	// Job was successfully terminated, remove it from the manager
	man.m.Lock()
	defer man.m.Unlock()
	delete(man.Jobs, job.Id)
	return nil
}

// createInstance creates and returns an EC2 instance.
func (man *Manager) createInstance() (*ec2.Instance, error) {
	params := &ec2.RunInstancesInput{
		ImageId:          aws.String(os.Getenv("IMG_ID")),
		InstanceType:     aws.String(os.Getenv("INST_TYPE")),
		MaxCount:         aws.Int64(1),
		MinCount:         aws.Int64(1),
		KeyName:          aws.String(os.Getenv("PEM_NAME")),
		SecurityGroupIds: []*string{aws.String(os.Getenv("SEC_GROUP"))},
	}
	res, err := man.EC2Svc.RunInstances(params)
	if err != nil {
		return nil, err
	}

	inst := res.Instances[0]
	fmt.Printf("%s: created new instance.\n", *inst.InstanceId)
	return inst, nil
}

// runCommand runs a command on a worker instance through SSH.
func (man *Manager) runCommand(worker *Worker, cmd string) error {
	inst, err := man.getWorkerInstance(worker)
	if err != nil {
		return err
	}

	// Open PEM file
	pemPath := os.Getenv("PEM_PATH")
	man.m.Lock()
	pemBytes, err := ioutil.ReadFile(pemPath)
	man.m.Unlock()
	if err != nil {
		return err
	}

	// Obtain private key
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return err
	}

	// Connect to the remote server and perform the SSH handshake
	config := &ssh.ClientConfig{
		User:    "ubuntu",
		Auth:    []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout: 5 * time.Second,
	}
	fmt.Printf("%s: executing command: %s\n", *inst.InstanceId, cmd)
	addr := fmt.Sprintf("%s:%d", *inst.PublicIpAddress, 22)

	// Store IP for later use
	worker.PublicIpAddress = *inst.PublicIpAddress

	// Retry SSH until successful
	var conn *ssh.Client
	try, max, interval := 1, 5, 10*time.Second
	for conn == nil && try <= max {
		conn, err = ssh.Dial("tcp", addr, config)
		if err != nil {
			// Timeout occurred
			fmt.Printf("%v (%d/%d), trying again in %v...\n", err, try, max, interval)
			time.Sleep(interval)
		}
		try++
	}
	defer conn.Close()

	var stdoutBuf bytes.Buffer
	session, err := conn.NewSession()
	for err != nil {
		session, err = conn.NewSession()
	}
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	for err != nil {
		err = session.Run(cmd)
		defer session.Close()
	}

	return nil
}

// getWorkerInstance returns the AWS instance corresponding to a worker
func (man *Manager) getWorkerInstance(w *Worker) (*ec2.Instance, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(w.Id)},
			},
		},
	}

	resp, err := man.EC2Svc.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			return inst, err
		}
	}

	return nil, fmt.Errorf("Could not find running instance.")
}

// Send a POST request with a portion of the job to the worker
func (man *Manager) submitWork(worker *Worker, share int, job *Job) error {
	ip := GetLocalIP()
	fmt.Printf("My IP: %s", ip)

	work := &Work{
		Id:       worker.Id,
		MasterIp: string(ip),
		Hash:     worker.Hash,
		HashType: worker.HashType,
		Share:    share,
		Capacity: worker.Capacity,
		Length:   worker.Length,
	}

	workJSON, err := json.Marshal(work)
	if err != nil {
		return err
	}

	fmt.Printf("Submitting work: %v", work)

	workerAddr := fmt.Sprintf("http://%s", worker.PublicIpAddress)
	resp, err := http.Post(workerAddr+"/start", "application/json", bytes.NewBuffer(workJSON))
	if err != nil {
		return err
	}
	fmt.Printf("Resp: %#v", resp)

	return nil
}

// JobFromRecord converts a DynamoDB record to a Job
func JobFromRecord(rec *Record) *Job {
	return &Job{
		Id:        rec.Id,
		Name:      rec.Name,
		Email:     rec.Email,
		Capacity:  rec.Capacity,
		Timelimit: rec.Timelimit,
		Workers:   make(map[string]*Worker),
		Hash:      rec.Hash,
		HashType:  rec.HashType,
		Length:    rec.Length,
	}
}

// Runs a list of commands on a Worker.
func (man *Manager) runCommands(worker *Worker, commands []string) error {
	for _, cmd := range commands {
		err := man.runCommand(worker, cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

// check panics if err is not nil
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
