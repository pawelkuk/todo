package repo

import (
	"bytes"
	"strings"

	model "github.com/pawelkuk/todo/pkg/user/model"
)

func applyFilter(filter model.QueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, filter.ID)
	}

	if filter.Email != nil {
		wc = append(wc, "email LIKE ?")
		*args = append(*args, filter.Email.Address)
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
