package task

import (
	"fmt"
	"strings"
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

type TaskOpt func(*Task) error

func Parse(title string, options ...TaskOpt) (*Task, error) {
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

func (t *Task) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%6d. ", t.ID))
	if t.Completed {
		builder.WriteString("[x] ")
	} else {
		builder.WriteString("[ ] ")
	}
	if !t.DueDate.IsZero() {
		builder.WriteString(fmt.Sprintf("due: %s", t.DueDate.Format(time.DateOnly)))
	} else {
		builder.WriteString("no due date")
	}
	builder.WriteString(fmt.Sprintf("| %s", truncateText(t.Title, 50)))
	return builder.String()
}

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndexAny(s[:max], " .,:;-")] + "…"
}
