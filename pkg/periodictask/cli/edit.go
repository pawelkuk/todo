/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/pawelkuk/todo/pkg/config"
	"github.com/pawelkuk/todo/pkg/periodictask/model"
	"github.com/pawelkuk/todo/pkg/periodictask/repo"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit task_id",
	Short: "Edit a periodic task in todo list",
	Long: `Edit a periodic task in todo list.

Opens the periodic task in yaml format in editor specified by $EDITOR env variable.
The id field is not editable. It serves an informational purpose only. If
any modified value does not match the required format the edit won't take effect.

Example:
todo pt edit 1  # edit periodic task with id = 1
`,
	RunE:              editHandler.RunE,
	ValidArgsFunction: editHandler.ListPeriodicTasks,
}

func initEdit(rootCmd *cobra.Command) {
	rootCmd.AddCommand(editCmd)
}

type EditHandler struct {
	Repo   repo.Repo
	Config config.Config
}

func (h *EditHandler) RunE(cmd *cobra.Command, args []string) error {
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
	err := h.Repo.Read(cmd.Context(), t)
	originalTask, err := marshalToYaml(t)
	if err != nil {
		return fmt.Errorf("could not get original content: %w", err)
	}
	tmpFile, err := os.CreateTemp("/tmp", "tmp-periodic-task.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temporary file
	if _, err := tmpFile.WriteString(originalTask); err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tmpFile.Close() // Close the file before opening it in the editor

	editor := h.Config.Editor
	if editor == "" {
		editor = "vi" // Default to vi if editor is not set
	}

	// Open the editor
	editCmd := exec.Command(editor, tmpFile.Name())
	editCmd.Stdin = os.Stdin
	editCmd.Stdout = os.Stdout
	editCmd.Stderr = os.Stderr

	if err := editCmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	// Read the modified file
	modifiedTaskStr, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read modified file: %w", err)
	}

	modifiedTask, err := unmarshalYaml(modifiedTaskStr)
	if err != nil {
		return fmt.Errorf("could not unmarshal task: %w", err)
	}
	if t.ID != modifiedTask.ID {
		fmt.Println("warning: can't change id")
	}
	t.Title = modifiedTask.Title
	t.Description = modifiedTask.Description
	t.Schedule = modifiedTask.Schedule
	t.UpdatedAt = time.Now()
	err = h.Repo.Update(cmd.Context(), t)
	if err != nil {
		return fmt.Errorf("could not update task: %w", err)
	}
	fmt.Printf("task %d updated\n", t.ID)
	return nil
}

func marshalToYaml(t *model.PeriodicTask) (string, error) {
	yamltask := &yamltask{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Schedule:    t.Schedule,
	}
	out, err := yaml.Marshal(yamltask)
	if err != nil {
		return "", fmt.Errorf("could not marshal task: %w", err)
	}
	return string(out), nil
}

func unmarshalYaml(taskStr []byte) (*model.PeriodicTask, error) {
	yamltask := &yamltask{}
	err := yaml.Unmarshal(taskStr, yamltask)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal task: %w", err)
	}
	t, err := model.Parse(
		yamltask.Title,
		yamltask.Schedule,
		model.WithDescription(yamltask.Description),
		model.WithID(yamltask.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("could not parse task: %w", err)
	}
	return t, nil

}

type yamltask struct {
	ID          int64  `yaml:"id"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Schedule    string `yaml:"schedule"`
}

func (h *EditHandler) ListPeriodicTasks(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return listPeriodicTasks(cmd, args, toComplete, h.Repo)
}
