package controller

import (
	"errors"
	"testing"

	"github.com/singhpranshu/cointracker/model"
)

// Mock client to simulate API responses
type MockClient struct{}

func (m *MockClient) FetchExternalTransfer(address, page string) ([]model.TransactionRecord, string, error) {
	if address == "valid_address" {
		return []model.TransactionRecord{
			{TokenID: "txn1", Amount: "100"},
			{TokenID: "txn2", Amount: "200"},
		}, "2", nil
	}
	return nil, "", errors.New("invalid address")
}

func (m *MockClient) FetchInternalTransfer(address, page string) ([]model.TransactionRecord, string, error) {
	return nil, "", nil
}

func (m *MockClient) FetchtokenTransfer(address, page string) ([]model.TransactionRecord, string, error) {
	return nil, "", nil
}

func TestFetchTxnData(t *testing.T) {
	mockClient := &MockClient{}
	handler := &Handler{Client: mockClient}

	// Test case: Valid external transfer
	txnType := model.TransferType("external")
	address := "valid_address"
	page := "1"

	transactions, nextPage, err := handler.FetchTxnData(txnType, address, page, mockClient)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(transactions))
	}

	if nextPage != "2" {
		t.Errorf("Expected nextPage to be '2', got %s", nextPage)
	}

	// Test case: Invalid address
	_, _, err = handler.FetchTxnData(txnType, "invalid_address", page, mockClient)
	if err == nil {
		t.Errorf("Expected an error for invalid address, got nil")
	}
}

func TestProcessHistoryForAddress(t *testing.T) {
	mockClient := &MockClient{}
	handler := &Handler{Client: mockClient}

	// Test case: Process history for a valid address
	address := "valid_address"
	err := handler.ProcessHistoryForAddress(address)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
