package repo

import (
	"context"

	task "github.com/pawelkuk/todo/pkg/task/model"
)

type FakeRepo struct {
	task    *task.Task
	taskErr error
}

func (r *FakeRepo) Create(_ context.Context, task *task.Task) error {
	if r.taskErr != nil {
		return r.taskErr
	}
	*task = *r.task
	return nil
}
func (r *FakeRepo) Read(_ context.Context, task *task.Task) error {
	if r.taskErr != nil {
		return r.taskErr
	}
	*task = *r.task
	return nil
}
func (r *FakeRepo) Update(_ context.Context, task *task.Task) error {
	if r.taskErr != nil {
		return r.taskErr
	}
	*task = *r.task
	return nil
}
func (r *FakeRepo) Delete(_ context.Context, task *task.Task) error {
	if r.taskErr != nil {
		return r.taskErr
	}
	*task = *r.task
	return nil
}
