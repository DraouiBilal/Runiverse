package queue

import (
	"fmt"
	"sync"
	"time"
)

// Job represents a unit of work
type Job struct {
	ID      int
	Payload string
}

// Queue manages jobs and workers
type Queue struct {
	jobs         chan Job       // Channel to hold jobs
	wg           sync.WaitGroup // WaitGroup to track active workers
	workersIdle  int            // Number of idle workers
	totalWorkers int            // Total number of workers
	mu           sync.Mutex     // Mutex to protect workersIdle
}

// NewQueue creates a new queue with a specified buffer size
func NewQueue(size int) *Queue {
	return &Queue{
		jobs: make(chan Job, size),
	}
}

// AddJob adds a job to the queue
func (q *Queue) AddJob(job Job) {
	q.jobs <- job
	fmt.Printf("Added job %d to queue\n", job.ID)
}

// Worker processes jobs from the queue
func (q *Queue) Worker(workerID int) {
	defer q.wg.Done() // Signal when this worker is done

	for {
		// Mark worker as idle before waiting for a job
		q.mu.Lock()
		q.workersIdle++
		q.mu.Unlock()

		job, ok := <-q.jobs // Wait for a job
		if !ok {            // Channel closed, exit
			return
		}

		// Mark worker as busy
		q.mu.Lock()
		q.workersIdle--
		q.mu.Unlock()

		fmt.Printf("Worker %d started job %d: %s\n", workerID, job.ID, job.Payload)
		time.Sleep(1 * time.Second) // Simulate work
		fmt.Printf("Worker %d finished job %d\n", workerID, job.ID)
	}
}

// StartWorkers starts a specified number of workers
func (q *Queue) StartWorkers(numWorkers int) {
	q.totalWorkers = numWorkers
	q.workersIdle = numWorkers // Initially, all workers are idle
	for i := 1; i <= numWorkers; i++ {
		q.wg.Add(1)
		go q.Worker(i)
	}
}

// Status returns the current queue status
func (q *Queue) Status() (jobsInQueue int, workersBusy int, isActive bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	jobsInQueue = len(q.jobs)                     // Number of jobs waiting in the channel
	workersBusy = q.totalWorkers - q.workersIdle  // Number of workers currently processing jobs
	isActive = jobsInQueue > 0 || workersBusy > 0 // Queue is active if there are jobs or busy workers
	return
}

// Close closes the queue
func (q *Queue) Close() {
	close(q.jobs)
	q.wg.Wait() // Wait for all workers to finish
}

func main() {
	// Create a new queue with a buffer size of 10
	queue := NewQueue(10)

	// Start 3 workers
	queue.StartWorkers(3)

	// Add some jobs
	for i := 1; i <= 5; i++ {
		job := Job{
			ID:      i,
			Payload: fmt.Sprintf("Job payload %d", i),
		}
		queue.AddJob(job)
	}

	// Monitor queue status periodically
	for i := 0; i < 10; i++ {
		jobsInQueue, workersBusy, isActive := queue.Status()
		fmt.Printf("Status: %d jobs in queue, %d workers busy, active: %v\n", jobsInQueue, workersBusy, isActive)
		time.Sleep(1 * time.Second)
	}

	// Close the queue and wait for workers to finish
	queue.Close()
	fmt.Println("Queue closed, all jobs completed")
}

