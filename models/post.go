package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostDoc struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content,omitempty" bson:"content,omitempty"`
	Author  string             `json:"author,omitempty" bson:"author,omitempty"`
}

func (p *PostDoc) CreatePost(mCl *mongo.Client, dbName, colName string) error {
	return createOneRecord(mCl, p, dbName, colName)
}

func (p *PostDoc) ReadPost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*PostDoc, error) {
	return readOneRecord(mCl, p, objId, dbName, colName)
}

func (p *PostDoc) UpdatePost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	return updateOneRecord(mCl, p, objId, dbName, colName)
}

func (p *PostDoc) DeletePost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	return deleteOneRecord(mCl, objId, dbName, colName)
}
