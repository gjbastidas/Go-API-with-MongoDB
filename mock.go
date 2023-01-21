package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockPost struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content,omitempty" bson:"content,omitempty"`
	Author  string             `json:"author,omitempty" bson:"author,omitempty"`
}

func (fP *mockPost) createPost(ctx context.Context, mCl *mongo.Client, dbName, colName string) (*mongo.InsertOneResult, error) {
	if dbName == "fakeDB" && colName == "fakeCollection" {
		res := &mongo.InsertOneResult{
			InsertedID: 1,
		}
		return res, nil
	}
	return &mongo.InsertOneResult{}, nil
}
