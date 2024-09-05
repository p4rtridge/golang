package usecases_test

import (
	"testing"
	"web_server/core/entities"
	"web_server/core/usecases"
)

type MockTaskRepo struct {
	tasks []*entities.Task
}

func (m *MockTaskRepo) GetTask(id int) (*entities.Task, error) {
	return nil, nil
}

func (m *MockTaskRepo) SaveTask(task *entities.Task) (*entities.Task, error) {
	m.tasks = append(m.tasks, task)
	return task, nil
}

func TestAddTask(t *testing.T) {
	mockRepo := &MockTaskRepo{}
	addTaskUseCase := usecases.AddTask{Repo: mockRepo}

	task, err := addTaskUseCase.Execute("New task")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if task.Title != "New task" {
		t.Errorf("Expected task title to be 'New task', got %s", task.Title)
	}

	if len(mockRepo.tasks) != 1 {
		t.Errorf("Expected 1 task in repository, got %d", len(mockRepo.tasks))
	}
}
