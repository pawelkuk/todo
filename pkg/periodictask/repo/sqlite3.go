package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	model "github.com/pawelkuk/todo/pkg/periodictask/model"
)

type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) Create(ctx context.Context, task *model.PeriodicTask) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	result, err := r.DB.ExecContext(ctx,
		`insert into periodic_task(title, description, created_at, updated_at, schedule) values(?, ?, ?, ?, ?)
		returning id`,
		task.Title, task.Description, task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339), task.Schedule,
	)
	if err != nil {
		return fmt.Errorf("could not create task: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("could not get last insert id: %w", err)
	}
	task.ID = id
	return nil
}
func (r *SQLiteRepo) Read(ctx context.Context, task *model.PeriodicTask) error {
	row := r.DB.QueryRowContext(ctx,
		"select title, description, created_at, updated_at, schedule from periodic_task where id = ?",
		task.ID)
	if row.Err() != nil {
		return fmt.Errorf("could not query row with id=%d: %w", task.ID, row.Err())
	}
	var title, description, createdAtStr, updatedAtStr, schedule string
	err := row.Scan(&title, &description, &createdAtStr, &updatedAtStr, &schedule)
	if err != nil {
		return fmt.Errorf("could not scan row: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return fmt.Errorf("could not parse timestamp: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return fmt.Errorf("could not parse timestamp: %w", err)
	}
	task.Title = title
	task.Description = description
	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt
	task.Schedule = schedule
	return nil
}
func (r *SQLiteRepo) Update(ctx context.Context, task *model.PeriodicTask) error {
	_, err := r.DB.ExecContext(ctx,
		`update periodic_task set title = ?, description = ?, created_at = ?, updated_at = ?, schedule = ? where id = ?`,
		task.Title, task.Description, task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339), task.Schedule, task.ID,
	)
	if err != nil {
		return fmt.Errorf("could not update task: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Delete(ctx context.Context, task *model.PeriodicTask) error {
	_, err := r.DB.ExecContext(ctx, `delete from periodic_task where id = ?`, task.ID)
	if err != nil {
		return fmt.Errorf("could not delete periodictask: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Query(ctx context.Context, filter model.QueryFilter) ([]model.PeriodicTask, error) {
	q := `
	select id, title, description, created_at, updated_at, schedule from periodic_task
	`
	buf := bytes.NewBufferString(q)
	args := []any{}
	applyFilter(filter, &args, buf)
	rows, err := r.DB.QueryContext(ctx, buf.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	tasks := []model.PeriodicTask{}
	for rows.Next() {
		t := &model.PeriodicTask{}
		var title, description string
		var createdAtStr, updatedAtStr, schedule string
		var id int64
		err := rows.Scan(&id, &title, &description, &createdAtStr, &updatedAtStr, &schedule)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse timestamp: %w", err)
		}
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse timestamp: %w", err)
		}
		t.ID = id
		t.Title = title
		t.Description = description
		t.CreatedAt = createdAt
		t.UpdatedAt = updatedAt
		t.Schedule = schedule
		tasks = append(tasks, *t)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return tasks, nil
}
