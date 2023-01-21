package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gjbastidas/GoSimpleAPIWithMongoDB/env"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog"
)

const (
	timeout = 10 * time.Second
)

type App struct {
	Router *mux.Router
	Db     *mongo.Client
}

var appCfg *env.AppConfig

func main() {
	var err error
	appCfg, err = env.Config()
	if err != nil {
		klog.Fatalf("bad application configuration. error: %v", err)
	}

	connStr := fmt.Sprintf("mongodb://%v:%v@%v:%v/?authSource=%v", appCfg.DbUsername, appCfg.DbPassword, appCfg.DbHost, appCfg.DbPort, appCfg.DbUsername)
	mClientOpts := options.Client().ApplyURI(connStr)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var a App
	a.Db, err = mongo.Connect(ctx, mClientOpts)
	if err != nil {
		klog.Fatalf("cannot set mongodb client. error: %v", err)
	}

	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/posts/", a.handleCreatePost(ctx, &Post{}, appCfg.DbName, "posts")).Methods("POST")

	klog.Info("App started...")
	klog.Fatal(http.ListenAndServe(":8088", a.Router))
}

func (a *App) handleCreatePost(ctx context.Context, post PostIface, dbName, colName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "cannot decode body")
			return
		}

		res, err := post.createPost(ctx, a.Db, dbName, colName)
		if err != nil {
			jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot create post")
			return
		}
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot encode body")
			return
		}
	}
}

func jsonPrintError(w http.ResponseWriter, code int, errMsj, consoleMsj string) {
	klog.Errorf(consoleMsj+". error: %v", errMsj)
	w.WriteHeader(code)
	w.Write([]byte(`{ "message": "` + errMsj + `" }`))
}
