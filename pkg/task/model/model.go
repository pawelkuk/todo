package task

import (
	"fmt"
	"time"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Completed   bool
}

func Parse(title string, options ...func(*Task) error) (*Task, error) {
	task := &Task{}
	if title == "" {
		return nil, fmt.Errorf("title can't be empty")
	}
	task.Title = title

	for _, o := range options {
		err := o(task)
		if err != nil {
			return nil, fmt.Errorf("could not construct task: %w", err)
		}
	}
	return task, nil
}

func WithDescription(description string) func(*Task) error {
	return func(t *Task) error {
		if description == "" {
			return fmt.Errorf("description can't be empty")
		}
		t.Description = description
		return nil
	}
}

func WithCompleted(completed bool) func(*Task) error {
	return func(t *Task) error {
		t.Completed = completed
		return nil
	}
}

func WithID(id int64) func(*Task) error {
	return func(t *Task) error {
		t.ID = id
		return nil
	}
}

func WithDueDate(dueDate string) func(*Task) error {
	return func(t *Task) error {
		dueDateTime, err := time.Parse(time.DateOnly, dueDate)
		if err != nil {
			return fmt.Errorf("invalid dueDate: %w", err)
		}
		t.DueDate = dueDateTime
		return nil
	}
}
