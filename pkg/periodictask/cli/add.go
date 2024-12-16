/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelkuk/todo/pkg/periodictask/model"
	"github.com/pawelkuk/todo/pkg/periodictask/repo"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add task_title",
	Short: "Add a task in your todo list",
	Long:  `Add a periodic task to your todo list.`,
	RunE:  addHandler.Handle,
	Args:  cobra.ArbitraryArgs,
}

func initAdd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("schedule", "S", "", "")
	addCmd.Flags().StringP("info", "i", "", "Additional information or details regarding the task. Default: none")
	err := addCmd.MarkFlagRequired("schedule")
	if err != nil {
		log.Fatalf("could not mark flag as required: %v", err)
	}
}

type AddHandler struct {
	Repo repo.Repo
}

func (h *AddHandler) Handle(cmd *cobra.Command, args []string) error {
	schedule, err := cmd.Flags().GetString("schedule")
	if err != nil {
		return fmt.Errorf("could not get schedule: %w", err)
	}
	_, err = cron.ParseStandard(schedule)
	if err != nil {
		return fmt.Errorf("could not parse schedule %w", err)
	}
	opts := []model.PeriodicTaskOpt{}
	info, err := cmd.Flags().GetString("info")
	if err != nil {
		return fmt.Errorf("could not get info: %w", err)
	}
	if info != "" {
		opts = append(opts, model.WithDescription(info))
	}
	title := strings.Join(args, " ")
	t, err := model.Parse(title, schedule, opts...)
	if err != nil {
		return fmt.Errorf("could not parse task: %w", err)
	}
	err = h.Repo.Create(cmd.Context(), t)
	if err != nil {
		return fmt.Errorf("could not create task: %w", err)
	}
	fmt.Printf("task %d:\n", t.ID)
	fmt.Printf("\ttitle:       %s\n", t.Title)
	fmt.Printf("\tschedule:    %s\n", t.Schedule)
	if t.Description != "" {
		fmt.Printf("\tdescription: %s\n", t.Description)
	}
	return nil
}
