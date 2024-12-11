package repo

import (
	"bytes"
	"strings"

	model "github.com/pawelkuk/todo/pkg/periodictask/model"
)

func applyFilter(filter model.QueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, filter.ID)
	}

	if filter.Description != nil {
		wc = append(wc, "description like ?")
		*args = append(*args, filter.Description)
	}

	if filter.Title != nil {
		wc = append(wc, "title like ?")
		*args = append(*args, filter.Title)
	}

	if filter.Schedule != nil {
		wc = append(wc, "schedule = ?")
		*args = append(*args, filter.Schedule)
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
