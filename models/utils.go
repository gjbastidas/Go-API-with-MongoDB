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

type Models struct {
	P        Post
	PColName string
	C        Comment
	CColName string
}

func NewModels() *Models {
	return &Models{
		P:        &PostDoc{},
		PColName: appConstants.PColl,
		C:        &CommentDoc{},
		CColName: appConstants.CColl,
	}
}

func NewClient(DbUsername, DbPassword, DbHost, DbPort string) (*mongo.Client, error) {
	connStr := fmt.Sprintf("mongodb://%v:%v@%v:%v/?authSource=%v", DbUsername, DbPassword, DbHost, DbPort, DbUsername)
	mClientOpts := options.Client().ApplyURI(connStr)
	return mongo.Connect(context.TODO(), mClientOpts)
}

type AnyDoc interface {
	*PostDoc | *CommentDoc
}

func createOneRecord[D AnyDoc](mCl *mongo.Client, d D, dbName, colName string) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	res, err := mCl.Database(dbName).Collection(colName).InsertOne(ctx, d)
	return res, err
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
