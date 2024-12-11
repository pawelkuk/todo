package repo

import (
	"context"

	model "github.com/pawelkuk/todo/pkg/periodictask/model"
)

type Repo interface {
	Create(context.Context, *model.PeriodicTask) error
	Read(context.Context, *model.PeriodicTask) error
	Update(context.Context, *model.PeriodicTask) error
	Delete(context.Context, *model.PeriodicTask) error
	Query(context.Context, model.QueryFilter) ([]model.PeriodicTask, error)
}
