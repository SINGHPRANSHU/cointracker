package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/singhpranshu/cointracker/client"
	"github.com/singhpranshu/cointracker/dto"
	"github.com/singhpranshu/cointracker/model"
	"github.com/singhpranshu/cointracker/queue"
)

var jobQueue = queue.NewJobQueue()

const (
	ExternalTransfer = "external"
	InternalTransfer = "internal"
	TokenTransfer    = "token"
)

// route to get transaction history for an address based on transfer type and fetches only 50 rows per page
// /{address}/history?type={external|internal|token}&page={page_number}
func (h *Handler) GetHistoryForAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	page := r.URL.Query().Get("page")
	txnType := r.URL.Query().Get("type")
	tnx, nextpage, err := h.FetchTxnData(model.TransferType(txnType), address, page, h.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(tnx)

	response := dto.GetHistoryApiResponse{}
	response.Transactions = tnx
	response.NextPage = nextpage
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

// route to add address in queue to process full history
// /{address}/history [POST]
func (h *Handler) AddAddressInQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	h.ProcessHistoryForAddress(address)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (h *Handler) FetchTxnData(TxnType model.TransferType, address string, page string, client client.Client) (tx []model.TransactionRecord, nextPage string, err error) {
	switch TxnType {
	case ExternalTransfer:
		return client.FetchExternalTransfer(address, page)
	case InternalTransfer:
		return client.FetchInternalTransfer(address, page)
	case TokenTransfer:
		return client.FetchtokenTransfer(address, page)
	}
	return []model.TransactionRecord{}, "", fmt.Errorf("invalid transaction type")
}

// ProcessHistoryForAddress processes the transaction history for a given address
// generates csv for all transactions
func (h *Handler) ProcessHistoryForAddress(address string) error {
	// signal completion of each goroutine
	types := []model.TransferType{ExternalTransfer, InternalTransfer, TokenTransfer}
	for _, txnType := range types {
		fmt.Println("Starting type:", txnType)
		// run in parallel to use multiple cores
		go h.EnqueueStart(address, txnType)
	}
	return nil
}

func (h *Handler) EnqueueStart(address string, txnType model.TransferType) {
	job := queue.Job{
		Address: address,
		Type:    txnType,
		Page:    "",
	}
	fmt.Println("Enqueuing Job:", job)
	jobQueue.Enqueue(job)

}
