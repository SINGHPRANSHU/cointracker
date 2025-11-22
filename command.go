package main

import (
	"github.com/singhpranshu/cointracker/client"
	"github.com/singhpranshu/cointracker/config"
	controller "github.com/singhpranshu/cointracker/controllers"
)

func command(address string) {
	config := config.NewConfig()
	config.Load()
	client := client.NewBlockScountClient(config.BlockScountClient.BaseURL, config.BlockScountClient.ExternalTxnEndpoint, config.BlockScountClient.InternalTxnEndpoint, config.BlockScountClient.TokenTransferEndpoint)
	h := controller.NewHandler(config, client, true)
	go h.ProcessHistoryForAddress(address)
	go h.StartWorkerPool()
	for i := 0; i < 3; i++ {
		<-h.CommandModeSignal
	}
}
