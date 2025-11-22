package controller

import (
	"runtime"

	"github.com/singhpranshu/cointracker/client"
	"github.com/singhpranshu/cointracker/config"
)

type Worker chan struct{}

func (w *Worker) Acquire() {
	*w <- struct{}{}
}
func (w *Worker) Release() {
	<-*w
}

type Handler struct {
	Client            client.Client
	Config            *config.Config
	Worker            Worker
	CommandModeSignal chan struct{}
	isCommandMode     bool
}

func NewHandler(config *config.Config, client client.Client, isCommandMode bool) *Handler {
	numCores := runtime.NumCPU()
	return &Handler{
		Config:            config,
		Client:            client,
		Worker:            make(chan struct{}, numCores), // limit to 5 concurrent workers
		CommandModeSignal: make(chan struct{}, 3),
		isCommandMode:     isCommandMode,
	}
}
