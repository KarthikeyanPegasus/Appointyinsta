package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type post struct {
	Caption string `json:"caption"`
	Id      string `json:"id"`
	ImgUrl  string `json:"imageUrl"`
	PtTime  string `json:"posttimestamp"`
	Userid  string `json:"userid"`
}

type Users []user

type UserHandler struct {
	sync.Mutex
	users Users
}
type PostHandler struct {
	sync.Mutex
	posts Posts
}

var (
	mongoURI = "mongodb://localhost:27017"
)

type Posts []post

///////// Serve HTTP

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uh.getusers(w, r)
		return
	case "POST":
		uh.postusers(w, r)
		return
	default:
		respondwithERROR(w, http.StatusMethodNotAllowed, "Method not Available")
	}

}

///// Get user /contains get /users, /users/<id>
func (uh *UserHandler) getusers(w http.ResponseWriter, r *http.Request) {
	defer uh.Unlock()
	uh.Lock()
	id, err := idFromURL(r)
	if err != nil {

		client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
		Appointydb := client.Database("AppointyInsta")

		usercollections := Appointydb.Collection("users")
		defer usercollections.Drop(ctx)
		cursor, err := usercollections.Find(ctx, bson.M{})
		if err != nil {
			respondwithERROR(w, http.StatusNotModified, err.Error())
		}

		var postings []bson.M
		if err = cursor.All(ctx, &postings); err != nil {
			respondwithERROR(w, http.StatusNotModified, err.Error())
		}

		respondwithJSON(w, http.StatusOK, postings)
		return

	}
	if id >= len(uh.users) || id < 0 {
		respondwithERROR(w, http.StatusNotFound, "not found")
		return
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}
	defer client.Disconnect(ctx)
	Appointydb := client.Database("AppointyInsta")

	usercollections := Appointydb.Collection("users")
	defer usercollections.Drop(ctx)
	var postings []bson.M
	Filter := bson.M{"ID": strconv.Itoa(id)}
	fCursor, err := usercollections.Find(ctx, Filter)
	if err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}
	if err = fCursor.All(ctx, &postings); err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}

	respondwithJSON(w, http.StatusOK, postings)

}

// Post new Users  into db /users using post method

func (uh *UserHandler) postusers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondwithERROR(w, http.StatusInternalServerError, err.Error())
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		respondwithERROR(w, http.StatusUnsupportedMediaType, "content type Application/json is required")
		return
	}
	var userses user
	err = json.Unmarshal(body, &userses)
	userses.Password = createHash(userses.Password)
	if err != nil {
		respondwithERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}
	defer client.Disconnect(ctx)
	Appointydb := client.Database("AppointyInsta")

	usercollections := Appointydb.Collection("users")
	defer usercollections.Drop(ctx)
	result, err := usercollections.InsertOne(ctx, bson.D{
		{Key: "ID", Value: userses.ID},
		{Key: "Name", Value: userses.Name},
		{Key: "Email", Value: userses.Email},
		{Key: "Password", Value: userses.Password},
	})
	if err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}

	defer uh.Unlock()
	uh.Lock()
	uh.users = append(uh.users, userses)
	respondwithJSON(w, http.StatusCreated, result)
}

//Encryption for the password

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

////////////// serve HTTP for posts
func (ph *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ph.getposts(w, r)
		return
	case "POST":
		ph.postposts(w, r)
		return
	default:
		respondwithERROR(w, http.StatusMethodNotAllowed, "Method not Available")
	}

}

func idwiththreeparts(r *http.Request) (int, bool, error) {
	var isconsistof bool = false
	parts := strings.Split(r.URL.String(), "/")

	if len(parts) == 4 {
		isconsistof = true
	}
	if len(parts) <= 2 {
		return 0, isconsistof, errors.New("not found")
	}
	id, err := strconv.Atoi(parts[len(parts)-1])

	if err != nil {
		return 0, isconsistof, errors.New("not found")
	}
	return id, isconsistof, nil
}

func idFromURL(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		return 0, errors.New("not found")
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, errors.New("not found")
	}
	return id, nil
}

//////////// gets posts // contains /posts ,/posts/<id>

func (ph *PostHandler) getposts(w http.ResponseWriter, r *http.Request) {
	defer ph.Unlock()
	ph.Lock()
	id, isconsistof, err := idwiththreeparts(r)
	// log.Fatal(id, isconsistof, err)
	if err != nil {

		client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
		Appointydb := client.Database("AppointyInsta")

		usercollections := Appointydb.Collection("posts")
		defer usercollections.Drop(ctx)
		cursor, err := usercollections.Find(ctx, bson.M{})
		if err != nil {
			respondwithERROR(w, http.StatusNotModified, err.Error())
		}

		var postings []bson.M
		if err = cursor.All(ctx, &postings); err != nil {
			respondwithERROR(w, http.StatusNotModified, err.Error())
		}

		respondwithJSON(w, http.StatusOK, postings)
		return
	}
	if id >= len(ph.posts) || id < 0 {
		respondwithERROR(w, http.StatusNotFound, "not found")
		return
	}
	if isconsistof {
		// userPosts := make([]post, 0, len(ph.posts))
		// // log.Fatal(userPosts)
		// for _, v := range ph.posts {

		// 	if v.Userid == strconv.Itoa(id) {

		// 		userPosts = append(userPosts, v)
		// 	}
		// }

		client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)
		Appointydb := client.Database("AppointyInsta")

		usercollections := Appointydb.Collection("posts")
		defer usercollections.Drop(ctx)
		var postings []bson.M
		Filter := bson.M{"UserId": strconv.Itoa(id)}
		fCursor, err := usercollections.Find(ctx, Filter)
		if err != nil {
			respondwithERROR(w, http.StatusNotModified, err.Error())
		}
		if err = fCursor.All(ctx, &postings); err != nil {
			respondwithERROR(w, http.StatusNotModified, err.Error())
		}

		respondwithJSON(w, http.StatusOK, postings)
		return

	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	Appointydb := client.Database("AppointyInsta")

	usercollections := Appointydb.Collection("posts")
	defer usercollections.Drop(ctx)
	var postings []bson.M
	Filter := bson.M{"ID": strconv.Itoa(id)}
	fCursor, err := usercollections.Find(ctx, Filter)
	if err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}
	if err = fCursor.All(ctx, &postings); err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}
	//if v.Id == strconv.Itoa(id) {
	respondwithJSON(w, http.StatusOK, postings)
	//return
	//}

}

//////////posting into posts Method POST
func (ph *PostHandler) postposts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondwithERROR(w, http.StatusInternalServerError, err.Error())
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		respondwithERROR(w, http.StatusUnsupportedMediaType, "content type Application/json is required")
		return
	}
	var postses post
	err = json.Unmarshal(body, &postses)
	if err != nil {
		respondwithERROR(w, http.StatusBadRequest, err.Error())
		return
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	Appointydb := client.Database("AppointyInsta")

	usercollections := Appointydb.Collection("posts")
	defer usercollections.Drop(ctx)
	result, err := usercollections.InsertOne(ctx, bson.D{
		{Key: "ID", Value: postses.Id},
		{Key: "Caption", Value: postses.Caption},
		{Key: "ImageURL", Value: postses.ImgUrl},
		{Key: "PostTimeStamp", Value: postses.PtTime},
		{Key: "UserId", Value: postses.Userid},
	})
	if err != nil {
		respondwithERROR(w, http.StatusNotModified, err.Error())
	}

	defer ph.Unlock()
	ph.Lock()
	ph.posts = append(ph.posts, postses)
	respondwithJSON(w, http.StatusCreated, result)
}

////////// reusable codes
func respondwithERROR(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"ERROR": msg})
}

func respondwithJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "Application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func newsUserhandler() *UserHandler {

	return &UserHandler{

		users: Users{
			user{ID: "1", Name: "Admin", Email: "someone@gmail.com", Password: "thisisthepassword"},
			user{ID: "2", Name: "Admin2", Email: "someone2@gmail.com", Password: "thisisthepassword2"},
		},
	}
}

func newspostHandler() *PostHandler {
	return &PostHandler{
		posts: Posts{
			post{Id: "1", Caption: "thisisthecaption", Userid: "1", ImgUrl: "thisistheiamgeurl", PtTime: "thistime"},
			post{Id: "2", Caption: "thisisthecaption2", Userid: "1", ImgUrl: "thisistheiamgeurl2", PtTime: "thistime2"},
			post{Id: "3", Caption: "thisisthecaption3", Userid: "1", ImgUrl: "thisistheiamgeurl3", PtTime: "thistime3"},
			post{Id: "4", Caption: "thisisthecaption4", Userid: "2", ImgUrl: "thisistheiamgeurl4", PtTime: "thistime4"},
		},
	}
}

///////////// Main Functions
func main() {

	port := ":8080"
	uh := newsUserhandler()
	ph := newspostHandler()
	http.Handle("/users", uh)
	http.Handle("/users/", uh)
	http.Handle("/posts", ph)
	http.Handle("/posts/", ph)
	http.Handle("/posts/users/", ph)

	log.Fatal(http.ListenAndServe(port, nil))

}
