package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//format to get inputs from console and write outputs to console

type Users struct{
  //ID, name, email ,password
  ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitmpty"`
  Name string `json:"name,omitempty" bson:"name,omitempty"`
  Email string `json:"email,omitempty" bson:"email,omitempty"`
  Password string `json:"password,omitempty" bson:"password,omitempty"`
}

var client *mongo.Client

//response of type http Response writer, reqwuest of type http request
func CreateUsersEndpoint (response http.ResponseWriter, request *http.Request){
	//user will aniticipate json
	response.Header().Add("content-type","application/json")
	var user Users
	json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("goLangTrial").Collection("user")
	ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
	result , _ := collection.InsertOne(ctx ,user)
	json.NewEncoder(response).Encode(result)
}

func main(){
	//run by go run main.go
	fmt.Println("Starting the application")
	//initialise router
	var ( 
		
		mongoURL = "mongodb+srv://rohanGolang:pliMmfICvzXjsXVy@cluster0.3kcv6.mongodb.net/goLangTrial?retryWrites=true&w=majority"
	)
	//initialise mongo client with options
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))

	//connect the mongo client to mongodb server
	ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	err = client.Connect(ctx)

	//ping mongodb
	// ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	if err = client.Ping(ctx,readpref.Primary()); err != nil {
		fmt.Println("couldnt ping to mongodb service:  %v",err)
		return
	}

	fmt.Println("connected to no sql database",mongoURL)


	// client, err := mongo.Connect(ctx,"mongodb+srv://rohanGolang:SZ6tp1ZODpp17I0F@cluster0.3kcv6.mongodb.net/goLangTrial?retryWrites=true&w=majority")
	//establish mongodb connection, for that define timeout, define context and pass that to connection

	router := mux.NewRouter()
		//arrange our route 
	router.HandleFunc("/users",CreateUsersEndpoint).Methods("POST")

	http.ListenAndServe(":12345", router)



}