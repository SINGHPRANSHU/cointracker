package client

import (
	"sync"

	"github.com/singhpranshu/cointracker/model"
)

// interface client where multiple client can implement like Blockscout, Etherscan etc
type Client interface {
	FetchExternalTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error)
	FetchtokenTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error)
	FetchInternalTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error)
}


// for testing if we are not processing multiple files
var mutex = sync.Mutex{}

type NewMap map[string]string

func (nm NewMap) Get(key string) (string, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	v, ok := nm[key]
	return v, ok
}
func (nm NewMap) Put(key string, v string) {
	mutex.Lock()
	defer mutex.Unlock()
	nm[key] = v
}
