package repository

import (
	"context"
	"os"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	client *mongo.Client
}

func NewUserRepository(c *mongo.Client) *UserRepository {
	return &UserRepository{
		client: c,
	}
}

func (ur UserRepository) CheckUserIsRegistered(ctx context.Context, user entity.User) (bool, error) {
	userColl := ur.client.Database(os.Getenv("MONGO_DB")).Collection("users")
	err := userColl.FindOne(ctx, bson.D{
		{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "email", Value: user.Email}},
				bson.D{{Key: "taxpayer_identification", Value: user.TaxpayerIdentification}},
			},
		},
	}).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (ur UserRepository) Save(ctx context.Context, user entity.User) error {
	session, err := ur.client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		userColl := ur.client.Database(os.Getenv("MONGO_DB")).Collection("users")
		if _, err := userColl.InsertOne(ctx, user); err != nil {
			return nil, err
		}

		wallet := entity.NewWallet(user.UserID, 0)

		walletColl := ur.client.Database(os.Getenv("MONGO_DB")).Collection("wallets")
		if _, err := walletColl.InsertOne(ctx, wallet); err != nil {
			return nil, err
		}

		if err := session.CommitTransaction(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}
