package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentDoc struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content,omitempty" bson:"content,omitempty"`
	Author  string             `json:"author,omitempty" bson:"author,omitempty"`
	PostId  primitive.ObjectID `json:"post,omitempty" bson:"post,omitempty"`
}

func (c *CommentDoc) CreateComment(mCl *mongo.Client, dbName, colName string) error {
	return createOneRecord(mCl, c, dbName, colName)
}

func (c *CommentDoc) ReadComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*CommentDoc, error) {
	return readOneRecord(mCl, c, objId, dbName, colName)
}

func (c *CommentDoc) UpdateComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	return updateOneRecord(mCl, c, objId, dbName, colName)
}

func (c *CommentDoc) DeleteComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	return deleteOneRecord(mCl, objId, dbName, colName)
}
