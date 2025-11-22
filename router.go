package main

import (
	"github.com/gorilla/mux"
	"github.com/singhpranshu/cointracker/client"
	"github.com/singhpranshu/cointracker/config"
	controller "github.com/singhpranshu/cointracker/controllers"
)

func router() *mux.Router {
	config := config.NewConfig()
	config.Load()
	client := client.NewBlockScountClient(config.BlockScountClient.BaseURL, config.BlockScountClient.ExternalTxnEndpoint, config.BlockScountClient.InternalTxnEndpoint, config.BlockScountClient.TokenTransferEndpoint)
	h := controller.NewHandler(config, client, false)
	go h.StartWorkerPool()
	mux := mux.NewRouter()
	mux.HandleFunc("/{address}/history", h.GetHistoryForAddress).Methods("GET")
	mux.HandleFunc("/{address}/history", h.GetHistoryForAddress).Methods("POST")
	return mux
}
