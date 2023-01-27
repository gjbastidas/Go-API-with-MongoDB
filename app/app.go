package app

import (
	"context"
	"encoding/json"
	"errors"
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
	"k8s.io/klog"
)

// TODO: create Comment handlers

type App struct {
	mCl *mongo.Client
	cfg *env.AppConfig
}

func New() *App {
	a := new(App)
	var err error

	// set environment variables
	a.cfg, err = env.Config()
	if err != nil {
		klog.Fatalf("bad application configuration. error: %v", err)
	}

	// set mongodb client
	a.mCl, err = appDb.NewClient(a.cfg.DbUsername, a.cfg.DbPassword, a.cfg.DbHost, a.cfg.DbPort)
	if err != nil {
		klog.Fatalf("cannot set mongodb client: %v", err)
	}

	// ping db
	ctx, cancel := context.WithTimeout(context.Background(), appConstants.RequestTimeout)
	defer cancel()
	err = a.mCl.Ping(ctx, nil)
	if err != nil {
		klog.Fatal(err)
	}

	a.serve()
	return a
}

// serve wires up routes and run server
func (a *App) serve() {
	// routing details
	r := mux.NewRouter()

	pSbr := r.PathPrefix("/post").Subrouter()
	pSbr.HandleFunc("/", a.handleCreatePost(new(appDb.PostDoc), appConstants.DbName, appConstants.PColl)).Methods("POST")
	pSbr.HandleFunc("/{id:[a-z0-9]+}", a.handleGetPost(new(appDb.PostDoc), appConstants.DbName, appConstants.PColl)).Methods("GET")
	pSbr.HandleFunc("/{id:[a-z0-9]+}", a.handlePutPost(new(appDb.PostDoc), appConstants.DbName, appConstants.PColl)).Methods("PUT")
	pSbr.HandleFunc("/{id:[a-z0-9]+}", a.handleDeletePost(new(appDb.PostDoc), appConstants.DbName, appConstants.PColl)).Methods("DELETE")

	// http server configs
	srv := &http.Server{
		Addr:    ":8088",
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
	klog.Info("app started at :8088")
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

		res, err := p.CreatePost(a.mCl, dbName, colName)
		if err != nil {
			jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot create post")
			return
		}

		jsonPrint(w, http.StatusCreated, res)
	}
}

func (a *App) handleGetPost(p appDb.Post, dbName, colName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		objId, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "invalid post id")
			return
		}

		res, err := p.ReadPost(a.mCl, objId, dbName, colName)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				jsonPrintError(w, http.StatusNotFound, err.Error(), "post not found")
				return
			default:
				jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot read post")
				return
			}
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
			return
		}

		_, err = p.ReadPost(a.mCl, objId, dbName, colName)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				jsonPrintError(w, http.StatusNotFound, err.Error(), "post not found")
				return
			default:
				jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot read post")
				return
			}
		}

		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			jsonPrintError(w, http.StatusBadRequest, err.Error(), "cannot decode body")
			return
		}

		err = p.UpdatePost(a.mCl, objId, dbName, colName)
		if err != nil {
			jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot update post")
			return
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
			return
		}

		_, err = p.ReadPost(a.mCl, objId, dbName, colName)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				jsonPrintError(w, http.StatusNotFound, err.Error(), "post not found")
				return
			default:
				jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot read post")
				return
			}
		}

		err = p.DeletePost(a.mCl, objId, dbName, colName)
		if err != nil {
			jsonPrintError(w, http.StatusInternalServerError, err.Error(), "cannot delete post")
			return
		}

		jsonPrint(w, http.StatusOK, map[string]string{"msj": "post deleted"})
	}
}
