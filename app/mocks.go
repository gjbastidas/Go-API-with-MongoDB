package app

import (
	"errors"

	appDb "github.com/gjbastidas/GoSimpleAPIWithMongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	fakeDbName          = "fakeDb"
	fakePostCol         = "fakePostCol"
	fakeCommentCol      = "fakeCommentCol"
	fakeObjIdHex        = "89372c88c133e1e4deb0e10a"
	fakeCommentObjIdHex = "bfc80a35195ed2079d97c43b"
)

func getObjId(hex string) primitive.ObjectID {
	out, _ := primitive.ObjectIDFromHex(hex)
	return out
}

func NewMockModels() *appDb.Models {
	return &appDb.Models{
		P: &MockPost{},
		C: &MockComment{},
	}
}

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
		if colName != fakePostCol && colName != fakeCommentCol {
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

type MockComment struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content,omitempty" bson:"content,omitempty"`
	Author  string             `json:"author,omitempty" bson:"author,omitempty"`
	PostId  primitive.ObjectID `json:"postId,omitempty" bson:"post,omitempty"`
}

func (mC *MockComment) CreateComment(mCl *mongo.Client, dbName, colName string) (*mongo.InsertOneResult, error) {
	if dbName == fakeDbName {
		out := &mongo.InsertOneResult{}
		if colName != fakeCommentCol {
			return out, errors.New("dummy error")
		}
		objId := getObjId(fakeCommentObjIdHex)
		out.InsertedID = objId
		return out, nil
	}
	return &mongo.InsertOneResult{}, nil
}

func (mC *MockComment) ReadComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*appDb.CommentDoc, error) {
	return &appDb.CommentDoc{}, nil
}

func (mC *MockComment) UpdateComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	return nil
}

func (mC *MockComment) DeleteComment(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) error {
	return nil
}

func (mC *MockComment) GetRelatedPostId() primitive.ObjectID {
	return getObjId(fakeObjIdHex)
}
