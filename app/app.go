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
	"time"

	appDb "github.com/gjbastidas/GoSimpleAPIWithMongoDB/db"
	"github.com/gjbastidas/GoSimpleAPIWithMongoDB/env"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog"
)

const (
	serverTimeout  = 1 * time.Second
	requestTimeout = 10 * time.Second
)

type App struct {
	db  *mongo.Client
	cfg *env.AppConfig
}

func New() *App {
	a := new(App)
	var err error
	a.cfg, err = env.Config()
	if err != nil {
		klog.Fatalf("bad application configuration. error: %v", err)
	}

	connStr := fmt.Sprintf("mongodb://%v:%v@%v:%v/?authSource=%v", a.cfg.DbUsername, a.cfg.DbPassword, a.cfg.DbHost, a.cfg.DbPort, a.cfg.DbUsername)
	mClientOpts := options.Client().ApplyURI(connStr)
	ctx, cancel := context.WithTimeout(context.TODO(), requestTimeout)
	defer cancel()
	a.db, err = mongo.Connect(ctx, mClientOpts)
	if err != nil {
		klog.Fatalf("cannot set mongodb client: %v", err)
	}
	err = a.db.Ping(ctx, nil)
	if err != nil {
		klog.Fatalf("cannot connect to mongodb: %v", err)
	}

	a.serve()
	return a
}

func (a *App) serve() {
	// routing details
	r := mux.NewRouter()

	pSbr := r.PathPrefix("/posts").Subrouter()
	pSbr.HandleFunc("/", a.handleCreatePost(&appDb.Post{}, a.cfg.DbName, "posts")).Methods("POST")

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

		ctx, cancel := context.WithTimeout(context.Background(), serverTimeout)
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

	<-done
	klog.Info("app stopped")
}

func (a *App) handleCreatePost(post appDb.PostIface, dbName, colName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "cannot decode body")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		res, err := post.CreatePost(ctx, a.db, dbName, colName)
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
