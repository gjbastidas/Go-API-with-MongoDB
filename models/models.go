package models

import (
	"context"

	appConstants "github.com/gjbastidas/GoSimpleAPIWithMongoDB/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Post interface {
	CreateOneRecord(mCl *mongo.Client, dbName, colName string) error
	ReadOneRecord(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*PostDoc, error)
	UpdateOneRecord(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error
	DeleteOneRecord(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error
}

type PostDoc struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content,omitempty" bson:"content,omitempty"`
	Author  string             `json:"author,omitempty" bson:"author,omitempty"`
}

func (p *PostDoc) CreateOneRecord(mCl *mongo.Client, dbName, colName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	_, err := mCl.Database(dbName).Collection(colName).InsertOne(ctx, p)
	return err
}

func (p *PostDoc) ReadOneRecord(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*PostDoc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	filter := bson.M{"_id": objId}
	pDoc := new(PostDoc)
	err := mCl.Database(dbName).Collection(colName).FindOne(ctx, filter).Decode(pDoc)
	return pDoc, err
}

func (p *PostDoc) UpdateOneRecord(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	filter := bson.M{"_id": objId}
	update := bson.M{"$set": p}
	_, err := mCl.Database(dbName).Collection(colName).UpdateOne(ctx, filter, update)
	return err
}

func (p *PostDoc) DeleteOneRecord(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	filter := bson.M{"_id": objId}
	_, err := mCl.Database(dbName).Collection(colName).DeleteOne(ctx, filter)
	return err
}
