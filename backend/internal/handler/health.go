package handler

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "sea-cucumber-trace",
		"go":      runtime.Version(),
	})
}
