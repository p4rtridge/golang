package usecases

import "web_server/core/entities"

// use cases layer
// belongs to the use cases layer and abstracts data storage, external dependencies must implement this.
type TaskRepo interface {
	GetTask(id int) (*entities.Task, error)
	SaveTask(task *entities.Task) (*entities.Task, error)
}

// application specific business rule
type GetTask struct {
	Repo TaskRepo
}

func (uc *GetTask) Execute(id int) (*entities.Task, error) {
	return uc.Repo.GetTask(id)
}

// application specific business rule
type AddTask struct {
	Repo TaskRepo
}

func (uc *AddTask) Execute(title string) (*entities.Task, error) {
	task := &entities.Task{
		Title: title,
	}

	newTask, err := uc.Repo.SaveTask(task)
	if err != nil {
		return nil, err
	}

	return newTask, nil
}

// application specific business rule
type UpdateTask struct {
	Repo TaskRepo
}

func (uc *UpdateTask) Execute(id int) (*entities.Task, error) {
	task, err := uc.Repo.GetTask(id)
	if err != nil {
		return nil, err
	}

	task.MarkAsCompleted()

	updatedTask, err := uc.Repo.SaveTask(task)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}
