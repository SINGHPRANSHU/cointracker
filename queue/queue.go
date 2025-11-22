package queue

import (
	"sync"

	"github.com/singhpranshu/cointracker/model"
)

var mutex = sync.Mutex{}

type Job struct {
	Address string
	Type    model.TransferType
	Page    string
}

type JobQueue []Job

func NewJobQueue() JobQueue {
	return []Job{}
}

func (jq *JobQueue) Enqueue(job Job) {
	mutex.Lock()
	defer mutex.Unlock()
	*jq = append(*jq, job)
}

func (jq *JobQueue) Dequeue() (Job, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	if len(*jq) == 0 {
		return Job{}, false
	}
	job := (*jq)[0]
	*jq = (*jq)[1:]
	return job, true
}
