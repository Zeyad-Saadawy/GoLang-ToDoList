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
	for i := range Todolists {
		if Todolists[i].ID == todoID {
			Todolists[i].bulletPointIDCounter++
			newBulletPoint := BulletPoint{
				ID:      Todolists[i].bulletPointIDCounter,
				Content: content,
			}
			Todolists[i].BulletPoints = append(Todolists[i].BulletPoints, newBulletPoint)
			fmt.Printf("Bullet point with ID %d content: %v added to todo with ID %d Title %v  \n", Todolists[i].bulletPointIDCounter, content, todoID, Todolists[i].Title)
			return
		}
	}
	fmt.Printf("Todo with ID %d not found\n", todoID)
}

func main() {
	NewTodo("First Todo")
	// NewTodo("Second Todo")
	AddBulletPoint(1, "11111111111111111")
	AddBulletPoint(1, "222222222222")
	for _, todo := range Todolists {
		fmt.Printf("Todo ID %d Title %v\n", todo.ID, todo.Title)
		for _, bulletPoint := range todo.BulletPoints {
			fmt.Printf("ID %d - Content %v\n", bulletPoint.ID, bulletPoint.Content)
		}
	}
}
