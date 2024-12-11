package repo

import (
	"context"

	model "github.com/pawelkuk/todo/pkg/periodictask/model"
)

type FakeRepo struct {
	periodictask    *model.PeriodicTask
	periodictaskErr error
}

func (r *FakeRepo) Create(_ context.Context, task *model.PeriodicTask) error {
	if r.periodictaskErr != nil {
		return r.periodictaskErr
	}
	*task = *r.periodictask
	return nil
}
func (r *FakeRepo) Read(_ context.Context, task *model.PeriodicTask) error {
	if r.periodictaskErr != nil {
		return r.periodictaskErr
	}
	*task = *r.periodictask
	return nil
}
func (r *FakeRepo) Update(_ context.Context, task *model.PeriodicTask) error {
	if r.periodictaskErr != nil {
		return r.periodictaskErr
	}
	*task = *r.periodictask
	return nil
}
func (r *FakeRepo) Delete(_ context.Context, task *model.PeriodicTask) error {
	if r.periodictaskErr != nil {
		return r.periodictaskErr
	}
	*task = *r.periodictask
	return nil
}
