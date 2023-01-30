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

func TestHandleCreatePost(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       fakePostCol,
			expectedResponse: `{"InsertedID":"` + fakePostObjIdHex + `"}`,
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
			subRouter.HandleFunc("/", a.handleCreatePost(new(MockPost), "fakeDb", st.collection)).Methods(http.MethodPost)

			w := httptest.NewRecorder()
			jsonBody := strings.NewReader(`{"content":"fake content", "author":"fake author"}`)
			r, err := http.NewRequest(http.MethodPost, "/post/", jsonBody)
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
			collection:       fakePostCol,
			postIdHex:        fakePostObjIdHex,
			expectedResponse: `{"id":"` + fakePostObjIdHex + `","content":"fake content","author":"fake author"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			postIdHex:        fakePostObjIdHex,
			expectedResponse: `{"error":"dummy error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "return-error-no-docs",
			collection:       "NoDocs",
			postIdHex:        fakePostObjIdHex,
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

			subRouter.HandleFunc("/{id:[a-z0-9]+}", a.handleGetPost(new(MockPost), "fakeDb", st.collection)).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/post/%v", st.postIdHex)
			r, err := http.NewRequest(http.MethodGet, url, nil)
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
			collection:       fakePostCol,
			postIdHex:        fakePostObjIdHex,
			expectedResponse: `{"msj":"post updated"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			postIdHex:        fakePostObjIdHex,
			expectedResponse: `{"error":"dummy error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "no-docs",
			collection:       "NoDocs",
			postIdHex:        fakePostObjIdHex,
			expectedResponse: `{"error":"mongo: no documents in result"}`,
			expectedCode:     http.StatusNotFound,
		},
		{
			name:             "return-error-invalid-hex-id",
			collection:       fakePostCol,
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
			collection:       fakePostCol,
			postIdHex:        fakePostObjIdHex,
			expectedResponse: `{"msj":"post deleted"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			postIdHex:        fakePostObjIdHex,
			expectedResponse: `{"error":"dummy error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "no-docs",
			collection:       "NoDocs",
			postIdHex:        fakePostObjIdHex,
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

func TestHandleCreateComment(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       fakeCommentCol,
			expectedResponse: `{"InsertedID":"` + fakeCommentObjIdHex + `"}`,
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
			subRouter := router.PathPrefix("/comment").Subrouter()

			mockModels := NewMockModels()
			mockModels.CColName = st.collection
			mockModels.PColName = fakePostCol
			subRouter.HandleFunc("/", a.handleCreateComment(mockModels, fakeDbName)).Methods(http.MethodPost)

			w := httptest.NewRecorder()
			jsonBody := strings.NewReader(`{"content":"fake content", "author":"fake author", "postId":"89372c88c133e1e4deb0e10a"}`)
			r, err := http.NewRequest(http.MethodPost, "/comment/", jsonBody)
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

func TestHandleGetComment(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		commentIdHex     string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       fakeCommentCol,
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"id":"` + fakeCommentObjIdHex + `","content":"fake content","author":"fake author","postId":"` + fakePostObjIdHex + `"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"error":"dummy ReadComment error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "return-error-no-docs",
			collection:       "NoDocs",
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"error":"mongo: no documents in result"}`,
			expectedCode:     http.StatusNotFound,
		},
		{
			name:             "return-error-invalid-hex-id",
			collection:       fakeCommentCol,
			commentIdHex:     "12345",
			expectedResponse: `{"error":"the provided hex string is not a valid ObjectID"}`,
			expectedCode:     http.StatusBadRequest,
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			a := new(App)

			router := mux.NewRouter()
			subRouter := router.PathPrefix("/comment").Subrouter()

			subRouter.HandleFunc("/{id:[a-z0-9]+}", a.handleGetComment(new(MockComment), fakeDbName, st.collection)).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/comment/%v", st.commentIdHex), nil)
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

func TestHandlePutComment(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		commentIdHex     string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       fakeCommentCol,
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"msj":"comment updated"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"error":"dummy ReadComment error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "no-docs",
			collection:       "NoDocs",
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"error":"mongo: no documents in result"}`,
			expectedCode:     http.StatusNotFound,
		},
		{
			name:             "return-error-invalid-hex-id",
			collection:       "fakePostCol",
			commentIdHex:     "12345",
			expectedResponse: `{"error":"the provided hex string is not a valid ObjectID"}`,
			expectedCode:     http.StatusBadRequest,
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			a := new(App)
			router := mux.NewRouter()
			subRouter := router.PathPrefix("/comment").Subrouter()
			subRouter.HandleFunc("/{id:[a-z0-9]+}", a.handlePutComment(new(MockComment), fakeDbName, st.collection)).Methods(http.MethodPut)

			w := httptest.NewRecorder()
			jsonBody := strings.NewReader(`{"content":"updated fake comment content", "author":"fake author"}`)
			url := fmt.Sprintf("/comment/%v", st.commentIdHex)
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

func TestHandleDeleteComment(t *testing.T) {
	subtests := []struct {
		name             string
		collection       string
		commentIdHex     string
		expectedResponse string
		expectedCode     int
	}{
		{
			name:             "happy-path",
			collection:       fakeCommentCol,
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"msj":"comment deleted"}`,
			expectedCode:     http.StatusOK,
		},
		{
			name:             "return-error",
			collection:       "fakeOtherCol",
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"error":"dummy ReadComment error"}`,
			expectedCode:     http.StatusInternalServerError,
		},
		{
			name:             "no-docs",
			collection:       "NoDocs",
			commentIdHex:     fakeCommentObjIdHex,
			expectedResponse: `{"error":"mongo: no documents in result"}`,
			expectedCode:     http.StatusNotFound,
		},
		{
			name:             "return-error-invalid-hex-id",
			collection:       fakeCommentCol,
			commentIdHex:     "12345",
			expectedResponse: `{"error":"the provided hex string is not a valid ObjectID"}`,
			expectedCode:     http.StatusBadRequest,
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			a := new(App)
			router := mux.NewRouter()
			subRouter := router.PathPrefix("/comment").Subrouter()

			subRouter.HandleFunc("/{id:[a-z0-9]+}", a.handleDeleteComment(new(MockComment), fakeDbName, st.collection)).Methods(http.MethodDelete)

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/comment/%v", st.commentIdHex)
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
