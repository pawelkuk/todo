package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	model "github.com/pawelkuk/todo/pkg/task/model"
)

type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) Create(ctx context.Context, task *model.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	result, err := r.DB.ExecContext(ctx,
		`insert into task(title, description, created_at, updated_at, due_date, completed) values(?, ?, ?, ?, ?, ?)
		returning id`,
		task.Title, task.Description, task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339),
		task.DueDate.Format(time.RFC3339), btoi(task.Completed),
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
func (r *SQLiteRepo) Read(ctx context.Context, task *model.Task) error {
	row := r.DB.QueryRowContext(ctx,
		"select title, description, created_at, updated_at, due_date, completed from task where id = ?",
		task.ID)
	if row.Err() != nil {
		return fmt.Errorf("could not query row with id=%d: %w", task.ID, row.Err())
	}
	var title, description string
	var createdAtStr, updatedAtStr, dueDateStr string
	var completedInt int
	err := row.Scan(&title, &description, &createdAtStr, &updatedAtStr, &dueDateStr, &completedInt)
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
	dueDate, err := time.Parse(time.RFC3339, dueDateStr)
	if err != nil {
		return fmt.Errorf("could not parse timestamp: %w", err)
	}
	task.Title = title
	task.Description = description
	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt
	task.DueDate = dueDate
	task.Completed = itob(completedInt)
	return nil
}
func (r *SQLiteRepo) Update(ctx context.Context, task *model.Task) error {
	_, err := r.DB.ExecContext(ctx,
		`update task set title = ?, description = ?, created_at = ?, updated_at = ?, due_date = ?, completed = ? where id = ?`,
		task.Title, task.Description, task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339),
		task.DueDate.Format(time.RFC3339), btoi(task.Completed), task.ID,
	)
	if err != nil {
		return fmt.Errorf("could not update task: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Delete(ctx context.Context, task *model.Task) error {
	_, err := r.DB.ExecContext(ctx, `delete from task where id = ?`, task.ID)
	if err != nil {
		return fmt.Errorf("could not delete task: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Query(ctx context.Context, filter model.QueryFilter) ([]model.Task, error) {
	q := `
	select id, title, description, created_at, updated_at, due_date, completed from task
	`
	buf := bytes.NewBufferString(q)
	args := []any{}
	applyFilter(filter, &args, buf)
	rows, err := r.DB.QueryContext(ctx, buf.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	tasks := []model.Task{}
	for rows.Next() {
		t := &model.Task{}
		var title, description string
		var createdAtStr, updatedAtStr, dueDateStr string
		var completedInt int
		var id int64
		err := rows.Scan(&id, &title, &description, &createdAtStr, &updatedAtStr, &dueDateStr, &completedInt)
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
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse timestamp: %w", err)
		}
		t.ID = id
		t.Title = title
		t.Description = description
		t.CreatedAt = createdAt
		t.UpdatedAt = updatedAt
		t.DueDate = dueDate
		t.Completed = itob(completedInt)
		tasks = append(tasks, *t)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return tasks, nil
}

func btoi(val bool) int {
	if val {
		return 1
	}
	return 0
}

func itob(val int) bool {
	if val == 0 {
		return false
	}
	return true
}
