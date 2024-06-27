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

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName string = "localhost:27017"
	dbName   string = "demo_todo"
	collName string = "todo"
	port     string = ":9000"
)

type (
	todoModel struct {
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"created_at"`
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
	sess, err := mgo.Dial(hostName)
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tmpl"}, nil)
	checkErr(err)
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	todos := []todoModel{}
	if err := db.C(collName).Find(bson.M{}).All(&todos); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{"message": "failed to fetch todos", "error": err})
		return
	}
	todolist := []todo{}
	for _, t := range todos {
		todolist = append(todolist, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}
	rnd.JSON(w, http.StatusOK, renderer.M{"data": todolist})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil { // take the body and put it in t
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "invalid request payload"})
		return
	}
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "title is required"})
		return
	}

	tm := todoModel{
		ID:        bson.NewObjectId(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	if err := db.C(collName).Insert(&tm); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{"message": "failed to create todo", "error": err})
		return
	}
	rnd.JSON(w, http.StatusCreated, renderer.M{"message": "todo created successfully", "todo_ID": tm.ID.Hex()})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "invalid todo ID"})
		return
	}

	if err := db.C(collName).RemoveId(bson.ObjectIdHex(id)); err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "failed to delete todo", "error": err})
		return
	}
	rnd.JSON(w, http.StatusOK, renderer.M{"message": "todo deleted successfully"})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if !bson.IsObjectIdHex(id) {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "invalid todo ID"})
		return
	}

	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil { // take the body and put it in t
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "invalid request payload"})
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{"message": "title is required"})
		return
	}

	if err := db.C(collName).
		Update(
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"title": t.Title, "completed": t.Completed},
		); err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "failed to update todo", "error": err})
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
	srv.Shutdown(ctx)
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
		log.Fatal(err)
	}
}
