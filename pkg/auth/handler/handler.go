package handler

import (
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	session "github.com/pawelkuk/todo/pkg/auth/model"
	"github.com/pawelkuk/todo/pkg/auth/repo"
	user "github.com/pawelkuk/todo/pkg/user/model"
	userrepo "github.com/pawelkuk/todo/pkg/user/repo"
)

type Handler struct {
	Repo     repo.Repo
	UserRepo userrepo.Repo
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	ul := UserLogin{}
	err := c.BindJSON(&ul)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email, err := mail.ParseAddress(ul.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := h.UserRepo.Query(c.Request.Context(), user.QueryFilter{Email: email})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(users) != 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	u := users[0]
	err = u.MatchPassword(ul.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	s := session.New(u.ID)
	err = h.Repo.Create(c.Request.Context(), &s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	setSessionCookie(c, &s)
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (h *Handler) Logout(c *gin.Context) {
	userIDAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	userID := userIDAny.(int64)
	sessions, err := h.Repo.Query(c.Request.Context(), session.QueryFilter{UserID: &userID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, s := range sessions {
		err := h.Repo.Delete(c.Request.Context(), &s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	resetSessionCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func setSessionCookie(c *gin.Context, s *session.Session) {
	c.SetCookie("session_token", s.Token.Value, int(s.Expiry.Unix()), "", "", true, false)
}

func resetSessionCookie(c *gin.Context) {
	// A.k.a logout
	c.SetCookie("session_token", "", 0, "", "", true, false)
}
