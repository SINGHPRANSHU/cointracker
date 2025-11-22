package model

type TransactionRecord struct {
	TxHash    string
	Timestamp string
	From      string
	To        string
	TxType    string
	Contract  string
	Symbol    string
	TokenID   string
	Amount    string
	GasFeeETH string
}

type TransferType string
