package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"sea-cucumber-trace/backend/internal/config"
	"sea-cucumber-trace/backend/internal/service"
)

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		raw := strings.TrimPrefix(h, "Bearer ")
		claims, err := service.ParseToken(cfg, raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func claimsFromCtx(c *gin.Context) *service.Claims {
	v, ok := c.Get("claims")
	if !ok {
		return nil
	}
	cl, _ := v.(*service.Claims)
	return cl
}
