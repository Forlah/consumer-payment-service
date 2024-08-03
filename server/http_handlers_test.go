package server

import (
	"bytes"
	"consumer-payment-service/client"
	"consumer-payment-service/environment"
	"consumer-payment-service/mocks"
	"consumer-payment-service/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_HttpHandler_PaymentCredit(t *testing.T) {
	const (
		success = iota
		errorGettingUser
		errorGettingAccount
		errorMakingDeposit
		errorCreatingTransaction
		errorUpdatingAccountBalance
	)

	testCases := []struct {
		name     string
		testType int
	}{
		{
			name:     "Test success",
			testType: success,
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	cfg := &environment.Config{
		THIRD_PARTY_SERVICE_BASE_URL: "http://example.domain.com/third-party",
	}

	mockDataStore := mocks.NewMockMongoDBStore(controller)
	mockThirdPartyClient := mocks.NewMockThirdPartyAPIClient(controller)

	handler := NewHTTPHandler(cfg, mockDataStore, mockThirdPartyClient)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockRequest := models.PaymentRequestPayload{
				UserId:    "usr-001",
				AccountId: "acc_001",
				Reference: "ref-001",
				Amount:    10,
			}
			mockPayload, err := json.Marshal(mockRequest)
			assert.NoError(t, err)

			switch testCase.testType {
			case success:
				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   1,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/credit", bytes.NewBuffer(mockPayload))

				mockDataStore.
					EXPECT().
					GetUserById(mockRequest.UserId).
					Return(&models.User{
						Id: mockRequest.UserId,
					}, nil)

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&mockAccount, nil)

				mockThirdPartyClient.
					EXPECT().
					MakeDeposit(mockRequest.AccountId, mockRequest.Reference, mockRequest.Amount).
					Return(&client.PaymentResponse{
						AccountId: mockRequest.AccountId,
						Reference: mockRequest.Reference,
						Amount:    mockRequest.Amount,
					}, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(nil)

				newBalance := mockAccount.Balance + mockRequest.Amount

				mockDataStore.EXPECT().UpdateAccountBalance(mockRequest.AccountId, newBalance).Return(nil)

				handler.PaymentCreditHandler(w, r)
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}

}
