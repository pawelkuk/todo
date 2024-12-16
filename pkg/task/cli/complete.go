/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pawelkuk/todo/pkg/task/model"
	"github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete task_id",
	Short: "Complete a task from your todo list",
	Long: `Complete a task from your todo list. The id of the task
can be obtained via tab completion or the list command.:

`,
	RunE:              completeHandler.Handle,
	ValidArgsFunction: completeHandler.ListIncompleteTasks,
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

type CompleteHandler struct {
	Repo repo.Repo
}

func (h *CompleteHandler) Handle(cmd *cobra.Command, args []string) error {
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
	t := &model.Task{ID: int64(taskID)}
	err := h.Repo.Read(cmd.Context(), t)
	if err != nil {
		return fmt.Errorf("could not read task: %w", err)
	}
	if t.Completed {
		return fmt.Errorf("task %d already completed...\n", taskID)
	}
	t.Completed = true
	err = h.Repo.Update(cmd.Context(), t)
	if err != nil {
		return fmt.Errorf("could not update task: %w\n", err)
	}
	fmt.Printf("task %d completed...\n", t.ID)
	return nil
}

func (h *CompleteHandler) ListIncompleteTasks(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return listIncompleteTasks(cmd, args, toComplete, h.Repo)
}
