package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	task "github.com/pawelkuk/todo/pkg/task/model"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestSqliteRepoTaskCRUD(t *testing.T) {
	skipCI(t)
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data/todo_test.db")
	if err != nil {
		panic(err)
	}
	repo := SQLiteRepo{DB: db}
	tt := task.Task{
		Title:       "take out trash",
		Description: "the keys to the bin are next to the door",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     time.Now().Add(24 * time.Hour),
	}
	err = repo.Create(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("task created:", tt.ID)
	tt.DueDate = tt.DueDate.Add(24 * time.Hour)
	err = repo.Update(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	taskRead := task.Task{ID: tt.ID}
	err = repo.Read(ctx, &taskRead)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = repo.Delete(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
