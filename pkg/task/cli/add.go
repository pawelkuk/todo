/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	task "github.com/pawelkuk/todo/pkg/task/model"
	taskrepo "github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task in your todo list",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("repo").(*taskrepo.SQLiteRepo)
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
		opts := []task.TaskOpt{}
		if !dueTime.IsZero() {
			opts = append(opts, task.WithDueDate(dueTime.Format(time.DateOnly)))
		}
		info, err := cmd.Flags().GetString("info")
		if err != nil {
			return fmt.Errorf("could not get info: %w", err)
		}
		if info != "" {
			opts = append(opts, task.WithDescription(info))
		}
		title := strings.Join(args, " ")
		t, err := task.Parse(title, opts...)
		if err != nil {
			return fmt.Errorf("could not parse task: %w", err)
		}
		err = repo.Create(cmd.Context(), t)
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
	},
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
Examples: 10d - due in 10 days, 1h2d3w4m - due in 1 hour 2 days 3 weeks and 4 months. Default: no due time`)
	addCmd.Flags().StringP("due-date", "D", "", "Date where the task is due. Format: yyyy-mm-dd. Default: no due date")
	addCmd.MarkFlagsMutuallyExclusive("due", "due-date")
	addCmd.Flags().StringP("info", "i", "", "Additional information or details regarding the task. Default: none")
}
