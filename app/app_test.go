package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TODO: Test Comment handlers

func TestHandleCreatePost(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       "fakePostCol",
			expectedResponse: `{"InsertedID":"89372c88c133e1e4deb0e10a"}`,
			expectedCode:     http.StatusCreated,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			expectedResponse: `{"error":"dummy error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			a := new(App)
			router := mux.NewRouter()
			subRouter := router.PathPrefix("/post").Subrouter()
			subRouter.HandleFunc("/", a.handleCreatePost(new(MockPost), "fakeDb", st.collection)).Methods("POST")

			w := httptest.NewRecorder()
			jsonBody := strings.NewReader(`{"content":"fake content", "author":"fake author"}`)
			r, err := http.NewRequest("POST", "/post/", jsonBody)
			router.ServeHTTP(w, r)

			if assert.NoError(t, err) {
				assert.EqualValues(t, st.expectedCode, w.Code)
			}

			b, err := io.ReadAll(w.Body)
			if assert.NoError(t, err) {
				assert.EqualValues(t, st.expectedResponse, strings.TrimSuffix(string(b), "\n"))
			}
		})
	}
}

func TestHandleGetPost(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		postIdHex        string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       "fakePostCol",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"id":"89372c88c133e1e4deb0e10a","content":"fake content","author":"fake author"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"error":"dummy error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "return-error-no-docs",
			collection:       "NoDocs",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"error":"mongo: no documents in result"}`,
			expectedCode:     http.StatusNotFound,
		},
		{
			name:             "return-error-invalid-hex-id",
			collection:       "fakePostCol",
			postIdHex:        "12345",
			expectedResponse: `{"error":"the provided hex string is not a valid ObjectID"}`,
			expectedCode:     http.StatusBadRequest,
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			a := new(App)
			router := mux.NewRouter()
			subRouter := router.PathPrefix("/post").Subrouter()

			subRouter.HandleFunc("/{id:[a-z0-9]+}", a.handleGetPost(new(MockPost), "fakeDb", st.collection)).Methods("GET")

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/post/%v", st.postIdHex)
			r, err := http.NewRequest("GET", url, nil)
			router.ServeHTTP(w, r)

			if assert.NoError(t, err) {
				assert.EqualValues(t, st.expectedCode, w.Code)
			}

			b, err := io.ReadAll(w.Body)
			if assert.NoError(t, err) {
				assert.EqualValues(t, st.expectedResponse, strings.TrimSuffix(string(b), "\n"))
			}
		})
	}
}

func TestHandlePutPost(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		postIdHex        string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       "fakePostCol",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"msj":"post updated"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"error":"dummy error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "no-docs",
			collection:       "NoDocs",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"error":"mongo: no documents in result"}`,
			expectedCode:     http.StatusNotFound,
		},
		{
			name:             "return-error-invalid-hex-id",
			collection:       "fakePostCol",
			postIdHex:        "12345",
			expectedResponse: `{"error":"the provided hex string is not a valid ObjectID"}`,
			expectedCode:     http.StatusBadRequest,
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			a := new(App)
			router := mux.NewRouter()
			subRouter := router.PathPrefix("/post").Subrouter()
			subRouter.HandleFunc("/{id:[a-z0-9]+}", a.handlePutPost(new(MockPost), "fakeDb", st.collection)).Methods(http.MethodPut)

			w := httptest.NewRecorder()
			jsonBody := strings.NewReader(`{"content":"updated fake content", "author":"fake author"}`)
			url := fmt.Sprintf("/post/%v", st.postIdHex)
			r, err := http.NewRequest(http.MethodPut, url, jsonBody)
			router.ServeHTTP(w, r)

			if assert.NoError(t, err) {
				assert.EqualValues(t, st.expectedCode, w.Code)
			}

			b, err := io.ReadAll(w.Body)
			if assert.NoError(t, err) {
				assert.EqualValues(t, st.expectedResponse, strings.TrimSuffix(string(b), "\n"))
			}
		})
	}
}

func TestHandleDeletePost(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		postIdHex        string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       "fakePostCol",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"msj":"post deleted"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"error":"dummy error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "no-docs",
			collection:       "NoDocs",
			postIdHex:        "89372c88c133e1e4deb0e10a",
			expectedResponse: `{"error":"mongo: no documents in result"}`,
			expectedCode:     http.StatusNotFound,
		},
		{
			name:             "return-error-invalid-hex-id",
			collection:       "fakePostCol",
			postIdHex:        "12345",
			expectedResponse: `{"error":"the provided hex string is not a valid ObjectID"}`,
			expectedCode:     http.StatusBadRequest,
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			a := new(App)
			router := mux.NewRouter()
			subRouter := router.PathPrefix("/post").Subrouter()

			subRouter.HandleFunc("/{id:[a-z0-9]+}", a.handleDeletePost(new(MockPost), "fakeDb", st.collection)).Methods(http.MethodDelete)

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/post/%v", st.postIdHex)
			r, err := http.NewRequest(http.MethodDelete, url, nil)
			router.ServeHTTP(w, r)

			if assert.NoError(t, err) {
				assert.EqualValues(t, st.expectedCode, w.Code)
			}

			b, err := io.ReadAll(w.Body)
			if assert.NoError(t, err) {
				assert.EqualValues(t, st.expectedResponse, strings.TrimSuffix(string(b), "\n"))
			}
		})
	}
}
