package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pawelkuk/todo/pkg/task/model"
	"github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/spf13/cobra"
)

func parseDue(d string) (time.Duration, error) {
	durationExpression := regexp.MustCompile(`^(?<hours>\d+h)?(?<days>\d+d)?(?<weeks>\d+w)?(?<months>\d+m)?(?<years>\d+y)?$`)
	var duration time.Duration
	match := durationExpression.FindStringSubmatch(d)
	for i, name := range durationExpression.SubexpNames() {
		if i != 0 && name != "" && len(match[i]) != 0 {
			num, err := strconv.Atoi(match[i][:len(match[i])-1])
			if err != nil {
				return 0, fmt.Errorf("could not parse number: %w", err)
			}
			switch name {
			case "hours":
				duration += time.Duration(num) * time.Hour
			case "days":
				duration += time.Duration(24*num) * time.Hour
			case "weeks":
				duration += time.Duration(7*24*num) * time.Hour
			case "months":
				duration += time.Duration(30*24*num) * time.Hour
			case "years":
				duration += time.Duration(365*24*num) * time.Hour
			}
		}
	}
	return duration, nil
}

func listIncompleteTasks(cmd *cobra.Command, _ []string, toComplete string, repo repo.Repo) ([]string, cobra.ShellCompDirective) {
	tasks, err := repo.Query(cmd.Context(), model.QueryFilter{})
	if err != nil {
		fmt.Printf("could not query for tasks: %v\n", err)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	res := []string{}
	for _, t := range tasks {
		if t.Completed {
			continue
		}
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
		if strings.Contains(t.DueDate.Format(time.DateOnly), toComplete) {
			res = append(res, t.String())
			continue
		}
	}
	return res, cobra.ShellCompDirectiveNoFileComp
}
