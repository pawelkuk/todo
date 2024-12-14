/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"strings"

	task "github.com/pawelkuk/todo/pkg/task/model"
	taskrepo "github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("repo").(*taskrepo.SQLiteRepo)
		qf := task.QueryFilter{}
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return fmt.Errorf("could not get flag: %w", err)
		}
		if !all {
			tmp := true
			qf.Completed = &tmp
		}
		tasks, err := repo.Query(cmd.Context(), qf)
		if err != nil {
			return fmt.Errorf("could not query tasks: %w", err)
		}
		res := lo.Map(tasks, func(t task.Task, idx int) string { return t.String() })
		fmt.Println(strings.Join(res, "\n"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "A", false, "Lists all tasks including incomplete onces. Default: false")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
