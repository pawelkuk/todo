package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pawelkuk/todo/pkg/periodictask/model"
	"github.com/pawelkuk/todo/pkg/periodictask/repo"
	"github.com/spf13/cobra"
)

func listPeriodicTasks(cmd *cobra.Command, _ []string, toComplete string, repo repo.Repo) ([]string, cobra.ShellCompDirective) {
	tasks, err := repo.Query(cmd.Context(), model.QueryFilter{})
	if err != nil {
		fmt.Printf("could not query for tasks: %v\n", err)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	res := []string{}
	for _, t := range tasks {
		if toComplete == "" {
			res = append(res, t.String())
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
		if strings.Contains(t.Schedule, toComplete) {
			res = append(res, t.String())
			continue
		}
	}
	return res, cobra.ShellCompDirectiveNoFileComp
}
