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
	Use:   "delete",
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
		err := repo.Delete(cmd.Context(), t)
		if err != nil {
			return fmt.Errorf("could not delete task: %w", err)
		}
		fmt.Printf("task %d deleted...\n", t.ID)
		return nil
	},
	ValidArgsFunction: listIncompleteTasks,
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
