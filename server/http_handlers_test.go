package server

import (
	"bytes"
	"consumer-payment-service/client"
	"consumer-payment-service/environment"
	"consumer-payment-service/mocks"
	"consumer-payment-service/models"
	"encoding/json"
	"errors"
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

		{
			name:     "Test error fetching user",
			testType: errorGettingUser,
		},

		{
			name:     "Test error fetching account",
			testType: errorGettingAccount,
		},

		{
			name:     "Test error making deposit on third party service",
			testType: errorMakingDeposit,
		},

		{
			name:     "Test error creating transaction record",
			testType: errorCreatingTransaction,
		},

		{
			name:     "Test error while updating account balance",
			testType: errorUpdatingAccountBalance,
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

			case errorGettingUser:
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/credit", bytes.NewBuffer(mockPayload))

				mockDataStore.
					EXPECT().
					GetUserById(mockRequest.UserId).
					Return(nil, errors.New("not found"))

				handler.PaymentCreditHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorGettingAccount:
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
					Return(nil, errors.New("not found"))

				handler.PaymentCreditHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorMakingDeposit:

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
					Return(&models.Account{
						AccountID: mockRequest.AccountId,
						Balance:   1,
					}, nil)

				mockThirdPartyClient.
					EXPECT().
					MakeDeposit(mockRequest.AccountId, mockRequest.Reference, mockRequest.Amount).
					Return(nil, errors.New(""))

				handler.PaymentCreditHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorCreatingTransaction:
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
					Return(errors.New(""))

				handler.PaymentCreditHandler(w, r)
				assert.Equal(t, http.StatusNotFound, w.Code)

			case errorUpdatingAccountBalance:

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

				mockDataStore.
					EXPECT().
					UpdateAccountBalance(mockRequest.AccountId, newBalance).
					Return(errors.New(""))

				handler.PaymentCreditHandler(w, r)
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})
	}

}

func Test_HttpHandler_PaymentDebit(t *testing.T) {
	const (
		success = iota
		errorGettingUser
		errorGettingAccount
		errorInsufficientBalance
		errorMakingWithdrawal
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

		{
			name:     "Test error fetching user",
			testType: errorGettingUser,
		},

		{
			name:     "Test error fetching account",
			testType: errorGettingAccount,
		},

		{
			name:     "Test error insufficient balance",
			testType: errorInsufficientBalance,
		},

		{
			name:     "Test error making withdrawal on third party service",
			testType: errorMakingWithdrawal,
		},

		{
			name:     "Test error creating transaction record",
			testType: errorCreatingTransaction,
		},

		{
			name:     "Test error while updating account balance",
			testType: errorUpdatingAccountBalance,
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
				Amount:    1.50,
			}
			mockPayload, err := json.Marshal(mockRequest)
			assert.NoError(t, err)

			switch testCase.testType {
			case success:
				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/debit", bytes.NewBuffer(mockPayload))

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
					MakeWithdrawal(mockRequest.AccountId, mockRequest.Reference, mockRequest.Amount).
					Return(&client.PaymentResponse{
						AccountId: mockRequest.AccountId,
						Reference: mockRequest.Reference,
						Amount:    mockRequest.Amount,
					}, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(nil)

				newBalance := mockAccount.Balance - mockRequest.Amount

				mockDataStore.EXPECT().UpdateAccountBalance(mockRequest.AccountId, newBalance).Return(nil)

				handler.PaymentDebitHandler(w, r)
				assert.Equal(t, http.StatusOK, w.Code)

			case errorGettingUser:
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/debit", bytes.NewBuffer(mockPayload))

				mockDataStore.
					EXPECT().
					GetUserById(mockRequest.UserId).
					Return(nil, errors.New(""))

				handler.PaymentDebitHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorGettingAccount:
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/debit", bytes.NewBuffer(mockPayload))

				mockDataStore.
					EXPECT().
					GetUserById(mockRequest.UserId).
					Return(&models.User{
						Id: mockRequest.UserId,
					}, nil)

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(nil, errors.New(""))

				handler.PaymentDebitHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorInsufficientBalance:

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/debit", bytes.NewBuffer(mockPayload))

				mockDataStore.
					EXPECT().
					GetUserById(mockRequest.UserId).
					Return(&models.User{
						Id: mockRequest.UserId,
					}, nil)

				mockDataStore.
					EXPECT().
					GetAccountByID(mockRequest.AccountId).
					Return(&models.Account{
						AccountID: mockRequest.AccountId,
						Balance:   0.50,
					}, nil)

				handler.PaymentDebitHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorMakingWithdrawal:
				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/debit", bytes.NewBuffer(mockPayload))

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
					MakeWithdrawal(mockRequest.AccountId, mockRequest.Reference, mockRequest.Amount).
					Return(nil, errors.New(""))

				handler.PaymentDebitHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorCreatingTransaction:
				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/debit", bytes.NewBuffer(mockPayload))

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
					MakeWithdrawal(mockRequest.AccountId, mockRequest.Reference, mockRequest.Amount).
					Return(&client.PaymentResponse{
						AccountId: mockRequest.AccountId,
						Reference: mockRequest.Reference,
						Amount:    mockRequest.Amount,
					}, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(errors.New(""))

				handler.PaymentDebitHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)

			case errorUpdatingAccountBalance:
				mockAccount := models.Account{
					AccountID: mockRequest.AccountId,
					Balance:   10,
				}

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/payments/debit", bytes.NewBuffer(mockPayload))

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
					MakeWithdrawal(mockRequest.AccountId, mockRequest.Reference, mockRequest.Amount).
					Return(&client.PaymentResponse{
						AccountId: mockRequest.AccountId,
						Reference: mockRequest.Reference,
						Amount:    mockRequest.Amount,
					}, nil)

				mockDataStore.
					EXPECT().
					CreateTransaction(gomock.Any()).
					Return(nil)

				newBalance := mockAccount.Balance - mockRequest.Amount

				mockDataStore.
					EXPECT().
					UpdateAccountBalance(mockRequest.AccountId, newBalance).
					Return(errors.New(""))

				handler.PaymentDebitHandler(w, r)
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			}
		})
	}
}
