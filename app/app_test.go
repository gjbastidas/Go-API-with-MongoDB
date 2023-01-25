package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	appMocks "github.com/gjbastidas/GoSimpleAPIWithMongoDB/mocks"
	appDb "github.com/gjbastidas/GoSimpleAPIWithMongoDB/models"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var a *App
var router *mux.Router
var subRouter *mux.Router

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	a = new(App)
	router = mux.NewRouter()
	subRouter = router.PathPrefix("/post").Subrouter()
}

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mP := appMocks.NewMockPost(ctrl)
	mP.EXPECT().CreateOneRecord(a.mCl, "fakeDB", "fakeCollection").Return(nil)

	subRouter.HandleFunc("/", a.handleCreatePost(mP, "fakeDB", "fakeCollection")).Methods("POST")

	w := httptest.NewRecorder()
	jsonBody := strings.NewReader(`{"content":"fake content", "author":"fake author"}`)
	r, err := http.NewRequest("POST", "/post/", jsonBody)
	router.ServeHTTP(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)

	b, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, `{"msj":"post created"}`, strings.TrimSuffix(string(b), "\n"))
}

func TestGetPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mP := appMocks.NewMockPost(ctrl)
	testIdHex := "63cefca7da1c1470664ec41c"
	objId, _ := primitive.ObjectIDFromHex(testIdHex)
	mP.EXPECT().ReadOneRecord(gomock.Any(), objId, gomock.Any(), gomock.Any()).DoAndReturn(func(mCl *mongo.Client, objId primitive.ObjectID, dbName, colName string) (*appDb.PostDoc, error) {
		p := &appDb.PostDoc{
			Id:      objId,
			Content: "fake content",
			Author:  "fake author",
		}
		return p, nil
	})

	subRouter.HandleFunc("/{id:[a-z0-9]+}", a.handleGetPost(mP, "fakeDB", "fakeCollection")).Methods("GET")

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/post/%v", testIdHex)
	r, err := http.NewRequest("GET", url, nil)
	router.ServeHTTP(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}
