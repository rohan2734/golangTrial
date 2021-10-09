package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//format to get inputs from console and write outputs to console

type Users struct{
  //ID, name, email ,password
//   ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitmpty"`
  Name string `json:"name,omitempty" bson:"name,omitempty"`
  Email string `json:"email,omitempty" bson:"email,omitempty"`
  Password string `json:"password,omitempty" bson:"password,omitempty"`
}

type Posts  struct{
	Caption string `json:"caption,omitempty" bson:"name,omitempty"`
	ImageURL string `json:"imageURL,omitempty" bson:"email,omitempty"`
	PostedTimestamp string `json:"postedTimestamp,omitempty" bson:"password,omitempty"`
}

func main(){
	//run by go run main.go
	fmt.Println("Starting the application")
	//initialise router
	var ( 
		client *mongo.Client
		mongoURL = "mongodb+srv://rohanGolang:pliMmfICvzXjsXVy@cluster0.3kcv6.mongodb.net/goLangTrial?retryWrites=true&w=majority"
	)
	//initialise mongo client with options
	client, _ = mongo.NewClient(options.Client().ApplyURI(mongoURL))

	//connect the mongo client to mongodb server
	ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	client.Connect(ctx)

	//ping mongodb
	// ctx, _ := context.WithTimeout(context.Background(),10*time.Second)
	if err := client.Ping(ctx,readpref.Primary()); err != nil {
		fmt.Println("couldnt ping to mongodb service:  %v",err)
		return
	}

	fmt.Println("connected to no sql database",mongoURL)


	// client, err := mongo.Connect(ctx,"mongodb+srv://rohanGolang:SZ6tp1ZODpp17I0F@cluster0.3kcv6.mongodb.net/goLangTrial?retryWrites=true&w=majority")
	//establish mongodb connection, for that define timeout, define context and pass that to connection

	router := mux.NewRouter()
		//arrange our route 
	// router.HandleFunc("/users",).Methods("POST")
    router.HandleFunc("/users",func  (response http.ResponseWriter, request *http.Request){
		//user will aniticipate json
		response.Header().Add("content-type","application/json")
		var user Users

		json.NewDecoder(request.Body).Decode(&user)
		// collection := client.Database("goLangTrial").Collection("user")
		// pass := user.Password
		// pass = "1"
		// user.Password = pass
		ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
		result , _ := client.Database("goLangTrial").Collection("user").InsertOne(ctx ,user)
		json.NewEncoder(response).Encode(result)
	}).Methods("POST")

	router.HandleFunc("/users/{ID}", func(response http.ResponseWriter, request *http.Request){
		response.Header().Add("content-type","application/json")
		// var userSlice []Users
		// collection = client.Database("goLangTrial").Collection("user")
		// ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
		// cursor, err := collection.Find(ctx , bson.M{}) 
		// if err != nil {
		// 	response.WriteHeader(http.StatusInternalServerError)
		// 	response.Write([]byte(`{"message" : "` + err.Error() + `"}`))
		// 	return
		// }
		// defer cursor.Close(ctx)

		// for cursor.Next(ctx){
		// 	var user Users
		// 	cursor.Decode(&user)
		// 	userSlice = append(userSlice,user)
		// }
		// var user bson.M
		params := mux.Vars(request)
		ID := params["ID"]

		objID , _ := primitive.ObjectIDFromHex(ID)

		fmt.Println(ID)
		usersCollection := client.Database("goLangTrial").Collection("user")
		ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
		filterCursor, err := usersCollection.Find(ctx , bson.M{"_id": objID }) 
		if err != nil {
			log.Fatal(err)
		}

		var usersFiltered []bson.M
		if err = filterCursor.All(ctx, &usersFiltered); err != nil {
			log.Fatal(err)
		}
		fmt.Println(usersFiltered)
		json.NewEncoder(response).Encode(usersFiltered)

		// result :=  client.Database("goLangTrial").Collection("user").FindOne(context.Background(),bson.M{"_id":ID})
		// var user Users
		// result.Decode(user)
		// fmt.Println(result)
		// json.NewEncoder(response).Encode(result)
	

	})
	
	router.HandleFunc("/posts",func(response http.ResponseWriter, request *http.Request) {
		response.Header().Add("content-type","application/json")
		var post Posts

		json.NewDecoder(request.Body).Decode(&post)

		ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
		result , _ := client.Database("goLangTrial").Collection("post").InsertOne(ctx ,post)
		json.NewEncoder(response).Encode(result)
	})
	http.ListenAndServe(":12345", router)



}