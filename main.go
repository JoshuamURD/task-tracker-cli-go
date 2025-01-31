package main

import (
	"fmt"
	"joshuamURD/go-tasks-cli/pkg/cli"
	"joshuamURD/go-tasks-cli/pkg/task"
	"log"
)


func main() {
	tasks := task.NewTasks()
	path := "tasks.json"
	task.LoadTasks(
		path,
		tasks,
	)
	cli.AddVerbHandler("add", func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no task provided. Usage: go run main.go add <task>")
		}
		if err := task.AddTask(args[0], path, tasks); err != nil {
			return err
		}
		return nil
	})
	cli.AddVerbHandler("list", func(args []string) error {
		return nil
	})
	cli.AddVerbHandler("delete", func(args []string) error {
		return nil
	})
	cli.AddVerbHandler("help", func(args []string) error {
		return nil
	})

	err := cli.HandleVerb()
	if err != nil {
		log.Fatal(err)
	}

}
