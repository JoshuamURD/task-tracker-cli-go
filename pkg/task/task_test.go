package task

import (
	"bytes"
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

// setupTasks returns a Tasks object and the path to the temporary file
func setupTasks(t *testing.T) (Tasks, string) {
	t.Helper()
	tasks := NewTasks()
	tempFile := setupTempFile(t)
	t.Cleanup(func() { os.Remove(tempFile) })
	return tasks, tempFile
}

func TestLoadTasks(t *testing.T) {
	t.Run("Loading tasks from empty file", func(t *testing.T) {
		tasks, tempFile := setupTasks(t)
		
		err := LoadTasks(tempFile, tasks)
		if err != nil {
			t.Fatalf("Failed to load tasks: %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("Expected 0 tasks, got %d", len(tasks))
		}
	})
}

func TestListTasks(t *testing.T) {
	var buf bytes.Buffer
	now := time.Now()
	tests := []struct {
		name  string
		tasks Tasks
		kind  string
		want  string
	}{
		{
			name: "List tasks",
			tasks: Tasks{
				1: Task{
					Description: "Task 1",
					Status:     "Pending",
					Createdate: now,
				},
			},
			kind: "all",
			want: "ID: 1, Description: Task 1, Status: Pending, CreatedAt: " + now.Format(time.RFC3339) + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf.Reset()
			err := test.tasks.ListTasks(test.kind, &buf)
			if err != nil {
				t.Fatalf("ListTasks() error = %v", err)
			}
			if buf.String() != test.want {
				t.Errorf("ListTasks() = %q, want %q", buf.String(), test.want)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	tasks, tempFile := setupTasks(t)

	err := AddTask("Test task", tempFile, tasks)
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	task, exists := tasks[1]
	if !exists {
		t.Fatal("Task with ID 1 not found")
	}
	if task.Description != "Test task" {
		t.Errorf("Expected task description 'Test task', got '%s'", task.Description)
	}
	if task.Status != "Pending" {
		t.Errorf("Expected status 'Pending', got '%s'", task.Status)
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name       string
		setupTasks Tasks
		taskID     int
		wantErr    bool
		wantLen    int
	}{
		{
			name: "Delete existing task",
			setupTasks: Tasks{
				1: Task{
					Description: "Task 1",
					Status:     "Pending",
					Createdate: time.Now(),
				},
			},
			taskID:  1,
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "Delete non-existent task",
			setupTasks: Tasks{
				1: Task{
					Description: "Task 1",
					Status:     "Pending",
					Createdate: time.Now(),
				},
			},
			taskID:  2,
			wantErr: true,
			wantLen: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tasks, tempFile := setupTasks(t)
			for k, v := range test.setupTasks {
				tasks[k] = v
			}

			err := tasks.DeleteTask(test.taskID, tempFile)
			
			if (err != nil) != test.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, test.wantErr)
			}

			if len(tasks) != test.wantLen {
				t.Errorf("Expected %d tasks after deletion, got %d", test.wantLen, len(tasks))
			}
		})
	}
}