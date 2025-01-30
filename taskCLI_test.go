package main

import (
	"os"
	"testing"
)

func TestTasks(t *testing.T) {
	tasks := NewTasks()
	t.Run("Loading tasks", func(t *testing.T) {
		_ = os.Remove("data.json")
		if err := tasks.LoadTasks("data.json"); err != nil {
			t.Errorf("Got an error %v", err)
		}

	})
	t.Run("Adding a task", func(t *testing.T) {
		tasks.AddTask("test")
		got := len(tasks)
		want := 1
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		tasks.AddTask("Another task")
		got = len(tasks)
		want = 2

		if got != want {
			t.Errorf("Have %v tasks, want %v tasks")
		}

	})
	t.Run("Deleting a task", func(t *testing.T) {
		tasks.RemoveTask(2)

		got := tasks[0]
		want := Task{}

	})
	// t.Run("Listing tasks", func(t *testing.T) {
	// 	buf := &bytes.Buffer{}

	// 	tasks.ListTasks("All", buf)
	// 	task := tasks[0]

	// 	got := buf.String()
	// 	want := fmt.Sprintf("ID: %v, Description: %s, Status: %s, CreatedAt: %s\n", task.ID, task.Description, task.Status, task.CreatedAt)

	// 	if got != want {
	// 		t.Errorf("got %v want %v", got, want)
	// 	}
	// })
	t.Run("Making sure added task is in json", func(t *testing.T) {
		tasks.LoadTasks("data.json")

	})
}
