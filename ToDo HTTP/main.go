package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thedevsaddam/renderer"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var rnd *renderer.Render
var client *mongo.Client
var db *mongo.Database

const (
	hostName string = "localhost:27017"
	dbName   string = "demo_todo"
	collName string = "todo"
	port     string = ":9000"
)

type (
	todoModel struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		Title     string             `bson:"title"`
		Completed bool               `bson:"completed"`
		CreatedAt time.Time          `bson:"created_at"`
	}
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func init() {
	rnd = renderer.New()
	// Set up MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://" + hostName)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the MongoDB server to check if the connection is established
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB server: %v", err)
	}

	db = client.Database(dbName)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	todos := []todoModel{}
	cursor, err := db.Collection(collName).Find(context.Background(), bson.D{})
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{"message": "failed to fetch todos", "error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var t todoModel
		if err := cursor.Decode(&t); err != nil {
			rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "error decoding todo", "error": err.Error()})
			return
		}
		todos = append(todos, t)
	}

	todolist := make([]todo, len(todos))
	for i, t := range todos {
		todolist[i] = todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		}
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"data": todolist})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "invalid request payload", "error": err.Error()})
		return
	}
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "title is required"})
		return
	}

	tm := todoModel{
		ID:        primitive.NewObjectID(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	_, err := db.Collection(collName).InsertOne(context.Background(), tm)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{"message": "failed to create todo", "error": err.Error()})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{"message": "todo created successfully", "todo_ID": tm.ID.Hex()})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "invalid todo ID", "error": err.Error()})
		return
	}

	res, err := db.Collection(collName).DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "failed to delete todo", "error": err.Error()})
		return
	}
	if res.DeletedCount == 0 {
		rnd.JSON(w, http.StatusNotFound, renderer.M{"message": "todo not found"})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"message": "todo deleted successfully"})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "invalid todo ID", "error": err.Error()})
		return
	}

	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "invalid request payload", "error": err.Error()})
		return
	}
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "title is required"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":     t.Title,
			"completed": t.Completed,
		},
	}

	res, err := db.Collection(collName).UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "failed to update todo", "error": err.Error()})
		return
	}
	if res.ModifiedCount == 0 {
		rnd.JSON(w, http.StatusNotFound, renderer.M{"message": "todo not found"})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"message": "todo updated successfully"})
}

func main() {
	stopchan := make(chan os.Signal)
	signal.Notify(stopchan, os.Interrupt)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())
	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("ListenAndServe: %s\n", err)
		}
	}()

	<-stopchan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s", err)
	}
	log.Println("Server gracefully stopped!")
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
