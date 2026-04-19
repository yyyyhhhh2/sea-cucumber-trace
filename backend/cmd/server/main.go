package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"sea-cucumber-trace/backend/internal/config"
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
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
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
	svc := service.New(cfg, repo)
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
			auth.GET("/batches", h.ListBatches)
			auth.POST("/batches", h.CreateBatch)
			auth.GET("/batches/:id", h.GetBatch)
			auth.POST("/batches/:id/events", h.AddEvent)
			auth.GET("/batches/:id/timeline", h.Timeline)
		}
	}

	addr := ":" + cfg.Port
	log.Printf("listening %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
