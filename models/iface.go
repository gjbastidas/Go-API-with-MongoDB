package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Post interface {
	CreatePost(mCl *mongo.Client, dbName, colName string) (*mongo.InsertOneResult, error)
	ReadPost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*PostDoc, error)
	UpdatePost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error
	DeletePost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error
}

type Comment interface {
	CreateComment(mCl *mongo.Client, dbName, colName string) error
	ReadComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*CommentDoc, error)
	UpdateComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error
	DeleteComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error
}
