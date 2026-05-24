package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Health(c *gin.Context) {
	health := h.svc.Health(context.Background())
	statusCode := http.StatusOK
	if health.Status != "ok" {
		statusCode = http.StatusServiceUnavailable
	}
	c.JSON(statusCode, health)
}
