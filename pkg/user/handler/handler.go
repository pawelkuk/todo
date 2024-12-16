package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pawelkuk/todo/pkg/user/model"
	"github.com/pawelkuk/todo/pkg/user/repo"
)

type userGET struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type userPOST struct {
	Email    string `json:"email,required"`
	Password string `json:"password,required"`
}

func parseUserGet(u *model.User) userGET {
	return userGET{Email: u.Email.Address, ID: u.ID}
}

type Handler struct {
	Repo repo.Repo
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := &model.User{ID: int64(intID)}
	err = h.Repo.Read(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, parseUserGet(u))
}

func (h *Handler) Post(c *gin.Context) {
	userpost := &userPOST{}
	err := c.BindJSON(userpost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := model.Parse(
		userpost.Email,
		model.WithPassword(userpost.Password),
	)
	err = h.Repo.Create(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := &model.User{ID: int64(intID)}
	err = h.Repo.Read(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	err = h.Repo.Delete(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) Patch(c *gin.Context) {

}
