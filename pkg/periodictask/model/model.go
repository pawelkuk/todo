package periodictask

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

type PeriodicTask struct {
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Schedule    string
}

func Parse(title string, schedule string, options ...func(*PeriodicTask) error) (*PeriodicTask, error) {
	task := &PeriodicTask{}
	if title == "" {
		return nil, fmt.Errorf("title can't be empty")
	}
	task.Title = title

	_, err := cron.ParseStandard(schedule)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule: %w", err)
	}
	task.Schedule = schedule

	for _, o := range options {
		err := o(task)
		if err != nil {
			return nil, fmt.Errorf("could not construct task: %w", err)
		}
	}
	return task, nil
}

func WithDescription(description string) func(*PeriodicTask) error {
	return func(t *PeriodicTask) error {
		if description == "" {
			return fmt.Errorf("description can't be empty")
		}
		t.Description = description
		return nil
	}
}

func WithID(id int64) func(*PeriodicTask) error {
	return func(t *PeriodicTask) error {
		t.ID = id
		return nil
	}
}
