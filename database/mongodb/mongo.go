package mongodb

import (
	"consumer-payment-service/database"
	"consumer-payment-service/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	AccountsCollectionName     = "accounts"
	TransactionsCollectionName = "transactions"
	UserCollection             = "users"
)

type mongodbStore struct {
	mongodbClient *mongo.Client
	databaseName  string
}

func (m *mongodbStore) collection(collectionName string) *mongo.Collection {
	return m.mongodbClient.Database(m.databaseName).Collection(collectionName)
}

// New returns a mongo instance that implements the mongodbstore
func New(connectUri, databaseName string) (database.MongoDBStore, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(connectUri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, nil, err
	}

	return &mongodbStore{mongodbClient: client, databaseName: databaseName}, client, nil
}

func (m *mongodbStore) GetAccountByID(accountId string) (*models.Account, error) {
	filter := bson.M{"account_id": accountId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	account := &models.Account{}

	err := m.collection(AccountsCollectionName).FindOne(ctx, filter).Decode(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *mongodbStore) UpdateAccountBalance(accountId string, amount float64) error {
	filter := bson.M{"account_id": accountId}
	update := bson.M{
		"$set": bson.M{
			"balance": amount,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.collection(AccountsCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongodbStore) CreateTransaction(transaction *models.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.collection(TransactionsCollectionName).InsertOne(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongodbStore) GetPaymentByReferenceId(reference string) (*models.Transaction, error) {
	filter := bson.M{"reference": reference}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	transaction := &models.Transaction{}

	err := m.collection(TransactionsCollectionName).FindOne(ctx, filter).Decode(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (m *mongodbStore) GetUserById(userId string) (*models.User, error) {
	filter := bson.M{"user_id": userId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{}

	err := m.collection(UserCollection).FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
