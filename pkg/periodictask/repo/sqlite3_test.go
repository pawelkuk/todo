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
	model "github.com/pawelkuk/todo/pkg/periodictask/model"

	"github.com/robfig/cron/v3"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestSqliteRepoPeriodicTaskCRUD(t *testing.T) {
	skipCI(t)
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data/todo_test.db")
	if err != nil {
		panic(err)
	}
	repo := SQLiteRepo{DB: db}
	schedule := "* * * * *"
	_, err = cron.ParseStandard(schedule)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	tt := model.PeriodicTask{
		Title:       "take out trash",
		Description: "the keys to the bin are next to the door",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Schedule:    schedule,
	}
	err = repo.Create(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("periodictask created:", tt.ID)
	tt.UpdatedAt = tt.UpdatedAt.Add(24 * time.Hour)
	err = repo.Update(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	taskRead := model.PeriodicTask{ID: tt.ID}
	err = repo.Read(ctx, &taskRead)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = repo.Delete(ctx, &tt)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
