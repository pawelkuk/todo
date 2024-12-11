package repo

import (
	"context"

	task "github.com/pawelkuk/todo/pkg/task/model"
)

type Repo interface {
	Create(context.Context, *task.Task) error
	Read(context.Context, *task.Task) error
	Update(context.Context, *task.Task) error
	Delete(context.Context, *task.Task) error
	Query(context.Context, task.QueryFilter) ([]task.Task, error)
}
