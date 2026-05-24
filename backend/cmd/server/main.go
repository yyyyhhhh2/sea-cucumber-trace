package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	appcache "sea-cucumber-trace/backend/internal/cache"
	"sea-cucumber-trace/backend/internal/config"
	"sea-cucumber-trace/backend/internal/database"
	"sea-cucumber-trace/backend/internal/handler"
	"sea-cucumber-trace/backend/internal/model"
	"sea-cucumber-trace/backend/internal/repository"
	"sea-cucumber-trace/backend/internal/service"
)

func main() {
	_ = godotenv.Load()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	cfg := config.Load()
	db, err := database.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := repository.PrepareForMigration(db); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(
		&model.User{},
		&model.Org{},
		&model.SeaCucumberBatch{},
		&model.TraceEvent{},
		&model.ChainRecord{},
	); err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db)
	cacheClient, err := appcache.New(cfg)
	if err != nil {
		log.Printf("redis disabled: %v", err)
	}
	defer func() {
		if cacheClient != nil {
			_ = cacheClient.Close()
		}
	}()

	svc := service.New(cfg, repo, cacheClient)
	h := handler.New(svc)

	if err := repo.SeedIfEmpty(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(handler.CORS())

	r.GET("/api/health", h.Health)

	api := r.Group("/api")
	{
		api.POST("/auth/login", h.Login)
		api.GET("/trace/:batchNo", h.PublicTimeline)
		api.GET("/verify/:batchNo", h.Verify)

		auth := api.Group("")
		auth.Use(handler.JWTAuth(cfg))
		{
			auth.GET("/me", h.Me)
			auth.GET("/orgs", h.ListOrgs)
			auth.GET("/dashboard/summary", h.DashboardSummary)
			auth.GET("/batches", h.ListBatches)
			auth.POST("/batches", h.CreateBatch)
			auth.GET("/batches/:id", h.GetBatch)
			auth.POST("/batches/:id/events", h.AddEvent)
			auth.GET("/batches/:id/timeline", h.Timeline)
			auth.POST("/admin/import-demo", h.ImportDemoData)
		}
	}

	if staticDir := strings.TrimSpace(cfg.StaticDir); staticDir != "" {
		indexPath := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			r.Static("/assets", filepath.Join(staticDir, "assets"))
			r.NoRoute(func(c *gin.Context) {
				if strings.HasPrefix(c.Request.URL.Path, "/api/") {
					c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
					return
				}
				c.File(indexPath)
			})
		}
	}

	addr := ":" + cfg.Port
	log.Printf("listening %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
