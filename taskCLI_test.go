package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

// Setup a temporary JSON file for each test.
func setupTempFile(t *testing.T) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestTasks(t *testing.T) {
	// Override time for consistent testing
	fixedTime := time.Date(2025, 1, 30, 12, 0, 0, 0, time.UTC)
	NowFunc = func() time.Time { return fixedTime }
	defer func() { NowFunc = time.Now }()

	// Use a temporary JSON file
	tempFile := setupTempFile(t)
	defer os.Remove(tempFile) // Cleanup after test

	tasks := NewTasks()

	t.Run("Loading tasks from empty file", func(t *testing.T) {
		err := tasks.LoadTasks(tempFile)
		if err != nil {
			t.Fatalf("Failed to load tasks: %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("Expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("Adding a task", func(t *testing.T) {
		tasks.AddTask("Test task")

		if len(tasks) != 1 {
			t.Errorf("Expected 1 task, got %d", len(tasks))
		}
		task := tasks[0]
		if task.Description != "Test task" {
			t.Errorf("Expected task description 'Test task', got '%s'", task.Description)
		}
		if !task.Createdate.Equal(fixedTime) {
			t.Errorf("Expected created time %v, got %v", fixedTime, task.Createdate)
		}
	})

	t.Run("Listing tasks", func(t *testing.T) {
		buf := &bytes.Buffer{}
		tasks.ListTasks("All", buf)

		task := tasks[0]
		want := fmt.Sprintf("ID: %d, Description: %s, Status: %s, CreatedAt: %s\n",
			task.ID, task.Description, task.Status, task.Createdate.Format(time.RFC3339))

		if buf.String() != want {
			t.Errorf("ListTasks output mismatch:\nGot: %s\nWant: %s", buf.String(), want)
		}
	})

	t.Run("Saving and Loading tasks", func(t *testing.T) {
		// Save tasks
		file, err := os.Create(tempFile)
		if err != nil {
			t.Fatalf("Failed to open temp file for writing: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		err = encoder.Encode(tasks)
		if err != nil {
			t.Fatalf("Failed to encode tasks: %v", err)
		}

		// Reload tasks from file
		loadedTasks := NewTasks()
		err = loadedTasks.LoadTasks(tempFile)
		if err != nil {
			t.Fatalf("Failed to load tasks from file: %v", err)
		}

		// Compare stored vs loaded tasks
		if len(loadedTasks) != len(tasks) {
			t.Errorf("Expected %d tasks after reload, got %d", len(tasks), len(loadedTasks))
		}
		if loadedTasks[0].Description != tasks[0].Description {
			t.Errorf("Mismatch in task description: Got %s, Want %s",
				loadedTasks[0].Description, tasks[0].Description)
		}
	})
}
