package dto

import "github.com/singhpranshu/cointracker/model"

type BlockscoutAPIResponse[T any] struct {
	Items          []T                    `json:"items"`
	NextPageParams map[string]interface{} `json:"next_page_params"`
}

type NextPageParams struct {
	Index            uint64 `json:"index"`
	Value            string `json:"value"`
	Hash             string `json:"hash"`
	InsertedAt       string `json:"inserted_at"`
	BlockNumber      uint64 `json:"block_number"`
	Fee              string `json:"fee"`
	ItemsCount       uint64 `json:"items_count"`
	TransactionIndex uint64 `json:"transaction_index"`
}

// "next_page_params": {
//     "index": 55,
//     "value": "0",
//     "hash": "0xbcb1d9f4c85a708b0eb37b50dd760048b185f20c806c3550c7ca101e09afe490",
//     "inserted_at": "2025-06-14T22:51:15.376529Z",
//     "block_number": 22705890,
//     "fee": "62653428273671",
//     "items_count": 50
//"transaction_index": 0
//   }

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
