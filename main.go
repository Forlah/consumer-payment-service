package main

import (
	"consumer-payment-service/client"
	"consumer-payment-service/database/mongodb"
	"consumer-payment-service/environment"
	"fmt"
	"log"
	"net/http"

	srv "consumer-payment-service/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load content of .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env ")
	}

	cfg := environment.LoadConfig()

	// Get mongodb instance
	store, _, err := mongodb.New(cfg.DatabaseURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("failed to establish MongoDB connection")
	}

	// Get instance of third party payment service client
	paymentClient := client.NewPaymentAPIClient(cfg)

	addr := fmt.Sprintf(":%s", cfg.PORT)
	router := srv.MountServer(cfg, store, paymentClient)
	// start HTTP server
	fmt.Println(fmt.Sprintf("starting HTTP service running on port %v", addr))
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("error starting http server")
	}
}
