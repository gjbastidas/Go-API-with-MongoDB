package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	var a App
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/posts/", a.handleCreatePost(context.TODO(), &mockPost{}, "fakeDB", "fakeCollection")).Methods("POST")

	w := httptest.NewRecorder()
	json := strings.NewReader(`{"content":"updated post", "author":"gus bast"}`)
	r, err := http.NewRequest("POST", "/posts/", json)
	a.Router.ServeHTTP(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
}
