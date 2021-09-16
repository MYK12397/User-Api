package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	//	"os/user"

	ds "UserAPI/UserDS"
	"UserAPI/middleware"

	godb "UserAPI/database"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "text/json")

	singleuser := r.Context().Value(middleware.KeyUser{}).(*ds.User)
	json.NewDecoder(r.Body).Decode(&singleuser)

	collection := godb.ConnectDB()

	ans, err := collection.InsertOne(context.TODO(), singleuser)

	if err != nil {
		godb.GetError(err, rw)
		return
	}
	json.NewEncoder(rw).Encode(ans)
}

func GetUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/json")

	collection := godb.ConnectDB()
	usr, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		http.Error(rw, "Unable to find Id", http.StatusNotFound)
		return
	}

	defer usr.Close(context.TODO())

	var Users []ds.User
	for usr.Next(context.TODO()) {
		var u ds.User

		err := usr.Decode(&u)
		if err != nil {
			log.Fatal(err)
		}
		if u.Active != false {
			Users = append(Users, u)
		}

	}
	if err := usr.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(rw).Encode(Users)
}

func GetID(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	var user ds.User
	vars := mux.Vars(r)

	Id, er := primitive.ObjectIDFromHex(vars["id"])
	if er != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	sorted := bson.M{"_id": Id}

	collection := godb.ConnectDB()
	err := collection.FindOne(context.TODO(), sorted).Decode(&user)

	if err != nil {
		godb.GetError(err, rw)
		return
	}

	json.NewEncoder(rw).Encode(user)
}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	collection := godb.ConnectDB()
	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{"$set": bson.M{"active": false}}

	ans, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		godb.GetError(err, rw)
		return
	}

	json.NewEncoder(rw).Encode(ans)
}
