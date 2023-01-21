package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostIface interface {
	createPost(ctx context.Context, mCl *mongo.Client, dbName, colName string) (*mongo.InsertOneResult, error)
}
