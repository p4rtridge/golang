package entities

// Core business
type Task struct {
	Title     string
	Id        int
	Completed bool
}

func (t *Task) MarkAsCompleted() {
	t.Completed = true
}
