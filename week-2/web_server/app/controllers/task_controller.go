package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"web_server/core/usecases"
	"web_server/utils/extract"
)

// interface adapters layer
type TaskController struct {
	GetTaskUseCase    usecases.GetTask
	AddTaskUseCase    usecases.AddTask
	UpdateTaskUseCase usecases.UpdateTask
}

func (tc *TaskController) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := extract.ExtractTaskID(r, "id")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	task, err := tc.GetTaskUseCase.Execute(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (tc *TaskController) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprint(w, "Invalid method\n")
		return
	}

	var request struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	task, err := tc.AddTaskUseCase.Execute(request.Title)
	if err != nil {
		http.Error(w, "error creating task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (tc *TaskController) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		fmt.Fprint(w, "Invalid method\n")
		return
	}

	id, err := extract.ExtractTaskID(r, "id")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	task, err := tc.UpdateTaskUseCase.Execute(id)
	if err != nil {
		http.Error(w, "error updating task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}
