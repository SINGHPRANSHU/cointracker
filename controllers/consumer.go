package controller

import (
	"fmt"
	"time"

	"github.com/singhpranshu/cointracker/csv"
	"github.com/singhpranshu/cointracker/queue"
)

func (h *Handler) StartWorkerPool() {
	for {
		job, exist := jobQueue.Dequeue()
		if !exist {
			// if no job wait for some time to poll again
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println("Dequeued Job:", job, exist)
		go h.FetchAndSaveAllTransactions(job)
	}

}

// FetchAllTransactions fetches and saves all transactions of a given type for the address
func (h *Handler) FetchAndSaveAllTransactions(job queue.Job) error {

	txnType := job.Type
	address := job.Address
	page := job.Page

	h.Worker.Acquire()
	txs, nextPage, err := h.FetchTxnData(txnType, address, page, h.Client)
	if err != nil {
		fmt.Println("Error fetching transactions:", err)
		jobQueue.Enqueue(job)
		return err
	}
	csvprocessor, err := csv.GetCSVProcessor(address, txnType)
	if err != nil {
		fmt.Println("Error getting CSV processor:", err)
		jobQueue.Enqueue(job)
		return err
	}
	err = csvprocessor.WriteCSV(txs)
	if err != nil {
		fmt.Println("Error writing CSV:", err)
		jobQueue.Enqueue(job)
		return err
	}
	if nextPage != "" {
		jobQueue.Enqueue(queue.Job{
			Address: address,
			Type:    txnType,
			Page:    nextPage,
		})
	} else if h.isCommandMode {
		h.CommandModeSignal <- struct{}{}
	}
	h.Worker.Release()

	// releasing for the last acquire
	return nil
}
