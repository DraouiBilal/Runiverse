package queue

import (
	"fmt"
	"github.com/DraouiBilal/Runiverse/cri"
	"github.com/DraouiBilal/Runiverse/invoker"
	"log"
	"strings"
	"sync"
)

var jobID int = 0

type Job struct {
	ID      int
	Payload *cri.InvocationRequest
}

type Queue struct {
	jobs chan Job       // Channel to hold jobs
	WG   sync.WaitGroup // WaitGroup to track active workers
	stop chan struct{}  // Channel to signal workers to stop
}

func newQueue(size int) *Queue {
	return &Queue{
		jobs: make(chan Job, size),
		stop: make(chan struct{}), // Initially no signal to stop
	}
}

func (q *Queue) AddJob(payload *cri.InvocationRequest) {
	job := Job{
		ID:      jobID,
		Payload: payload,
	}

	jobID++

	q.jobs <- job
	log.Println("Running container with image " + job.Payload.Image + " and command [" + strings.Join(job.Payload.Command, ", ") + "]")
}

// Worker processes jobs from the queue
func (q *Queue) worker(workerID int) {
	defer q.WG.Done() // Ensure that the worker is done when exiting
	for {
		job, ok := <-q.jobs // Wait for a job
		if !ok { // If the channel is closed, exit the worker
			fmt.Printf("Worker %d stopping, no more jobs.\n", workerID)
			return
		}

		// Process the job
		fmt.Printf("Worker %d started job %d: %s\n", workerID, job.ID, job.Payload)
		invoker.Invoke(job.Payload.Image, job.Payload.Command) // Assuming invoker.Invoke is valid
		fmt.Printf("Worker %d finished job %d\n", workerID, job.ID)
	}
}

// StartWorkers starts a specified number of workers
func (q *Queue) startWorkers(numWorkers int) {
	for i := 1; i <= numWorkers; i++ {
		q.WG.Add(1)
		go q.worker(i)
	}
}

// StopWorkers signals all workers to stop
func (q *Queue) StopWorkers() {
	close(q.stop) // Send stop signal to workers
}

func InitQueue() *Queue {
	// Create a new queue with a buffer size of 10
	queue := newQueue(10)

	// Start 3 workers
	queue.startWorkers(3)

	return queue
}
