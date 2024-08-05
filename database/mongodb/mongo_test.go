package mongodb

import (
	"consumer-payment-service/models"
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"math/rand"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

const (
	databaseName = "banking-app"
)

var mongoDbPort = ""

// func makeRandomString() string {
// 	rand.Seed(time.Now().Unix())
// 	length := 4

// 	ran_str := make([]byte, length)

// 	// Generating Random string
// 	for i := 0; i < length; i++ {
// 		ran_str[i] = (65 + rand.Intn(25))
// 	}

// 	// Displaying the random string
// 	return string(ran_str)
// }

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}

	resource, err := pool.Run("mongo", "6.0.6", []string{
		"MONGO_INITDB_DATABASE=" + databaseName,
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	mongoDbPort = resource.GetPort("27017/tcp")
	if err := pool.Retry(func() error {
		var err error
		connectURL := fmt.Sprintf("mongodb://localhost:%s", mongoDbPort)

		_, _, err = New(connectURL, databaseName)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	code := m.Run()
	err = pool.Purge(resource)
	if err != nil {
		log.Fatal(err)
	}

	rand.NewSource(time.Now().UnixNano())
	os.Exit(code)
}

func TestMongoStore_GetAccountByID(t *testing.T) {
	const (
		success = iota
		errorGetAccount
	)

	var tests = []struct {
		name      string
		accountId string
		testType  int
	}{
		{
			name:      "Test get account successfully",
			accountId: "accountId",
			testType:  success,
		},
		{
			name:      "Test error get account",
			accountId: "invalid_id",
			testType:  errorGetAccount,
		},
	}

	for _, testCase := range tests {

		mockAccount := &models.Account{
			AccountID: testCase.accountId,
			Balance:   19.33,
			CreatedAt: time.Now().Unix(),
		}

		t.Run(testCase.name, func(t *testing.T) {
			connectUri := "mongodb://localhost:" + mongoDbPort
			dbStore, client, errRt := New(connectUri, databaseName)
			if errRt != nil {
				assert.Nil(t, errRt)
				t.Fail()
			}
			assert.NotNil(t, client)
			ctx := context.Background()

			switch testCase.testType {
			case success:
				_, err := client.Database(databaseName).Collection(AccountsCollectionName).InsertOne(ctx, mockAccount)
				if err != nil {
					assert.NoError(t, err)
					t.Fail()
				}

				acc, err := dbStore.GetAccountByID(testCase.accountId)
				assert.NoError(t, err)
				assert.NotNil(t, acc)
				assert.Equal(t, mockAccount, acc)

			case errorGetAccount:
				acc, err := dbStore.GetAccountByID(testCase.accountId)
				assert.Error(t, err)
				assert.Nil(t, acc)
			}
		})
	}
}

func TestMongoStore_UpdateBalance(t *testing.T) {
	const (
		success = iota
		errorOccurred
	)

	var tests = []struct {
		name      string
		accountId string
		testType  int
	}{
		{
			name:      "Test update account balance successfully",
			accountId: "accountId",
			testType:  success,
		},
		{
			name:      "Test error updating balance",
			accountId: "invalid_id",
			testType:  errorOccurred,
		},
	}

	for _, testCase := range tests {

		mockAccount := &models.Account{
			AccountID: testCase.accountId,
			Balance:   19.33,
			CreatedAt: time.Now().Unix(),
		}

		t.Run(testCase.name, func(t *testing.T) {
			connectUri := "mongodb://localhost:" + mongoDbPort
			dbStore, client, errRt := New(connectUri, databaseName)
			if errRt != nil {
				assert.Nil(t, errRt)
				t.Fail()
			}
			assert.NotNil(t, client)
			ctx := context.Background()

			switch testCase.testType {
			case success:
				_, err := client.Database(databaseName).Collection(AccountsCollectionName).InsertOne(ctx, mockAccount)
				if err != nil {
					assert.NoError(t, err)
					t.Fail()
				}

				updateErr := dbStore.UpdateAccountBalance(testCase.accountId, 20)
				acc, accErr := dbStore.GetAccountByID(testCase.accountId)

				assert.NoError(t, updateErr)
				assert.NoError(t, accErr)
				assert.NotNil(t, acc)
				assert.Equal(t, float64(20), acc.Balance)

			case errorOccurred:
				_ = client.Disconnect(ctx)
				err := dbStore.UpdateAccountBalance(testCase.accountId, 20)
				assert.Error(t, err)
			}
		})
	}
}

func TestMongoStore_CreateTransaction(t *testing.T) {
	const (
		success = iota
		errorOccurred
	)

	var tests = []struct {
		name     string
		testType int
	}{
		{
			name:     "Test save transaction successfully",
			testType: success,
		},
		{
			name:     "Test error saving transaction",
			testType: errorOccurred,
		},
	}

	for _, testCase := range tests {

		mockTransaction := &models.Transaction{
			UserID:    "usr-0001",
			AccountID: "acc-0001",
			Amount:    2.5,
			Reference: "rand-ref",
			Type:      models.CREDIT,
			Status:    models.SUCCESS,
			CreatedAt: time.Now().Unix(),
		}

		t.Run(testCase.name, func(t *testing.T) {
			connectUri := "mongodb://localhost:" + mongoDbPort
			dbStore, client, errRt := New(connectUri, databaseName)
			if errRt != nil {
				assert.Nil(t, errRt)
				t.Fail()
			}
			assert.NotNil(t, client)
			ctx := context.Background()

			switch testCase.testType {
			case success:
				err := dbStore.CreateTransaction(mockTransaction)
				assert.NoError(t, err)

			case errorOccurred:
				_ = client.Disconnect(ctx)
				err := dbStore.CreateTransaction(mockTransaction)
				assert.Error(t, err)
			}
		})
	}
}

func TestMongoStore_GetTransactionByReference(t *testing.T) {
	const (
		success = iota
		errorNotFound
	)

	var tests = []struct {
		name      string
		reference string
		testType  int
	}{
		{
			name:      "Test get transaction successfully",
			reference: "trans-ref-001",
			testType:  success,
		},
		{
			name:      "Test error get transaction",
			reference: "invalid-ref",
			testType:  errorNotFound,
		},
	}

	for _, testCase := range tests {

		t.Run(testCase.name, func(t *testing.T) {
			connectUri := "mongodb://localhost:" + mongoDbPort
			dbStore, client, errRt := New(connectUri, databaseName)
			if errRt != nil {
				assert.Nil(t, errRt)
				t.Fail()
			}
			assert.NotNil(t, client)
			ctx := context.Background()

			mockTransaction := &models.Transaction{
				UserID:    "usr-0001",
				AccountID: "acc-0001",
				Amount:    2.5,
				Reference: testCase.reference,
				Type:      models.CREDIT,
				Status:    models.SUCCESS,
				CreatedAt: time.Now().Unix(),
			}

			switch testCase.testType {
			case success:
				_, err := client.Database(databaseName).Collection(TransactionsCollectionName).InsertOne(ctx, mockTransaction)
				if err != nil {
					assert.NoError(t, err)
					t.Fail()
				}

				transaction, err := dbStore.GetPaymentByReferenceId(testCase.reference)
				assert.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.Equal(t, mockTransaction, transaction)

			case errorNotFound:
				transaction, err := dbStore.GetPaymentByReferenceId(testCase.reference)
				assert.Error(t, err)
				assert.Nil(t, transaction)
			}
		})
	}
}

func TestMongoStore_GetUserById(t *testing.T) {
	const (
		success = iota
		errorNotFound
	)

	var tests = []struct {
		name     string
		userID   string
		testType int
	}{
		{
			name:     "Test get user successfully",
			userID:   "usr-001",
			testType: success,
		},
		{
			name:     "Test error get user",
			userID:   "invalid-usr",
			testType: errorNotFound,
		},
	}

	for _, testCase := range tests {

		t.Run(testCase.name, func(t *testing.T) {
			connectUri := "mongodb://localhost:" + mongoDbPort
			dbStore, client, errRt := New(connectUri, databaseName)
			if errRt != nil {
				assert.Nil(t, errRt)
				t.Fail()
			}
			assert.NotNil(t, client)
			ctx := context.Background()

			mockUser := &models.User{
				Id:        testCase.userID,
				Name:      "name",
				CreatedAt: time.Now().Unix(),
			}

			switch testCase.testType {
			case success:
				_, err := client.Database(databaseName).Collection(UserCollection).InsertOne(ctx, mockUser)
				if err != nil {
					assert.NoError(t, err)
					t.Fail()
				}

				user, err := dbStore.GetUserById(testCase.userID)
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, mockUser, user)

			case errorNotFound:
				user, err := dbStore.GetUserById(testCase.userID)
				assert.Error(t, err)
				assert.Nil(t, user)
			}
		})
	}
}
