package csv

import (
	"encoding/csv"
	"os"
	"sync"

	"github.com/singhpranshu/cointracker/model"
)

type CSVProcessor struct {
	FileName string
	mutex    *sync.Mutex
	file     *os.File
}

var CSVProcessorMap map[string]*CSVProcessor = make(map[string]*CSVProcessor)

// GetCSVProcessor is thread safe and returns a singleton CSVProcessor for the given address and transaction type
func GetCSVProcessor(address string, txnType model.TransferType) (*CSVProcessor, error) {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	filename := GetFileName(address, txnType)
	_, exist := CSVProcessorMap[filename]
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	if !exist {
		CSVProcessorMap[filename] = &CSVProcessor{
			FileName: filename,
			mutex:    &sync.Mutex{},
			file:     file,
		}
	}
	return CSVProcessorMap[filename], nil
}
func GetFileName(address string, txnType model.TransferType) string {
	return address + "_" + string(txnType) + ".csv"
}

// WriteCSV is thread safe and only one writes will happen at a time writes transaction records to a CSV file
func (csvProcessor *CSVProcessor) WriteCSV(records []model.TransactionRecord) error {
	csvProcessor.mutex.Lock()
	defer csvProcessor.mutex.Unlock()
	writer := csv.NewWriter(csvProcessor.file)
	defer writer.Flush()

	// CSV Header
	header := []string{"TxHash", "Timestamp", "From", "To", "TxType", "Contract", "Symbol", "TokenID", "Amount", "GasFeeETH"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range records {
		row := []string{
			r.TxHash, r.Timestamp, r.From, r.To, r.TxType,
			r.Contract, r.Symbol, r.TokenID, r.Amount, r.GasFeeETH,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
