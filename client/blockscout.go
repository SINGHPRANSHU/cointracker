package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/singhpranshu/cointracker/dto"
	"github.com/singhpranshu/cointracker/model"
)

type BlockScountClient struct {
	BaseURL               string
	ExternalTxnEndpoint   string
	InternalTxnEndpoint   string
	TokenTransferEndpoint string
}

func NewBlockScountClient(baseUrl string, externalTxnEndpoint string, InternalTxnEndpoint string, TokenTransferEndpoint string) *BlockScountClient {
	return &BlockScountClient{
		BaseURL:               baseUrl,
		ExternalTxnEndpoint:   "transactions",
		InternalTxnEndpoint:   "internal-transactions",
		TokenTransferEndpoint: "token-transfers",
	}
}

// FetchInBatches fetches paginated results from Blockscout API in batches
func Fetch[T any](baseURL, address string, endpoint string, page string) (res []T, nextpage string, err error) {
	var allItems []T
	url := fmt.Sprintf("%s/addresses/%s/%s?", baseURL, address, endpoint)
	url += page

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != http.StatusOK {
		return []T{}, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var apiResp dto.BlockscoutAPIResponse[T]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, "", err
	}

	allItems = append(allItems, apiResp.Items...)

	// Build next page URL
	if apiResp.NextPageParams == nil || len(apiResp.NextPageParams) == 0 {
		return allItems, "", nil
	}

	url = fmt.Sprintf("%s/addresses/%s/%s?", baseURL, address, endpoint)
	for k, v := range apiResp.NextPageParams {
		nextpage = fmt.Sprintf("%s=%v", k, v)
		url += fmt.Sprintf("%s", nextpage)
	}

	return allItems, nextpage, nil
}

func (blockScountClient BlockScountClient) FetchExternalTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error) {
	txn, nextPage, err := Fetch[dto.Transaction](blockScountClient.BaseURL, address, blockScountClient.ExternalTxnEndpoint, page)
	if err != nil {
		return nil, "", err
	}
	var records []model.TransactionRecord
	for _, t := range txn {
		records = append(records, ConvertEthTx(t))
	}
	return records, nextPage, nil
}
func (blockScountClient BlockScountClient) FetchtokenTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error) {
	txn, nextPage, err := Fetch[dto.TokenTransfer](blockScountClient.BaseURL, address, blockScountClient.ExternalTxnEndpoint, page)
	if err != nil {
		return nil, "", err
	}
	var records []model.TransactionRecord
	for _, t := range txn {
		records = append(records, ConvertTokenTx(t))
	}
	return records, nextPage, nil
}
func (blockScountClient BlockScountClient) FetchInternalTransfer(address string, page string) (tx []model.TransactionRecord, nextPage string, err error) {
	txn, nextPage, err := Fetch[dto.InternalTransaction](blockScountClient.BaseURL, address, blockScountClient.InternalTxnEndpoint, page)
	if err != nil {
		return nil, "", err
	}
	var records []model.TransactionRecord
	for _, t := range txn {
		records = append(records, ConvertInternalTx(t))
	}
	return records, nextPage, nil
}

func ConvertEthTx(tx dto.Transaction) model.TransactionRecord {
	gasFee := ""
	if tx.Fee != nil {
		gasFee = tx.Fee.Value
	}
	return model.TransactionRecord{
		TxHash:    tx.Hash,
		Timestamp: tx.Timestamp,
		From:      tx.From.Hash,
		To:        tx.To.Hash,
		TxType:    "ETH",
		Contract:  "",
		Symbol:    "ETH",
		TokenID:   "",
		Amount:    tx.Value,
		GasFeeETH: gasFee,
	}
}

// Convert Token transfer to unified record
func ConvertTokenTx(tx dto.TokenTransfer) model.TransactionRecord {
	return model.TransactionRecord{
		TxHash:    tx.TransactionHash,
		Timestamp: tx.Timestamp,
		From:      tx.From.Hash,
		To:        tx.To.Hash,
		TxType:    tx.Token.Type,
		Contract:  tx.Token.AddressHash,
		Symbol:    tx.Token.Symbol,
		TokenID:   "",
		Amount:    tx.Total.Value,
		GasFeeETH: "",
	}
}

// Convert Internal transaction to unified record
func ConvertInternalTx(tx dto.InternalTransaction) model.TransactionRecord {
	return model.TransactionRecord{
		TxHash:    tx.TransactionHash,
		Timestamp: tx.Timestamp,
		From:      tx.From.Hash,
		To:        tx.To.Hash,
		TxType:    "INTERNAL",
		Contract:  "",
		Symbol:    "",
		TokenID:   "",
		Amount:    tx.Value,
		GasFeeETH: "",
	}
}
