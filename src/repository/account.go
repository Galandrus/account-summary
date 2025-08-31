package repository

import (
	"account-summary/src/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountRepository interface {
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	GetAccountById(ctx context.Context, id string) (*models.Account, error)
	UpsertAccount(ctx context.Context, account models.Account) error
}

type accountRepository struct {
	collection *mongo.Collection
}

func NewAccountRepository(mongoClient *mongo.Client) AccountRepository {
	collection := mongoClient.Database(dbName).Collection(collectionAccounts)

	return &accountRepository{collection: collection}
}

func (r *accountRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	account := &models.Account{}
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return account, nil
}

func (r *accountRepository) GetAccountById(ctx context.Context, id string) (*models.Account, error) {
	account := &models.Account{}
	err := r.collection.FindOne(context.Background(), bson.M{"id": id}).Decode(account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return account, nil
}

func (r *accountRepository) UpsertAccount(ctx context.Context, account models.Account) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"id": account.ID}, bson.M{"$set": account}, opts)
	if err != nil {
		return err
	}

	return nil
}
