package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	neturl "net/url"

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
var m NewMap = make(map[string]string)

func Fetch[T any](baseURL, address string, endpoint string, page string) (res []T, nextpage string, err error) {
	// fmt.Println("Fetching page:", page)

	var allItems []T
	url := fmt.Sprintf("%s/addresses/%s/%s?", baseURL, address, endpoint)
	url += page
	// Check for duplicate requests
	// use this in dev as it may reduce performance and parallelism due to mutex
	// if val, ok := m.Get(url); ok {
	// 	panic("Duplicate url request detected: " + val + " and " + page)
	// }
	// m.Put(url, url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// fmt.Println("Error response status:", string(body))
		return []T{}, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp dto.BlockscoutAPIResponse[T]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, "", err
	}

	allItems = append(allItems, apiResp.Items...)

	// Build next page URL
	if apiResp.NextPageParams == nil || len(apiResp.NextPageParams) == 0 {
		return allItems, "", nil
	}

	queryParmas := neturl.Values{}
	for k, v := range apiResp.NextPageParams {
		switch val := v.(type) {
		case int:
			queryParmas.Add(k, fmt.Sprintf("%d", val))
		case string:
			queryParmas.Add(k, val)
		case float64:
			queryParmas.Add(k, fmt.Sprintf("%.0f", val))
		default:
			fmt.Println("Unknown type for next page param:", k, v, val)
			panic("unknown type in next page params")
		}
	}

	nextpage = queryParmas.Encode()
	// fmt.Println("Next page params:", nextpage)

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
	txn, nextPage, err := Fetch[dto.TokenTransfer](blockScountClient.BaseURL, address, blockScountClient.TokenTransferEndpoint, page)
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
