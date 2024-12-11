package repo

import (
	"bytes"
	"strings"
	"time"

	model "github.com/pawelkuk/todo/pkg/task/model"
)

func applyFilter(filter model.QueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, filter.ID)
	}

	if filter.Completed != nil {
		wc = append(wc, "completed = ?")
		*args = append(*args, filter.Completed)
	}

	if filter.Description != nil {
		wc = append(wc, "description like ?")
		*args = append(*args, filter.Description)
	}

	if filter.Title != nil {
		wc = append(wc, "title like ?")
		*args = append(*args, filter.Title)
	}

	if filter.DueDate != nil {
		wc = append(wc, "due_date = ?")
		*args = append(*args, filter.DueDate.Format(time.DateOnly))
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
