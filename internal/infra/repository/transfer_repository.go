package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransferRepository struct {
	client *mongo.Client
}

func NewTransferRepository(c *mongo.Client) *TransferRepository {
	return &TransferRepository{
		client: c,
	}
}

func (t TransferRepository) List(ctx context.Context, userID string, page, limit int64, transfers *[]entity.Transfer) error {
	transferColl := t.client.Database(os.Getenv("MONGO_DB")).Collection("transfers")
	options := new(options.FindOptions)
	options.SetSkip(page*limit - limit)
	options.SetLimit(limit)

	cur, err := transferColl.Find(ctx, bson.D{}, options)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var transfer entity.Transfer
		if err := cur.Decode(&transfer); err != nil {
			return err
		}
		fmt.Println(transfer)
		*transfers = append(*transfers, transfer)
	}

	return nil
}
