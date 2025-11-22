package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	serverMode := flag.Bool("servermode", false, "Run in server mode")
	address := flag.String("address", "", "Run in server mode")
	flag.Parse()
	if serverMode != nil && *serverMode {
		mux := router()
		fmt.Println("Server starting at port", os.Getenv("APP_PORT"))
		err := http.ListenAndServe(os.Getenv("APP_PORT"), mux)
		if err != nil {
			log.Fatal(err)
		}
	} else if address != nil && *address != "" {
		fmt.Println("Fetching history for address:", *address)
		command(*address)
		// Here you can add code to fetch and display history for the given address
	} else {
		fmt.Println("Please provide either -servermode to run in server mode or -address to fetch history for a specific address.")

	}

}
