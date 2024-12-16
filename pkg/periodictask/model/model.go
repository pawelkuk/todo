package model

import (
	"fmt"
	"strings"
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

func Parse(title string, schedule string, options ...PeriodicTaskOpt) (*PeriodicTask, error) {
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

type PeriodicTaskOpt func(*PeriodicTask) error

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

func (t *PeriodicTask) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%6d. ", t.ID))
	builder.WriteString(t.Schedule)
	builder.WriteString(fmt.Sprintf("| %s", truncateText(t.Title, 50)))
	if t.Description != "" {
		builder.WriteString(fmt.Sprintf(" | %s", truncateText(t.Description, 50)))
	}
	return builder.String()
}

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndexAny(s[:max], " .,:;-")] + "…"
}
