/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pawelkuk/todo/pkg/task/model"
	"github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add task_title",
	Short: "Add a task in your todo list",
	Long: `Add a task to your todo list. 
	
Examples:
  todo add implement a new feature --due 1w -i "feature regards our web app and it is due one week from now"
  todo add test feature --due-date "2040-12-31" -i "we have a long time to test this feature"`,
	RunE: addHandler.Handle,
	Args: cobra.ArbitraryArgs,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().StringP("due", "d", "", `Time from now when the task is due.
Units: h - hours, d - days, w - weeks , m - months, y - years.
Examples: 10d - due in 10 days, 1h2d3w4m - due in 1 hour 2 days 3 weeks and 4 months.
Default: no due time`)
	addCmd.Flags().StringP("due-date", "D", "", "Date where the task is due. Format: yyyy-mm-dd. Default: no due date")
	addCmd.MarkFlagsMutuallyExclusive("due", "due-date")
	addCmd.Flags().StringP("info", "i", "", "Additional information or details regarding the task. Default: none")
}

type AddHandler struct {
	Repo repo.Repo
}

func (h *AddHandler) Handle(cmd *cobra.Command, args []string) error {
	due, err := cmd.Flags().GetString("due")
	if err != nil {
		return fmt.Errorf("could not get due: %w", err)
	}
	dueDate, err := cmd.Flags().GetString("due-date")
	if err != nil {
		return fmt.Errorf("could not get due-date: %w", err)
	}
	var dueTime time.Time
	if due != "" {
		dur, err := parseDue(due)
		if err != nil {
			return fmt.Errorf("could not parse due time: %w", err)
		}
		dueTime = time.Now().Add(dur)
	}
	if dueDate != "" {
		dueTime, err = time.Parse(time.DateOnly, dueDate)
		if err != nil {
			return fmt.Errorf("could not parse due date: %w", err)
		}
	}
	opts := []model.TaskOpt{}
	if !dueTime.IsZero() {
		opts = append(opts, model.WithDueDate(dueTime.Format(time.DateOnly)))
	}
	info, err := cmd.Flags().GetString("info")
	if err != nil {
		return fmt.Errorf("could not get info: %w", err)
	}
	if info != "" {
		opts = append(opts, model.WithDescription(info))
	}
	title := strings.Join(args, " ")
	t, err := model.Parse(title, opts...)
	if err != nil {
		return fmt.Errorf("could not parse task: %w", err)
	}
	err = h.Repo.Create(cmd.Context(), t)
	if err != nil {
		return fmt.Errorf("could not create task: %w", err)
	}
	fmt.Printf("task %d:\n", t.ID)
	fmt.Printf("\ttitle:       %s\n", t.Title)
	if t.Description != "" {
		fmt.Printf("\tdescription: %s\n", t.Description)
	}
	if !t.DueDate.IsZero() {
		fmt.Printf("\tdue date:    %s\n", t.DueDate.Format(time.DateOnly))
	}
	return nil
}
