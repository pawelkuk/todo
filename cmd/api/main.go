package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	auth "github.com/pawelkuk/todo/pkg/auth/handler"
	authmiddleware "github.com/pawelkuk/todo/pkg/auth/middleware"
	authrepo "github.com/pawelkuk/todo/pkg/auth/repo"
	"github.com/pawelkuk/todo/pkg/config"
	periodictask "github.com/pawelkuk/todo/pkg/periodictask/handler"
	periodictaskrepo "github.com/pawelkuk/todo/pkg/periodictask/repo"
	task "github.com/pawelkuk/todo/pkg/task/handler"
	taskrepo "github.com/pawelkuk/todo/pkg/task/repo"
	user "github.com/pawelkuk/todo/pkg/user/handler"
	userrepo "github.com/pawelkuk/todo/pkg/user/repo"
)

func main() {
	err := serve()
	log.Fatal(err)
}

func serve() error {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("could not parse config: %v", err)
	}
	db, err := sql.Open("sqlite3", cfg.DBPath)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("could not connect to db: %v", err)
	}
	r := gin.Default()

	authMiddleware := authmiddleware.Middleware{Repo: &authrepo.SQLiteRepo{DB: db}}
	authHandler := auth.Handler{
		Repo:     &authrepo.SQLiteRepo{DB: db},
		UserRepo: &userrepo.SQLiteRepo{DB: db}}
	a := r.Group("/auth")
	a.POST("/login", authHandler.Login)
	a.POST("/logout", authMiddleware.Handle, authHandler.Logout)

	taskHandler := task.Handler{Repo: &taskrepo.SQLiteRepo{DB: db}}
	t := r.Group("/tasks").Use(authMiddleware.Handle)
	t.GET("/:id", taskHandler.Get)
	t.GET("/", taskHandler.GetList)
	t.POST("/", taskHandler.Post)
	t.POST("/:id/*complete", taskHandler.PostComplete)
	t.PATCH("/:id", taskHandler.Patch)
	t.DELETE("/:id", taskHandler.Delete)

	periodicTaskHandler := periodictask.Handler{Repo: &periodictaskrepo.SQLiteRepo{DB: db}}
	pt := r.Group("/periodic-tasks").Use(authMiddleware.Handle)
	pt.GET("/", periodicTaskHandler.List)
	pt.GET("/:id", periodicTaskHandler.Get)
	pt.POST("/", periodicTaskHandler.Post)
	pt.PATCH("/:id", periodicTaskHandler.Patch)
	pt.DELETE("/:id", periodicTaskHandler.Delete)

	userHandler := user.Handler{Repo: &userrepo.SQLiteRepo{DB: db}}
	u := r.Group("/users")
	u.GET("/:id", authMiddleware.Handle, userHandler.Get)
	u.POST("/", userHandler.Post)
	u.PATCH("/:id", authMiddleware.Handle, userHandler.Patch)
	u.DELETE("/:id", authMiddleware.Handle, userHandler.Delete)

	err = r.Run(":8080")
	return err
}
