/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"strings"

	"github.com/pawelkuk/todo/pkg/periodictask/model"
	"github.com/pawelkuk/todo/pkg/periodictask/repo"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List periodic tasks",
	Long: `List periodic tasks. Example:

todo pt list   # list all periodic tasks `,
	RunE: listHandler.Handle,
}

func initList(rootCmd *cobra.Command) {
	rootCmd.AddCommand(listCmd)
}

type ListHandler struct {
	Repo repo.Repo
}

func (h *ListHandler) Handle(cmd *cobra.Command, args []string) error {
	qf := model.QueryFilter{}
	tasks, err := h.Repo.Query(cmd.Context(), qf)
	if err != nil {
		return fmt.Errorf("could not query tasks: %w", err)
	}
	res := lo.Map(tasks, func(t model.PeriodicTask, idx int) string { return t.String() })
	if len(res) != 0 {
		fmt.Println(strings.Join(res, "\n"))
	} else {
		fmt.Println("no tasks with matching criteria found")
	}
	return nil
}
