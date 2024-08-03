package client

import (
	"consumer-payment-service/environment"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_PaymentClient_MakeDeposit(t *testing.T) {
	const (
		success = iota
		requestError
		errorMakeDeposit
	)

	testCases := []struct {
		name      string
		inputArgs struct {
			accountId string
			reference string
			amount    float64
		}
		testType int
	}{
		{
			name: "Test success",
			inputArgs: struct {
				accountId string
				reference string
				amount    float64
			}{
				accountId: "account_id",
				reference: "reference",
				amount:    10.0,
			},
			testType: success,
		},

		{
			name: "Test error with request",
			inputArgs: struct {
				accountId string
				reference string
				amount    float64
			}{
				accountId: "account_id",
				reference: "reference",
				amount:    10.0,
			},
			testType: requestError,
		},

		{
			name: "Test error response when making deposit",
			inputArgs: struct {
				accountId string
				reference string
				amount    float64
			}{
				accountId: "account_id",
				reference: "reference",
				amount:    10.0,
			},
			testType: errorMakeDeposit,
		},
	}

	cfg := &environment.Config{
		THIRD_PARTY_SERVICE_BASE_URL: "http://example.domain.com/third-party",
	}

	httpClient := resty.New()

	paymentAPIClient := &paymentAPIClient{
		restClient: httpClient,
		config:     cfg,
	}

	httpmock.ActivateNonDefault(httpClient.GetClient())

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockUrl := fmt.Sprintf("%s/payments?type=credit", cfg.THIRD_PARTY_SERVICE_BASE_URL)
			switch testCase.testType {

			case success:
				mockPaymentResponse := PaymentResponse{
					AccountId: testCase.inputArgs.accountId,
					Reference: testCase.inputArgs.reference,
					Amount:    testCase.inputArgs.amount,
				}

				httpmock.RegisterResponder("POST", mockUrl, func(r *http.Request) (*http.Response, error) {
					response, err := httpmock.NewJsonResponse(http.StatusOK, mockPaymentResponse)
					if err != nil {
						return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
					}
					return response, nil
				})

				resp, err := paymentAPIClient.MakeDeposit(testCase.inputArgs.accountId, testCase.inputArgs.reference, testCase.inputArgs.amount)
				assert.NoError(t, err)
				assert.EqualValues(t, resp.AccountId, testCase.inputArgs.accountId)
				assert.EqualValues(t, resp.Amount, testCase.inputArgs.amount)
				assert.EqualValues(t, resp.Reference, testCase.inputArgs.reference)

			case requestError:
				httpmock.RegisterResponder("POST", mockUrl, httpmock.ConnectionFailure)
				resp, err := paymentAPIClient.MakeDeposit(testCase.inputArgs.accountId, testCase.inputArgs.reference, testCase.inputArgs.amount)
				assert.Error(t, err)
				assert.Nil(t, resp)

			case errorMakeDeposit:
				mockErrorResponse := ErrorResponse{
					ErrorMessage: "transaction failed",
				}

				httpmock.RegisterResponder("POST", mockUrl, func(r *http.Request) (*http.Response, error) {
					response, err := httpmock.NewJsonResponse(http.StatusInternalServerError, mockErrorResponse)
					if err != nil {
						return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
					}
					return response, nil
				})

				resp, err := paymentAPIClient.MakeDeposit(testCase.inputArgs.accountId, testCase.inputArgs.reference, testCase.inputArgs.amount)
				assert.Error(t, err)
				assert.Nil(t, resp)
			}

		})
	}

}

func Test_PaymentClient_MakeWithdrawal(t *testing.T) {
	const (
		success = iota
		requestError
		errorMakeWithdrawal
	)

	testCases := []struct {
		name      string
		inputArgs struct {
			accountId string
			reference string
			amount    float64
		}
		testType int
	}{
		{
			name: "Test success",
			inputArgs: struct {
				accountId string
				reference string
				amount    float64
			}{
				accountId: "account_id",
				reference: "reference",
				amount:    10.0,
			},
			testType: success,
		},

		{
			name: "Test error with request",
			inputArgs: struct {
				accountId string
				reference string
				amount    float64
			}{
				accountId: "account_id",
				reference: "reference",
				amount:    10.0,
			},
			testType: requestError,
		},

		{
			name: "Test error response when posting credit",
			inputArgs: struct {
				accountId string
				reference string
				amount    float64
			}{
				accountId: "account_id",
				reference: "reference",
				amount:    10.0,
			},
			testType: errorMakeWithdrawal,
		},
	}

	cfg := &environment.Config{
		THIRD_PARTY_SERVICE_BASE_URL: "http://example.domain.com/third-party",
	}

	httpClient := resty.New()

	paymentAPIClient := &paymentAPIClient{
		restClient: httpClient,
		config:     cfg,
	}

	httpmock.ActivateNonDefault(httpClient.GetClient())

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockUrl := fmt.Sprintf("%s/payments?type=debit", cfg.THIRD_PARTY_SERVICE_BASE_URL)
			switch testCase.testType {

			case success:
				mockPaymentResponse := PaymentResponse{
					AccountId: testCase.inputArgs.accountId,
					Reference: testCase.inputArgs.reference,
					Amount:    testCase.inputArgs.amount,
				}

				httpmock.RegisterResponder("POST", mockUrl, func(r *http.Request) (*http.Response, error) {
					response, err := httpmock.NewJsonResponse(http.StatusOK, mockPaymentResponse)
					if err != nil {
						return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
					}
					return response, nil
				})

				resp, err := paymentAPIClient.MakeWithdrawal(testCase.inputArgs.accountId, testCase.inputArgs.reference, testCase.inputArgs.amount)
				assert.NoError(t, err)
				assert.EqualValues(t, resp.AccountId, testCase.inputArgs.accountId)
				assert.EqualValues(t, resp.Amount, testCase.inputArgs.amount)
				assert.EqualValues(t, resp.Reference, testCase.inputArgs.reference)

			case requestError:
				httpmock.RegisterResponder("POST", mockUrl, httpmock.ConnectionFailure)
				resp, err := paymentAPIClient.MakeWithdrawal(testCase.inputArgs.accountId, testCase.inputArgs.reference, testCase.inputArgs.amount)
				assert.Error(t, err)
				assert.Nil(t, resp)

			case errorMakeWithdrawal:
				mockErrorResponse := ErrorResponse{
					ErrorMessage: "transaction failed",
				}

				httpmock.RegisterResponder("POST", mockUrl, func(r *http.Request) (*http.Response, error) {
					response, err := httpmock.NewJsonResponse(http.StatusInternalServerError, mockErrorResponse)
					if err != nil {
						return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
					}
					return response, nil
				})

				resp, err := paymentAPIClient.MakeWithdrawal(testCase.inputArgs.accountId, testCase.inputArgs.reference, testCase.inputArgs.amount)
				assert.Error(t, err)
				assert.Nil(t, resp)
			}

		})
	}

}

func Test_PaymentClient_RetrievePayment(t *testing.T) {
	const (
		success = iota
		requestError
		errorMakeWithdrawal
	)

	testCases := []struct {
		name      string
		reference string
		testType  int
	}{
		{
			name:      "Test success",
			reference: "ref-001",
			testType:  success,
		},

		{
			name:      "Test error with request",
			reference: "reference",
			testType:  requestError,
		},

		{
			name:      "Test error response with invalid reference",
			reference: "invalid_ref",
			testType:  errorMakeWithdrawal,
		},
	}

	cfg := &environment.Config{
		THIRD_PARTY_SERVICE_BASE_URL: "http://example.domain.com/third-party",
	}

	httpClient := resty.New()

	paymentAPIClient := &paymentAPIClient{
		restClient: httpClient,
		config:     cfg,
	}

	httpmock.ActivateNonDefault(httpClient.GetClient())

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockUrl := fmt.Sprintf("%s/payments/%s", cfg.THIRD_PARTY_SERVICE_BASE_URL, testCase.reference)
			switch testCase.testType {

			case success:
				mockPaymentResponse := PaymentResponse{
					AccountId: "account_id",
					Reference: testCase.reference,
					Amount:    10.0,
				}

				httpmock.RegisterResponder("GET", mockUrl, func(r *http.Request) (*http.Response, error) {
					response, err := httpmock.NewJsonResponse(http.StatusOK, mockPaymentResponse)
					if err != nil {
						return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
					}
					return response, nil
				})

				resp, err := paymentAPIClient.RetrieveTransaction(testCase.reference)
				assert.NoError(t, err)
				assert.Equal(t, testCase.reference, resp.Reference)

			case requestError:
				httpmock.RegisterResponder("GET", mockUrl, httpmock.ConnectionFailure)
				resp, err := paymentAPIClient.RetrieveTransaction(testCase.reference)
				assert.Error(t, err)
				assert.Nil(t, resp)

			case errorMakeWithdrawal:
				mockErrorResponse := ErrorResponse{
					ErrorMessage: "transaction not found",
				}

				httpmock.RegisterResponder("GET", mockUrl, func(r *http.Request) (*http.Response, error) {
					response, err := httpmock.NewJsonResponse(http.StatusInternalServerError, mockErrorResponse)
					if err != nil {
						return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
					}
					return response, nil
				})

				resp, err := paymentAPIClient.RetrieveTransaction(testCase.reference)
				assert.Error(t, err)
				assert.Nil(t, resp)
			}

		})
	}

}
