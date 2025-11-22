package controller

import (
	"runtime"

	"github.com/singhpranshu/cointracker/client"
	"github.com/singhpranshu/cointracker/config"
)

type Worker chan struct{}

func (w Worker) Acquire() {
	w <- struct{}{}
}
func (w Worker) Release() {
	<-w
}

type IWorker interface {
	Acquire()
	Release()
}

type Handler struct {
	Client            client.Client
	Config            *config.Config
	Worker            IWorker
	CommandModeSignal chan struct{}
	isCommandMode     bool
}

func NewHandler(config *config.Config, client client.Client, isCommandMode bool) *Handler {
	numCores := runtime.NumCPU()
	// limit to 5 concurrent workers
	worker := make(chan struct{}, numCores)
	return &Handler{
		Config:            config,
		Client:            client,
		Worker:            Worker(worker),
		CommandModeSignal: make(chan struct{}, 3),
		isCommandMode:     isCommandMode,
	}
}
