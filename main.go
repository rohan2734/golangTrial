package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"os"

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
	PostedTimestamp primitive.Timestamp `json:"lastUpdate" bson:"lastUpdate"`
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
	
		//create user
    router.HandleFunc("/users",func  (response http.ResponseWriter, request *http.Request){
		//user will aniticipate json
		response.Header().Add("content-type","application/json")
		var user Users

		json.NewDecoder(request.Body).Decode(&user)
		// collection := client.Database("goLangTrial").Collection("user")
		// pass := user.Password
		// pass = "1"
		// user.Password = pass


		var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

		plainPass := []byte(user.Password)
		// aes encryption string
		key_text := "astaxie12798akljzmknm.ahkjkljl;k"

		c, err := aes.NewCipher([]byte(key_text))
		if err != nil {
			fmt.Println("Error: NewCipher(%d bytes) = %s", len(key_text), err)
			os.Exit(-1)
		}

		// Encrypted string
		cfb := cipher.NewCFBEncrypter(c, commonIV)
		cipherPass := make([]byte, len(plainPass))
		cfb.XORKeyStream(cipherPass, plainPass)
		// fmt.Printf("%s=>%x\n", plainPass, ciphertext)
		user.Password = string( cipherPass )
		fmt.Println("%s",user.Password)


		ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
		result , _ := client.Database("goLangTrial").Collection("user").InsertOne(ctx ,user)
		json.NewEncoder(response).Encode(result)
	}).Methods("POST")

	//get user by ID
	router.HandleFunc("/users/{ID}", func(response http.ResponseWriter, request *http.Request){
		response.Header().Add("content-type","application/json")
		
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


	})
	
	//createpost
	router.HandleFunc("/posts",func(response http.ResponseWriter, request *http.Request) {
		response.Header().Add("content-type","application/json")
		var post Posts

		json.NewDecoder(request.Body).Decode(&post)

		ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
		result , _ := client.Database("goLangTrial").Collection("post").InsertOne(ctx ,post)
		json.NewEncoder(response).Encode(result)
	})
	
	//get post by post ID
	router.HandleFunc("/posts/{ID}", func(response http.ResponseWriter, request *http.Request){
		response.Header().Add("content-type","application/json")
		
		params := mux.Vars(request)
		ID := params["ID"]

		objID , _ := primitive.ObjectIDFromHex(ID)

		fmt.Println(ID)
		postsCollection := client.Database("goLangTrial").Collection("post")
		ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
		filterCursor, err := postsCollection.Find(ctx , bson.M{"_id": objID }) 
		if err != nil {
			log.Fatal(err)
		}

		var postsFiltered []bson.M
		if err = filterCursor.All(ctx, &postsFiltered); err != nil {
			log.Fatal(err)
		}
		fmt.Println(postsFiltered)
		json.NewEncoder(response).Encode(postsFiltered)


	})

	//list all posts of a user
	router.HandleFunc("/posts/{userID}/{pageNumber}", func(response http.ResponseWriter, request *http.Request){
		response.Header().Add("content-type","application/json")
		
		params := mux.Vars(request)
		ID := params["ID"]

		objID , _ := primitive.ObjectIDFromHex(ID)

		fmt.Println(ID)
		
		postsCollection := client.Database("goLangTrial").Collection("post")
		ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
		filterCursor, err := postsCollection.Find(ctx , bson.M{"_id": objID }) 
		if err != nil {
			log.Fatal(err)
		}

		var postsFiltered []bson.M
		if err = filterCursor.All(ctx, &postsFiltered); err != nil {
			log.Fatal(err)
		}
		fmt.Println(postsFiltered)
		json.NewEncoder(response).Encode(postsFiltered)


	})
	http.ListenAndServe(":12345", router)



}