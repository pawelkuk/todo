package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	session "github.com/pawelkuk/todo/pkg/auth/model"
	"github.com/pawelkuk/todo/pkg/auth/repo"
)

type Middleware struct {
	Repo repo.Repo
}

func (m *Middleware) Handle(c *gin.Context) {
	// runtime.Breakpoint()
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	s := &session.Session{Token: session.SessionToken{Value: sessionToken}}
	err = m.Repo.Read(c.Request.Context(), s)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	if !time.Now().Before(s.Expiry) {
		c.AbortWithError(http.StatusUnauthorized, errors.New("session expired"))
		return
	}
	s.Refresh()
	err = m.Repo.Update(c.Request.Context(), s)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("could not refresh session: %w", err))
	}
	c.Set("user_id", s.UserID)
}
