package server

import (
	"consumer-payment-service/client"
	"consumer-payment-service/database"
	"consumer-payment-service/environment"
	"net/http"

	"github.com/go-chi/chi"
)

func MountServer(cfg *environment.Config, mongodbStore database.MongoDBStore, paymentClient client.ThirdPartyAPIClient) *chi.Mux {
	router := chi.NewRouter()

	httpHandler := NewHTTPHandler(cfg, mongodbStore, paymentClient)

	// service check
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("a simple banking app service"))
	})

	router.Post("/payments/debit", httpHandler.PaymentDebitHandler)

	router.Post("/payments/credit", httpHandler.PaymentCreditHandler)

	return router
}
