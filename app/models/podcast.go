package models

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Podcast struct describing podcast model
type Podcast struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Author string             `json:"author,omitempty" bson:"author,omitempty"`
	Tags   []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}

//DeletePodcast handler
func DeletePodcast(podcasts *mongo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("content-type", "application/json")
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])
		result, err := podcasts.DeleteOne(context.Background(), Podcast{ID: id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		respondWithJSON(w, http.StatusOK, result.DeletedCount)
	}
}

//UpdatePodcast handler
func UpdatePodcast(podcasts *mongo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("content-type", "application/json")
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])
		var podcast Podcast
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&podcast)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		result, err := podcasts.UpdateOne(context.Background(), Podcast{ID: id}, bson.D{{"$set", &podcast}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		respondWithJSON(w, http.StatusOK, result.UnmarshalBSON)
	}
}

//InsertPodcast handler
func InsertPodcast(podcasts *mongo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var podcast Podcast
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&podcast)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		insertResult, err := podcasts.InsertOne(context.Background(), podcast)
		if err != nil {
			panic(err)
		}
		respondWithJSON(w, http.StatusOK, insertResult.InsertedID)
	}
}

//GetPodcast handler
func GetPodcast(podcasts *mongo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("content-type", "application/json")
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])
		var podcast Podcast
		err := podcasts.FindOne(context.Background(), Podcast{ID: id}).Decode(&podcast)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, &podcast)
	}
}

//GetPodcasts handler
func GetPodcasts(podcasts *mongo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var podcastsRes []Podcast
		cursor, err := podcasts.Find(context.Background(), bson.D{})
		if err != nil {
			panic(err)
		}
		if err = cursor.All(context.Background(), &podcastsRes); err != nil {
			panic(err)
		}
		respondWithJSON(w, http.StatusOK, &podcastsRes)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
