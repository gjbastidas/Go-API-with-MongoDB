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
	fakeDbName  = "fakeDb"
	fakePostCol = "fakePostCol"
)

type MockPost struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content    string             `json:"content,omitempty" bson:"content,omitempty"`
	Author     string             `json:"author,omitempty" bson:"author,omitempty"`
	Database   string
	Collection string
}

func (mP *MockPost) CreatePost(mCl *mongo.Client, dbName, colName string) error {
	if dbName == fakeDbName && colName != fakePostCol {
		return errors.New("dummy error")
	}
	return nil
}

func (mP *MockPost) ReadPost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*appDb.PostDoc, error) {
	if dbName == fakeDbName {
		out := new(appDb.PostDoc)
		if colName != fakePostCol {
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
		return errors.New("dummy error")
	}
	return nil
}

func (mP *MockPost) DeletePost(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	if dbName == fakeDbName && colName != fakePostCol {
		return errors.New("dummy error")
	}
	return nil
}
