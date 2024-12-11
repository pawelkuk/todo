package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	model "github.com/pawelkuk/todo/pkg/periodictask/model"
	repo "github.com/pawelkuk/todo/pkg/periodictask/repo"
	"github.com/robfig/cron"
)

type Handler struct {
	Repo repo.Repo
}

type periodicTaskPOST struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Schedule    string `json:"schedule" binding:"required"`
}

type periodicTaskGET struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Schedule    string `json:"schedule"`
}

type periodicTaskList struct {
	Items []periodicTaskGET `json:"items"`
}

type periodicTaskPATCH struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
}

func parsePeriodicTaskGet(task model.PeriodicTask) periodicTaskGET {
	return periodicTaskGET{
		ID:          int(task.ID),
		Title:       task.Title,
		Description: task.Description,
		Schedule:    task.Schedule,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *Handler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := &model.PeriodicTask{ID: id}
	err = h.Repo.Read(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, parsePeriodicTaskGet(*task))
}

func (h *Handler) List(c *gin.Context) {
	tasks, err := h.Repo.Query(c.Request.Context(), model.QueryFilter{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	res := periodicTaskList{}
	for _, t := range tasks {
		res.Items = append(res.Items, parsePeriodicTaskGet(t))
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Post(c *gin.Context) {
	task := &periodicTaskPOST{}
	err := c.BindJSON(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := model.Parse(
		task.Title,
		task.Schedule,
		model.WithDescription(task.Description),
	)
	err = h.Repo.Create(c.Request.Context(), t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": t.ID})
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := &model.PeriodicTask{ID: id}
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
	taskPatch := &periodicTaskPATCH{}
	err = c.BindJSON(taskPatch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := &model.PeriodicTask{ID: id}
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
	if taskPatch.Schedule != "" {
		_, err := cron.ParseStandard(taskPatch.Schedule)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task.Schedule = taskPatch.Schedule
	}
	task.UpdatedAt = time.Now()
	err = h.Repo.Update(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": task.ID})
}
