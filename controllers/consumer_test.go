package controller

import (
	"testing"

	"github.com/singhpranshu/cointracker/queue"
)

type MockWorker struct{}

func (w MockWorker) Acquire() {

}
func (w MockWorker) Release() {

}

func TestFetchAndSaveAllTransactions(t *testing.T) {
	mockClient := &MockClient{}
	Worker := MockWorker{}
	handler := &Handler{Client: mockClient, Worker: Worker}

	// Test case: Valid address
	err := handler.FetchAndSaveAllTransactions(queue.Job{Address: "valid_address", Type: "external", Page: "1"})
	if err != nil {
		t.Errorf("Expected no error for valid address, got %v", err)
	}

	// Test case: Invalid address
	err = handler.FetchAndSaveAllTransactions(queue.Job{Address: "invalid", Type: "", Page: "1"})
	if err == nil {
		t.Errorf("Expected an error for invalid address, got nil")
	}
}
