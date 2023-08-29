package repository

import (
	"context"
	"os"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WalleteRepository struct {
	client *mongo.Client
}

func NewWalleteRepository(c *mongo.Client) *WalleteRepository {
	return &WalleteRepository{
		client: c,
	}
}

func (wr WalleteRepository) Load(ctx context.Context, userID string) (*entity.Wallet, error) {
	walletColl := wr.client.Database(os.Getenv("MONGO_DB")).Collection("wallets")

	var wallet entity.Wallet

	err := walletColl.FindOne(ctx, bson.D{
		{Key: "user_id", Value: userID},
	}).Decode(&wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (wr WalleteRepository) Transfer(ctx context.Context, transfer entity.Transfer) (userPayer, userPayee *entity.User, err error) {
	session, err := wr.client.StartSession()
	if err != nil {
		return nil, nil, err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		userColl := wr.client.Database(os.Getenv("MONGO_DB")).Collection("users")
		if err := userColl.FindOne(ctx, bson.D{{Key: "user_id", Value: transfer.Payer}}).Decode(&userPayer); err != nil {
			return nil, err
		}
		if err := userColl.FindOne(ctx, bson.D{{Key: "user_id", Value: transfer.Payee}}).Decode(&userPayee); err != nil {
			return nil, err
		}

		var wallet entity.Wallet

		walletColl := wr.client.Database(os.Getenv("MONGO_DB")).Collection("wallets")
		if err := walletColl.FindOne(ctx, bson.D{{Key: "user_id", Value: userPayer.UserID}}).Decode(&wallet); err != nil {
			return nil, err
		}

		walletID, err := primitive.ObjectIDFromHex(wallet.ID)
		if err != nil {
			return nil, err
		}

		walletColl.UpdateByID(ctx, walletID, bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "balance", Value: wallet.Balance - transfer.Amount},
			}},
		})

		if err := walletColl.FindOne(ctx, bson.D{{Key: "user_id", Value: userPayee.UserID}}).Decode(&wallet); err != nil {
			return nil, err
		}

		walletID, err = primitive.ObjectIDFromHex(wallet.ID)
		if err != nil {
			return nil, err
		}

		walletColl.UpdateByID(ctx, walletID, bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "balance", Value: wallet.Balance + transfer.Amount},
			}},
		})

		if err := session.CommitTransaction(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return
}
