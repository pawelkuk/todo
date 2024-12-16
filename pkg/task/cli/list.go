/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/pawelkuk/todo/pkg/task/model"
	"github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long: `List tasks. Example:

todo list --today  # list all incomplete tasks due today
todo list --today -A  # list all tasks due today
todo list --before 2024-12-31 --after 2024-12-01  # list all incomplete tasks for December`,
	RunE: listHandler.Handle,
}

func initList(rootCmd *cobra.Command) {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "A", false, "Lists tasks including incomplete onces. Default: false")
	listCmd.Flags().BoolP("today", "t", false, "Lists tasks for today. Default: false")
	listCmd.Flags().StringP("before", "b", "", "Lists tasks due before given date. Format: yyyy-mm-dd")
	listCmd.Flags().StringP("after", "a", "", "Lists tasks due after given date. Format: yyyy-mm-dd")
	listCmd.MarkFlagsMutuallyExclusive("today", "before")
	listCmd.MarkFlagsMutuallyExclusive("today", "after")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type ListHandler struct {
	Repo repo.Repo
}

func (h *ListHandler) Handle(cmd *cobra.Command, args []string) error {
	qf := model.QueryFilter{}
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return fmt.Errorf("could not get flag: %w", err)
	}
	if !all {
		tmp := false
		qf.Completed = &tmp
	}
	tasks, err := h.Repo.Query(cmd.Context(), qf)
	if err != nil {
		return fmt.Errorf("could not query tasks: %w", err)
	}
	before, err := cmd.Flags().GetString("before")
	if err != nil {
		return fmt.Errorf("could not get flag: %w", err)
	}
	if before != "" {
		beforeTime, err := time.Parse(time.DateOnly, before)
		if err != nil {
			return fmt.Errorf("could not parse due date: %w", err)
		}
		tasks = lo.Filter(tasks, func(t model.Task, idx int) bool { return t.DueDate.Before(beforeTime) })
	}
	after, err := cmd.Flags().GetString("after")
	if err != nil {
		return fmt.Errorf("could not get flag: %w", err)
	}
	if after != "" {
		afterTime, err := time.Parse(time.DateOnly, after)
		if err != nil {
			return fmt.Errorf("could not parse due date: %w", err)
		}
		tasks = lo.Filter(tasks, func(t model.Task, idx int) bool { return t.DueDate.After(afterTime) })
	}
	today, err := cmd.Flags().GetBool("today")
	if err != nil {
		return fmt.Errorf("could not get flag: %w", err)
	}
	if today {
		todayStr := time.Now().Format(time.DateOnly)
		tasks = lo.Filter(tasks, func(t model.Task, idx int) bool { return t.DueDate.Format(time.DateOnly) == todayStr })
	}
	res := lo.Map(tasks, func(t model.Task, idx int) string { return t.String() })
	if len(res) != 0 {
		fmt.Println(strings.Join(res, "\n"))
	} else {
		fmt.Println("no tasks with matching criteria found")
	}
	return nil
}
