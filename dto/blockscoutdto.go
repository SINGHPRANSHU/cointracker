package dto

import "github.com/singhpranshu/cointracker/model"

type BlockscoutAPIResponse[T any] struct {
	Items          []T                    `json:"items"`
	NextPageParams map[string]interface{} `json:"next_page_params"`
}

// Blockscout ETH transaction
type Transaction struct {
	Hash      string `json:"hash"`
	Timestamp string `json:"timestamp"`
	From      struct {
		Hash string `json:"hash"`
	} `json:"from"`
	To struct {
		Hash string `json:"hash"`
	} `json:"to"`
	Value string `json:"value"`
	Fee   *struct {
		Value string `json:"value"`
	} `json:"fee"`
}

// Blockscout Token transfer
type TokenTransfer struct {
	TransactionHash string `json:"transaction_hash"`
	Timestamp       string `json:"timestamp"`
	From            struct {
		Hash string `json:"hash"`
	} `json:"from"`
	To struct {
		Hash string `json:"hash"`
	} `json:"to"`
	Token struct {
		AddressHash string `json:"address_hash"`
		Symbol      string `json:"symbol"`
		Name        string `json:"name"`
		Type        string `json:"type"`
	} `json:"token"`
	Total struct {
		Value string `json:"value"`
	} `json:"total"`
}

// Blockscout Internal transaction
type InternalTransaction struct {
	TransactionHash string `json:"transaction_hash"`
	Timestamp       string `json:"timestamp"`
	From            struct {
		Hash string `json:"hash"`
	} `json:"from"`
	To struct {
		Hash string `json:"hash"`
	} `json:"to"`
	Value string `json:"value"`
}

type GetHistoryApiResponse struct {
	Transactions []model.TransactionRecord `json:"transactions"`
	NextPage     string                    `json:"nextpage"`
}
