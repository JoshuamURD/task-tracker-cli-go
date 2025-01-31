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

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening json file: %v", err)
	}
	defer file.Close()
	file.Truncate(0)
	file.Seek(0, 0)

	if err = json.NewEncoder(file).Encode(*t); err != nil {
		return fmt.Errorf("error writing task to JSON: %v", err)
	}

	fmt.Printf("Task added successfully (ID: %v)", len(*t))
	return nil
}

// ListTasks writes thes tasks to an io.Writer
func (t Tasks) ListTasks(kind string, w io.Writer) error {
	for _, task := range t {
		_, err := fmt.Fprintf(w, "ID: %v, Description: %s, Status: %s, CreatedAt: %s\n",
			task.ID, task.Description, task.Status, task.Createdate.Format(time.RFC3339))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteTask takes a taskID and deletes that task
func (t *Tasks) DeleteTask(taskID int, path string) error {
	for i, task := range *t {
		if task.ID == taskID {
			*t = append((*t)[:i], (*t)[i+1:]...)
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("error opening json file: %v", err)
			}
			defer file.Close()
			file.Truncate(0)
			file.Seek(0, 0)

			if err = json.NewEncoder(file).Encode(*t); err != nil {
				return fmt.Errorf("error writing task to JSON: %v", err)
			}
			return nil
		}
	}
	return fmt.Errorf("task with ID %v not found", taskID)
}
