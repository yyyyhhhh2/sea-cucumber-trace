package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	tok, u, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tok,
		"user": gin.H{
			"id":          u.ID,
			"username":    u.Username,
			"displayName": u.DisplayName,
			"role":        u.Role,
			"orgId":       u.OrgID,
		},
	})
}

func (h *Handler) Me(c *gin.Context) {
	cl := claimsFromCtx(c)
	if cl == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	u, err := h.svc.GetUser(cl.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":          u.ID,
		"username":    u.Username,
		"displayName": u.DisplayName,
		"role":        u.Role,
		"orgId":       u.OrgID,
	})
}
