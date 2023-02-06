package app

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDDAO struct {
	c *mongo.Collection
}

func NewIDDAO(c *mongo.Client) *IDDAO {
	return &IDDAO{
		c: c.Database(database).Collection("ids"),
	}
}

func (dao *IDDAO) Insert(ctx context.Context, id *UrlID) error {
	_, err := dao.c.InsertOne(ctx, id)
	return err
}

func (dao *IDDAO) FindUnused(ctx context.Context) ([]*UrlID, error) {
	q := bson.D{{"used", bson.D{{"$ne", true}}}}
	cursor, err := dao.c.Find(ctx, q, options.Find().SetLimit(1000))
	if err != nil {
		return nil, err
	}

	var result []*UrlID
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (dao *IDDAO) ReserveID(ctx context.Context, id string) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"used", true}}}}
	updateResult, err := dao.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return ErrIDNotFound
	}
	return nil
}
