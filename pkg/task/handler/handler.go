package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	model "github.com/pawelkuk/todo/pkg/task/model"
	repo "github.com/pawelkuk/todo/pkg/task/repo"
)

type Handler struct {
	Repo repo.Repo
}

type taskPOST struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate" binding:"required"`
}

type taskGET struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Completed   bool   `json:"completed"`
}

type taskList struct {
	Items []taskGET `json:"items"`
}

type TaskPATCH struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"dueDate,omitempty"`
	Completed   *bool  `json:"completed"`
}

type TaskFilter struct {
	After       time.Time `form:"after" time_format:"2006-01-02" time_utc:"1"`
	DueDate     time.Time `form:"dueDate" time_format:"2006-01-02" time_utc:"1"`
	Before      time.Time `form:"before" time_format:"2006-01-02" time_utc:"1"`
	Completed   *bool     `form:"completed"`
	Title       string    `form:"title"`
	Description string    `form:"description"`
}

func parseTaskGet(task model.Task) taskGET {
	return taskGET{
		ID:          int(task.ID),
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate.Format(time.DateOnly),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
		Completed:   task.Completed,
	}
}

func (h *Handler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := &model.Task{ID: id}
	err = h.Repo.Read(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, parseTaskGet(*task))
}

func (h *Handler) GetList(c *gin.Context) {
	ulrQueryParams := TaskFilter{}
	err := c.ShouldBind(&ulrQueryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	qf := parseQueryFilter(ulrQueryParams)
	tasks, err := h.Repo.Query(c.Request.Context(), qf)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	res := taskList{Items: make([]taskGET, 0)}
	for _, t := range tasks {
		res.Items = append(res.Items, parseTaskGet(t))
	}
	c.JSON(http.StatusOK, res)
}

func parseQueryFilter(ulrQueryParams TaskFilter) model.QueryFilter {
	qf := model.QueryFilter{}
	if ulrQueryParams.Title != "" {
		qf.Title = &ulrQueryParams.Title
	}
	if ulrQueryParams.Description != "" {
		qf.Title = &ulrQueryParams.Description
	}
	if !ulrQueryParams.DueDate.IsZero() {
		qf.DueDate = &ulrQueryParams.DueDate
	}
	if !ulrQueryParams.After.IsZero() {
		qf.StartDueDate = &ulrQueryParams.After
	}
	if !ulrQueryParams.Before.IsZero() {
		qf.EndDueDate = &ulrQueryParams.Before
	}
	if ulrQueryParams.Completed != nil {
		qf.Completed = ulrQueryParams.Completed
	}
	return qf
}

func (h *Handler) Post(c *gin.Context) {
	taskpost := &taskPOST{}
	err := c.BindJSON(taskpost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := model.Parse(
		taskpost.Title,
		model.WithDescription(taskpost.Description),
		model.WithDueDate(taskpost.DueDate),
	)
	err = h.Repo.Create(c.Request.Context(), t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": t.ID})
}

func (h *Handler) PostComplete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := &model.Task{ID: id}
	err = h.Repo.Read(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if task.Completed {
		c.JSON(http.StatusAccepted, gin.H{"id": task.ID})
		return
	}
	task.Completed = true
	task.UpdatedAt = time.Now()
	err = h.Repo.Update(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": task.ID})
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := &model.Task{ID: id}
	err = h.Repo.Read(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	err = h.Repo.Delete(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) Patch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskPatch := &TaskPATCH{}
	err = c.BindJSON(taskPatch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := &model.Task{ID: id}
	err = h.Repo.Read(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// no need to check if null because of omit_empty tag
	if taskPatch.Title != "" {
		task.Title = taskPatch.Title
	}
	if taskPatch.Description != "" {
		task.Description = taskPatch.Description
	}
	if taskPatch.Completed != nil {
		task.Completed = *taskPatch.Completed
	}
	if taskPatch.DueDate != "" {
		dueDate, err := time.Parse(time.DateOnly, taskPatch.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task.DueDate = dueDate
	}
	task.UpdatedAt = time.Now()
	err = h.Repo.Update(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": task.ID})
}
