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
	status        bool
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
				status:  false,
			}
			Todolists[i].BulletPoints = append(Todolists[i].BulletPoints, newBulletPoint)
			fmt.Printf("Bullet point with ID %d content: %v status: %v ,added to todo with ID %d Title %v  \n", Todolists[i].bulletPointIDCounter, content, newBulletPoint.status, todoID, Todolists[i].Title)
			return
		}
	}
	fmt.Printf("Todo with ID %d not found\n", todoID)
}

// delete a bullet point from a todo list
func DeleteBulletPoint(todoID int, bulletPointID int) {
	for i := range Todolists {
		if Todolists[i].ID == todoID {
			for j := range Todolists[i].BulletPoints {
				if Todolists[i].BulletPoints[j].ID == bulletPointID {
					Todolists[i].BulletPoints = append(Todolists[i].BulletPoints[:j], Todolists[i].BulletPoints[j+1:]...)
					fmt.Printf("Bullet point with ID %d deleted from todo with ID %d Title %v\n", bulletPointID, todoID, Todolists[i].Title)
					return
				}
			}
			fmt.Printf("Bullet point with ID %d not found in todo with ID %d\n", bulletPointID, todoID)
			return
		}
	}
	fmt.Printf("Todo with ID %d not found\n", todoID)

}

// mark a bullet point as completed
func BulletPointCompleted(todoID int, bulletPointID int) {
	for i := range Todolists {
		if Todolists[i].ID == todoID {
			for j := range Todolists[i].BulletPoints {
				if Todolists[i].BulletPoints[j].ID == bulletPointID {
					t := time.Now()
					Todolists[i].BulletPoints[j].CompletedTime = &t
					Todolists[i].BulletPoints[j].status = true
					fmt.Printf("Bullet point with ID %d marked as completed in todo with ID %d Title %v\n", bulletPointID, todoID, Todolists[i].Title)
					return
				}
			}
			fmt.Printf("Bullet point with ID %d not found in todo with ID %d\n", bulletPointID, todoID)
			return
		}
		fmt.Printf("Todo with ID %d not found\n", todoID)
	}
}
func main() {
	NewTodo("First Todo")
	AddBulletPoint(1, "11111111111111111")
	AddBulletPoint(1, "222222222222")
	AddBulletPoint(1, "333333333333")
	NewTodo("Second Todo")
	AddBulletPoint(2, "44444444444")
	AddBulletPoint(2, "555555555")
	fmt.Println("-----------------")
	for _, todo := range Todolists {
		fmt.Printf("Todo ID %d Title %v\n", todo.ID, todo.Title)
		for _, bulletPoint := range todo.BulletPoints {
			fmt.Printf("ID: %d - Content: %v status: %v", bulletPoint.ID, bulletPoint.Content, bulletPoint.status)
			if bulletPoint.CompletedTime != nil {
				fmt.Printf(" Completed at %v\n", bulletPoint.CompletedTime)
			} else {
				fmt.Println()
			}
		}
	}
	fmt.Println("-----------------")
	DeleteBulletPoint(1, 1)
	BulletPointCompleted(1, 2)
	for _, todo := range Todolists {
		fmt.Printf("Todo ID %d Title %v\n", todo.ID, todo.Title)
		for _, bulletPoint := range todo.BulletPoints {
			fmt.Printf("ID: %d - Content: %v status: %v", bulletPoint.ID, bulletPoint.Content, bulletPoint.status)
			if bulletPoint.CompletedTime != nil {
				fmt.Printf(" Completed at %v\n", bulletPoint.CompletedTime)
			} else {
				fmt.Println()
			}
		}
	}
}
