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
				bson.D{{Key: "taxpayeer_identification", Value: user.TaxpayeerIdentification}},
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
	userColl := ur.client.Database(os.Getenv("MONGO_DB")).Collection("users")
	_, err := userColl.InsertOne(ctx, user)
	return err
}
