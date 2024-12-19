package client

import (
	"consumer-payment-service/environment"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

//go:generate mockgen -source=client.go -destination=../mocks/client_mock.go -package=mocks
type ThirdPartyAPIClient interface {
	MakeDeposit(accountId, reference string, amount float64) (*PaymentResponse, error)
	MakeWithdrawal(accountId, reference string, amount float64) (*PaymentResponse, error)
	RetrieveTransaction(reference string) (*PaymentResponse, error)
}

type paymentAPIClient struct {
	restClient *resty.Client
	config     *environment.Config
}

func NewPaymentAPIClient(config *environment.Config) ThirdPartyAPIClient {
	restClient := resty.New()
	restClient.SetDebug(true)
	return &paymentAPIClient{
		restClient: restClient,
		config:     config,
	}
}

func (p *paymentAPIClient) MakeDeposit(accountId, reference string, amount float64) (*PaymentResponse, error) {
	payload := PaymentRequest{
		AccountId: accountId,
		Reference: reference,
		Amount:    amount,
	}
	url := fmt.Sprintf("%s/payments?type=credit", p.config.THIRD_PARTY_SERVICE_BASE_URL)
	resp, err := p.restClient.
		R().
		SetResult(&PaymentResponse{}).
		SetError(&ErrorResponse{}).
		SetBody(payload).
		Post(url)
	if err != nil {
		log.Printf("error making deposit %v", err)
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("make deposit error occurred with httpCode: %d, message: %s ", resp.StatusCode(), resp.Body())
	}

	return resp.Result().(*PaymentResponse), nil
}

func (p *paymentAPIClient) MakeWithdrawal(accountId, reference string, amount float64) (*PaymentResponse, error) {
	payload := PaymentRequest{
		AccountId: accountId,
		Reference: reference,
		Amount:    amount,
	}
	url := fmt.Sprintf("%s/payments?type=debit", p.config.THIRD_PARTY_SERVICE_BASE_URL)
	resp, err := p.restClient.
		R().
		SetResult(&PaymentResponse{}).
		SetError(&ErrorResponse{}).
		SetBody(payload).
		Post(url)
	if err != nil {
		log.Printf("error making withdrawal %v", err)
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("make withdrawal error occurred with httpCode: %d, message: %s ", resp.StatusCode(), resp.Body())
	}

	return resp.Result().(*PaymentResponse), nil
}

func (p *paymentAPIClient) RetrieveTransaction(reference string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/payments/%s", p.config.THIRD_PARTY_SERVICE_BASE_URL, reference)
	resp, err := p.restClient.
		R().
		SetResult(&PaymentResponse{}).
		SetError(&ErrorResponse{}).
		Get(url)
	if err != nil {
		log.Printf("error retrieving payment %v", err)
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("retrieve payment error httpCode: %d, message: %s ", resp.StatusCode(), resp.Body())
	}

	return resp.Result().(*PaymentResponse), nil
}
