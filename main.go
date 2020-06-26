package main

import (
	"context"
	"log"
	"mongo-crud/app/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DbCollections stores db collections
type dbCollections struct {
	podcasts *mongo.Collection
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("quickstart")
	startMux(&dbCollections{podcasts: db.Collection("podcasts")})
}

func startMux(cols *dbCollections) {
	r := mux.NewRouter()

	r.HandleFunc("/podcasts/{id}", models.GetPodcast(cols.podcasts)).Methods("GET")
	r.HandleFunc("/podcasts", models.GetPodcasts(cols.podcasts)).Methods("GET")
	r.HandleFunc("/podcasts", models.InsertPodcast(cols.podcasts)).Methods("POST")
	r.HandleFunc("/podcasts/{id}", models.UpdatePodcast(cols.podcasts)).Methods("PUT")
	r.HandleFunc("/podcasts/{id}", models.DeletePodcast(cols.podcasts)).Methods("DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
