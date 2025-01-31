package task

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

// Count is a counter for the number of tasks
var count = 0

// ID is a type alias for int
type ID = int

// tasks is an in memory cache of tasks loaded from json
type Task struct {
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Createdate  time.Time `json:"createdate"`
	Updatedate  time.Time `json:"updatedate"`
}

// Tasks stores a map of tasks with ID as the key
type Tasks map[ID]Task

// NewTasks creates a new Tasks map
func NewTasks() Tasks {
	return make(Tasks)
}

// LoadTasks loads tasks from the json file, or creates a json file to store tasks if it does not exist
func LoadTasks(path string, tasks Tasks) error {
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
	err = decoder.Decode(&tasks)
	if err != nil {
		return err
	}
	
	// Update count to the highest ID to maintain ID sequence
	for id := range tasks {
		if id > count {
			count = id
		}
	}
	
	log.Printf("%v tasks were loaded", len(tasks))
	return nil
}

// AddTask adds a task
func AddTask(task, path string, tasks Tasks) error {
	now := NowFunc()
	count++
	tasks[count] = Task{
		Description: task,
		Status:      "Pending",
		Createdate:  now,
		Updatedate:  now,
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening json file: %v", err)
	}
	defer file.Close()
	file.Truncate(0)
	file.Seek(0, 0)

	if err = json.NewEncoder(file).Encode(tasks); err != nil {
		return fmt.Errorf("error writing task to JSON: %v", err)
	}

	fmt.Printf("Task added successfully (ID: %v)", count)
	return nil
}

// ListTasks writes the tasks to an io.Writer
func (t Tasks) ListTasks(kind string, w io.Writer) error {
	for id, task := range t {
		_, err := fmt.Fprintf(w, "ID: %v, Description: %s, Status: %s, CreatedAt: %s\n",
			id, task.Description, task.Status, task.Createdate.Format(time.RFC3339))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteTask takes a taskID and deletes that task
func (t Tasks) DeleteTask(taskID int, path string) error {
	if _, exists := t[taskID]; !exists {
		return fmt.Errorf("task with ID %v not found", taskID)
	}

	delete(t, taskID)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening json file: %v", err)
	}
	defer file.Close()
	file.Truncate(0)
	file.Seek(0, 0)

	if err = json.NewEncoder(file).Encode(t); err != nil {
		return fmt.Errorf("error writing task to JSON: %v", err)
	}
	return nil
}
