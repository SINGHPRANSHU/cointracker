package client

import "github.com/singhpranshu/cointracker/model"

// interface client where multiple client can implement like Blockscout, Etherscan etc
type Client interface {
	FetchExternalTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error)
	FetchtokenTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error)
	FetchInternalTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error)
}
