package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	appConstants "github.com/gjbastidas/GoSimpleAPIWithMongoDB/constants"
	"github.com/gjbastidas/GoSimpleAPIWithMongoDB/env"
	appDb "github.com/gjbastidas/GoSimpleAPIWithMongoDB/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog"
)

// App defines the application
type App struct {
	mCl *mongo.Client
	cfg *env.AppConfig
}

// New App's constructor
func New() *App {
	a := new(App)
	var err error

	// set environment variables
	a.cfg, err = env.Config()
	if err != nil {
		klog.Fatalf("bad application configuration. error: %v", err)
	}

	// set mongodb client and ping db
	connStr := fmt.Sprintf("mongodb://%v:%v@%v:%v/?authSource=%v", a.cfg.DbUsername, a.cfg.DbPassword, a.cfg.DbHost, a.cfg.DbPort, a.cfg.DbUsername)
	mClientOpts := options.Client().ApplyURI(connStr)
	ctx, cancel := context.WithTimeout(context.TODO(), appConstants.RequestTimeout)
	defer cancel()
	a.mCl, err = mongo.Connect(ctx, mClientOpts)
	if err != nil {
		klog.Fatalf("cannot set mongodb client: %v", err)
	}

	a.serve()
	return a
}

// serve wires up routes and run server
func (a *App) serve() {
	// routing details
	r := mux.NewRouter()

	pSbr := r.PathPrefix("/post").Subrouter()
	pSbr.HandleFunc("/", a.handleCreatePost(&appDb.PostDoc{}, a.cfg.DbName, "posts")).Methods("POST")
	pSbr.HandleFunc("/{id:[a-z0-9]+}", a.handleGetPost(&appDb.PostDoc{}, a.cfg.DbName, "posts")).Methods("GET")
	pSbr.HandleFunc("/{id:[a-z0-9]+}", a.handlePutPost(&appDb.PostDoc{}, a.cfg.DbName, "posts")).Methods("PUT")
	pSbr.HandleFunc("/{id:[a-z0-9]+}", a.handleGetPost(&appDb.PostDoc{}, a.cfg.DbName, "posts")).Methods("DELETE")

	// http server configs
	srv := &http.Server{
		Addr:    a.cfg.SvrAddr,
		Handler: r,
	}

	// graceful server shutdown
	done := make(chan struct{})
	go func() {
		osSigs := make(chan os.Signal, 1)
		signal.Notify(osSigs, syscall.SIGINT, syscall.SIGTERM)
		<-osSigs
		klog.Info("os interrupt signal received")

		ctx, cancel := context.WithTimeout(context.Background(), appConstants.ServerTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			klog.Errorf("server shutdown error: %v", err)
		}
		klog.Info("server shutdown complete")

		close(done)
	}()

	// start http server
	klog.Infof("app started at %v", a.cfg.SvrAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		klog.Fatalf("server failed to start: %v", err)
	}

	// waits until any SIGINT or SIGTERM os signal is sent
	<-done
	klog.Info("app stopped")
}

func (a *App) handleCreatePost(p appDb.Post, dbName, colName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "cannot decode body")
			return
		}

		err = p.CreateOneRecord(a.mCl, dbName, colName)
		if err != nil {
			jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot create post")
			return
		}

		jsonPrint(w, http.StatusCreated, map[string]string{"msj": "post created"})
	}
}

func (a *App) handleGetPost(p appDb.Post, dbName, colName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		objId, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "invalid post id")
		}

		res, err := p.ReadOneRecord(a.mCl, objId, dbName, colName)
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "cannot decode post")
		}

		jsonPrint(w, http.StatusOK, res)
	}
}

func (a *App) handlePutPost(p appDb.Post, dbName, colName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		objId, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "invalid post id")
		}

		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "cannot decode body")
			return
		}

		err = p.UpdateOneRecord(a.mCl, objId, dbName, colName)
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "cannot decode post")
		}

		jsonPrint(w, http.StatusOK, map[string]string{"msj": "post updated"})
	}
}

func (a *App) handleDeletePost(p appDb.Post, dbName, colName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		objId, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "invalid post id")
		}

		err = p.DeleteOneRecord(a.mCl, objId, dbName, colName)
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "cannot decode post")
		}
		jsonPrint(w, http.StatusOK, map[string]string{"msj": "post deleted"})
	}
}

// jsonPrint prints output in json format
func jsonPrint(w http.ResponseWriter, code int, res any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		klog.Errorf("cannot encode response: %v", err)
	}
}

// jsonPrintError error log to server console and prints out error in json format
func jsonPrintError(w http.ResponseWriter, code int, errMsj, consoleMsj string) {
	klog.Errorf(consoleMsj+" : %v", errMsj)
	jsonPrint(w, code, map[string]string{"error": errMsj})
}
