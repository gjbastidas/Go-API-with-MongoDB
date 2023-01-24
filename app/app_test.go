package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appMocks "github.com/gjbastidas/GoSimpleAPIWithMongoDB/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreatePost(t *testing.T) {
	var a App
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/post").Subrouter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPost := appMocks.NewMockPostIface(ctrl)
	mockPost.EXPECT().CreatePost(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.InsertOneResult{InsertedID: 123456}, nil).Times(1)

	subRouter.HandleFunc("/", a.handleCreatePost(mockPost, "fakeDB", "fakeCollection")).Methods("POST")

	w := httptest.NewRecorder()
	json := strings.NewReader(`{"content":"updated post", "author":"gus bast"}`)
	r, err := http.NewRequest("POST", "/post/", json)
	router.ServeHTTP(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetPost(t *testing.T) {
	// TODO
	var a App
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/post").Subrouter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPost := appMocks.NewMockPostIface(ctrl)
	mockPost.EXPECT().GetPost(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&mongo.SingleResult{}).Times(1)

	subRouter.HandleFunc("/", a.handleGetPost(mockPost, "fakeDB", "fakeCollection")).Methods("GET")

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/posts/123456", nil)
	router.ServeHTTP(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}
