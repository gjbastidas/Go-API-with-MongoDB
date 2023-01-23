package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostIface interface {
	CreatePost(ctx context.Context, mCl *mongo.Client, dbName, colName string) (*mongo.InsertOneResult, error)
}
