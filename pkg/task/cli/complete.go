/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	task "github.com/pawelkuk/todo/pkg/task/model"
	taskrepo "github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("repo").(*taskrepo.SQLiteRepo)
		taskExpr := regexp.MustCompile(`^(\s+)?(?<taskid>\d+)`)
		match := taskExpr.FindStringSubmatch(strings.Join(args, " "))
		var taskID int
		for i, name := range taskExpr.SubexpNames() {
			if i != 0 && name != "" && len(match[i]) != 0 {
				tmpTaskID, err := strconv.Atoi(match[i])
				if err != nil {
					return fmt.Errorf("could not parse number: %w", err)
				}
				taskID = tmpTaskID
			}
		}
		if taskID == 0 {
			return fmt.Errorf("could not much provided args: %s", strings.Join(args, " "))
		}
		t := &task.Task{ID: int64(taskID)}
		err := repo.Read(cmd.Context(), t)
		if err != nil {
			return fmt.Errorf("could not read task: %w", err)
		}
		if t.Completed {
			return fmt.Errorf("task %d already completed...\n", taskID)
		}
		t.Completed = true
		err = repo.Update(cmd.Context(), t)
		if err != nil {
			return fmt.Errorf("could not update task: %w\n", err)
		}
		fmt.Printf("task %d completed...\n", t.ID)
		return nil
	},
	ValidArgsFunction: listIncompleteTasks,
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listIncompleteTasks(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	repo := cmd.Context().Value("repo").(*taskrepo.SQLiteRepo)
	tasks, err := repo.Query(cmd.Context(), task.QueryFilter{})
	if err != nil {
		fmt.Printf("could not query for tasks: %v\n", err)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	res := []string{}
	for _, t := range tasks {
		if t.Completed {
			continue
		}
		if strings.Contains(strconv.Itoa(int(t.ID)), toComplete) {
			res = append(res, t.String())
			continue
		}
		if strings.Contains(t.Title, toComplete) {
			res = append(res, t.String())
			continue
		}
		if strings.Contains(t.Description, toComplete) {
			res = append(res, t.String())
			continue
		}
		if strings.Contains(t.DueDate.Format(time.DateOnly), toComplete) {
			res = append(res, t.String())
			continue
		}
	}
	return res, cobra.ShellCompDirectiveNoFileComp
}
