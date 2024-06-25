package main

import (
	"fmt"
	"time"
)

// global variables to generate unique IDs for todos
var todoIDCounter int

// global variable to keep track of all the todolists
var Todolists []Todolist

type Todolist struct {
	ID                   int
	Title                string
	BulletPoints         []BulletPoint
	bulletPointIDCounter int
}

type BulletPoint struct {
	ID            int
	Content       string
	CompletedTime *time.Time
}

// Create a new Todo
func NewTodo(title string) {
	todoIDCounter++
	newtodo := Todolist{
		ID:           todoIDCounter,
		Title:        title,
		BulletPoints: []BulletPoint{},
	}
	Todolists = append(Todolists, newtodo)
	fmt.Printf("Todo with ID %d and title %v created\n", todoIDCounter, title)
}

// add a bullet point to a todo list
func AddBulletPoint(todoID int, content string) {
	for i, todo := range Todolists {
		if todo.ID == todoID {
			todo.bulletPointIDCounter++
			newBulletPoint := BulletPoint{
				ID:      todo.bulletPointIDCounter,
				Content: content,
			}
			Todolists[i].BulletPoints = append(Todolists[i].BulletPoints, newBulletPoint)
			fmt.Printf("Bullet point with ID %d added to todo with ID %d Title %v \n content: %v", todo.bulletPointIDCounter, todoID, todo.Title, content)
			return
		}
	}
	fmt.Printf("Todo with ID %d not found\n", todoID)
}

func main() {

}
