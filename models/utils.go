package models

import (
	"context"
	"fmt"

	appConstants "github.com/gjbastidas/GoSimpleAPIWithMongoDB/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnyDoc interface {
	*PostDoc | *CommentDoc
}

func New(DbUsername, DbPassword, DbHost, DbPort string) (*mongo.Client, error) {
	connStr := fmt.Sprintf("mongodb://%v:%v@%v:%v/?authSource=%v", DbUsername, DbPassword, DbHost, DbPort, DbUsername)
	mClientOpts := options.Client().ApplyURI(connStr)
	ctx, cancel := context.WithTimeout(context.TODO(), appConstants.RequestTimeout)
	defer cancel()
	return mongo.Connect(ctx, mClientOpts)
}

func createOneRecord[D AnyDoc](mCl *mongo.Client, d D, dbName, colName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	_, err := mCl.Database(dbName).Collection(colName).InsertOne(ctx, d)
	return err
}

func readOneRecord[D AnyDoc](mCl *mongo.Client, d D, objId primitive.ObjectID, dbName, colName string) (D, error) {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	filter := bson.M{"_id": objId}
	err := mCl.Database(dbName).Collection(colName).FindOne(ctx, filter).Decode(d)
	return d, err
}

func updateOneRecord[D AnyDoc](mCl *mongo.Client, d D, objId primitive.ObjectID, dbName, colName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	filter := bson.M{"_id": objId}
	update := bson.M{"$set": d}
	_, err := mCl.Database(dbName).Collection(colName).UpdateOne(ctx, filter, update)
	return err
}

func deleteOneRecord(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	filter := bson.M{"_id": objId}
	_, err := mCl.Database(dbName).Collection(colName).DeleteOne(ctx, filter)
	return err
}
