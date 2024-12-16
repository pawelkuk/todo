/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	task "github.com/pawelkuk/todo/pkg/task/model"
	taskrepo "github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete task_id",
	Short: "Delete a task from todo list",
	Long: `Delete a task from todo list. Task id can be obtained
via tab completion or via the list command. Example:

todo delete 1  # delete task with id 1`,
	RunE:              deleteHandler.Handle,
	ValidArgsFunction: deleteHandler.ListIncompleteTasks,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type DeleteHandler struct {
	Repo taskrepo.Repo
}

func (h *DeleteHandler) Handle(cmd *cobra.Command, args []string) error {
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
	err := h.Repo.Delete(cmd.Context(), t)
	if err != nil {
		return fmt.Errorf("could not delete task: %w", err)
	}
	fmt.Printf("task %d deleted...\n", t.ID)
	return nil
}

func (h *DeleteHandler) ListIncompleteTasks(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return listIncompleteTasks(cmd, args, toComplete, h.Repo)
}
