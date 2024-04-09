package main

import (
	"cube/manager"
	"cube/node"
	"cube/task"
	"cube/worker"
	"fmt"
	"os"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	d := task.NewDocker(&c)

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Error: %v\n", result.Error)
		return nil, &result
	}

	fmt.Printf("Container is %s is running with config: %v\n", result.ContainerID, c)
	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Printf("Error: %v\n", result.Error)
		return &result
	}

	fmt.Printf("Container %s is stopped and removed\n", id)
	return &result
}

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "Task-1",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}
	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Timestamp: time.Now(),
		Task:      t,
	}
	fmt.Printf("Task: %v\n", t)
	fmt.Printf("TaskEvent: %v\n", te)

	w := worker.Worker{
		Name:      "Worker-1",
		Queue:     *queue.New(),
		Db:        make(map[uuid.UUID]*task.Task),
		TaskCount: 0,
	}
	fmt.Printf("Worker: %v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending:     *queue.New(),
		TaskDb:      make(map[string]*task.Task),
		TaskEventDb: make(map[string]*task.TaskEvent),
		Workers:     []string{w.Name},
	}
	fmt.Printf("Manager: %v\n", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()

	n := node.Node{
		Name:   "Node-1",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 4096,
		Disk:   100,
		Role:   "Worker",
	}
	fmt.Printf("Node: %v\n", n)

	fmt.Println("Creating a test container")
	dockerTask, createResult := createContainer()
	if createResult.Error != nil {
		fmt.Printf("Error creating container: %v\n", createResult.Error)
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)
	fmt.Printf("Stopping container %s\n", createResult.ContainerID)
	stopResult := stopContainer(dockerTask, createResult.ContainerID)
	if stopResult.Error != nil {
		fmt.Printf("Error stopping container: %v\n", stopResult.Error)
		os.Exit(1)
	}
}