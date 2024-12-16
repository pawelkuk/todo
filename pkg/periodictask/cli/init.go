package cli

import (
	"database/sql"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/pawelkuk/todo/pkg/config"
	"github.com/pawelkuk/todo/pkg/periodictask/repo"
	"github.com/spf13/cobra"
)

var editHandler EditHandler
var listHandler ListHandler
var deleteHandler DeleteHandler
var addHandler AddHandler

func Initialize(rootCmd *cobra.Command) {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("could not parse config: %v", err)
	}
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	editHandler = EditHandler{
		Config: cfg,
		Repo:   &repo.SQLiteRepo{DB: db},
	}
	listHandler = ListHandler{
		Repo: &repo.SQLiteRepo{DB: db},
	}
	deleteHandler = DeleteHandler{
		Repo: &repo.SQLiteRepo{DB: db},
	}
	addHandler = AddHandler{
		Repo: &repo.SQLiteRepo{DB: db},
	}
	rootCmd.AddCommand(periodicCmd)
	// We want them to be subcommands of periodicCmd in order not to have name clashes with normal tasks
	initAdd(periodicCmd)
	initList(periodicCmd)
	initDelete(periodicCmd)
	initEdit(periodicCmd)
}
