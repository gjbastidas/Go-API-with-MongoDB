package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appMocks "github.com/gjbastidas/GoSimpleAPIWithMongoDB/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	var a App
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/posts").Subrouter()
	subRouter.HandleFunc("/", a.handleCreatePost(context.TODO(), &appMocks.MockPost{}, "fakeDB", "fakeCollection")).Methods("POST")

	w := httptest.NewRecorder()
	json := strings.NewReader(`{"content":"updated post", "author":"gus bast"}`)
	r, err := http.NewRequest("POST", "/posts/", json)
	router.ServeHTTP(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
}
