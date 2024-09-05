package entities_test

import (
	"testing"
	"web_server/core/entities"
)

func TestMarkAsCompleted(t *testing.T) {
	task := entities.Task{Title: "New task"}
	task.MarkAsCompleted()

	if !task.Completed {
		t.Errorf("Expected task to be marked as completed")
	}
}
