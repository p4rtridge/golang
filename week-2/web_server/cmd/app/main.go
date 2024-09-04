package main

import (
	"fmt"
	"log"
	"net/http"
	"web_server/app/controllers"
	"web_server/config"
	"web_server/core/usecases"
	"web_server/infras/repos"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("[Error]: config error: %v", err)
	}

	repo, err := repos.NewPostgresRepo(cfg, 5)
	if err != nil {
		log.Fatalf("[Error]: repo error: %v", err)
	}
	defer repo.DB.Close()

	getTaskUseCase := usecases.GetTask{Repo: repo}
	addTaskUseCase := usecases.AddTask{Repo: repo}
	updateTaskUseCase := usecases.UpdateTask{Repo: repo}

	taskController := &controllers.TaskController{
		GetTaskUseCase:    getTaskUseCase,
		AddTaskUseCase:    addTaskUseCase,
		UpdateTaskUseCase: updateTaskUseCase,
	}

	http.HandleFunc("/tasks", taskController.AddTaskHandler)
	http.HandleFunc("/tasks/{id}", taskController.GetTaskHandler)
	http.HandleFunc("/tasks/{id}/update", taskController.UpdateTaskHandler)

	tcpListener := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	fmt.Printf("Listening on %s...\n", tcpListener)
	log.Fatalln(http.ListenAndServe(tcpListener, nil))
}
