package app

import (
	"errors"

	appDb "github.com/gjbastidas/GoSimpleAPIWithMongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: create MockComment

const (
	fakeDbName   = "fakeDb"
	fakePostCol  = "fakePostCol"
	fakeObjIdHex = "89372c88c133e1e4deb0e10a"
)

type MockPost struct {
}

func (mP *MockPost) CreatePost(mCl *mongo.Client, dbName, colName string) (*mongo.InsertOneResult, error) {
	if dbName == fakeDbName {
		out := &mongo.InsertOneResult{}
		if colName != fakePostCol {
			return out, errors.New("dummy error")
		}
		objId, _ := primitive.ObjectIDFromHex(fakeObjIdHex)
		out.InsertedID = objId
		return out, nil
	}
	return &mongo.InsertOneResult{}, nil
}

func (mP *MockPost) ReadPost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*appDb.PostDoc, error) {
	if dbName == fakeDbName {
		out := new(appDb.PostDoc)
		if colName != fakePostCol {
			if colName == "NoDocs" {
				return out, mongo.ErrNoDocuments
			}
			return out, errors.New("dummy error")
		}
		res, _ := bson.Marshal(bson.M{"_id": objId, "content": "fake content", "author": "fake author"})
		_ = bson.Unmarshal(res, out)
		return out, nil
	}
	return &appDb.PostDoc{}, nil
}

func (mP *MockPost) UpdatePost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	if dbName == fakeDbName && colName != fakePostCol {
		if colName == "NoDocs" {
			return mongo.ErrNoDocuments
		}
		return errors.New("dummy error")
	}
	return nil
}

func (mP *MockPost) DeletePost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	if dbName == fakeDbName && colName != fakePostCol {
		if colName == "NoDocs" {
			return mongo.ErrNoDocuments
		}
		return errors.New("dummy error")
	}
	return nil
}
