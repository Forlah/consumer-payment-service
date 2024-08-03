package server

import (
	"consumer-payment-service/client"
	"consumer-payment-service/database"
	"consumer-payment-service/environment"
	"consumer-payment-service/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type HttpHandler struct {
	config        *environment.Config
	mongodbStore  database.MongoDBStore
	paymentClient client.ThirdPartyAPIClient
}

func NewHTTPHandler(config *environment.Config, store database.MongoDBStore, paymentClient client.ThirdPartyAPIClient) *HttpHandler {
	return &HttpHandler{config: config, mongodbStore: store, paymentClient: paymentClient}
}

func (handler *HttpHandler) responseWriter(w http.ResponseWriter, codes ...int) {
	statusCode := http.StatusOK
	if len(codes) > 0 {
		statusCode = codes[0]
	}

	w.WriteHeader(statusCode)
}

func (handler *HttpHandler) PaymentCreditHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("Error closing request body")
		}
	}()

	var payload models.PaymentRequestPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate user exist
	if _, err = handler.mongodbStore.GetUserById(payload.UserId); err != nil {
		log.Printf("error getting user %v", err)
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	// validate account exist
	account, err := handler.mongodbStore.GetAccountByID(payload.AccountId)
	if err != nil {
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	// make credit API call to third party service
	resp, err := handler.paymentClient.MakeDeposit(payload.AccountId, payload.Reference, payload.Amount)
	if err != nil {
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	transaction := &models.Transaction{
		Reference: resp.Reference,
		UserID:    payload.UserId,
		AccountID: resp.AccountId,
		Amount:    resp.Amount,
		Type:      models.CREDIT,
		Status:    models.SUCCESS,
		CreatedAt: time.Now().Unix(),
	}

	err = handler.mongodbStore.CreateTransaction(transaction)
	if err != nil {
		log.Println("error creating transaction")
		handler.responseWriter(w, http.StatusNotFound)
		return
	}

	newBalance := account.Balance + payload.Amount

	if err = handler.mongodbStore.UpdateAccountBalance(payload.AccountId, newBalance); err != nil {
		handler.responseWriter(w, http.StatusNotFound)
		return
	}

	handler.responseWriter(w)
}

func (handler *HttpHandler) PaymentDebitHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("Error closing request body")
		}
	}()

	var payload models.PaymentRequestPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate user exist
	if _, err = handler.mongodbStore.GetUserById(payload.UserId); err != nil {
		log.Printf("error getting user %v", err)
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	// validate account exist
	account, err := handler.mongodbStore.GetAccountByID(payload.AccountId)
	if err != nil {
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	// check balance
	if payload.Amount > float64(account.Balance) {
		log.Println("insufficient balance")
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	// make API call to third party payment service for debit
	resp, err := handler.paymentClient.MakeWithdrawal(payload.AccountId, payload.Reference, payload.Amount)
	if err != nil {
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	newBalance := account.Balance - payload.Amount

	transaction := &models.Transaction{
		Reference: resp.Reference,
		UserID:    payload.UserId,
		AccountID: resp.AccountId,
		Amount:    resp.Amount,
		Type:      models.DEBIT,
		Status:    models.SUCCESS,
		CreatedAt: time.Now().Unix(),
	}

	err = handler.mongodbStore.CreateTransaction(transaction)
	if err != nil {
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	if err = handler.mongodbStore.UpdateAccountBalance(payload.AccountId, newBalance); err != nil {
		handler.responseWriter(w, http.StatusInternalServerError)
		return
	}

	handler.responseWriter(w)
}
