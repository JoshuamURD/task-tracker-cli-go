package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// NowFunc acts as a DI for time.Now for testing purposes
var NowFunc = time.Now

// tasks is an in memory cache of tasks loaded from json
type Task struct {
	ID          int       `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Createdate  time.Time `json:"createdat"`
	Updatedate  time.Time `json:"updatedat"`
}

// Tasks stores a slice of the in memory tasks
type Tasks []Task

// NewTasks initialises an empty Tasks slice
func NewTasks() Tasks {
	return Tasks{}
}

// LoadTasks loads tasks from the json file, or creates a json file to store tasks if it does not exist
func (t *Tasks) LoadTasks(path string) error {

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		log.Println("No tasks were loaded.")
		return nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(t)
	if err != nil {
		return err
	}
	log.Printf("%v tasks were loaded in ", len(*t))

	return nil
}

// AddTask adds a task
func (t *Tasks) AddTask(task, path string) error {
	now := NowFunc()
	*t = append(*t, Task{
		ID:          len(*t) + 1,
		Description: task,
		Status:      "todo",
		Createdate:  now,
		Updatedate:  now,
	})

	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("An error occurred opening the json file: %v", err)
	}
	defer file.Close()

	json.NewEncoder(file).Encode(*t[len(*t)])

	fmt.Printf("Task added successfully (ID: %v)", len(*t))
	return nil
}

// ListTasks writes thes tasks to an io.Writer
func (t Tasks) ListTasks(kind string, w io.Writer) error {
	for _, task := range t {
		_, err := fmt.Fprintf(w, "ID: %v, Description: %s, Status: %s, CreatedAt: %s\n",
			task.ID, task.Description, task.Status, task.Createdate)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

}
