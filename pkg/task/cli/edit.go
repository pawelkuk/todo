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
	"github.com/pawelkuk/todo/pkg/task/model"
	"github.com/pawelkuk/todo/pkg/task/repo"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit task_id",
	Short: "Edit a task in todo list",
	Long: `Edit a task in todo list.

Opens the task in yaml format in editor specified by $EDITOR env variable.
The id field is not editable. It serves an informational purpose only. If
any modified value does not match the required format the edit won't take effect.

Example:
todo edit 1  # edit task with id = 1
`,
	RunE: editHandler.RunE,
}

func init() {
	rootCmd.AddCommand(editCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
	t := &model.Task{ID: int64(taskID)}
	err := h.Repo.Read(cmd.Context(), t)
	originalTask, err := marshalToYaml(t)
	if err != nil {
		return fmt.Errorf("could not get original content: %w", err)
	}
	tmpFile, err := os.CreateTemp("/tmp", "tmp-task.yaml")
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
	t.DueDate = modifiedTask.DueDate
	t.Completed = modifiedTask.Completed
	t.UpdatedAt = time.Now()
	err = h.Repo.Update(cmd.Context(), t)
	if err != nil {
		return fmt.Errorf("could not update task: %w", err)
	}
	fmt.Printf("task %d updated\n", t.ID)
	return nil
}

func marshalToYaml(t *model.Task) (string, error) {
	yamltask := &yamltask{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate.Format(time.DateOnly),
		Completed:   t.Completed,
	}
	out, err := yaml.Marshal(yamltask)
	if err != nil {
		return "", fmt.Errorf("could not marshal task: %w", err)
	}
	return string(out), nil
}

func unmarshalYaml(taskStr []byte) (*model.Task, error) {
	yamltask := &yamltask{}
	err := yaml.Unmarshal(taskStr, yamltask)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal task: %w", err)
	}
	t, err := model.Parse(
		yamltask.Title,
		model.WithCompleted(yamltask.Completed),
		model.WithDescription(yamltask.Description),
		model.WithID(yamltask.ID),
		model.WithDueDate(yamltask.DueDate),
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
	DueDate     string `yaml:"dueDate"`
	Completed   bool   `yaml:"completed"`
}
