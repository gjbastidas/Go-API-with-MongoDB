package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Post struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content,omitempty" bson:"content,omitempty"`
	Author  string             `json:"author,omitempty" bson:"author,omitempty"`
}

func (p *Post) createPost(ctx context.Context, mCl *mongo.Client, dbName, colName string) (*mongo.InsertOneResult, error) {
	res, err := mCl.Database(dbName).Collection(colName).InsertOne(ctx, *p)
	return res, err
}

// func (p *post) readPost(mCl *mongo.Client) error {
// 	return errors.New("Not implemented")
// }

// func (p *post) updatePost(mCl *mongo.Client) error {
// 	return errors.New("Not implemented")
// }

// func (p *post) deletePost(mCl *mongo.Client) error {
// 	return errors.New("Not implemented")
// }
