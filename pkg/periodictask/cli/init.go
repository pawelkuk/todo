package cli

// import (
// 	"database/sql"
// 	"log"

// 	"github.com/caarlos0/env/v11"
// 	"github.com/pawelkuk/todo/pkg/config"
// )

// var editHandler EditHandler
// var listHandler ListHandler
// var deleteHandler DeleteHandler
// var completeHandler CompleteHandler
// var addHandler AddHandler

// func Initialize() {
// 	var cfg config.Config
// 	err := env.Parse(&cfg)
// 	if err != nil {
// 		log.Fatalf("could not parse config: %v", err)
// 	}
// 	db, err := sql.Open("sqlite3", cfg.DBPath)
// 	if err != nil {
// 		log.Fatalf("could not open database: %v", err)
// 	}
// 	editHandler = EditHandler{
// 		Config: cfg,
// 		Repo:   &taskrepo.SQLiteRepo{DB: db},
// 	}
// 	listHandler = ListHandler{
// 		Repo: &taskrepo.SQLiteRepo{DB: db},
// 	}
// 	deleteHandler = DeleteHandler{
// 		Repo: &taskrepo.SQLiteRepo{DB: db},
// 	}
// 	completeHandler = CompleteHandler{
// 		Repo: &taskrepo.SQLiteRepo{DB: db},
// 	}
// 	addHandler = AddHandler{
// 		Repo: &taskrepo.SQLiteRepo{DB: db},
// 	}
// }
