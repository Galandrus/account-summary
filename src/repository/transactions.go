package repository

import (
	"account-summary/src/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository interface {
	GetTransactions(ctx context.Context, accountId string) ([]models.Transaction, error)
	CreateTransactions(ctx context.Context, transactions []models.Transaction) error
}

type transactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(mongoClient *mongo.Client) TransactionRepository {
	collection := mongoClient.Database(dbName).Collection(collectionTransactions)

	return &transactionRepository{collection: collection}
}

func (r *transactionRepository) GetTransactions(ctx context.Context, accountId string) ([]models.Transaction, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"accountId": accountId})
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := cursor.All(context.Background(), &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) CreateTransactions(ctx context.Context, transactions []models.Transaction) error {
	transactionBson := make([]interface{}, len(transactions))
	for i, transaction := range transactions {
		transactionBson[i] = transaction
	}

	log.Println("Inserting transactions:", len(transactionBson))

	_, err := r.collection.InsertMany(context.Background(), transactionBson)
	if err != nil {
		return err
	}

	return nil
}
