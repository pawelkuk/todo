/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pawelkuk/todo/pkg/periodictask/model"
	"github.com/pawelkuk/todo/pkg/periodictask/repo"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete task_id",
	Short: "Delete a periodic task from todo list",
	Long: `Delete a periodic task from todo list. Periodic task id can be obtained
via tab completion or via the list command. Example:

todo pt delete 1  # delete task with id 1`,
	RunE:              deleteHandler.Handle,
	ValidArgsFunction: deleteHandler.ListPeriodicTasks,
}

func initDelete(rootCmd *cobra.Command) {
	rootCmd.AddCommand(deleteCmd)
}

type DeleteHandler struct {
	Repo repo.Repo
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
	t := &model.PeriodicTask{ID: int64(taskID)}
	err := h.Repo.Delete(cmd.Context(), t)
	if err != nil {
		return fmt.Errorf("could not delete task: %w", err)
	}
	fmt.Printf("task %d deleted...\n", t.ID)
	return nil
}

func (h *DeleteHandler) ListPeriodicTasks(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return listPeriodicTasks(cmd, args, toComplete, h.Repo)
}
